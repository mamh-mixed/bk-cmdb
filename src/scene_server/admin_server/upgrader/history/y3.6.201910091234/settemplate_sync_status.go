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

package y3_6_201910091234

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

// SetTemplateSyncStatusMigrate TODO
func SetTemplateSyncStatusMigrate(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	tableNames := []string{"cc_SetTemplateSyncStatus", "cc_SetTemplateSyncHistory"}
	for _, tableName := range tableNames {
		existTable, err := db.HasTable(ctx, tableName)
		if err != nil {
			blog.Errorf("check table %s exist failed, err:%s", tableName, err.Error())
			return err
		}

		if !existTable {
			err := db.CreateTable(ctx, tableName)
			if err != nil {
				blog.Errorf("create table %s failed, err:%s", tableName, err.Error())
				return err
			}
		}

		setIDUnique := false
		if tableName == "cc_SetTemplateSyncStatus" {
			setIDUnique = true
		}
		taskIDUnique := false
		if tableName == "cc_SetTemplateSyncHistory" {
			taskIDUnique = true
		}
		indexArr := []types.Index{
			{
				Keys:       bson.D{{"task_id", 1}},
				Name:       "idx_taskID",
				Unique:     taskIDUnique,
				Background: true,
			},
			{
				Keys:       bson.D{{"bk_set_id", 1}},
				Name:       "idx_setID",
				Unique:     setIDUnique,
				Background: true,
			},
			{
				Keys:       bson.D{{"last_time", 1}, {"create_time", 1}},
				Name:       "idx_createLastTime",
				Background: true,
			},
			{
				Keys:       bson.D{{"status", 1}},
				Name:       "idx_status",
				Unique:     false,
				Background: true,
			},
		}

		existIndexes, err := db.Table(tableName).Indexes(ctx)
		existIndexNames := make([]string, 0)
		for _, item := range existIndexes {
			existIndexNames = append(existIndexNames, item.Name)
		}
		for _, index := range indexArr {
			if util.InStrArr(existIndexNames, index.Name) {
				continue
			}
			err := db.Table(tableName).CreateIndex(ctx, index)
			if err != nil && !db.IsDuplicatedError(err) {
				blog.ErrorJSON("add index %s for table %s failed, err:%s", index, tableName, err.Error())
				return err
			}
		}
	}

	return nil
}
