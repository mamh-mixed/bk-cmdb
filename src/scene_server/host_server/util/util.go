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

// Package util TODO
package util

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
)

// NewOperation TODO
func NewOperation() *operation {
	return &operation{
		op: make(map[string]interface{}),
	}
}

type operation struct {
	op map[string]interface{}
}

// Data TODO
func (o *operation) Data() map[string]interface{} {
	return o.op
}

// MapStr TODO
func (o *operation) MapStr() mapstr.MapStr {
	return mapstr.NewFromMap(o.op)
}

// WithHostID TODO
func (o *operation) WithHostID(hostID int64) *operation {
	o.op[common.BKHostIDField] = hostID
	return o
}

// WithAppID TODO
func (o *operation) WithAppID(appID int64) *operation {
	o.op[common.BKAppIDField] = appID
	return o
}

// WithOwnerID TODO
func (o *operation) WithOwnerID(ownerID string) *operation {
	o.op[common.BKOwnerIDField] = ownerID
	return o
}

// WithDefaultField TODO
func (o *operation) WithDefaultField(d int64) *operation {
	o.op[common.BKDefaultField] = d
	return o
}

// WithInstID TODO
func (o *operation) WithInstID(instID int64) *operation {
	o.op[common.BKInstIDField] = instID
	return o
}

// WithObjID TODO
func (o *operation) WithObjID(objID string) *operation {
	o.op[common.BKObjIDField] = objID
	return o
}

// WithPropertyID TODO
func (o *operation) WithPropertyID(id string) *operation {
	o.op[common.BKObjAttIDField] = id
	return o
}

// WithModuleName TODO
func (o *operation) WithModuleName(name string) *operation {
	o.op[common.BKModuleNameField] = name
	return o
}

// WithModuleIDs TODO
func (o *operation) WithModuleIDs(id []int64) *operation {
	o.op[common.BKModuleIDField] = id
	return o
}

// WithModuleID TODO
func (o *operation) WithModuleID(id int64) *operation {
	o.op[common.BKModuleIDField] = id
	return o
}

// WithAssoObjID TODO
func (o *operation) WithAssoObjID(id string) *operation {
	o.op[common.BKAsstObjIDField] = id
	return o
}

// WithAssoInstID TODO
func (o *operation) WithAssoInstID(id map[string]interface{}) *operation {
	o.op[common.BKAsstInstIDField] = id
	return o
}

// WithHostInnerIP TODO
func (o *operation) WithHostInnerIP(ip string) *operation {
	o.op[common.BKHostInnerIPField] = ip
	return o
}

// WithCloudID TODO
func (o *operation) WithCloudID(id int64) *operation {
	o.op[common.BKCloudIDField] = id
	return o
}
