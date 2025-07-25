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
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful/v3"
)

// SearchReportSummary TODO
func (s *Service) SearchReportSummary(req *restful.Request, resp *restful.Response) {
	var err error
	pHeader := req.Request.Header
	rid := httpheader.GetRid(pHeader)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pHeader))

	param := metadata.ParamSearchNetcollectReport{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] decode body failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result, err := s.logics.SearchReportSummary(pHeader, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReportSummary failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusInternalServerError,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetReportSearchFail)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(result))
	return

}

// SearchReport TODO
func (s *Service) SearchReport(req *restful.Request, resp *restful.Response) {
	var err error
	pHeader := req.Request.Header
	rid := httpheader.GetRid(pHeader)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pHeader))

	param := metadata.ParamSearchNetcollectReport{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] decode body failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	count, result, err := s.logics.SearchReport(pHeader, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReportSummary failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusOK,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetReportSearchFail)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(metadata.RspNetcollectReport{
		Count: count,
		Info:  result,
	}))
	return
}

// ConfirmReport TODO
func (s *Service) ConfirmReport(req *restful.Request, resp *restful.Response) {
	var err error
	pHeader := req.Request.Header
	rid := httpheader.GetRid(pHeader)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pHeader))

	param := metadata.ParamNetcollectComfirm{}
	if err = json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][ConfirmReport] decode body failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result := s.logics.ConfirmReport(pHeader, param.Reports)
	if len(result.Errors) > 0 {
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:  defErr.Error(common.CCErrCollectNetReportConfirmFail),
			Data: result,
		})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(result))
	return

}

// SearchHistory TODO
func (s *Service) SearchHistory(req *restful.Request, resp *restful.Response) {
	var err error
	pHeader := req.Request.Header
	rid := httpheader.GetRid(pHeader)
	defErr := s.engine.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(pHeader))

	param := metadata.ParamSearchNetcollectReport{}
	if err := json.NewDecoder(req.Request.Body).Decode(&param); err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] decode body failed, err: %v, rid: %s", err, rid)
		_ = resp.WriteError(http.StatusBadRequest,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	count, result, err := s.logics.SearchHistory(pHeader, param)
	if err != nil {
		blog.Errorf("[NetDevice][SearchReportSummary] SearchReportSummary failed, err: %v", err)
		_ = resp.WriteError(http.StatusInternalServerError,
			&metadata.RespError{Msg: defErr.Error(common.CCErrCollectNetHistorySearchFail)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(metadata.RspNetcollectHistory{
		Count: count,
		Info:  result,
	}))
	return
}
