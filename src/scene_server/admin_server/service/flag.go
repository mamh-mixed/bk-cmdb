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

package service

import (
	"net/http"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful/v3"
)

var (
	// 允许用户设置的key
	userConfigKeyMap = map[string]bool{
		"blueking_modify": true,
	}
	// 过期时间
	userConfigDefaultExpireHour = 6
)

// UserConfigSwitch update blueking modify flag
func (s *Service) UserConfigSwitch(req *restful.Request, resp *restful.Response) {
	kit := rest.NewKitFromHeader(req.Request.Header, s.CCErr)

	canModify := strings.ToLower(req.PathParameter("can"))
	key := req.PathParameter("key")
	blCanModify := false

	if _, ok := userConfigKeyMap[key]; !ok {
		result := &metadata.RespError{
			Msg: kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, key),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	switch canModify {
	case "true":
		blCanModify = true
	case "false":
		blCanModify = false
	default:
		result := &metadata.RespError{
			Msg: kit.CCError.Errorf(common.CCErrCommParamsNeedBool, "can"),
		}
		resp.WriteError(http.StatusBadRequest, result)
		return
	}
	cond := map[string]interface{}{
		"type": metadata.CCSystemUserConfigSwitch,
	}
	data := map[string]metadata.SysUserConfigItem{
		key: {
			Flag:     blCanModify,
			ExpireAt: time.Now().Unix() + int64(userConfigDefaultExpireHour*3600),
		},
	}

	err := s.db.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Upsert(s.ctx, cond, data)
	if err != nil {
		blog.ErrorJSON("set key %s value %s failed, err: %v, rid: %s", key, canModify, err, kit.Rid)
		resp.WriteError(http.StatusBadGateway, kit.CCError.Error(common.CCErrCommDBUpdateFailed))
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp("modify system user config success"))

}
