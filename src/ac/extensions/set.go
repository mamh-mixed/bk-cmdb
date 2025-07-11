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
 * set instance
 */

func (a *AuthManager) collectSetBySetIDs(ctx context.Context, header http.Header, setIDs ...int64) ([]SetSimplify,
	error) {

	rid := util.ExtractRequestIDFromContext(ctx)

	cond := metadata.QueryCondition{
		Condition: condition.CreateCondition().Field(common.BKSetIDField).In(setIDs).ToMapStr(),
	}
	result, err := a.clientSet.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDSet, &cond)
	if err != nil {
		blog.V(3).Infof("get sets by id failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get sets by id failed, err: %+v", err)
	}

	sets := make([]SetSimplify, 0)
	for _, cls := range result.Info {
		set := SetSimplify{}
		_, err = set.Parse(cls)
		if err != nil {
			return nil, fmt.Errorf("get sets by object failed, err: %+v", err)
		}
		sets = append(sets, set)
	}
	return sets, nil
}

func (a *AuthManager) extractBusinessIDFromSets(sets ...SetSimplify) (int64, error) {
	var businessID int64
	for idx, set := range sets {
		bizID := set.BKAppIDField
		// we should ignore metadata.LabelBusinessID field not found error
		if idx > 0 && bizID != businessID {
			return 0, fmt.Errorf("authorization failed, get multiple business ID from sets")
		}
		businessID = bizID
	}
	return businessID, nil
}

// MakeResourcesBySet TODO
func (a *AuthManager) MakeResourcesBySet(header http.Header, action meta.Action, businessID int64,
	sets ...SetSimplify) []meta.ResourceAttribute {
	resources := make([]meta.ResourceAttribute, 0)
	for _, set := range sets {
		resource := meta.ResourceAttribute{
			Basic: meta.Basic{
				Action:     action,
				Type:       meta.ModelSet,
				Name:       set.BKSetNameField,
				InstanceID: set.BKSetIDField,
			},
			TenantID:   httpheader.GetTenantID(header),
			BusinessID: businessID,
		}

		resources = append(resources, resource)
	}
	return resources
}

// AuthorizeBySetID TODO
func (a *AuthManager) AuthorizeBySetID(ctx context.Context, header http.Header, action meta.Action,
	ids ...int64) error {
	if !a.Enabled() {
		return nil
	}

	if len(ids) == 0 {
		return nil
	}
	if !a.RegisterSetEnabled {
		return nil
	}

	sets, err := a.collectSetBySetIDs(ctx, header, ids...)
	if err != nil {
		return fmt.Errorf("collect set by id failed, err: %+v", err)
	}
	return a.AuthorizeBySet(ctx, header, action, sets...)
}

// AuthorizeBySet TODO
func (a *AuthManager) AuthorizeBySet(ctx context.Context, header http.Header, action meta.Action,
	sets ...SetSimplify) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	if !a.Enabled() {
		return nil
	}

	if a.SkipReadAuthorization && (action == meta.Find || action == meta.FindMany) {
		blog.V(4).Infof("skip authorization for reading, sets: %+v, rid: %s", sets, rid)
		return nil
	}
	if !a.RegisterSetEnabled {
		return nil
	}

	// extract business id
	bizID, err := a.extractBusinessIDFromSets(sets...)
	if err != nil {
		return fmt.Errorf("authorize sets failed, extract business id from sets failed, err: %+v", err)
	}

	// make auth resources
	resources := a.MakeResourcesBySet(header, action, bizID, sets...)

	return a.batchAuthorize(ctx, header, resources...)
}
