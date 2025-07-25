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
	registerIndexes(common.BKTableNameModelQuoteRelation, commModelQuoteRelationIndexes)
}

var commModelQuoteRelationIndexes = []types.Index{
	{
		Name: common.CCLogicIndexNamePrefix + "destModel_bkSupplierAccount",
		Keys: bson.D{
			{
				common.BKDestModelField, 1,
			},
			{
				common.BKOwnerIDField, 1,
			},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "srcModel_bkPropertyID_bkSupplierAccount",
		Keys: bson.D{
			{
				common.BKSrcModelField, 1,
			},
			{
				common.BKPropertyIDField, 1,
			},
			{
				common.BKOwnerIDField, 1,
			},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "srcModel_bkSupplierAccount",
		Keys: bson.D{
			{
				common.BKSrcModelField, 1,
			},
			{
				common.BKOwnerIDField, 1,
			},
		},
		Background: true,
	},
}
