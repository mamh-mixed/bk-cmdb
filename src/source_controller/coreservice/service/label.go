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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/selector"
)

// AddLabels TODO
func (s *coreService) AddLabels(ctx *rest.Contexts) {
	inputData := selector.LabelAddRequest{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	if err := s.core.LabelOperation().AddLabel(ctx.Kit, inputData.TableName, inputData.Option); err != nil {
		blog.Errorf("AddLabels failed, table: %s, option: %+v, err: %s, rid: %s", inputData.TableName, inputData.Option, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// RemoveLabels TODO
func (s *coreService) RemoveLabels(ctx *rest.Contexts) {
	inputData := selector.LabelRemoveRequest{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	if err := s.core.LabelOperation().RemoveLabel(ctx.Kit, inputData.TableName, inputData.Option); err != nil {
		blog.Errorf("RemoveLabels failed, table: %s, option: %+v, err: %s, rid: %s", inputData.TableName, inputData.Option, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateLabels update service instance tag.
func (s *coreService) UpdateLabels(ctx *rest.Contexts) {
	inputData := selector.LabelUpdateRequest{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	if err := s.core.LabelOperation().UpdateLabel(ctx.Kit, inputData.TableName, inputData.Option); err != nil {
		blog.Errorf("update labels failed, table: %s, option: %+v, err: %v, rid: %s", inputData.TableName,
			inputData.Option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
