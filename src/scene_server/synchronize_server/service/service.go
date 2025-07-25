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
	"net/http"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/synchronize"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/synchronize_server/app/options"
	"configcenter/src/scene_server/synchronize_server/logics"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service TODO
type Service struct {
	*options.Config
	*backbone.Engine
	disc           discovery.DiscoveryInterface
	CacheDB        redis.Client
	synchronizeSrv synchronize.SynchronizeClientInterface
}

type srvComm struct {
	header        http.Header
	rid           string
	ccErr         errors.DefaultCCErrorIf
	ccLang        language.DefaultCCLanguageIf
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	user          string
	ownerID       string
	lgc           *logics.Logics
}

// SetSynchronizeServer TODO
func (s *Service) SetSynchronizeServer(synchronizeSrv synchronize.SynchronizeClientInterface) {
	s.synchronizeSrv = synchronizeSrv
}

func (s *Service) newSrvComm(header http.Header) *srvComm {
	lang := httpheader.GetLanguage(header)
	ctx, cancel := s.Engine.CCCtx.WithCancel()
	return &srvComm{
		header:        header,
		rid:           httpheader.GetRid(header),
		ccErr:         s.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:        s.Language.CreateDefaultCCLanguageIf(lang),
		ctx:           ctx,
		ctxCancelFunc: cancel,
		user:          httpheader.GetUser(header),
		ownerID:       httpheader.GetSupplierAccount(header),
		lgc:           logics.NewLogics(s.Engine, header, s.CacheDB, s.synchronizeSrv),
	}
}

// WebService TODO
func (s *Service) WebService() *restful.Container {
	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	ws := new(restful.WebService)
	ws.Path("/synchronize/{version}").Filter(s.Engine.Metric().RestfulMiddleWare).Filter(rdapi.HTTPRequestIDFilter()).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/search").To(s.Find))
	ws.Route(ws.POST("/set/identifier/flag").To(s.SetIdentifierFlag))

	container.Add(ws)

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

	// coreservice
	coreSrv := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_CORESERVICE}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_CORESERVICE); err != nil {
		coreSrv.IsHealthy = false
		coreSrv.Message = err.Error()
	}
	meta.Items = append(meta.Items, coreSrv)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "cloud server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_SYNCHRONZESERVER,
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

// InitBackground  initialization background task
func (s *Service) InitBackground() {
	header := make(http.Header, 0)
	if httpheader.GetSupplierAccount(header) == "" {
		httpheader.SetSupplierAccount(header, common.BKSuperOwnerID)
		httpheader.SetUser(header, common.BKSynchronizeDataTaskDefaultUser)
	}

	srvData := s.newSrvComm(header)
	go srvData.lgc.TriggerSynchronize(srvData.ctx, s.Config)
}
