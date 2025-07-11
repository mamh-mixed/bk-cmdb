// Package metadata TODO
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
package metadata

// ProcessModule TODO
type ProcessModule struct {
	AppID      int64  `json:"bk_biz_id" bson:"bk_biz_id"`
	ModuleName string `json:"bk_module_name" bson:"bk_module_name"`
	ProcessID  int64  `json:"bk_process_id" bson:"bk_process_id"`
}

// ListProcessRelatedInfoResponse TODO
type ListProcessRelatedInfoResponse struct {
	BaseResp `json:",inline"`
	Data     struct {
		Count int                            `json:"count"`
		Info  []ListProcessRelatedInfoResult `json:"info"`
	} `json:"data"`
}
