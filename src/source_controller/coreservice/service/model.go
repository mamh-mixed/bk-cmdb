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
	"fmt"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateManyModelClassification TODO
func (s *coreService) CreateManyModelClassification(ctx *rest.Contexts) {
	inputDatas := metadata.CreateManyModelClassifiaction{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateManyModelClassification(ctx.Kit, inputDatas))
}

// CreateOneModelClassification TODO
func (s *coreService) CreateOneModelClassification(ctx *rest.Contexts) {
	inputData := metadata.CreateOneModelClassification{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateOneModelClassification(ctx.Kit, inputData))
}

// SetOneModelClassification TODO
func (s *coreService) SetOneModelClassification(ctx *rest.Contexts) {
	inputData := metadata.SetOneModelClassification{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().SetOneModelClassification(ctx.Kit, inputData))
}

// SetManyModelClassification TODO
func (s *coreService) SetManyModelClassification(ctx *rest.Contexts) {
	inputDatas := metadata.SetManyModelClassification{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SetManyModelClassification(ctx.Kit, inputDatas))
}

// UpdateModelClassification TODO
func (s *coreService) UpdateModelClassification(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelClassification(ctx.Kit, inputData))
}

// DeleteModelClassification TODO
func (s *coreService) DeleteModelClassification(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelClassification(ctx.Kit, inputData))
}

// SearchModelClassification TODO
func (s *coreService) SearchModelClassification(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelClassification(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate language
	lang := s.Language(ctx.Kit.Header)
	defaultIDMap := map[string]bool{
		metadata.ClassificationHostManageID:    true,
		metadata.ClassificationBizTopoID:       true,
		metadata.ClassificationOrganizationID:  true,
		metadata.ClassificationNetworkID:       true,
		metadata.ClassificationUncategorizedID: true,
	}
	nameMap := map[string]string{
		metadata.ClassificationHostManageID:    metadata.ClassificationHostManage,
		metadata.ClassificationBizTopoID:       metadata.ClassificationTopo,
		metadata.ClassificationOrganizationID:  metadata.ClassificationOrganization,
		metadata.ClassificationNetworkID:       metadata.ClassificationNet,
		metadata.ClassificationUncategorizedID: metadata.ClassificationUncategorized,
	}

	for index := range dataResult.Info {
		result := dataResult.Info[index]
		if defaultIDMap[result.ClassificationID] && result.ClassificationName == nameMap[result.ClassificationID] {
			dataResult.Info[index].ClassificationName = s.TranslateClassificationName(lang, &result)
		}
	}
	ctx.RespEntity(dataResult)
}

// CreateModel TODO
func (s *coreService) CreateModel(ctx *rest.Contexts) {
	inputData := metadata.CreateModel{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateModel(ctx.Kit, inputData))
}

// CreateTableModel create inner table model
func (s *coreService) CreateTableModel(ctx *rest.Contexts) {
	inputData := metadata.CreateModel{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateTableModel(ctx.Kit, inputData))
}

// SetModel TODO
func (s *coreService) SetModel(ctx *rest.Contexts) {
	inputData := metadata.SetModel{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SetModel(ctx.Kit, inputData))
}

// UpdateModel TODO
func (s *coreService) UpdateModel(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModel(ctx.Kit, inputData))
}

// DeleteModel TODO
func (s *coreService) DeleteModel(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModel(ctx.Kit, inputData))
}

// CascadeDeleteModel TODO
func (s *coreService) CascadeDeleteModel(ctx *rest.Contexts) {
	idStr := ctx.Request.PathParameter(common.BKFieldID)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID))
		return
	}
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CascadeDeleteModel(ctx.Kit, id))
}

// CascadeDeleteTableModel delete table model related resources in a cascading manner.
func (s *coreService) CascadeDeleteTableModel(ctx *rest.Contexts) {
	inputData := metadata.DeleteTableOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	if err := s.core.ModelOperation().CascadeDeleteTableModel(ctx.Kit, inputData); err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)

}

// SearchModel TODO
func (s *coreService) SearchModel(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModel(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate
	lang := s.Language(ctx.Kit.Header)
	for modelIdx := range dataResult.Info {
		if needTranslateObjMap[dataResult.Info[modelIdx].ObjectID] {
			dataResult.Info[modelIdx].ObjectName = s.TranslateObjectName(lang, &dataResult.Info[modelIdx])
		}
	}

	ctx.RespEntity(dataResult)
}

// ListModel list object.
func (s *coreService) ListModel(cts *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	filter = util.SetQueryOwner(filter, cts.Kit.SupplierAccount)

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(common.BKTableNameObjDes).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count objects failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.QueryModelDataResult{Count: int64(count)})
		return
	}

	objects := make([]metadata.Object, 0)
	err = mongodb.Client().Table(common.BKTableNameObjDes).Find(filter).Start(uint64(opt.Page.Start)).
		Limit(uint64(opt.Page.Limit)).Sort(opt.Page.Sort).Fields(opt.Fields...).All(cts.Kit.Ctx, &objects)
	if err != nil {
		blog.Errorf("list objects failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	result := metadata.QueryModelDataResult{Info: objects}

	// translate
	lang := s.Language(cts.Kit.Header)
	for modelIdx := range result.Info {
		if needTranslateObjMap[result.Info[modelIdx].ObjectID] {
			result.Info[modelIdx].ObjectName = s.TranslateObjectName(lang, &result.Info[modelIdx])
		}
	}

	cts.RespEntity(result)
}

// SearchModelWithAttribute search model with attribute
func (s *coreService) SearchModelWithAttribute(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelWithAttribute(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate
	lang := s.Language(ctx.Kit.Header)
	for modelIdx := range dataResult.Info {
		if needTranslateObjMap[dataResult.Info[modelIdx].Spec.ObjectID] {
			dataResult.Info[modelIdx].Spec.ObjectName = s.TranslateObjectName(lang, &dataResult.Info[modelIdx].Spec)
		}
		for attributeIdx := range dataResult.Info[modelIdx].Attributes {
			if dataResult.Info[modelIdx].Attributes[attributeIdx].IsPre || dataResult.Info[modelIdx].Spec.IsPre ||
				needTranslateObjMap[dataResult.Info[modelIdx].Spec.ObjectID] {
				dataResult.Info[modelIdx].Attributes[attributeIdx].PropertyName =
					s.TranslatePropertyName(lang, &dataResult.Info[modelIdx].Attributes[attributeIdx])
				dataResult.Info[modelIdx].Attributes[attributeIdx].Placeholder =
					s.TranslatePlaceholder(lang, &dataResult.Info[modelIdx].Attributes[attributeIdx])
				switch dataResult.Info[modelIdx].Attributes[attributeIdx].PropertyType {
				case common.FieldTypeEnum, common.FieldTypeEnumMulti:
					dataResult.Info[modelIdx].Attributes[attributeIdx].Option =
						s.TranslateEnumName(ctx.Kit.Ctx, lang, &dataResult.Info[modelIdx].Attributes[attributeIdx],
							dataResult.Info[modelIdx].Attributes[attributeIdx].Option)
				}
			}
		}
	}

	ctx.RespEntity(dataResult)
}

// GetModelStatistics 用于统计各个模型的实例数(Web页面展示需要)
func (s *coreService) GetModelStatistics(ctx *rest.Contexts) {
	// statistics data include all object model statistics.
	statistics := []metadata.ObjectIDCount{}

	// stat set count.
	filter := map[string]interface{}{}
	setCount, err := mongodb.Client().Table(common.BKTableNameBaseSet).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count set model instances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	statistics = append(statistics, metadata.ObjectIDCount{ObjID: common.BKInnerObjIDSet, Count: int64(setCount)})

	// stat module count.
	moduleCount, err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count module model instances failed, err: %+v, rid: %s", err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	statistics = append(statistics, metadata.ObjectIDCount{ObjID: common.BKInnerObjIDModule, Count: int64(moduleCount)})

	// stat host count.
	hostCount, err := mongodb.Client().Table(common.BKTableNameBaseHost).Find(filter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count host model instances failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	statistics = append(statistics, metadata.ObjectIDCount{ObjID: common.BKInnerObjIDHost, Count: int64(hostCount)})

	// stat biz count.
	appFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultAppFlag,
		},
		common.BKDataStatusField: map[string]interface{}{
			common.BKDBNE: common.DataStatusDisabled,
		},
	}
	bizCount, err := mongodb.Client().Table(common.BKTableNameBaseApp).Find(appFilter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetModelStatistics failed, count application model instances failed, err: %+v, rid: %s", err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	statistics = append(statistics, metadata.ObjectIDCount{ObjID: common.BKInnerObjIDApp, Count: int64(bizCount)})

	// stat common object counts.
	allObjects := []metadata.ObjectIDCount{}
	commonObjects := []metadata.ObjectIDCount{}

	objectFilter := []map[string]interface{}{
		{
			common.BKDBGroup: map[string]interface{}{
				"_id": "$bk_obj_id",
				"count": map[string]interface{}{
					common.BKDBSum: 1,
				},
			},
		},
	}
	err = mongodb.Client().Table(common.BKTableNameObjDes).AggregateAll(ctx.Kit.Ctx, objectFilter, &allObjects)
	if err != nil {
		blog.Errorf("get all object models failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// only stat common object models.
	for _, object := range allObjects {
		if metadata.IsCommon(object.ObjID) {
			commonObjects = append(commonObjects, object)
		}
	}

	// stat common object counts in sharding tables.
	for _, object := range commonObjects {
		// stat object sharding data one by one.
		data := []metadata.ObjectIDCount{}

		// sharding table name.
		tableName := common.GetObjectInstTableName(object.ObjID, ctx.Kit.SupplierAccount)

		if err := mongodb.Client().Table(tableName).AggregateAll(ctx.Kit.Ctx, objectFilter, &data); err != nil {
			blog.Errorf("get object %s instances count failed, err: %+v, rid: %s", object.ObjID, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		statistics = append(statistics, data...)
	}

	ctx.RespEntity(statistics)
}

// CreateModelAttributeGroup TODO
func (s *coreService) CreateModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.CreateModelAttributeGroup{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().CreateModelAttributeGroup(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// SetModelAttributeGroup TODO
func (s *coreService) SetModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.SetModelAttributeGroup{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().SetModelAttributeGroup(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// UpdateModelAttributeGroup TODO
func (s *coreService) UpdateModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributeGroup(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// UpdateModelAttributeGroupByCondition TODO
func (s *coreService) UpdateModelAttributeGroupByCondition(ctx *rest.Contexts) {
	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributeGroupByCondition(ctx.Kit, inputData))
}

// SearchModelAttributeGroup TODO
func (s *coreService) SearchModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributeGroup(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsDefault {
			dataResult.Info[index].GroupName = s.TranslatePropertyGroupName(lang, &dataResult.Info[index])
		}
	}
	ctx.RespEntity(dataResult)
}

// SearchModelAttributeGroupByCondition TODO
func (s *coreService) SearchModelAttributeGroupByCondition(ctx *rest.Contexts) {
	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributeGroupByCondition(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsDefault {
			dataResult.Info[index].GroupName = s.TranslatePropertyGroupName(lang, &dataResult.Info[index])
		}
	}
	ctx.RespEntity(dataResult)
}

// DeleteModelAttributeGroup TODO
func (s *coreService) DeleteModelAttributeGroup(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttributeGroup(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// DeleteModelAttributeGroupByCondition TODO
func (s *coreService) DeleteModelAttributeGroupByCondition(ctx *rest.Contexts) {
	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttributeGroupByCondition(ctx.Kit, inputData))
}

// CreateModelAttributes TODO
func (s *coreService) CreateModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.CreateModelAttributes{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().CreateModelAttributes(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// SetModelAttributes TODO
func (s *coreService) SetModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.SetModelAttributes{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SetModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"),
		inputData))
}

// UpdateTableModelAttrByCondition update properties of form fields based on conditions.
func (s *coreService) UpdateTableModelAttrByCondition(ctx *rest.Contexts) {

	inputData := metadata.UpdateTableOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	err := s.core.ModelOperation().UpdateTableModelAttributes(ctx.Kit, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
	return
}

// UpdateModelAttributes TODO
func (s *coreService) UpdateModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	updateCount, err := s.core.ModelOperation().UpdateModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"),
		inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(updateCount)
}

// UpdateModelAttributesIndex update model attribute
func (s *coreService) UpdateModelAttributesIndex(ctx *rest.Contexts) {
	inputData := new(metadata.UpdateAttrIndexInput)
	if err := ctx.DecodeInto(inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	objID := ctx.Request.PathParameter(common.BKObjIDField)
	idStr := ctx.Request.PathParameter("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		blog.Errorf("parse id from path params failed, err: %v, id: %s, rid: %s", err, idStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID))
		return
	}

	err = s.core.ModelOperation().UpdateModelAttributeIndex(ctx.Kit, objID, id, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// UpdateModelAttributesByCondition TODO
func (s *coreService) UpdateModelAttributesByCondition(ctx *rest.Contexts) {

	inputData := metadata.UpdateOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().UpdateModelAttributesByCondition(ctx.Kit, inputData))
}

// DeleteModelAttribute TODO
func (s *coreService) DeleteModelAttribute(ctx *rest.Contexts) {

	inputData := metadata.DeleteOption{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttributes(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), inputData))
}

// SearchModelAttrsWithTableByCondition querying model properties containing table types
// NOTICE: include table attributes
func (s *coreService) SearchModelAttrsWithTableByCondition(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ModelOperation().SearchModelAttrsWithTableByCondition(ctx.Kit, inputData)
	if err != nil {
		ctx.RespEntityWithError(result, err)
		return
	}

	// translate
	lang := s.Language(ctx.Kit.Header)
	for index := range result.Info {
		if result.Info[index].IsPre || needTranslateObjMap[result.Info[index].ObjectID] {
			result.Info[index].PropertyName = s.TranslatePropertyName(lang, &result.Info[index])
			result.Info[index].Placeholder = s.TranslatePlaceholder(lang, &result.Info[index])
			switch result.Info[index].PropertyType {
			case common.FieldTypeEnum, common.FieldTypeEnumMulti:
				result.Info[index].Option = s.TranslateEnumName(ctx.Kit.Ctx, lang, &result.Info[index],
					result.Info[index].Option)
			case common.FieldTypeTable:
				option, err := metadata.ParseSubAttribute(ctx.Kit.Ctx, result.Info[index].Option)
				if err != nil {
					blog.Errorf("marshal table type option failed, err: %v, option: %v, rid: %s, ", err,
						result.Info[index].Option, ctx.Kit.Rid)
					ctx.RespAutoError(err)
					return
				}

				for i, attr := range option {
					option[i].Placeholder = s.TranslateTablePlaceholder(lang, result.Info[index].ObjectID,
						result.Info[index].PropertyID, attr.PropertyID)
				}
				result.Info[index].Option = option
			}
		}
	}

	ctx.RespEntity(result)
}

// SearchModelAttributesByCondition query for model attributes that do not contain table types
// NOTICE: exclude table attributes
func (s *coreService) SearchModelAttributesByCondition(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributesByCondition(ctx.Kit, inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsPre || needTranslateObjMap[dataResult.Info[index].ObjectID] {
			dataResult.Info[index].PropertyName = s.TranslatePropertyName(lang, &dataResult.Info[index])
			dataResult.Info[index].Placeholder = s.TranslatePlaceholder(lang, &dataResult.Info[index])
			switch dataResult.Info[index].PropertyType {
			case common.FieldTypeEnum, common.FieldTypeEnumMulti:
				dataResult.Info[index].Option =
					s.TranslateEnumName(ctx.Kit.Ctx, lang, &dataResult.Info[index], dataResult.Info[index].Option)
			}
		}
	}

	ctx.RespEntity(dataResult)
}

// SearchModelAttributes TODO
func (s *coreService) SearchModelAttributes(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	dataResult, err := s.core.ModelOperation().SearchModelAttributes(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"),
		inputData)
	if nil != err {
		ctx.RespEntityWithError(dataResult, err)
		return
	}

	// translate 主机内置字段bk_state不做翻译
	lang := s.Language(ctx.Kit.Header)
	for index := range dataResult.Info {
		if dataResult.Info[index].IsPre || needTranslateObjMap[dataResult.Info[index].ObjectID] {
			dataResult.Info[index].PropertyName = s.TranslatePropertyName(lang, &dataResult.Info[index])
			dataResult.Info[index].Placeholder = s.TranslatePlaceholder(lang, &dataResult.Info[index])

			switch dataResult.Info[index].PropertyType {
			case common.FieldTypeEnum, common.FieldTypeEnumMulti:
				dataResult.Info[index].Option =
					s.TranslateEnumName(ctx.Kit.Ctx, lang, &dataResult.Info[index], dataResult.Info[index].Option)
			}
		}
	}

	ctx.RespEntity(dataResult)
}

// SearchModelAttrUnique TODO
func (s *coreService) SearchModelAttrUnique(ctx *rest.Contexts) {

	inputData := metadata.QueryCondition{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(s.core.ModelOperation().SearchModelAttrUnique(ctx.Kit, inputData))
}

// CreateModelAttrUnique TODO
func (s *coreService) CreateModelAttrUnique(ctx *rest.Contexts) {
	inputData := metadata.CreateModelAttrUnique{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}

	objID := ctx.Request.PathParameter("bk_obj_id")
	res, err := s.core.ModelOperation().CreateModelAttrUnique(ctx.Kit, objID, inputData)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// UpdateModelAttrUnique TODO
func (s *coreService) UpdateModelAttrUnique(ctx *rest.Contexts) {
	inputDatas := metadata.UpdateModelAttrUnique{}
	if err := ctx.DecodeInto(&inputDatas); nil != err {
		ctx.RespAutoError(err)
		return
	}

	id, err := strconv.ParseUint(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}

	objID := ctx.Request.PathParameter("bk_obj_id")
	res, err := s.core.ModelOperation().UpdateModelAttrUnique(ctx.Kit, objID, id, inputDatas)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(res)
}

// DeleteModelAttrUnique TODO
func (s *coreService) DeleteModelAttrUnique(ctx *rest.Contexts) {
	id, err := strconv.ParseUint(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, "id"))
		return
	}

	ctx.RespEntityWithError(s.core.ModelOperation().DeleteModelAttrUnique(ctx.Kit,
		ctx.Request.PathParameter("bk_obj_id"), id))
}

// CreateModelTables TODO
func (s *coreService) CreateModelTables(ctx *rest.Contexts) {
	inputData := metadata.CreateModelTable{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(nil, s.core.ModelOperation().CreateModelTables(ctx.Kit, inputData))
}

// CreateTableModelTables create table model table
func (s *coreService) CreateTableModelTables(ctx *rest.Contexts) {
	inputData := metadata.CreateModelTable{}
	if err := ctx.DecodeInto(&inputData); nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntityWithError(nil, s.core.ModelOperation().CreateTableModelTables(ctx.Kit, inputData))
}

// UpdateIDGenerator update id generator
func (s *coreService) UpdateIDGenerator(ctx *rest.Contexts) {
	opt := new(metadata.UpdateIDGenOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := mapstr.MapStr{common.BKFieldDBID: opt.Type}
	result := make(map[string]interface{})
	err := mongodb.Client().Table(common.BKTableNameIDgenerator).Find(cond).Fields(common.BKFieldSeqID).One(ctx.Kit.Ctx,
		&result)
	if err != nil {
		blog.Errorf("find id generator failed, err: %v, filter: %+v, rid: %s", err, cond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrObjectSelectInstFailed))
		return
	}

	curID, err := util.GetInt64ByInterface(result[common.BKFieldSeqID])
	if err != nil {
		blog.Errorf("get id generator sequence failed, data: %v, err: %v, rid: %s", result, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrObjectSelectInstFailed))
		return
	}

	if curID > opt.SequenceID {
		blog.Errorf("sequence can not be less than current value, cur: %d, opt: %+v, rid: %s", curID, opt, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, fmt.Sprintf("sequence can not "+
			"be less than current value: %d", curID)))
		return
	}

	cond = mapstr.MapStr{
		common.BKFieldDBID:  opt.Type,
		common.BKFieldSeqID: mapstr.MapStr{common.BKDBLT: opt.SequenceID},
	}
	data := mapstr.MapStr{
		common.BKFieldSeqID:  opt.SequenceID,
		common.LastTimeField: time.Now(),
	}
	if err := mongodb.Client().Table(common.BKTableNameIDgenerator).Update(ctx.Kit.Ctx, cond, data); err != nil {
		blog.Errorf("update id generator failed, err: %v, filter: %+v, rid: %s", err, cond, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrObjectDBOpErrno))
		return
	}

	ctx.RespEntity(nil)
}
