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

package y3_9_202107011154

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {

	upgrader.RegistUpgrader("y3.9.202107011154", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	err = dropVersionColumn(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.9.202107011154] drop set template table version column error  %s", err.Error())
		return err
	}

	err = dropSetTplVersionColumn(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.9.202107011154] drop set table set_template_version column error  %s", err.Error())
		return err
	}

	return
}
