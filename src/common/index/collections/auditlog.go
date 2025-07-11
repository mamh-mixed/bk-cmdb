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

package collections

import (
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	registerIndexes(common.BKTableNameAuditLog, commAuditLogIndexes)
}

var commAuditLogIndexes = []types.Index{
	{
		Name:                    common.CCLogicIndexNamePrefix + "ID",
		Keys:                    bson.D{{common.BKFieldID, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "operationTime",
		Keys:                    bson.D{{"operation_time", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "user",
		Keys:                    bson.D{{"user", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "resourceName",
		Keys:                    bson.D{{"resource_name", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "auditType_resourceName_action_operationTime",
		Keys:                    bson.D{{"audit_type", 1}, {"resource_type", 1}, {"action", 1}, {"operation_time", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
