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

// Package dal is the data access layer
package dal

import (
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/sharding"
)

// RDB rename the RDB into DB
// Compatible stock code
// Deprecated: do not use anymore.
type RDB DB

// DB db operation interface
// TODO remove this later
type DB local.DB

// Dal is the data access layer
type Dal sharding.ShardingDB
