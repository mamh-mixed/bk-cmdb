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

package auditlog

import (
	"configcenter/src/apimachinery/coreservice"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
)

type cloudAccountAuditLog struct {
	audit
}

// GenerateAuditLog generate audit log of cloud account, if data is nil, will auto get data by accountID.
func (h *cloudAccountAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, accountID int64,
	data *metadata.CloudAccount) (*metadata.AuditLog, error) {
	kit := parameter.kit

	if data == nil {
		// get data by accountID.
		cond := metadata.SearchCloudOption{
			Condition: mapstr.MapStr{common.BKCloudAccountID: accountID},
		}

		res, err := h.clientSet.Cloud().SearchAccount(kit.Ctx, kit.Header, &cond)
		if err != nil {
			blog.Errorf("generate audit log of cloud account, failed to read cloud account, err: %v, rid: %s",
				err.Error(), kit.Rid)
			return nil, err
		}
		if len(res.Info) <= 0 {
			blog.Errorf("generate audit log of cloud account failed, not find cloud account, rid: %s",
				kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCloudAccountIDNoExistFail)
		}

		data = &res.Info[0].CloudAccount
	}

	secretKey := data.SecretKey
	data.SecretKey = ""

	dataMap, err := mapstruct.Struct2Map(data)
	if err != nil {
		blog.Errorf("convert cloud account(%+v) to map failed, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	data.SecretKey = secretKey

	return &metadata.AuditLog{
		AuditType:    metadata.CloudResourceType,
		ResourceType: metadata.CloudAccountRes,
		Action:       parameter.action,
		ResourceID:   accountID,
		ResourceName: data.AccountName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: parameter.NewBasicContent(dataMap),
		},
	}, nil
}

// NewCloudAccountAuditLog TODO
func NewCloudAccountAuditLog(clientSet coreservice.CoreServiceClientInterface) *cloudAccountAuditLog {
	return &cloudAccountAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
