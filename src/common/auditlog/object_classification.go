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

type objectClsAuditLog struct {
	audit
}

// GenerateAuditLog generate audit of model classification, if data is nil, will auto get current model classification data by id.
func (h *objectClsAuditLog) GenerateAuditLog(parameter *generateAuditCommonParameter, id int64,
	data *metadata.Classification) (*metadata.AuditLog, error) {

	kit := parameter.kit
	if data == nil {
		// get current model classification data by id.
		query := mapstr.MapStr{metadata.ClassificationFieldID: id}
		rsp, err := h.clientSet.Model().ReadModelClassification(kit.Ctx, kit.Header,
			&metadata.QueryCondition{Condition: query})
		if err != nil {
			blog.Errorf("generate audit log of model classification failed, failed to read model classification,"+
				" err: %v, rid: %s", err.Error(), kit.Rid)
			return nil, err
		}

		if len(rsp.Info) <= 0 {
			blog.Errorf("generate audit log of model classification failed, failed to read model classification,"+
				" err: %s, rid: %s", kit.CCError.CCError(common.CCErrorModelNotFound), kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorModelNotFound)
		}

		data = &rsp.Info[0]
	}

	dataMap, err := mapstruct.Struct2Map(data)
	if err != nil {
		blog.Errorf("convert model classification(%+v) to map failed, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	return &metadata.AuditLog{
		AuditType:    metadata.ModelType,
		ResourceType: metadata.ModelGroupRes,
		Action:       parameter.action,
		ResourceID:   id,
		ResourceName: data.ClassificationName,
		OperateFrom:  parameter.operateFrom,
		OperationDetail: &metadata.BasicOpDetail{
			Details: parameter.NewBasicContent(dataMap),
		},
	}, nil
}

// NewObjectClsAuditLog TODO
func NewObjectClsAuditLog(clientSet coreservice.CoreServiceClientInterface) *objectClsAuditLog {
	return &objectClsAuditLog{
		audit: audit{
			clientSet: clientSet,
		},
	}
}
