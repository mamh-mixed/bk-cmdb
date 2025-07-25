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

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// CreateProcessTemplateBatch TODO
// create a process template for a service template.
func (ps *ProcServer) CreateProcessTemplateBatch(ctx *rest.Contexts) {
	input := new(metadata.CreateProcessTemplateBatchInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.Processes) == 0 {
		ctx.RespEntity([]int64{})
		blog.Infof("no process to create, return")
		return
	}

	if len(input.Processes) > common.BKMaxUpdateOrCreatePageSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "create process template",
			common.BKMaxUpdateOrCreatePageSize))
		return
	}

	// authorize
	if err := ps.AuthManager.AuthorizeByServiceTemplateID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, input.ServiceTemplateID); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCheckAuthorizeFailed, "authorize by service template id failed, id: %d, err: %+v", input.ServiceTemplateID, err)
		return
	}

	ids := make([]int64, 0)
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		for _, process := range input.Processes {
			t := &metadata.ProcessTemplate{
				BizID:             input.BizID,
				ServiceTemplateID: input.ServiceTemplateID,
				Property:          process.Spec,
			}

			temp, err := ps.CoreAPI.CoreService().Process().CreateProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, t)
			if err != nil {
				blog.Errorf("create process template failed, template: %+v", *t)
				return err
			}

			ids = append(ids, temp.ID)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(ids)
}

// DeleteProcessTemplateBatch TODO
func (ps *ProcServer) DeleteProcessTemplateBatch(ctx *rest.Contexts) {
	input := new(metadata.DeleteProcessTemplateBatchInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.ProcessTemplates) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "process_templates"))
		return
	}
	if len(input.ProcessTemplates) > common.BKMaxDeletePageSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "delete process template",
			common.BKMaxDeletePageSize))
		return
	}

	// authorize by service template
	listOption := &metadata.ListProcessTemplatesOption{
		BusinessID:         input.BizID,
		ProcessTemplateIDs: input.ProcessTemplates,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	serviceTemplateIDs := make([]int64, 0)
	for _, processTemplate := range processTemplates.Info {
		serviceTemplateIDs = append(serviceTemplateIDs, processTemplate.ServiceTemplateID)
	}

	if err := ps.AuthManager.AuthorizeByServiceTemplateID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, serviceTemplateIDs...); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCheckAuthorizeFailed, "authorize by service template id failed, id: %+v, err: %+v", serviceTemplateIDs, err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := ps.CoreAPI.CoreService().Process().DeleteProcessTemplateBatch(ctx.Kit.Ctx, ctx.Kit.Header, input.ProcessTemplates)
		if err != nil {
			blog.Errorf("delete process template: %v failed", input.ProcessTemplates)
			return ctx.Kit.CCError.CCError(common.CCErrProcDeleteTemplateFail)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// UpdateProcessTemplate TODO
func (ps *ProcServer) UpdateProcessTemplate(ctx *rest.Contexts) {
	input := new(metadata.UpdateProcessTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if input.Property == nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update process template, but property empty, input: %+v", input)
		return
	}
	if input.ProcessTemplateID <= 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update process template, but get nil process template, input: %+v", input)
		return
	}

	listOption := &metadata.ListProcessTemplatesOption{
		BusinessID:         input.BizID,
		ProcessTemplateIDs: []int64{input.ProcessTemplateID},
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	serviceTemplateIDs := make([]int64, 0)
	for _, processTemplate := range processTemplates.Info {
		serviceTemplateIDs = append(serviceTemplateIDs, processTemplate.ServiceTemplateID)
	}

	if err := ps.AuthManager.AuthorizeByServiceTemplateID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, serviceTemplateIDs...); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCheckAuthorizeFailed, "authorize by service template id failed, id: %+v, err: %+v", serviceTemplateIDs, err)
		return
	}

	var template *metadata.ProcessTemplate
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		template, err = ps.CoreAPI.CoreService().Process().UpdateProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, input.ProcessTemplateID, input.Property)
		if err != nil {
			blog.Errorf("update process template: %v failed.", input)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(template)
}

// GetProcessTemplate TODO
func (ps *ProcServer) GetProcessTemplate(ctx *rest.Contexts) {
	input := &struct {
		BizID int64 `json:"bk_biz_id"`
	}{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	templateID, err := strconv.ParseInt(ctx.Request.PathParameter("processTemplateID"), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get process template, but get process template id failed, err: %v", err)
		return
	}

	template, err := ps.CoreAPI.CoreService().Process().GetProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get process template: %v failed, err: %v.", input, err)
		return
	}
	ctx.RespEntity(template)
}

// ListProcessTemplate TODO
func (ps *ProcServer) ListProcessTemplate(ctx *rest.Contexts) {

	input := new(metadata.ListProcessTemplateWithServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	option := &metadata.ListProcessTemplatesOption{
		BusinessID: input.BizID,
		Page:       input.Page,
	}

	if input.ServiceTemplateID > 0 {
		option.ServiceTemplateIDs = []int64{input.ServiceTemplateID}
	}

	if input.ProcessTemplatesIDs != nil {
		option.ProcessTemplateIDs = input.ProcessTemplatesIDs
	}
	template, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplateFailed, "get process template: %v failed", input)
		return
	}
	ctx.RespEntity(template)
}
