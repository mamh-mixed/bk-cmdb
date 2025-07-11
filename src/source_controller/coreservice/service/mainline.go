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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
)

// SearchMainlineModelTopo TODO
func (s *coreService) SearchMainlineModelTopo(ctx *rest.Contexts) {
	detail := struct {
		WithDetail bool `json:"with_detail"`
	}{}
	if err := ctx.DecodeInto(&detail); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.TopoOperation().SearchMainlineModelTopo(ctx.Kit, detail.WithDetail)
	if err != nil {
		blog.Errorf("search mainline model topo failed, %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoMainlineSelectFailed))
		return
	}
	ctx.RespEntity(result)
}

// SearchMainlineInstanceTopo TODO
func (s *coreService) SearchMainlineInstanceTopo(ctx *rest.Contexts) {
	bkBizID := ctx.Request.PathParameter(common.BKAppIDField)
	if len(bkBizID) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField))
		return
	}
	bizID, err := strconv.ParseInt(bkBizID, 10, 64)
	if err != nil {
		blog.Errorf("field %s with value:%s invalid, %v, rid: %s", common.BKAppIDField, bkBizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	detail := struct {
		WithDetail bool `json:"with_detail"`
	}{}
	if err := ctx.DecodeInto(&detail); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.TopoOperation().SearchMainlineInstanceTopo(ctx.Kit, bizID, detail.WithDetail)
	if err != nil {
		blog.Errorf("search mainline instance topo by business:%d failed, %+v, rid: %s", bizID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrTopoMainlineSelectFailed))
		return
	}
	ctx.RespEntity(result)
}
