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

package parser

import (
	"net/http"
	"regexp"

	"configcenter/src/ac/meta"
)

var (
	updateGlobalConfigRegexp   = regexp.MustCompile(`^/api/v3/admin/update/config/global_config/[^\s/]+/?$`)
	findPlatformConfigRegexp   = regexp.MustCompile(`^/api/v3/admin/find/config/platform_config/[^\s/]+/?$`)
	updatePlatformConfigRegexp = regexp.MustCompile(`^/api/v3/admin/update/config/platform_config/[^\s/]+/?$`)
)

func (ps *parseStream) adminRelated() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.PlatformSettingConfigAuth()
	return ps
}

// PlatformSettingConfig TODO
var PlatformSettingConfig = []AuthConfig{
	{
		Name:           "findPlatformSettingConfig",
		Description:    "查询平台配置管理",
		Pattern:        "/api/v3/admin/find/config/global_config",
		HTTPMethod:     http.MethodGet,
		ResourceType:   meta.ConfigAdmin,
		ResourceAction: meta.Find,
	}, {
		Name:           "updatePlatformSettingConfig",
		Description:    "更新平台配置管理",
		Regex:          updateGlobalConfigRegexp,
		HTTPMethod:     http.MethodPut,
		ResourceType:   meta.ConfigAdmin,
		ResourceAction: meta.Update,
	}, {
		Name:           "findPlatformConfig",
		Description:    "查询平台管理配置",
		Regex:          findPlatformConfigRegexp,
		HTTPMethod:     http.MethodGet,
		ResourceType:   meta.ConfigAdmin,
		ResourceAction: meta.Find,
	}, {
		Name:           "updatePlatformConfig",
		Description:    "更新平台管理配置",
		Regex:          updatePlatformConfigRegexp,
		HTTPMethod:     http.MethodPut,
		ResourceType:   meta.ConfigAdmin,
		ResourceAction: meta.Update,
	},
}

// PlatformSettingConfigAuth platform auth
func (ps *parseStream) PlatformSettingConfigAuth() *parseStream {
	return ParseStreamWithFramework(ps, PlatformSettingConfig)

}
