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

// Package openapi TODO
package openapi

import (
	"context"
	"net/http"

	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// GetProcessPortByApplicationID TODO
func (p *openapi) GetProcessPortByApplicationID(ctx context.Context, appID string, h http.Header,
	dat []mapstr.MapStr) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/GetProcessPortByApplicationID/%s"

	err = p.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, appID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// GetProcessPortByIP TODO
func (p *openapi) GetProcessPortByIP(ctx context.Context, h http.Header,
	dat map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/openapi/GetProcessPortByIP"

	err = p.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
