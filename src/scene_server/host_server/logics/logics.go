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

// Package logics TODO
package logics

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/common/backbone"
	"configcenter/src/storage/dal/redis"
)

// Logics TODO
type Logics struct {
	*backbone.Engine
	cache       redis.Client
	AuthManager *extensions.AuthManager
}

// NewLogics get logics handle
func NewLogics(b *backbone.Engine, cache redis.Client, authManager *extensions.AuthManager) *Logics {
	return &Logics{
		Engine:      b,
		cache:       cache,
		AuthManager: authManager,
	}
}
