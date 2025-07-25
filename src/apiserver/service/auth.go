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

package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"github.com/emicklei/go-restful/v3"
)

// AuthVerify verify if user is authorized.
func (s *service) AuthVerify(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pheader))
	ownerID := httpheader.GetSupplierAccount(pheader)
	rid := httpheader.GetRid(pheader)

	if auth.EnableAuthorize() == false {
		blog.Errorf("inappropriate calling, auth is disabled, rid: %s", rid)
		s.RespError(req, resp, http.StatusBadRequest, defErr.CCError(common.CCErrCommInappropriateVisitToIAM))
		return
	}

	body := metadata.AuthBathVerifyRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&body); err != nil {
		blog.Errorf("get user's resource auth verify status, but decode body failed, err: %v, rid: %s", err, rid)
		s.RespError(req, resp, http.StatusBadRequest, defErr.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	user := meta.UserInfo{
		UserName:        httpheader.GetUser(pheader),
		SupplierAccount: ownerID,
	}

	resources := make([]metadata.AuthBathVerifyResult, len(body.Resources), len(body.Resources))

	attrs := make([]meta.ResourceAttribute, 0)
	needExactAuthAttrs := make([]meta.ResourceAttribute, 0)
	needExactAuthMap := make(map[int]bool)

	for i, res := range body.Resources {
		resources[i].AuthResource = res
		attr := meta.ResourceAttribute{
			Basic: meta.Basic{
				Type:         meta.ResourceType(res.ResourceType),
				Action:       meta.Action(res.Action),
				InstanceID:   res.ResourceID,
				InstanceIDEx: res.ResourceIDEx,
			},
			SupplierAccount: ownerID,
			BusinessID:      res.BizID,
		}
		for _, item := range res.ParentLayers {
			attr.Layers = append(attr.Layers, meta.Item{Type: meta.ResourceType(item.ResourceType),
				InstanceID: item.ResourceID, InstanceIDEx: item.ResourceIDEx})
		}
		// contains exact resource info, need exact authorize
		if res.ResourceID > 0 || res.ResourceIDEx != "" || res.BizID > 0 || len(res.ParentLayers) > 0 {
			needExactAuthMap[i] = true
			needExactAuthAttrs = append(needExactAuthAttrs, attr)
		} else {
			attrs = append(attrs, attr)
		}
	}

	ctx := context.WithValue(req.Request.Context(), common.ContextRequestIDField, rid)

	if len(needExactAuthAttrs) > 0 {
		verifyResults, err := s.authorizer.AuthorizeBatch(ctx, pheader, user, needExactAuthAttrs...)
		if err != nil {
			blog.Errorf("authorize batch failed, err: %v, attrs: %+v, rid: %s", err, needExactAuthAttrs, rid)
			s.RespError(req, resp, http.StatusInternalServerError,
				defErr.CCErrorf(common.CCErrAPIGetUserResourceAuthStatusFailed))
			return
		}
		index := 0
		resourceLen := len(body.Resources)
		for i := 0; i < resourceLen; i++ {
			if needExactAuthMap[i] {
				resources[i].Passed = verifyResults[index].Authorized
				index++
			}
		}
	}

	if len(attrs) > 0 {
		verifyResults, err := s.authorizer.AuthorizeAnyBatch(ctx, pheader, user, attrs...)
		if err != nil {
			blog.Errorf("authorize any batch failed, err: %v, attrs: %+v, rid: %s", err, attrs, rid)
			s.RespError(req, resp, http.StatusInternalServerError,
				defErr.CCErrorf(common.CCErrAPIGetUserResourceAuthStatusFailed))
			return
		}
		index := 0
		resourceLen := len(body.Resources)
		for i := 0; i < resourceLen; i++ {
			if !needExactAuthMap[i] {
				resources[i].Passed = verifyResults[index].Authorized
				index++
			}
		}
	}

	resp.WriteEntity(metadata.NewSuccessResp(resources))
}

// GetAnyAuthorizedAppList TODO
func (s *service) GetAnyAuthorizedAppList(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pheader))
	rid := httpheader.GetRid(pheader)

	if auth.EnableAuthorize() == false {
		blog.Errorf("inappropriate calling, auth is disabled, rid: %s", rid)
		s.RespError(req, resp, http.StatusBadRequest, defErr.CCError(common.CCErrCommInappropriateVisitToIAM))
		return
	}

	userInfo := meta.UserInfo{
		UserName:        httpheader.GetUser(pheader),
		SupplierAccount: httpheader.GetSupplierAccount(pheader),
	}

	authInput := meta.ListAuthorizedResourcesParam{
		UserName:     httpheader.GetUser(pheader),
		ResourceType: meta.Business,
		Action:       meta.ViewBusinessResource,
	}
	authorizedResources, err := s.authorizer.ListAuthorizedResources(req.Request.Context(), pheader, authInput)
	if err != nil {
		blog.Errorf("get user: %s authorized business list failed, err: %v, rid: %s", userInfo.UserName, err, rid)
		s.RespError(req, resp, http.StatusInternalServerError,
			defErr.CCError(common.CCErrAPIGetAuthorizedAppListFromAuthFailed))
		return
	}
	input := params.SearchParams{}
	appIDList := make([]int64, 0)
	// if any Flag is false, we should parse the appIds, else we find all.
	if !authorizedResources.IsAny {
		appIDList := make([]int64, 0)
		for _, resourceID := range authorizedResources.Ids {
			bizID, err := strconv.ParseInt(resourceID, 10, 64)
			if err != nil {
				blog.Errorf("parse bizID(%s) failed, err: %v, rid: %s", bizID, err, rid)
				s.RespError(req, resp, http.StatusInternalServerError,
					defErr.CCErrorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
				return
			}
			appIDList = append(appIDList, bizID)
		}

		if len(appIDList) == 0 {
			resp.WriteEntity(metadata.NewSuccessResp(metadata.InstResult{Info: make([]mapstr.MapStr, 0)}))
			return
		}

		input = params.SearchParams{
			Condition: mapstr.MapStr{common.BKAppIDField: mapstr.MapStr{"$in": appIDList}},
		}
	}

	result, err := s.engine.CoreAPI.TopoServer().Instance().SearchApp(req.Request.Context(), userInfo.SupplierAccount,
		req.Request.Header, &input)
	if err != nil {
		blog.Errorf("get authorized business list, auth anyFlag is: %v, but get apps[%v] failed, err: %v, rid: %s",
			authorizedResources.IsAny, appIDList, err, rid)
		s.RespError(req, resp, http.StatusInternalServerError,
			defErr.CCError(common.CCErrAPIGetAuthorizedAppListFromAuthFailed))
		return
	}

	if !result.Result {
		blog.Errorf("get authorized business list,auth anyFlag is: %v, but get apps[%v] failed, err: %v, rid: %s",
			authorizedResources.IsAny, appIDList, result.ErrMsg, rid)
		s.RespError(req, resp, http.StatusBadRequest, defErr.CCError(common.CCErrAPIGetAuthorizedAppListFromAuthFailed))
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(result.Data))
}

// GetUserNoAuthSkipURL returns the url which can helps to launch the bk-auth-center when a user do not
// have the authorize to access resource(s).
func (s *service) GetUserNoAuthSkipURL(req *restful.Request, resp *restful.Response) {
	reqHeader := req.Request.Header
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(reqHeader))
	rid := httpheader.GetRid(reqHeader)

	p := new(metadata.IamPermission)
	err := json.NewDecoder(req.Request.Body).Decode(p)
	if err != nil {
		blog.Errorf("get user's skip url when no auth, but decode request failed, err: %v, rid: %s", err, rid)
		s.RespError(req, resp, http.StatusBadRequest, defErr.CCError(common.CCErrCommJSONUnmarshalFailed))
		return
	}

	url, err := s.authorizer.GetNoAuthSkipUrl(req.Request.Context(), reqHeader, p)
	if err != nil {
		blog.Errorf("get user's skip url when no auth, but request to auth center failed, err: %v, rid: %s", err, rid)
		s.RespError(req, resp, http.StatusBadRequest, defErr.CCError(common.CCErrGetNoAuthSkipURLFailed))
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(url))
}
