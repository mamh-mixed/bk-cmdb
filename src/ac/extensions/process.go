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

package extensions

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * process
 */

func (a *AuthManager) collectProcessesByIDs(ctx context.Context, header http.Header, ids ...int64) ([]ProcessSimplify,
	error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	ids = util.IntArrayUnique(ids)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKProcIDField).In(ids).ToMapStr(),
	}
	result, err := a.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDProc, &cond)
	if err != nil {
		blog.Errorf("get processes by id %+v failed, err: %+v, rid: %s", ids, err, rid)
		return nil, fmt.Errorf("get processes by id failed, err: %+v", err)
	}
	processes := make([]ProcessSimplify, 0)
	for _, item := range result.Info {
		process := ProcessSimplify{}
		_, err = process.Parse(item)
		if err != nil {
			blog.Errorf("collectProcessesByIDs by id %+v failed, parse process %+v failed, err: %+v, rid: %s", ids,
				item, err, rid)
			return nil, fmt.Errorf("parse process from db data failed, err: %+v", err)
		}
		processes = append(processes, process)
	}
	return processes, nil
}

// MakeResourcesByProcesses TODO
func (a *AuthManager) MakeResourcesByProcesses(header http.Header, action meta.Action, businessID int64,
	processes ...ProcessSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, process := range processes {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Process,
				Name:       process.ProcessName,
				InstanceID: process.ProcessID,
			},
			TenantID:   httpheader.GetTenantID(header),
			BusinessID: businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

// GenProcessNoPermissionResp TODO
func (a *AuthManager) GenProcessNoPermissionResp(ctx context.Context, header http.Header,
	businessID int64) (*metadata.BaseResp, error) {
	// process read authorization is skipped
	resp := metadata.NewNoPermissionResp(nil)
	return &resp, nil
}

func (a *AuthManager) extractBusinessIDFromProcesses(processes ...ProcessSimplify) (int64, error) {
	var businessID int64
	for idx, process := range processes {
		bizID := process.BKAppIDField
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("get multiple business ID from processes")
		}
		businessID = bizID
	}
	return businessID, nil
}

// AuthorizeByProcesses TODO
func (a *AuthManager) AuthorizeByProcesses(ctx context.Context, header http.Header, action meta.Action,
	processes ...ProcessSimplify) error {
	if !a.Enabled() {
		return nil
	}

	// extract business id
	bizID, err := a.extractBusinessIDFromProcesses(processes...)
	if err != nil {
		return fmt.Errorf("authorize processes failed, extract business id from processes failed, err: %+v", err)
	}

	// make auth resources
	resources := a.MakeResourcesByProcesses(header, action, bizID, processes...)

	return a.batchAuthorize(ctx, header, resources...)
}

// AuthorizeByProcessID TODO
func (a *AuthManager) AuthorizeByProcessID(ctx context.Context, header http.Header, action meta.Action,
	ids ...int64) error {
	if !a.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	processes, err := a.collectProcessesByIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("authorize processes failed, collect process by id failed, id: %+v, err: %+v", ids, err)
	}

	return a.AuthorizeByProcesses(ctx, header, action, processes...)
}
