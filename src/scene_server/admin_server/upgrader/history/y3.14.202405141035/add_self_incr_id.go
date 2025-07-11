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

package y3_14_202405141035

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func addSelfIncrID(ctx context.Context, db dal.RDB, conf *history.Config) error {
	objects := make([]metadata.Object, 0)
	err := db.Table(common.BKTableNameObjDes).Find(mapstr.MapStr{}).Fields(common.BKObjIDField).All(ctx, &objects)
	if err != nil {
		blog.Errorf("find all objects failed, err: %v", err)
		return err
	}

	ids := make([]string, 0)
	for _, object := range objects {
		ids = append(ids, GetIDRule(object.ObjectID))
	}
	ids = append(ids, GetIDRule(common.GlobalIDRule))

	cond := mapstr.MapStr{common.BKFieldDBID: mapstr.MapStr{common.BKDBIN: ids}}
	data := make([]map[string]interface{}, 0)
	err = db.Table(common.BKTableNameIDgenerator).Find(cond).Fields(common.BKFieldDBID).All(ctx, &data)
	if err != nil {
		blog.Errorf("find id generator data failed, cond: %+v, err: %v", cond, err)
		return err
	}

	dbIDMap := make(map[string]struct{})
	for _, val := range data {
		dbIDMap[util.GetStrByInterface(val[common.BKFieldDBID])] = struct{}{}
	}

	needAddIDs := make([]map[string]interface{}, 0)
	now := time.Now()
	for _, id := range ids {
		if _, ok := dbIDMap[id]; ok {
			continue
		}

		addID := map[string]interface{}{
			common.BKFieldDBID:     id,
			common.BKFieldSeqID:    0,
			common.CreateTimeField: now,
			common.LastTimeField:   now,
		}
		needAddIDs = append(needAddIDs, addID)
	}

	if len(needAddIDs) == 0 {
		return nil
	}

	if err = db.Table(common.BKTableNameIDgenerator).Insert(ctx, needAddIDs); err != nil {
		blog.Errorf("add id generator data failed, data: %+v, err: %v", needAddIDs, err)
		return err
	}

	return nil
}

// GetIDRule 获取对应id rule自增id的唯一标识，目前bk_obj_id唯一，后续涉及到多租户，可能需要调整
func GetIDRule(flag string) string {
	return fmt.Sprintf("%s%s", common.IDRulePrefix, flag)
}
