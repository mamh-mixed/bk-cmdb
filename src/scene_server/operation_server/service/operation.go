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
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateOperationChart TODO
func (o *OperationServer) CreateOperationChart(ctx *rest.Contexts) {
	chartInfo := new(metadata.ChartConfig)
	if err := ctx.DecodeInto(chartInfo); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// 图表是否已经存在
	filterCondition := mapstr.MapStr{}
	filterCondition[common.BKObjIDField] = chartInfo.ObjID
	filterCondition[common.OperationReportType] = chartInfo.ReportType
	filterCondition["field"] = chartInfo.Field
	exist, err := o.CoreAPI.CoreService().Operation().SearchChartCommon(ctx.Kit.Ctx, ctx.Kit.Header, filterCondition)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail,
			"new add operation chart fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		return
	}
	if exist.Data.Count > 0 {
		ctx.RespErrorCodeOnly(common.CCErrOperationChartAlreadyExist,
			"create operation chart fail, err: chart already exist, rid: %v", ctx.Kit.Rid)
		return
	}

	var id uint64
	resp := new(metadata.SearchChartCommon)

	defer func() {
		if id != 0 {
			opt := mapstr.MapStr{"config_id": id}
			resp, err = o.Engine.CoreAPI.CoreService().Operation().SearchChartCommon(ctx.Kit.Ctx, ctx.Kit.Header, opt)
			if err != nil {
				ctx.RespErrorCodeOnly(common.CCErrOperationSearchChartFail,
					"search operation chart fail, err: %v, rid: %v", err, ctx.Kit.Rid)
				return
			}
		}
		ctx.RespEntity(resp.Data)
		return
	}()

	// 自定义报表
	if chartInfo.ReportType == common.OperationCustom {
		result, err := o.Engine.CoreAPI.CoreService().Operation().CreateOperationChart(ctx.Kit.Ctx, ctx.Kit.Header,
			chartInfo)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail,
				"create operation chart fail, err: %v, rid: %v", err, ctx.Kit.Rid)
			return
		}

		id = result.Data
		return
	}

	// 内置报表
	srvData := o.newSrvComm(ctx.Kit.Header)
	configID, err := srvData.lgc.CreateInnerChart(ctx.Kit, chartInfo)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationNewAddStatisticFail, "create operation chart fail, err: %v, rid: %v",
			err, ctx.Kit.Rid)
		return
	}
	id = configID
}

// DeleteOperationChart TODO
func (o *OperationServer) DeleteOperationChart(ctx *rest.Contexts) {
	id := ctx.Request.PathParameter("id")
	_, err := o.Engine.CoreAPI.CoreService().Operation().DeleteOperationChart(ctx.Kit.Ctx, ctx.Kit.Header, id)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationDeleteChartFail, "delete operation chart fail, err: %v, rid: %v",
			err, ctx.Kit.Rid)
		return
	}

	ctx.RespEntity(nil)
}

// SearchOperationChart TODO
func (o *OperationServer) SearchOperationChart(ctx *rest.Contexts) {
	opt := make(map[string]interface{})

	ctx.SetReadPreference(common.SecondaryPreferredMode)
	result, err := o.Engine.CoreAPI.CoreService().Operation().SearchOperationCharts(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationSearchChartFail, "search operation chart fail, err: %v, rid: %v",
			err, ctx.Kit.Rid)
		return
	}

	ctx.RespEntity(result.Data)
}

// UpdateOperationChart TODO
func (o *OperationServer) UpdateOperationChart(ctx *rest.Contexts) {
	opt := mapstr.MapStr{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if _, err := o.Engine.CoreAPI.CoreService().Operation().UpdateOperationChart(ctx.Kit.Ctx, ctx.Kit.Header,
		opt); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationUpdateChartFail,
			"update operation chart fail, err: %v, chartInfo: %v, rid: %v", err, opt, ctx.Kit.Rid)
		return
	}

	ctx.RespEntity(opt["config_id"])
}

// SearchChartData search data of chart
func (o *OperationServer) SearchChartData(ctx *rest.Contexts) {
	srvData := o.newSrvComm(ctx.Kit.Header)
	inputParams := mapstr.MapStr{}
	if err := ctx.DecodeInto(&inputParams); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.SetReadPreference(common.SecondaryPreferredMode)
	chart, err := o.CoreAPI.CoreService().Operation().SearchChartCommon(ctx.Kit.Ctx, ctx.Kit.Header, inputParams)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v, cond: %v, "+
			"rid: %v", err, inputParams, ctx.Kit.Rid)
		return
	}

	innerChart := []string{
		common.BizModuleHostChart,
		common.ModelAndInstCount,
		common.HostChangeBizChart,
		common.ModelInstChart,
		common.ModelInstChangeChart,
	}

	if util.InStrArr(innerChart, chart.Data.Info.ReportType) {
		data, err := srvData.lgc.InnerChartData(ctx.Kit, chart.Data.Info)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, cond: %v, err: %v, "+
				"rid: %v", chart.Data.Info, err, ctx.Kit.Rid)
			return
		}
		ctx.RespEntity(data)
		return
	}

	// 判断模型是否存在，不存在返回nil
	cond := make(map[string]interface{}, 0)
	cond[common.BKObjIDField] = chart.Data.Info.ObjID
	query := metadata.QueryCondition{Condition: cond}
	models, err := o.CoreAPI.CoreService().Model().ReadModel(ctx.Kit.Ctx, ctx.Kit.Header, &query)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, err: %v, rid: %v", err,
			ctx.Kit.Rid)
		return
	}
	if models.Count <= 0 {
		ctx.RespEntity(nil)
		return
	}

	result, err := o.CoreAPI.CoreService().Operation().SearchChartData(ctx.Kit.Ctx, ctx.Kit.Header, chart.Data.Info)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationGetChartDataFail, "search chart data fail, cond: %v, err: %v, "+
			"rid: %v", chart.Data.Info, err, ctx.Kit.Rid)
		return
	}
	ctx.RespEntity(result.Data)
}

// UpdateChartPosition TODO
func (o *OperationServer) UpdateChartPosition(ctx *rest.Contexts) {
	opt := metadata.ChartPosition{}
	if err := ctx.DecodeInto(&opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := o.CoreAPI.CoreService().Operation().UpdateChartPosition(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrOperationUpdateChartPositionFail,
			"update chart position fail, err: %v, rid: %v", err, ctx.Kit.Rid)
		return
	}

	ctx.RespEntity(nil)
}
