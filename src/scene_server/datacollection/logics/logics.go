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
	"context"

	"configcenter/src/common/backbone"
	"configcenter/src/storage/dal"
	"configcenter/src/thirdparty/esbserver"
)

// Logics framwork need
type Logics struct {
	*backbone.Engine
	db  dal.RDB
	ESB esbserver.EsbClientInterface
	ctx context.Context
}

// NewLogics TODO
func NewLogics(ctx context.Context, engine *backbone.Engine, mgoCli dal.RDB, esb esbserver.EsbClientInterface) *Logics {
	return &Logics{ctx: ctx, db: mgoCli, Engine: engine, ESB: esb}
}
