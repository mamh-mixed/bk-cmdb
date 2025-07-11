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

package y3_9_202011201146

import (
	"context"

	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

const (
	tableNameSubscription = "cc_Subscription"
	subscriptionNameField = "subscription_name"
)

func removeRedundantItems(ctx context.Context, db dal.RDB, conf *history.Config) error {
	subscriptionFilter := map[string]interface{}{
		subscriptionNameField: "process instance refresh [Do not remove it]",
		"system_name":         "cmdb",
		"bk_supplier_account": "0",
	}

	if err := db.Table(tableNameSubscription).Delete(ctx, subscriptionFilter); err != nil {
		blog.Errorf("remove redundant subscription failed, err: %v", err)
		return err
	}

	if err := db.DropTable(ctx, "cc_Proc2Module"); err != nil {
		blog.Errorf("drop table cc_Proc2Module failed, err: %v", err)
		return err
	}

	if err := db.DropTable(ctx, "cc_Proc2Template"); err != nil {
		blog.Errorf("drop table cc_Proc2Template failed, err: %v", err)
		return err
	}

	if err := db.DropTable(ctx, "cc_ProcInstanceModel"); err != nil {
		blog.Errorf("drop table cc_ProcInstanceModel failed, err: %v", err)
		return err
	}

	if err := db.DropTable(ctx, "cc_ProcInstanceDetail"); err != nil {
		blog.Errorf("drop table cc_ProcInstanceDetail failed, err: %v", err)
		return err
	}

	if err := db.DropTable(ctx, "cc_ProcOpTask"); err != nil {
		blog.Errorf("drop table cc_ProcOpTask failed, err: %v", err)
		return err
	}

	return nil
}
