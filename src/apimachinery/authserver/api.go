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

package authserver

import (
	"context"
	"net/http"

	"configcenter/src/ac/meta"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/types"
)

type authorizeBatchResp struct {
	metadata.BaseResp `json:",inline"`
	Data              []types.Decision `json:"data"`
}

// AuthorizeBatch TODO
func (a *authServer) AuthorizeBatch(ctx context.Context, h http.Header,
	input *types.AuthBatchOptions) ([]types.Decision, error) {
	subPath := "/authorize/batch"
	response := new(authorizeBatchResp)

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}

// AuthorizeAnyBatch TODO
func (a *authServer) AuthorizeAnyBatch(ctx context.Context, h http.Header,
	input *types.AuthBatchOptions) ([]types.Decision, error) {
	subPath := "/authorize/any/batch"
	response := new(authorizeBatchResp)

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}

// ListAuthorizedResources 获取有权限的资源列表
func (a *authServer) ListAuthorizedResources(ctx context.Context, h http.Header,
	input meta.ListAuthorizedResourcesParam) (*types.AuthorizeList, error) {
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              *types.AuthorizeList `json:"data"`
	})
	subPath := "/findmany/authorized_resource"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}

// GetNoAuthSkipUrl TODO
func (a *authServer) GetNoAuthSkipUrl(ctx context.Context, h http.Header, input *metadata.IamPermission) (string,
	error) {
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              string `json:"data"`
	})
	subPath := "/find/no_auth_skip_url"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return "", errors.CCHttpError
	}
	if response.Code != 0 {
		return "", response.CCError()
	}

	return response.Data, nil
}

// GetPermissionToApply TODO
func (a *authServer) GetPermissionToApply(ctx context.Context, h http.Header,
	input []meta.ResourceAttribute) (*metadata.IamPermission, error) {
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              *metadata.IamPermission `json:"data"`
	})
	subPath := "/find/permission_to_apply"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}

// RegisterResourceCreatorAction TODO
func (a *authServer) RegisterResourceCreatorAction(ctx context.Context, h http.Header,
	input metadata.IamInstanceWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              []metadata.IamCreatorActionPolicy `json:"data"`
	})
	subPath := "/register/resource_creator_action"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}

// BatchRegisterResourceCreatorAction TODO
func (a *authServer) BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header,
	input metadata.IamInstancesWithCreator) (
	[]metadata.IamCreatorActionPolicy, error) {
	response := new(struct {
		metadata.BaseResp `json:",inline"`
		Data              []metadata.IamCreatorActionPolicy `json:"data"`
	})
	subPath := "/register/batch_resource_creator_action"

	err := a.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(response)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if response.Code != 0 {
		return nil, response.CCError()
	}

	return response.Data, nil
}
