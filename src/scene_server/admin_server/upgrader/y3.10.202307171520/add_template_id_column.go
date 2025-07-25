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

package y3_10_202307171520

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addTemplateIDColumnAndIndex(ctx context.Context, db dal.RDB) error {
	collections := []string{common.BKTableNameObjAttDes, common.BKTableNameObjUnique}
	index := types.Index{
		Name: common.CCLogicIndexNamePrefix + "bkTemplateID_bkSupplierAccount",
		Keys: bson.D{
			{
				common.BKTemplateID, 1,
			},
			{
				common.BKOwnerIDField, 1,
			},
		},
		Background: true,
	}

	for _, collection := range collections {
		if err := db.Table(collection).AddColumn(ctx, common.BKTemplateID, 0); err != nil {
			blog.Errorf("add %s column to table %s failed, err: %v", common.BKTemplateID, collection, err)
			return err
		}

		if err := addIndexIfNotExist(ctx, db, collection, []types.Index{index}); err != nil {
			blog.Errorf("add index failed, table: %s, index: %+v, err: %v", collection, index, err)
			return err
		}
	}

	return nil
}
