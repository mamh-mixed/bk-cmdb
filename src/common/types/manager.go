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

package types

import "configcenter/src/common/auth"

var (
	// needDiscoveryServiceName 服务依赖的第三方服务名字的配置
	needDiscoveryServiceName = make(map[string]struct{}, 0)
)

// DiscoveryAllService 发现所有定义的服务
func DiscoveryAllService() {
	for name := range AllModule {
		needDiscoveryServiceName[name] = struct{}{}
	}
}

// AddDiscoveryService 新加需要发现服务的名字
func AddDiscoveryService(name ...string) {
	for _, name := range name {
		needDiscoveryServiceName[name] = struct{}{}
	}
}

// GetDiscoveryService TODO
func GetDiscoveryService() map[string]struct{} {
	// compatible 如果没配置,发现所有的服务
	if len(needDiscoveryServiceName) == 0 {
		DiscoveryAllService()
	}
	// 如果没有开启鉴权，则不需要发现auth节点
	if !auth.EnableAuthorize() {
		delete(needDiscoveryServiceName, CC_MODULE_AUTH)
	}
	return needDiscoveryServiceName
}
