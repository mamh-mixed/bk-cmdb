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

package y3_9_202011251014

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func init() {
	history.RegistUpgrader("y3.9.202011251014", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *history.Config) (err error) {
	blog.Infof("start execute y3.9.202011251014")

	if err = addProcBindInfo(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.9.202011251014] change process bind attr, error  %s", err.Error())
		return err
	}

	if err = migrateProcTempBindInfo(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.9.202011251014] migrate process template bind info, error  %s", err.Error())
		return err
	}

	if err = migrateProcBindInfo(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.9.202011251014] migrate process bind info, error  %s", err.Error())
		return err
	}

	// upgrate data 之后才可以删除
	if err = clearProcAttrAndGroup(ctx, db, conf); err != nil {
		blog.Errorf("[upgrade y3.9.202011251014] clean process bind attr, error  %s", err.Error())
		return err
	}

	return nil
}
