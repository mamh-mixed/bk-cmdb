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

package x19_09_03_04

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

// ChangeSetPropertyName TODO
func ChangeSetPropertyName(ctx context.Context, db dal.RDB, conf *history.Config) error {
	filter := map[string]interface{}{
		common.BKPropertyIDField: "bk_set_name",
		common.BKObjIDField:      common.BKInnerObjIDSet,
	}
	doc := map[string]interface{}{
		"bk_property_name": "集群名",
	}
	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, filter, doc); err != nil {
		blog.Errorf("ChangeSetPropertyName failed, err: %+v", err)
		return fmt.Errorf("ChangeSetPropertyName failed, err: %v", err)
	}
	return nil
}
