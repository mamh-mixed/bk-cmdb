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

package logics

import (
	"configcenter/src/ac/iam"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/types"
)

// ListAttr list enumeration and list type attributes of instance type resource
func (lgc *Logics) ListAttr(kit *rest.Kit, resourceType iamtypes.TypeID) ([]types.AttrResource, error) {
	attrs := make([]types.AttrResource, 0)
	objID := getInstanceResourceObjID(resourceType)
	if objID == "" && !iam.IsIAMSysInstance(resourceType) {
		return attrs, nil
	}

	param := metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKPropertyTypeField: map[string]interface{}{
				common.BKDBIN: []interface{}{
					common.FieldTypeEnum,
					common.FieldTypeList,
				},
			},
		},
		Fields: []string{common.BKPropertyIDField, common.BKPropertyNameField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
	}

	var err error
	if iam.IsIAMSysInstance(resourceType) {
		obj, err := lgc.GetObjFromResourceType(kit.Ctx, kit.Header, resourceType)
		if err != nil {
			blog.ErrorJSON("get object id from resource type failed, error: %s, resource type: %s, rid: %s",
				err, resourceType, kit.Rid)
			return nil, err
		}
		objID = obj.ObjectID
	}

	res, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, &param)
	if err != nil {
		blog.ErrorJSON("read model attribute failed, error: %s, param: %s, rid: %s", err.Error(), param, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		return attrs, nil
	}

	for _, attr := range res.Info {
		displayName := attr.PropertyName
		attrs = append(attrs, types.AttrResource{
			ID:          attr.PropertyID,
			DisplayName: displayName,
		})
	}
	return attrs, nil
}
