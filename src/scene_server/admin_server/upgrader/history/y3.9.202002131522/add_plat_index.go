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

package y3_9_202002131522

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func addPlatIndex(ctx context.Context, db dal.RDB, conf *history.Config) error {
	tableName := common.BKTableNameBasePlat
	index := types.Index{
		Keys:       bson.D{{"bk_vpc_id", 1}},
		Name:       "vpcID",
		Unique:     false,
		Background: true,
	}

	existIndexes, err := db.Table(tableName).Indexes(ctx)
	if err != nil {
		blog.ErrorJSON("get exist indexes for table %s failed, err:%s", tableName, err)
		return err
	}
	existIndexNames := make([]string, 0)
	for _, item := range existIndexes {
		existIndexNames = append(existIndexNames, item.Name)
	}

	if util.InStrArr(existIndexNames, index.Name) {
		return nil
	}

	err = db.Table(tableName).CreateIndex(ctx, index)
	if err != nil && !db.IsDuplicatedError(err) {
		blog.ErrorJSON("add index %s for table %s failed, err:%s", index, tableName, err)
		return err
	}

	return nil
}
