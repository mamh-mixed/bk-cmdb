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

	// 先注册未规范化的索引，如果索引出现冲突旧，删除未规范化的索引
	registerIndexes(common.BKTableNameObjDes, deprecatedObjDesIndexes)
	registerIndexes(common.BKTableNameObjDes, commObjDesIndexes)

}

//  新加和修改后的索引,索引名字一定要用对应的前缀，CCLogicUniqueIdxNamePrefix|common.CCLogicIndexNamePrefix

var commObjDesIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkObjID",
		Keys: bson.D{{
			common.BKObjIDField, 1},
		},
		Unique: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKObjIDField: map[string]string{common.BKDBType: "string"},
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkObjName",
		Keys: bson.D{{
			common.BKObjNameField, 1},
		},
		Unique: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKObjNameField: map[string]string{common.BKDBType: "string"},
		},
	},
}

// deprecated 未规范化前的索引，只允许删除不允许新加和修改，
var deprecatedObjDesIndexes = []types.Index{
	{
		Name: "bk_classification_id_1",
		Keys: bson.D{{
			"bk_classification_id", 1},
		},
		Background: true,
	},
	{
		Name: "bk_supplier_account_1",
		Keys: bson.D{{
			"bk_supplier_account", 1},
		},
		Background: true,
	},
	{
		Name: "idx_unique_id",
		Keys: bson.D{{
			"id", 1},
		},
		Unique:     true,
		Background: true,
	},
}
