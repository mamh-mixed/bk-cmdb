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

package y3_13_202407191507

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func init() {
	upgrader.RegistUpgrader("y3.13.202407191507", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *upgrader.Config) (err error) {
	blog.Infof("start execute y3.13.202407191507")

	if err = upgradeContainer(ctx, db, conf); err != nil {
		blog.Errorf("upgrade y3.13.202407191507 upgrade container table failed, err: %v", err)
		return err
	}

	if err = upgradeDelArchive(ctx, db, conf); err != nil {
		blog.Errorf("upgrade y3.13.202407191507 upgrade del archive table failed, err: %v", err)
		return err
	}

	blog.Infof("upgrade y3.13.202407191507 upgrade container and del archive table success")
	return nil
}
