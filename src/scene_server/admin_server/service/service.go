/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package service TODO
package service

import (
	"context"
	"fmt"

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/configures"
	"configcenter/src/scene_server/admin_server/logics"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/apigw"
	"configcenter/src/thirdparty/dataid"
	"configcenter/src/thirdparty/logplatform/opentelemetry"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"github.com/emicklei/go-restful/v3"
)

// Service TODO
type Service struct {
	*backbone.Engine
	*logics.Logics
	db           dal.RDB
	watchDB      dal.RDB
	cache        redis.Client
	ctx          context.Context
	Config       options.Config
	iam          *iam.IAM
	ConfigCenter *configures.ConfCenter
	GseClient    dataid.DataIDInterface
}

// NewService TODO
func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

// SetDB TODO
func (s *Service) SetDB(db dal.RDB) {
	s.db = db
}

// SetWatchDB TODO
func (s *Service) SetWatchDB(watchDB dal.RDB) {
	s.watchDB = watchDB
}

// SetCache TODO
func (s *Service) SetCache(cache redis.Client) {
	s.cache = cache
}

// SetIam TODO
func (s *Service) SetIam(iam *iam.IAM) {
	s.iam = iam
}

// WebService TODO
func (s *Service) WebService() *restful.Container {
	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}
	api.Path("/migrate/v3")
	api.Filter(s.Engine.Metric().RestfulMiddleWare)
	api.Filter(rdapi.AllGlobalFilter(getErrFunc))
	api.Produces(restful.MIME_JSON)

	api.Route(api.POST("/authcenter/init").To(s.InitAuthCenter))
	api.Route(api.POST("/authcenter/register").To(s.RegisterAuthAccount))
	api.Route(api.POST("/migrate/{distribution}/{ownerID}").To(s.migrate))
	api.Route(api.POST("/migrate/system/hostcrossbiz/{ownerID}").To(s.SetSystemConfiguration))
	api.Route(api.POST("/migrate/system/user_config/{key}/{can}").To(s.UserConfigSwitch))
	api.Route(api.GET("/find/system/config_admin").To(s.SearchConfigAdmin))
	api.Route(api.PUT("/update/system/config_admin").To(s.UpdateConfigAdmin))

	api.Route(api.PUT("/update/system_config/platform_setting").To(s.UpdatePlatformSettingConfig))
	api.Route(api.GET("/find/system_config/platform_setting/{type}").To(s.SearchPlatformSettingConfig))

	api.Route(api.POST("/migrate/specify/version/{distribution}/{ownerID}").To(s.migrateSpecifyVersion))
	api.Route(api.POST("/migrate/config/refresh").To(s.refreshConfig))
	api.Route(api.POST("/migrate/dataid").To(s.migrateDataID))
	api.Route(api.POST("/migrate/old/dataid").To(s.migrateOldDataID))
	api.Route(api.POST("/delete/auditlog").To(s.DeleteAuditLog))
	api.Route(api.POST("/migrate/sync/db/index").To(s.RunSyncDBIndex))
	api.Route(api.GET("/healthz").To(s.Healthz))
	api.Route(api.GET("/monitor_healthz").To(s.MonitorHealth))

	container.Add(api)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

// Healthz TODO
func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// mongodb
	healthItem := metric.NewHealthItem(types.CCFunctionalityMongo, s.db.Ping())
	meta.Items = append(meta.Items, healthItem)

	// redis
	redisItem := metric.NewHealthItem(types.CCFunctionalityRedis, s.cache.Ping(context.Background()).Err())
	meta.Items = append(meta.Items, redisItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "admin server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_MIGRATE,
		HealthMeta: meta,
		AtTime:     metadata.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	answer.SetCommonResponse()
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}

// MonitorHealth TODO
func (s *Service) MonitorHealth(req *restful.Request, resp *restful.Response) {
	rid := util.GenerateRID()
	alam := &meta.Alarm{
		RequestID: rid,
		Type:      meta.EventTestInfo,
		Detail:    fmt.Sprintf("test event link connectivity"),
		Module:    types.CC_MODULE_MIGRATE,
	}
	monitor.Collect(alam)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(metadata.NewSuccessResp(alam))

}

// InitGseClient init gse apiGW client
func (s *Service) InitGseClient() error {
	switch s.Config.DataIdMigrateWay {
	case options.MigrateWayESB, "":
		s.GseClient = esb.EsbClient().GseSrv()
		return nil

	case options.MigrateWayApiGW:
		err := apigwcli.Init("apiGW", s.Engine.Metric().Registry(), []apigw.ClientType{apigw.Gse})
		if err != nil {
			blog.Errorf("init gse api gateway client failed, err: %v", err)
			return err
		}
		s.GseClient = apigwcli.Client().Gse()
		return nil

	default:
		return fmt.Errorf("init gse client failed, unknow migrate dataid way")
	}
}
