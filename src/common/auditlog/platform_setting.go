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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// PlatFormSettingAuditLog is audit log handler for platform.
type PlatFormSettingAuditLog struct {
	// audit base audit handler.
	audit
}

// NewPlatFormSettingAuditLog creates a new platform object.
func NewPlatFormSettingAuditLog(clientSet coreservice.CoreServiceClientInterface) *PlatFormSettingAuditLog {
	return &PlatFormSettingAuditLog{audit: audit{clientSet: clientSet}}
}

// GenerateAuditLog generates an audit log for platform operations.
func (l *PlatFormSettingAuditLog) GenerateAuditLog(param *generateAuditCommonParameter,
	oldConf interface{}, operationType string) ([]metadata.AuditLog, error) {

	if oldConf == nil {
		return make([]metadata.AuditLog, 0), nil
	}

	content := mapstr.MapStr{
		operationType: oldConf,
	}
	content[common.LastTimeField] = metadata.Now()
	logs := []metadata.AuditLog{{
		AuditType:       metadata.PlatFormSettingType,
		ResourceType:    metadata.PlatFormSettingRes,
		Action:          param.action,
		ResourceID:      common.ConfigAdminID,
		ResourceName:    common.ConfigAdminValueField,
		OperationDetail: &metadata.BasicOpDetail{Details: param.NewBasicContent(content)},
	}}

	return logs, nil
}
