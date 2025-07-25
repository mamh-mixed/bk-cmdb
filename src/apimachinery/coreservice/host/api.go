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

package host

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// TransferToInnerModule  transfer host to inner module  eg:idle module and fault module
func (h *host) TransferToInnerModule(ctx context.Context, header http.Header,
	input *metadata.TransferHostToInnerModule) ([]metadata.ExceptionResult, errors.CCErrorCoder) {

	resp := new(metadata.OperaterException)
	subPath := "/set/module/host/relation/inner/module"

	err := h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.Data, resp.CCError()
	}
	return resp.Data, nil
}

// TransferToNormalModule  transfer host to normal module
func (h *host) TransferToNormalModule(ctx context.Context, header http.Header, input *metadata.HostsModuleRelation) (
	[]metadata.ExceptionResult, errors.CCErrorCoder) {

	resp := new(metadata.OperaterException)
	subPath := "/set/module/host/relation/module"

	err := h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}
	if resp.CCError() != nil {
		return resp.Data, resp.CCError()
	}
	return resp.Data, nil
}

// RemoveFromModule 将主机从模块中移出
// 如果主机属于n+1个模块（n>0），操作之后，主机属于n个模块
// 如果主机属于1个模块, 且非空闲机模块，操作之后，主机属于空闲机模块
// 如果主机属于空闲机模块，操作失败
// 如果主机属于故障机模块，操作失败
// 如果主机不在参数指定的模块中，操作失败
func (h *host) RemoveFromModule(ctx context.Context, header http.Header,
	input *metadata.RemoveHostsFromModuleOption) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/delete/host/host_module_relations"

	err = h.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// TransferToAnotherBusiness  transfer host to other business module
func (h *host) TransferToAnotherBusiness(ctx context.Context, header http.Header,
	input *metadata.TransferHostsCrossBusinessRequest) (resp *metadata.OperaterException, err error) {
	resp = new(metadata.OperaterException)
	subPath := "/set/module/host/relation/cross/business"

	err = h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteHostFromSystem delete host from cc system
func (h *host) DeleteHostFromSystem(ctx context.Context, header http.Header,
	input *metadata.DeleteHostRequest) errors.CCErrorCoder {

	resp := new(metadata.BaseResp)
	subPath := "/delete/host"

	err := h.client.Delete().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	return resp.CCError()
}

// GetHostModuleRelation get host module relation
func (h *host) GetHostModuleRelation(ctx context.Context, header http.Header,
	input *metadata.HostModuleRelationRequest) (*metadata.HostConfigData, error) {

	resp := new(metadata.HostConfig)
	subPath := "/read/module/host/relation"

	err := h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// FindIdentifier  query host identifier
func (h *host) FindIdentifier(ctx context.Context, header http.Header, input *metadata.SearchHostIdentifierParam) (
	*metadata.SearchHostIdentifierData, error) {

	resp := new(metadata.SearchHostIdentifierResult)
	subPath := "/read/host/indentifier"

	err := h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetHostByID TODO
func (h *host) GetHostByID(ctx context.Context, header http.Header, hostID int64) (resp *metadata.HostInstanceResult,
	err error) {
	resp = new(metadata.HostInstanceResult)
	subPath := "/find/host/%d"

	err = h.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, hostID).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

// GetHosts search hosts
func (h *host) GetHosts(ctx context.Context, header http.Header, opt *metadata.QueryInput) (*metadata.HostInfo, error) {

	resp := new(metadata.GetHostsResult)
	subPath := "/findmany/hosts/search"

	err := h.client.Post().
		Body(opt).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)

	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// LockHost TODO
func (h *host) LockHost(ctx context.Context, header http.Header,
	input *metadata.HostLockRequest) (resp *metadata.HostLockResponse, err error) {
	resp = new(metadata.HostLockResponse)
	subPath := "/find/host/lock"

	err = h.client.Post().
		Body(input).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

// UnlockHost TODO
func (h *host) UnlockHost(ctx context.Context, header http.Header,
	input *metadata.HostLockRequest) (resp *metadata.HostLockResponse, err error) {
	resp = new(metadata.HostLockResponse)
	subPath := "/delete/host/lock"

	err = h.client.Delete().
		Body(input).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

// QueryHostLock TODO
func (h *host) QueryHostLock(ctx context.Context, header http.Header,
	input *metadata.QueryHostLockRequest) (resp *metadata.HostLockQueryResponse, err error) {
	resp = new(metadata.HostLockQueryResponse)
	subPath := "/findmany/host/lock/search"

	err = h.client.Post().
		Body(input).
		WithContext(ctx).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return resp, err
}

// CreateDynamicGroup is dynamic group query datas base on conditions action api machinery.
func (h *host) CreateDynamicGroup(ctx context.Context, header http.Header,
	data *metadata.DynamicGroup) (resp *metadata.IDResult, err error) {

	resp = new(metadata.IDResult)
	subPath := "/create/dynamicgroup"

	err = h.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// UpdateDynamicGroup is dynamic group update action api machinery.
func (h *host) UpdateDynamicGroup(ctx context.Context, bizID, id string,
	header http.Header, data map[string]interface{}) (resp *metadata.BaseResp, err error) {

	resp = new(metadata.BaseResp)
	subPath := "/update/dynamicgroup/%s/%s"

	err = h.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteDynamicGroup is dynamic group delete action api machinery.
func (h *host) DeleteDynamicGroup(ctx context.Context, bizID, id string,
	header http.Header) (resp *metadata.BaseResp, err error) {

	resp = new(metadata.BaseResp)
	subPath := "/delete/dynamicgroup/%s/%s"

	err = h.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetDynamicGroup is dynamic group query detail action api machinery.
func (h *host) GetDynamicGroup(ctx context.Context, bizID, id string,
	header http.Header) (resp *metadata.GetDynamicGroupResult, err error) {

	resp = new(metadata.GetDynamicGroupResult)
	subPath := "/find/dynamicgroup/%s/%s"

	err = h.client.Get().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, bizID, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// SearchDynamicGroup is dynamic group search action api machinery.
func (h *host) SearchDynamicGroup(ctx context.Context, header http.Header,
	opt *metadata.QueryCondition) (resp *metadata.SearchDynamicGroupResult, err error) {

	resp = new(metadata.SearchDynamicGroupResult)
	subPath := "/findmany/dynamicgroup/search"

	err = h.client.Post().
		WithContext(ctx).
		Body(opt).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// AddUserCustom TODO
func (h *host) AddUserCustom(ctx context.Context, user string, header http.Header,
	dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/create/usercustom/%s"

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, user).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// UpdateUserCustomByID TODO
func (h *host) UpdateUserCustomByID(ctx context.Context, user string, id string, header http.Header,
	dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/update/usercustom/%s/%s"

	err = h.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, user, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetUserCustomByUser TODO
func (h *host) GetUserCustomByUser(ctx context.Context, user string,
	header http.Header) (resp *metadata.GetUserCustomResult, err error) {
	resp = new(metadata.GetUserCustomResult)
	subPath := "/find/usercustom/user/search/%s"

	err = h.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, user).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetDefaultUserCustom TODO
func (h *host) GetDefaultUserCustom(ctx context.Context, header http.Header) (resp *metadata.GetUserCustomResult,
	err error) {
	resp = new(metadata.GetUserCustomResult)
	subPath := "/find/usercustom/default"

	err = h.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// UpdateDefaultUserCustom TODO
func (h *host) UpdateDefaultUserCustom(ctx context.Context, header http.Header,
	dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)

	err = h.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef("/update/usercustom/default").
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// AddHostFavourite TODO
func (h *host) AddHostFavourite(ctx context.Context, user string, header http.Header,
	dat *metadata.FavouriteParms) (resp *metadata.IDResult, err error) {
	resp = new(metadata.IDResult)
	subPath := "/create/hosts/favorites/%s"

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, user).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// UpdateHostFavouriteByID TODO
func (h *host) UpdateHostFavouriteByID(ctx context.Context, user string, id string, header http.Header,
	dat map[string]interface{}) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/update/hosts/favorites/%s/%s"

	err = h.client.Put().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, user, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// DeleteHostFavouriteByID TODO
func (h *host) DeleteHostFavouriteByID(ctx context.Context, user string, id string,
	header http.Header) (resp *metadata.BaseResp, err error) {
	resp = new(metadata.BaseResp)
	subPath := "/delete/hosts/favorites/%s/%s"

	err = h.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, user, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// ListHostFavourites TODO
func (h *host) ListHostFavourites(ctx context.Context, user string, header http.Header,
	dat *metadata.QueryInput) (resp *metadata.GetHostFavoriteResult, err error) {
	resp = new(metadata.GetHostFavoriteResult)
	subPath := "/findmany/hosts/favorites/search/%s"

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath, user).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetHostFavouriteByID TODO
func (h *host) GetHostFavouriteByID(ctx context.Context, user string, id string,
	header http.Header) (resp *metadata.GetHostFavoriteWithIDResult, err error) {
	resp = new(metadata.GetHostFavoriteWithIDResult)
	subPath := "/find/hosts/favorites/search/%s/%s"

	err = h.client.Get().
		WithContext(ctx).
		SubResourcef(subPath, user, id).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// GetHostModulesIDs TODO
func (h *host) GetHostModulesIDs(ctx context.Context, header http.Header,
	dat *metadata.ModuleHostConfigParams) (resp *metadata.GetHostModuleIDsResult, err error) {
	resp = new(metadata.GetHostModuleIDsResult)
	subPath := "/findmany/meta/hosts/modules/search"

	err = h.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// ListHosts get hosts list
func (h *host) ListHosts(ctx context.Context, header http.Header, option *metadata.ListHosts) (*metadata.ListHostResult,
	error) {

	result := new(metadata.ListHostResp)
	subPath := "/findmany/hosts/list_hosts"

	err := h.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(result)
	if err != nil {
		return nil, errors.CCHttpError
	}

	if err = result.CCError(); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// UpdateHostCloudAreaField update host's cloud area
func (h *host) UpdateHostCloudAreaField(ctx context.Context, header http.Header,
	option metadata.UpdateHostCloudAreaFieldOption) errors.CCErrorCoder {

	result := metadata.BaseResp{}
	subPath := "/updatemany/hosts/cloudarea_field"

	err := h.client.Put().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(&result)
	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// FindCloudAreaHostCount TODO
func (h *host) FindCloudAreaHostCount(ctx context.Context, header http.Header,
	option metadata.CloudAreaHostCount) (resp *metadata.CloudAreaHostCountResult, err error) {
	resp = new(metadata.CloudAreaHostCountResult)
	subPath := "/findmany/cloudarea/hostcount"

	err = h.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	return
}

// TransferHostResourceDirectory transfer host resource directory
func (h *host) TransferHostResourceDirectory(ctx context.Context, header http.Header,
	option *metadata.TransferHostResourceDirectory) errors.CCErrorCoder {

	result := metadata.BaseResp{}
	subPath := "/host/transfer/resource/directory"

	err := h.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(&result)
	if err != nil {
		return errors.CCHttpError
	}

	if ccErr := result.CCError(); ccErr != nil {
		return ccErr
	}

	return nil
}

// GetDistinctHostIDByTopology get distinct host id by topology relation
func (h *host) GetDistinctHostIDByTopology(ctx context.Context, header http.Header,
	input *metadata.DistinctHostIDByTopoRelationRequest) ([]int64, errors.CCErrorCoder) {

	resp := new(metadata.DistinctIDResponse)
	subPath := "/read/distinct/host_id/topology/relation"

	err := h.client.Post().
		WithContext(ctx).
		Body(input).
		SubResourcef(subPath).
		WithHeaders(header).
		Do().
		Into(resp)
	if err != nil {
		return nil, errors.CCHttpError
	}
	if err := resp.CCError(); err != nil {
		return nil, err
	}
	return resp.Data.IDArr, nil
}
