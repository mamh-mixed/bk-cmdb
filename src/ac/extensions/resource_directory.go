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
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
 * module instance
 */

func (a *AuthManager) collectResourceDirectoryByDirectoryIDs(ctx context.Context, header http.Header,
	directoryIDs ...int64) ([]ModuleSimplify, error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	// unique ids so that we can be aware of invalid id if query result length not equal ids's length
	directoryIDs = util.IntArrayUnique(directoryIDs)

	cond := metadata.QueryCondition{
		Condition: map[string]interface{}{common.BKModuleIDField: map[string]interface{}{common.BKDBIN: directoryIDs}},
	}
	result, err := a.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDModule, &cond)
	if err != nil {
		blog.V(3).Infof("get directory by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get directory by id failed, err: %+v", err)
	}

	directoryArr := make([]ModuleSimplify, 0)
	for _, cls := range result.Info {
		directory := ModuleSimplify{}
		_, err = directory.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("parse directory failed, err: %+v", err)
		}
		directoryArr = append(directoryArr, directory)
	}
	return directoryArr, nil
}

// MakeResourcesByResourceDirectory TODO
func (a *AuthManager) MakeResourcesByResourceDirectory(header http.Header, action meta.Action,
	directoryArr ...ModuleSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, directory := range directoryArr {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ResourcePoolDirectory,
				Name:       directory.BKModuleNameField,
				InstanceID: directory.BKModuleIDField,
			},
			TenantID: httpheader.GetTenantID(header),
		}

		resources = append(resources, resource)
	}
	return resources
}

// AuthorizeByResourceDirectoryID TODO
func (a *AuthManager) AuthorizeByResourceDirectoryID(ctx context.Context, header http.Header, action meta.Action,
	ids ...int64) error {
	if !a.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	directoryArr, err := a.collectResourceDirectoryByDirectoryIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("update registered directoryArr failed, get directoryArr by id failed, err: %+v", err)
	}
	return a.AuthorizeByResourceDirectory(ctx, header, action, directoryArr...)
}

// AuthorizeByResourceDirectory TODO
func (a *AuthManager) AuthorizeByResourceDirectory(ctx context.Context, header http.Header, action meta.Action,
	directoryArr ...ModuleSimplify) error {
	if !a.Enabled() {
		return nil
	}

	if len(directoryArr) == 0 {
		return nil
	}

	// make auth resources
	resources := a.MakeResourcesByResourceDirectory(header, action, directoryArr...)

	return a.batchAuthorize(ctx, header, resources...)
}
