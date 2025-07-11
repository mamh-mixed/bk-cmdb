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
	registerIndexes(common.BKTableNameBaseApp, commApplicationBaseIndexes)
}

var commApplicationBaseIndexes = []types.Index{
	{
		Name:                    common.CCLogicIndexNamePrefix + "default",
		Keys:                    bson.D{{common.BKDefaultField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizName",
		Keys:                    bson.D{{common.BKAppNameField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: map[string]interface{}{common.BKAppNameField: map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkBizID",
		Keys:                    bson.D{{common.BKAppIDField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
