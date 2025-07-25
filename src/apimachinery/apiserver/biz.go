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

package apiserver

import (
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CreateBiz TODO
func (a *apiServer) CreateBiz(ctx context.Context, ownerID string, h http.Header,
	params map[string]interface{}) (resp *metadata.CreateInstResult, err error) {
	resp = new(metadata.CreateInstResult)
	subPath := "/biz/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(params).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// UpdateBiz TODO
func (a *apiServer) UpdateBiz(ctx context.Context, ownerID string, bizID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/biz/%s/%s"
	err = a.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, ownerID, bizID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// UpdateBizDataStatus update biz data status
func (a *apiServer) UpdateBizDataStatus(ctx context.Context, ownerID string, flag common.DataStatusFlag, bizID int64,
	h http.Header) errors.CCErrorCoder {

	resp := new(metadata.Response)
	subPath := "/biz/status/%s/%s/%d"

	err := a.client.Put().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, flag, ownerID, bizID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.CCError()
	}

	return nil
}

// SearchBiz TODO
func (a *apiServer) SearchBiz(ctx context.Context, ownerID string, h http.Header,
	param *metadata.QueryBusinessRequest) (resp *metadata.SearchInstResult, err error) {
	resp = new(metadata.SearchInstResult)
	subPath := "/biz/search/%s"
	err = a.client.Post().
		WithContext(ctx).
		Body(param).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// UpdateBizPropertyBatch batch update business properties
func (a *apiServer) UpdateBizPropertyBatch(ctx context.Context, h http.Header,
	param metadata.UpdateBizPropertyBatchParameter) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/updatemany/biz/property"
	err = a.client.Put().
		WithContext(ctx).
		Body(param).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)
	return
}

// DeleteBiz delete archived businesses
func (a *apiServer) DeleteBiz(ctx context.Context, h http.Header, param metadata.DeleteBizParam) error {
	resp := new(metadata.Response)
	subPath := "/deletemany/biz"

	err := a.client.Post().
		WithContext(ctx).
		Body(param).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		IntoCmdbResp(resp)

	if err != nil {
		return errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.CCError()
	}
	return nil
}
