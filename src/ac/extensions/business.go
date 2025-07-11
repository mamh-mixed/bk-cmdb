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
 * business related auth interface
 */

func (a *AuthManager) collectBusinessByIDs(ctx context.Context, header http.Header, businessIDs ...int64) (
	[]BusinessSimplify, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	businessIDs = util.IntArrayUnique(businessIDs)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKAppIDField).In(businessIDs).ToMapStr(),
	}
	result, err := a.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDApp, &cond)
	if err != nil {
		blog.V(3).Infof("get businesses by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get businesses by id failed, err: %+v", err)
	}
	blog.V(5).Infof("get businesses by id result: %+v", result)
	instances := make([]BusinessSimplify, 0)
	for _, cls := range result.Info {
		instance := BusinessSimplify{}
		_, err = instance.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("parse business from db data failed, err: %+v", err)
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

// MakeResourcesByBusiness TODO
func (a *AuthManager) MakeResourcesByBusiness(header http.Header, action meta.Action,
	businesses ...BusinessSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, business := range businesses {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.Business,
				Name:       business.BKAppNameField,
				InstanceID: business.BKAppIDField,
			},
			TenantID: httpheader.GetTenantID(header),
		}

		resources = append(resources, resource)
	}
	return resources
}

// AuthorizeByBusiness authorize by business
func (a *AuthManager) AuthorizeByBusiness(ctx context.Context, header http.Header, action meta.Action,
	businesses ...BusinessSimplify) error {

	if !a.Enabled() {
		return nil
	}

	resourcePoolBusinessID, err := a.getResourcePoolBusinessID(ctx, header)
	if err != nil {
		return err
	}

	bizArr := make([]BusinessSimplify, 0)
	if action == meta.ViewBusinessResource {
		for _, biz := range businesses {
			if biz.BKAppIDField == resourcePoolBusinessID {
				continue
			}
			bizArr = append(bizArr, biz)
		}
	} else {
		bizArr = businesses
	}

	// make auth resources
	resources := a.MakeResourcesByBusiness(header, action, bizArr...)

	return a.batchAuthorize(ctx, header, resources...)
}

// AuthorizeByBusinessID TODO
func (a *AuthManager) AuthorizeByBusinessID(ctx context.Context, header http.Header, action meta.Action,
	businessIDs ...int64) error {
	if !a.Enabled() {
		return nil
	}

	businesses, err := a.collectBusinessByIDs(ctx, header, businessIDs...)
	if err != nil {
		return fmt.Errorf("authorize businesses failed, get business by id failed, err: %+v", err)
	}

	return a.AuthorizeByBusiness(ctx, header, action, businesses...)
}

// GenBizBatchNoPermissionResp TODO
func (a *AuthManager) GenBizBatchNoPermissionResp(ctx context.Context, header http.Header, action meta.Action,
	bizIDs []int64) (*metadata.BaseResp, error) {
	businesses, err := a.collectBusinessByIDs(ctx, header, bizIDs...)
	if err != nil {
		return nil, err
	}

	// make auth resources
	resources := a.MakeResourcesByBusiness(header, action, businesses...)

	rid := util.ExtractRequestIDFromContext(ctx)
	permission, err := a.Authorizer.GetPermissionToApply(ctx, header, resources)
	if err != nil {
		blog.Errorf("get permission to apply failed, err: %v, rid: %s", err, rid)
		return nil, err
	}
	resp := metadata.NewNoPermissionResp(permission)
	return &resp, nil
}
