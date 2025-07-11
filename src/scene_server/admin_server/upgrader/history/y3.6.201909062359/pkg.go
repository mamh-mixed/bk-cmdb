// Package y3_6_201909062359 TODO
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
package y3_6_201909062359

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func init() {
	history.RegistUpgrader("y3.6.201909062359", upgrade)
}

func upgrade(ctx context.Context, db dal.RDB, conf *history.Config) (err error) {
	err = createSetTemplateTables(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.6.201909062359] create tables failed, err: %s", err.Error())
		return err
	}

	err = addSetProperty(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.6.201909062359] addSetProperty failed, err: %s", err.Error())
		return err
	}

	err = addModuleProperty(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.6.201909062359] addModuleProperty failed, err: %s", err.Error())
		return err
	}

	err = initializeSetTemplateField(ctx, db, conf)
	if err != nil {
		blog.Errorf("[upgrade y3.6.201909062359] initializeSetTemplateField failed, err: %s", err.Error())
		return err
	}

	return nil
}
