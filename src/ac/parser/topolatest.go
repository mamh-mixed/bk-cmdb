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

package parser

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"

	"github.com/tidwall/gjson"
)

// this package's topology filter is the latest api version
// for these resources, it also has a elder resource api version.
// TODO: if the elder api has been removed, delete their resource
// filter at the same time.

func (ps *parseStream) topologyLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	ps.objectUniqueLatest().
		associationTypeLatest().
		objectAssociationLatest().
		objectInstanceAssociationLatest().
		objectInstanceLatest().
		objectLatest().
		objectClassificationLatest().
		objectAttributeGroupLatest().
		objectAttributeLatest().
		mainlineLatest().
		setTemplate().
		modelQuote().
		fieldTemplate()

	return ps
}

var (
	createObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/create/objectunique/object/[^\s/]+/?$`)
	updateObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectunique/object/[^\s/]+/unique/[0-9]+/?$`)
	deleteObjectUniqueLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectunique/object/[^\s/]+/unique/[0-9]+/?$`)
	findObjectUniqueLatestRegexp   = regexp.MustCompile(`^/api/v3/find/objectunique/object/[^\s/]+/?$`)
)

func (ps *parseStream) objectUniqueLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// TODO: add business id for these filter rules to resources.
	// add object unique operation.
	if ps.hitRegexp(createObjectUniqueLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("add object unique, but got invalid url")
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelUnique,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// update object unique operation.
	if ps.hitRegexp(updateObjectUniqueLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update object unique, but got invalid url")
			return ps
		}

		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[7])
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelUnique,
					Action:     meta.Update,
					InstanceID: uniqueID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// delete object unique operation.
	if ps.hitRegexp(deleteObjectUniqueLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("delete object unique, but got invalid url")
			return ps
		}

		uniqueID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object unique, but got invalid unique id %s", ps.RequestCtx.Elements[7])
			return ps
		}
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: ps.RequestCtx.Elements[5]})
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.ModelUnique,
					Action:     meta.Delete,
					InstanceID: uniqueID,
				},
				Layers:     []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				BusinessID: bizID,
			},
		}
		return ps
	}

	// find model unique operation
	if ps.hitRegexp(findObjectUniqueLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelUnique,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findManyAssociationKindLatestPattern = "/api/v3/find/associationtype"
	createAssociationKindLatestPattern   = "/api/v3/create/associationtype"
)

var (
	updateAssociationKindLatestRegexp = regexp.MustCompile(`^/api/v3/update/associationtype/[0-9]+/?$`)
	deleteAssociationKindLatestRegexp = regexp.MustCompile(`^/api/v3/delete/associationtype/[0-9]+/?$`)
)

func (ps *parseStream) associationTypeLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find association kind operation
	if ps.hitPattern(findManyAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AssociationType,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// create association kind operation
	if ps.hitPattern(createAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.AssociationType,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update association kind operation
	if ps.hitRegexp(updateAssociationKindLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update association kind, but got invalid url")
			return ps
		}

		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.AssociationType,
					Action:     meta.Update,
					InstanceID: kindID,
				},
			},
		}

		return ps
	}

	// delete association kind operation
	if ps.hitRegexp(deleteAssociationKindLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete association kind, but got invalid url")
			return ps
		}

		kindID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete association kind, but got invalid kind id %s", ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:       meta.AssociationType,
					Action:     meta.Delete,
					InstanceID: kindID,
				},
			},
		}

		return ps
	}

	return ps
}

const (
	findObjectAssociationLatestPattern                     = "/api/v3/find/objectassociation"
	createObjectAssociationLatestPattern                   = "/api/v3/create/objectassociation"
	findObjectAssociationWithAssociationKindLatestPattern  = "/api/v3/find/topoassociationtype"
	countObjectAssociationWithAssociationKindLatestPattern = "/api/v3/count/topoassociationtype"
	// excel 导入主机专用接口
	findModelAssociationPattern = "/api/v3/find/instassociation/model"
)

var (
	updateObjectAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectassociation/[0-9]+/?$`)
	deleteObjectAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectassociation/[0-9]+/?$`)
	// excel 导入关联关系专用接口
	findAssociationByObjectAssociationIDLatestRegexp = regexp.MustCompile(
		`^/api/v3/topo/find/object/[^\s/]+/association/by/bk_obj_asst_id$`)
	// excel 导入关联关系专用接口
	importAssociationByObjectAssociationIDLatestRegexp = regexp.MustCompile(
		`^/api/v3/import/instassociation/[^\s/]+$`)
)

// NOCC:golint/fnsize(设计如此)
func (ps *parseStream) objectAssociationLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// search object association operation
	if ps.hitPattern(findObjectAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// create object association operation
	if ps.hitPattern(createObjectAssociationLatestPattern, http.MethodPost) {
		objVal, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}

		asstVal, err := ps.RequestCtx.getValueFromBody(common.BKAsstObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}

		res, err := ps.getModelAsstUpdateResAttr(objVal.Value(), asstVal.Value())
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = append(ps.Attribute.Resources, res...)
		return ps
	}

	// update/delete object association operation
	if ps.hitRegexp(updateObjectAssociationLatestRegexp, http.MethodPut) ||
		ps.hitRegexp(deleteObjectAssociationLatestRegexp, http.MethodDelete) {

		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object association, but got invalid url")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object association id %s is invalid", ps.RequestCtx.Elements[4])
			return ps
		}

		asst, err := ps.getModelAssociation(mapstr.MapStr{common.BKFieldID: assoID})
		if err != nil {
			ps.err = err
			return ps
		}

		res, err := ps.getModelAsstUpdateResAttr(asst[0].ObjectID, asst[0].AsstObjID)
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = append(ps.Attribute.Resources, res...)
		return ps
	}

	// find object association with a association kind list.
	if ps.hitPattern(findObjectAssociationWithAssociationKindLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// count object association with a association kind list.
	if ps.hitPattern(countObjectAssociationWithAssociationKindLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// excel 导入关联关系专用接口, 跳过鉴权
	if ps.hitRegexp(findAssociationByObjectAssociationIDLatestRegexp, http.MethodPost) ||
		ps.hitRegexp(importAssociationByObjectAssociationIDLatestRegexp, http.MethodPost) {

		ps.Attribute.Resources = []meta.ResourceAttribute{{Basic: meta.Basic{Action: meta.SkipAction}}}
		return ps
	}

	// excel 导入主机专用接口, 跳过鉴权
	if ps.hitPattern(findModelAssociationPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAssociation,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

func (ps *parseStream) getModelAsstUpdateResAttr(objID, asstObjID interface{}) ([]meta.ResourceAttribute, error) {
	filter := mapstr.MapStr{
		common.BKObjIDField: mapstr.MapStr{
			common.BKDBIN: []interface{}{
				objID,
				asstObjID,
			},
		},
	}
	models, err := ps.searchModels(filter)
	if err != nil {
		return nil, err
	}

	bizID, err := ps.RequestCtx.getBizIDFromBody()
	if err != nil {
		return nil, err
	}

	resources := make([]meta.ResourceAttribute, 0, len(models))
	for _, model := range models {
		resources = append(resources, meta.ResourceAttribute{
			BusinessID: bizID,
			Basic: meta.Basic{
				Type:       meta.Model,
				Action:     meta.Update,
				InstanceID: model.ID,
			},
		})
	}
	return resources, nil
}

const (
	findObjectInstanceAssociationLatestPattern        = "/api/v3/find/instassociation"
	findObjectInstanceAssociationRelatedLatestPattern = "/api/v3/find/instassociation/related"
	createObjectInstanceAssociationLatestPattern      = "/api/v3/create/instassociation"
	createObjectManyInstanceAssociationLatestPattern  = "/api/v3/createmany/instassociation"
)

var (
	deleteObjectInstanceAssociationLatestRegexp      = regexp.MustCompile(`^/api/v3/delete/instassociation/[^\s/]+/[0-9]+/?$`)
	deleteObjectInstanceAssociationBatchLatestRegexp = regexp.MustCompile("^/api/v3/delete/instassociation/batch")
	findObjectInstanceTopologyUILatestRegexp         = regexp.MustCompile(`^/api/v3/findmany/inst/association/object/[^\s/]+/inst_id/[0-9]+/offset/[0-9]+/limit/[0-9]+/web$`)
	findInstAssociationObjInstInfoLatestRegexp       = regexp.MustCompile(`^/api/v3/findmany/inst/association/association_object/inst_base_info$`)
	searchInstAssociationAndInstDetailLatestRegexp   = regexp.MustCompile(
		`^/api/v3/find/instassociation/object/[^\s/]+/inst/detail`)

	searchInstanceAssociationsRegexp = regexp.MustCompile(`^/api/v3/search/instance_associations/object/[^\s/]+/?$`)
	countInstanceAssociationsRegexp  = regexp.MustCompile(`^/api/v3/count/instance_associations/object/[^\s/]+/?$`)

	findObjectInstanceAssociationWithBizIDRegexp = regexp.MustCompile(`^/api/v3/find/instassociation/biz/([0-9]+)/?$`)
)

func (ps *parseStream) objectInstanceAssociationLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// find instance's association operation.
	if ps.hitPattern(findObjectInstanceAssociationLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find instance's association with bizID operation.
	if ps.hitRegexp(findObjectInstanceAssociationWithBizIDRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// search instance associations operation.
	if ps.hitRegexp(searchInstanceAssociationsRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// count instance associations operation.
	if ps.hitRegexp(countInstanceAssociationsRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find instance's association related info operation.
	if ps.hitPattern(findObjectInstanceAssociationRelatedLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// create instance association operation.
	if ps.hitPattern(createObjectInstanceAssociationLatestPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.AssociationObjAsstIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		associationObjAsstID := val.String()
		filter := mapstr.MapStr{
			common.AssociationObjAsstIDField: associationObjAsstID,
		}
		asst, err := ps.getModelAssociation(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		modelFilter := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: []interface{}{
					asst[0].ObjectID,
					asst[0].AsstObjID,
				},
			},
		}
		models, err := ps.searchModels(modelFilter)
		if err != nil {
			ps.err = err
			return ps
		}

		val, err = ps.RequestCtx.getValueFromBody(common.BKInstIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		instanceID := val.Int()
		if instanceID <= 0 {
			ps.err = errors.New("invalid bk_inst_id value")
			return ps
		}

		val, err = ps.RequestCtx.getValueFromBody(common.BKAsstInstIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		asstInstID := val.Int()
		if asstInstID <= 0 {
			ps.err = errors.New("invalid bk_asst_inst_id value")
			return ps
		}

		// 处理模型自关联的情况
		if len(models) == 1 {
			instRes, err := ps.generateUpdateInstanceResource(&models[0], instanceID)
			if err != nil {
				ps.err = err
				return ps
			}
			asstInstRes, err := ps.generateUpdateInstanceResource(&models[0], asstInstID)
			if err != nil {
				ps.err = err
				return ps
			}

			ps.Attribute.Resources = []meta.ResourceAttribute{*instRes, *asstInstRes}
			return ps
		}

		for _, model := range models {
			var instID int64
			if model.ObjectID == asst[0].ObjectID {
				instID = instanceID
			} else {
				instID = asstInstID
			}
			instRes, err := ps.generateUpdateInstanceResource(&model, instID)
			if err != nil {
				ps.err = err
				return ps
			}

			ps.Attribute.Resources = append(ps.Attribute.Resources, *instRes)
		}
		return ps
	}

	if ps.hitPattern(createObjectManyInstanceAssociationLatestPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.AssociationObjAsstIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		associationObjAsstID := val.String()
		filter := mapstr.MapStr{
			common.AssociationObjAsstIDField: associationObjAsstID,
		}
		asst, err := ps.getModelAssociation(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		modelFilter := mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{
				common.BKDBIN: []interface{}{
					asst[0].ObjectID,
					asst[0].AsstObjID,
				},
			},
		}
		models, err := ps.searchModels(modelFilter)
		if err != nil {
			ps.err = err
			return ps
		}

		instances, err := ps.RequestCtx.getValueFromBody("details")
		if err != nil {
			ps.err = err
			return ps
		}
		for _, instance := range instances.Array() {
			instanceMap := instance.Map()
			instID := instanceMap[common.BKInstIDField]
			asstInstanceID := instanceMap[common.BKAsstInstIDField]
			instanceID := instID.Int()
			if instanceID <= 0 {
				ps.err = errors.New("invalid bk_inst_id value")
				return ps
			}

			asstInstID := asstInstanceID.Int()
			if asstInstID <= 0 {
				ps.err = errors.New("invalid bk_asst_inst_id value")
				return ps
			}

			// 处理模型自关联的情况
			if len(models) == 1 {
				instRes, err := ps.generateUpdateInstanceResource(&models[0], instanceID)
				if err != nil {
					ps.err = err
					return ps
				}
				asstInstRes, err := ps.generateUpdateInstanceResource(&models[0], asstInstID)
				if err != nil {
					ps.err = err
					return ps
				}

				ps.Attribute.Resources = []meta.ResourceAttribute{*instRes, *asstInstRes}
			} else {
				for _, model := range models {
					var instID int64
					if model.ObjectID == asst[0].ObjectID {
						instID = instanceID
					} else {
						instID = asstInstID
					}
					instRes, err := ps.generateUpdateInstanceResource(&model, instID)
					if err != nil {
						ps.err = err
						return ps
					}

					ps.Attribute.Resources = append(ps.Attribute.Resources, *instRes)
				}
			}
		}
		return ps
	}

	// delete object's instance association operation. for web
	if ps.hitRegexp(deleteObjectInstanceAssociationLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object's instance association, but got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[4]
		if len(objID) == 0 {
			ps.err = fmt.Errorf("delete object instance association, but got empty object id")
			return ps
		}

		assoID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance association, but got invalid association id %s",
				ps.RequestCtx.Elements[5])
			return ps
		}

		asst, err := ps.getInstAssociation(objID, mapstr.MapStr{common.BKFieldID: assoID})
		if err != nil {
			ps.err = err
			return ps
		}
		models, err := ps.searchModels(mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: []interface{}{
			asst.ObjectID,
			asst.AsstObjectID,
		}}})
		if err != nil {
			ps.err = err
			return ps
		}

		// 处理模型自关联的情况
		if len(models) == 1 {
			instRes, err := ps.generateUpdateInstanceResource(&models[0], asst.InstID)
			if err != nil {
				ps.err = err
				return ps
			}
			asstInstRes, err := ps.generateUpdateInstanceResource(&models[0], asst.AsstInstID)
			if err != nil {
				ps.err = err
				return ps
			}

			ps.Attribute.Resources = []meta.ResourceAttribute{*instRes, *asstInstRes}
			return ps
		}

		for _, model := range models {
			var instID int64
			if model.ObjectID == asst.ObjectID {
				instID = asst.InstID
			} else {
				instID = asst.AsstInstID
			}

			instRes, err := ps.generateUpdateInstanceResource(&model, instID)
			if err != nil {
				ps.err = err
				return ps
			}

			ps.Attribute.Resources = append(ps.Attribute.Resources, *instRes)
		}

		return ps
	}

	// delete object instance's association batch operation.
	if ps.hitRegexp(deleteObjectInstanceAssociationBatchLatestRegexp, http.MethodDelete) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		val, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		objID := val.Value()
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.DeleteMany,
				},
			},
		}
		return ps
	}

	// find object instance topology ui operation.
	if ps.hitRegexp(findObjectInstanceTopologyUILatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find object instance's association object instance info operation.
	if ps.hitRegexp(findInstAssociationObjInstInfoLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	if ps.hitRegexp(searchInstAssociationAndInstDetailLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	findObjectInstancesNamesPattern = "/api/v3/findmany/object/instances/names"
)

var (
	createObjectInstanceLatestRegexp             = regexp.MustCompile(`^/api/v3/create/instance/object/[^\s/]+/?$`)
	createObjectManyInstanceByImportLatestRegexp = regexp.MustCompile(
		`^/api/v3/create/instance/object/[^\s/]+/by_import/?$`)
	createObjectManyInstanceLatestRegexp      = regexp.MustCompile(`^/api/v3/createmany/instance/object/[^\s/]+/?$`)
	findObjectInstanceAssociationLatestRegexp = regexp.MustCompile(`^/api/v3/find/instassociation/object/[^\s/]+/?$`)
	updateObjectInstanceLatestRegexp          = regexp.MustCompile(
		`^/api/v3/update/instance/object/[^\s/]+/inst/[0-9]+/?$`)
	updateObjectInstanceBatchLatestRegexp = regexp.MustCompile(`^/api/v3/updatemany/instance/object/[^\s/]+/?$`)
	deleteObjectInstanceBatchLatestRegexp = regexp.MustCompile(`^/api/v3/deletemany/instance/object/[^\s/]+/?$`)
	deleteObjectInstanceLatestRegexp      = regexp.MustCompile(
		`^/api/v3/delete/instance/object/[^\s/]+/inst/[0-9]+/?$`)
	// TODO remove it
	findObjectInstanceSubTopologyLatestRegexp = regexp.MustCompile(
		`^/api/v3/find/insttopo/object/[^\s/]+/inst/[0-9]+/?$`)
	findObjectInstanceTopologyLatestRegexp = regexp.MustCompile(
		`^/api/v3/find/instassttopo/object/[^\s/]+/inst/[0-9]+/?$`)
	findObjectInstancesLatestRegexp       = regexp.MustCompile(`^/api/v3/find/instance/object/[^\s/]+/?$`)
	findObjectInstancesUniqueFieldsRegexp = regexp.MustCompile(
		`^/api/v3/find/instance/object/[^\s/]+/unique_fields/by/unique/[0-9]+/?$`)

	searchObjectInstancesRegexp = regexp.MustCompile(`^/api/v3/search/instances/object/[^\s/]+/?$`)
	countObjectInstancesRegexp  = regexp.MustCompile(`^/api/v3/count/instances/object/[^\s/]+/?$`)
	// excel 导入主机专用接口
	findObjectInstancesForExcelRegexp = regexp.MustCompile(`^/api/v3/find/instance/[^\s/]+/?$`)
)

// NOCC:golint/fnsize(设计如此)
func (ps *parseStream) objectInstanceLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create instance operation
	if ps.hitRegexp(createObjectInstanceLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("create instance, but got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.Create,
				},
			},
		}

		return ps
	}

	// create instance operation
	if ps.hitRegexp(createObjectManyInstanceByImportLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("create instance, but got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.Create,
				},
			},
		}

		return ps
	}

	if ps.hitRegexp(createObjectManyInstanceLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("create instance, but got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.Create,
				},
			},
		}

		return ps
	}

	// search instance association
	if ps.hitRegexp(findObjectInstanceAssociationLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// update instance operation
	if ps.hitRegexp(updateObjectInstanceLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update object instance, but got invalid url")
			return ps
		}

		instID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object instance, but got invalid instance id %s", ps.RequestCtx.Elements[5])
			return ps
		}

		objectID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objectID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       instanceType,
					Action:     meta.Update,
					InstanceID: instID,
				},
			},
		}
		return ps
	}

	// batch update instance operation
	if ps.hitRegexp(updateObjectInstanceBatchLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("update object instance batch, but got invalid url")
			return ps
		}

		objectID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objectID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		ids := make([]int64, 0)
		val, err := ps.RequestCtx.getValueFromBody("update.#.inst_id")
		if err != nil {
			ps.err = err
			return ps
		}
		val.ForEach(
			func(key, value gjson.Result) bool {
				ids = append(ids, value.Int())
				return true
			})

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		for _, id := range ids {
			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       instanceType,
					Action:     meta.UpdateMany,
					InstanceID: id,
				},
			})
		}

		return ps
	}

	// batch delete instance operation
	if ps.hitRegexp(deleteObjectInstanceBatchLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("delete object instance batch, but got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		isMainline, err := ps.isMainlineModel(objID)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		if isMainline {
			// special logic for mainline object's instance authorization.
			instanceType = meta.MainlineInstance
		}

		instIDs := make([]int64, 0)
		val, err := ps.RequestCtx.getValueFromBody("delete.inst_ids")
		if err != nil {
			ps.err = err
			return ps
		}
		val.ForEach(
			func(key, value gjson.Result) bool {
				instIDs = append(instIDs, value.Int())
				return true
			})

		for _, instID := range instIDs {
			ps.Attribute.Resources = append(ps.Attribute.Resources, meta.ResourceAttribute{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       instanceType,
					Action:     meta.Delete,
					InstanceID: instID,
				},
			})
		}
		return ps
	}

	// delete instance operation.
	if ps.hitRegexp(deleteObjectInstanceLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("delete object instance, but got invalid url")
			return ps
		}

		e7 := ps.RequestCtx.Elements[7]
		instID, err := strconv.ParseInt(e7, 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object instance, but got invalid instance id %s", e7)
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       instanceType,
					Action:     meta.Delete,
					InstanceID: instID,
				},
			},
		}

		return ps
	}

	// find object instance sub topology operation
	if ps.hitRegexp(findObjectInstanceSubTopologyLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find object instance fully topology operation.
	if ps.hitRegexp(findObjectInstanceTopologyLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("find object instance topology, but got invalid url")
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find object's instance list operation
	if ps.hitRegexp(findObjectInstancesLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// search object instances operation.
	if ps.hitRegexp(searchObjectInstancesRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// count object instances operation.
	if ps.hitRegexp(countObjectInstancesRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("count object instances, got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		if len(objID) == 0 {
			ps.err = fmt.Errorf("find object instance topology ui failed, got empty object id")
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find object's instances' unique fields operation
	if ps.hitRegexp(findObjectInstancesUniqueFieldsRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 10 {
			ps.err = errors.New("find object's instances' unique fields, but got invalid url")
			return ps
		}

		objID := ps.RequestCtx.Elements[5]
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find object instances' names to get set/module name for host advanced filter
	if ps.hitPattern(findObjectInstancesNamesPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		objID := val.Value()

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: objID})
		if err != nil {
			ps.err = err
			return ps
		}
		instanceType, err := ps.getInstanceTypeByObject(model.ObjectID, model.ID)
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   instanceType,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// excel 导入主机专用接口, 跳过鉴权
	if ps.hitRegexp(findObjectInstancesForExcelRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectLatestPattern      = "/api/v3/create/object"
	createManyObjectLatestPattern  = "/api/v3/createmany/object/by_import"
	createObjectBatchLatestPattern = "/api/v3/createmany/object"
	findObjectsLatestPattern       = "/api/v3/find/object"
	findObjectsModelLatestPattern  = "/api/v3/find/object/model"

	findObjectBatchLatestPattern         = "/api/v3/findmany/object"
	findObjectWithTotalInfoLatestPattern = "/api/v3/findmany/object/total/info"
	findObjectTopologyLatestPattern      = "/api/v3/find/objecttopology"
)

var (
	deleteObjectLatestRegexp = regexp.MustCompile(`^/api/v3/delete/object/[0-9]+/?$`)
	updateObjectLatestRegexp = regexp.MustCompile(`^/api/v3/update/object/[0-9]+/?$`)

	// TODO remove it
	// 获取模型拓扑图及位置信息-Web
	findObjectTopologyGraphicLatestRegexp = regexp.MustCompile(`^/api/v3/find/objecttopo/scope_type/[^\s/]+/scope_id/[^\s/]+/?$`)
	// 设置模型拓扑图及位置信息-Web
	updateObjectTopologyGraphicLatestRegexp = regexp.MustCompile(`^/api/v3/update/objecttopo/scope_type/[^\s/]+/scope_id/[^\s/]+/?$`)
)

func (ps *parseStream) objectLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create common object operation.
	if ps.hitPattern(createObjectLatestPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.BKClassificationIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		classID := val.String()
		filter := map[string]interface{}{
			common.BKClassificationIDField: classID,
		}

		classification, err := ps.getOneClassification(filter)
		if err != nil {
			ps.err = fmt.Errorf("create object, but get classification failed, err:%v", err)
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.ModelClassification, InstanceID: classification.ID}},
			},
		}
		return ps
	}

	// create many common object by yaml import.
	if ps.hitPattern(createManyObjectLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// create or update models' attributes batch
	if ps.hitPattern(createObjectBatchLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// delete object operation
	if ps.hitRegexp(deleteObjectLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object, but got invalid url")
			return ps
		}

		id, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object, but got invalid object's id %s", ps.RequestCtx.Elements[3])
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Delete,
					InstanceID: id,
				},
			},
		}
		return ps
	}

	// update object operation.
	if ps.hitRegexp(updateObjectLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object, but got invalid url")
			return ps
		}

		id, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object, but got invalid object's id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.Model,
					Action:     meta.Update,
					InstanceID: id,
				},
			},
		}
		return ps
	}

	// get object operation.
	if ps.hitPattern(findObjectsLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// get object model operation.
	if ps.hitPattern(findObjectsModelLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.FindMany,
				},
			},
		}
		return ps
	}

	// get object in batch operation.
	if ps.hitPattern(findObjectBatchLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	if ps.hitPattern(findObjectWithTotalInfoLatestPattern, http.MethodPost) {

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// find object's topology operation.
	if ps.hitPattern(findObjectTopologyLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// find object's topology graphic operation.
	if ps.hitRegexp(findObjectTopologyGraphicLatestRegexp, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.ModelTopology,
					// Action: meta.Find,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// update object's topology graphic operation.
	if ps.hitRegexp(updateObjectTopologyGraphicLatestRegexp, http.MethodPost) {

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelTopology,
					Action: meta.Update,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectClassificationLatestPattern   = "/api/v3/create/objectclassification"
	findObjectClassificationListLatestPattern = "/api/v3/find/objectclassification"
	// 查找模型分组及分组下的模型列表
	findObjectsBelongsToClassificationLatestPattern = `/api/v3/find/classificationobject`
)

var (
	deleteObjectClassificationLatestRegexp = regexp.MustCompile("^/api/v3/delete/objectclassification/[0-9]+/?$")
	updateObjectClassificationLatestRegexp = regexp.MustCompile("^/api/v3/update/objectclassification/[0-9]+/?$")
)

func (ps *parseStream) objectClassificationLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's classification operation.
	if ps.hitPattern(createObjectClassificationLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete object's classification operation.
	if ps.hitRegexp(deleteObjectClassificationLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object classification, but got invalid url")
			return ps
		}

		classID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object classification, but got invalid object's id %s",
				ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelClassification,
					Action:     meta.Delete,
					InstanceID: classID,
				},
			},
		}
		return ps
	}

	// update object's classification operation.
	if ps.hitRegexp(updateObjectClassificationLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object classification, but got invalid url")
			return ps
		}

		classID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object classification, but got invalid object's  classification id %s",
				ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelClassification,
					Action:     meta.Update,
					InstanceID: classID,
				},
			},
		}
		return ps
	}

	// find object's classification list operation.
	if ps.hitPattern(findObjectClassificationListLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelClassification,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}
	// find all the objects belongs to a classification
	if ps.hitPattern(findObjectsBelongsToClassificationLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.Model,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectAttributeGroupLatestPattern   = "/api/v3/create/objectattgroup"
	updateObjectAttributeGroupLatestPattern   = "/api/v3/update/objectattgroup"
	exchangeObjectAttributeGroupLatestPattern = "/api/v3/update/objectattgroup/groupindex"
)

var (
	findObjectAttributeGroupLatestRegexp   = regexp.MustCompile(`^/api/v3/find/objectattgroup/object/[^\s/]+/?$`)
	deleteObjectAttributeGroupLatestRegexp = regexp.MustCompile(`^/api/v3/delete/objectattgroup/[0-9]+/?$`)
)

func (ps *parseStream) objectAttributeGroupLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}
	// create object's attribute group operation.
	if ps.hitPattern(createObjectAttributeGroupLatestPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		filter := mapstr.MapStr{
			common.BKObjIDField: val.Value(),
		}
		model, err := ps.getOneModel(filter)
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				// 业务ID的解释
				// case  0: 创建公共的属性分组
				// case ~0: 创建业务私有的属性分组
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// find object's attribute group operation.
	if ps.hitRegexp(findObjectAttributeGroupLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelAttributeGroup,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	// update object's attribute group operation.
	if ps.hitPattern(updateObjectAttributeGroupLatestPattern, http.MethodPut) {
		val, err := ps.RequestCtx.getValueFromBody("condition")
		if err != nil {
			ps.err = err
			return ps
		}

		groups, err := ps.getAttributeGroup(val.Value())
		if err != nil {
			ps.err = err
			return ps
		}

		for _, group := range groups {
			filter := mapstr.MapStr{
				common.BKObjIDField: group.ObjectID,
			}
			model, err := ps.getOneModel(filter)
			if err != nil {
				ps.err = err
				return ps
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					BusinessID: group.BizID,
					Basic: meta.Basic{
						Type:       meta.ModelAttributeGroup,
						Action:     meta.Update,
						InstanceID: group.ID,
					},
					Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				})
		}
		return ps
	}

	// exchange objects' attribute group_index
	if ps.hitPattern(exchangeObjectAttributeGroupLatestPattern, http.MethodPut) {
		val, err := ps.RequestCtx.getValueFromBody("condition.id")
		if err != nil {
			ps.err = err
			return ps
		}

		cond := mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: val.Value()}}
		groups, err := ps.getAttributeGroup(cond)
		if err != nil {
			ps.err = err
			return ps
		}

		for _, group := range groups {
			filter := mapstr.MapStr{
				common.BKObjIDField: group.ObjectID,
			}
			model, err := ps.getOneModel(filter)
			if err != nil {
				ps.err = err
				return ps
			}
			ps.Attribute.Resources = append(ps.Attribute.Resources,
				meta.ResourceAttribute{
					BusinessID: group.BizID,
					Basic: meta.Basic{
						Type:       meta.ModelAttributeGroup,
						Action:     meta.Update,
						InstanceID: group.ID,
					},
					Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
				})
		}
		return ps
	}

	// delete object's attribute group operation.
	if ps.hitRegexp(deleteObjectAttributeGroupLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object's attribute group, but got invalid url")
			return ps
		}

		groupID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object's attribute group, but got invalid group's id %s",
				ps.RequestCtx.Elements[4])
			return ps
		}

		groups, err := ps.getAttributeGroup(mapstr.MapStr{"id": groupID})
		if err != nil {
			ps.err = err
			return ps
		}

		if len(groups) == 0 {
			ps.err = fmt.Errorf("object attribute group(id: %d) is not exist", groupID)
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: groups[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: groups[0].BizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttributeGroup,
					Action:     meta.Delete,
					InstanceID: groupID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	return ps
}

const (
	createObjectAttributeLatestPattern  = "/api/v3/create/objectattr"
	findObjectAttributeLatestPattern    = "/api/v3/find/objectattr"
	findObjectAttributeForLatestPattern = "/api/v3/find/objectattr/web"
)

var (
	deleteObjectAttributeLatestRegexp      = regexp.MustCompile(`^/api/v3/delete/objectattr/[0-9]+/?$`)
	updateObjectAttributeLatestRegexp      = regexp.MustCompile(`^/api/v3/update/objectattr/[0-9]+/?$`)
	updateObjectAttributeIndexLatestRegexp = regexp.MustCompile(`^/api/v3/update/objectattr/index/[^\s/]+/[0-9]+/?$`)
	createBizCustomFieldLatestRegexp       = regexp.MustCompile(`^/api/v3/create/objectattr/biz/[0-9]+/?$`)
	updateBizCustomFieldLatestRegexp       = regexp.MustCompile(`^/api/v3/update/objectattr/biz/[0-9]+/id/[0-9]+/?$`)
)

func (ps *parseStream) objectAttributeLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create object's attribute operation.
	if ps.hitPattern(createObjectAttributeLatestPattern, http.MethodPost) {
		val, err := ps.RequestCtx.getValueFromBody(common.BKObjIDField)
		if err != nil {
			ps.err = err
			return ps
		}
		modelEn := val.String()
		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: modelEn})
		if err != nil {
			ps.err = err
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				// 注意业务ID是否为0表示创建两种不同的属性
				// case 0: 创建公共属性，这种属性相比业务私有属性，所有业务都可见
				// case ~0: 创建业务私有属性，业务私有属性，其它业务不可见
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.Create,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}

		return ps
	}

	// delete object's attribute operation.
	if ps.hitRegexp(deleteObjectAttributeLatestRegexp, http.MethodDelete) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("delete object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		attr, err := ps.getModelAttribute(bizID, mapstr.MapStr{common.BKFieldID: attrID})
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but fetch attribute by %v failed %v",
				mapstr.MapStr{common.BKFieldID: attrID}, err)
			return ps
		}

		if len(attr) == 0 {
			ps.err = errors.New("can not find attribute detail")
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: attr[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				// 对属性操作的鉴权，依赖于属性是公有属性，还是业务私有属性
				BusinessID: attr[0].BizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Delete,
					InstanceID: attrID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// update object attribute operation
	if ps.hitRegexp(updateObjectAttributeLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 5 {
			ps.err = errors.New("update object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[4], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		attr, err := ps.getModelAttribute(0, mapstr.MapStr{common.BKFieldID: attrID})
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but fetch attribute by %v failed %v",
				mapstr.MapStr{common.BKFieldID: attrID}, err)
			return ps
		}

		if len(attr) == 0 {
			ps.err = errors.New("can not find attribute detail")
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: attr[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				// 对属性操作的鉴权，依赖于属性是公有属性，还是业务私有属性
				BusinessID: attr[0].BizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Update,
					InstanceID: attrID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// update object attribute index operation
	if ps.hitRegexp(updateObjectAttributeIndexLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("update object attribute, but got invalid url")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[6], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update object attribute, but got invalid attribute id %s", ps.RequestCtx.Elements[4])
			return ps
		}

		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		attr, err := ps.getModelAttribute(bizID, mapstr.MapStr{common.BKFieldID: attrID})
		if err != nil {
			ps.err = fmt.Errorf("delete object attribute, but fetch attribute by %v failed %v",
				mapstr.MapStr{common.BKFieldID: attrID}, err)
			return ps
		}

		if len(attr) == 0 {
			ps.err = errors.New("can not find attribute detail")
			return ps
		}

		model, err := ps.getOneModel(mapstr.MapStr{common.BKObjIDField: attr[0].ObjectID})
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				// 对属性操作的鉴权，依赖于属性是公有属性，还是业务私有属性
				BusinessID: attr[0].BizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Update,
					InstanceID: attrID,
				},
				Layers: []meta.Item{{Type: meta.Model, InstanceID: model.ID}},
			},
		}
		return ps
	}

	// get object's attribute operation.
	if ps.hitPattern(findObjectAttributeLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = append(ps.Attribute.Resources,
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.SkipAction,
				},
			})
		return ps
	}

	// get object's attribute operation.
	if ps.hitPattern(findObjectAttributeForLatestPattern, http.MethodPost) {
		ps.Attribute.Resources = append(ps.Attribute.Resources,
			meta.ResourceAttribute{
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.SkipAction,
				},
			})
	}

	// create business custom field operation
	if ps.hitRegexp(createBizCustomFieldLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("create business custom field, but got invalid url")
			return ps
		}

		// 创建业务自定义属性，业务ID不可为0
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			blog.Errorf("get business id failed, err: %v", err)
			ps.err = err
			return ps
		}

		if bizID == 0 {
			blog.Error("biz custom field business id can't be 0")
			ps.err = fmt.Errorf("create biz custom field failed, business id is 0")
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelAttribute,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// update business custom field operation
	if ps.hitRegexp(updateBizCustomFieldLatestRegexp, http.MethodPut) {
		if len(ps.RequestCtx.Elements) != 8 {
			ps.err = errors.New("update business custom field, but got invalid url")
			return ps
		}

		// 更新业务自定义属性，业务ID不可为0
		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			blog.Errorf("get business id failed, err: %v", err)
			ps.err = err
			return ps
		}

		if bizID == 0 {
			blog.Error("biz custom field business id can't be 0")
			ps.err = fmt.Errorf("update biz custom field failed, business id is 0")
			return ps
		}

		attrID, err := strconv.ParseInt(ps.RequestCtx.Elements[7], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("update business custom field, but got invalid attribute id %s",
				ps.RequestCtx.Elements[4])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:       meta.ModelAttribute,
					Action:     meta.Update,
					InstanceID: attrID,
				},
			},
		}
		return ps
	}

	return ps
}

const (
	createMainlineObjectLatestPattern   = "/api/v3/create/topomodelmainline"
	findMainlineObjectTopoLatestPattern = "/api/v3/find/topomodelmainline"
)

var (
	deleteMainlineObjectLatestRegexp                       = regexp.MustCompile(`^/api/v3/delete/topomodelmainline/object/[^\s/]+/?$`)
	findBusinessInstanceTopologyLatestRegexp               = regexp.MustCompile(`^/api/v3/find/topoinst/biz/[0-9]+/?$`)
	findBusinessInstanceTopologyPathRegexp                 = regexp.MustCompile(`^/api/v3/find/topopath/biz/[0-9]+/?$`)
	findHostApplyRelatedObjectTopologyRegex                = regexp.MustCompile(`^/api/v3/find/topoinst/bk_biz_id/([0-9]+)/host_apply_rule_related/?$`)
	findBusinessInstanceTopologyWithStatisticsLatestRegexp = regexp.MustCompile(`^/api/v3/find/topoinst_with_statistics/biz/[0-9]+/?$`)
	findTopoNodeHostAndServiceInstCountLatestRegexp        = regexp.MustCompile(
		`^/api/v3/find/topoinstnode/host_serviceinst_count/[0-9]+/?$`)
)

func (ps *parseStream) mainlineLatest() *parseStream {
	if ps.shouldReturn() {
		return ps
	}

	// create mainline object operation.
	if ps.hitPattern(createMainlineObjectLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Create,
				},
			},
		}
		return ps
	}

	// delete mainline object operation
	if ps.hitRegexp(deleteMainlineObjectLatestRegexp, http.MethodDelete) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineModel,
					Action: meta.Delete,
				},
			},
		}

		return ps
	}

	// get mainline object operation
	if ps.hitPattern(findMainlineObjectTopoLatestPattern, http.MethodPost) {
		bizID, err := ps.RequestCtx.getBizIDFromBody()
		if err != nil {
			ps.err = err
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type: meta.MainlineModelTopology,
					// Action: meta.Find,
					Action: meta.SkipAction,
				},
			},
		}

		return ps
	}

	// 根据主机属性自动应用规则查找拓扑节点
	if ps.hitRegexp(findHostApplyRelatedObjectTopologyRegex, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 7 {
			ps.err = errors.New("find host apply rule related topo node, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("find host apply rule related topo node, but got invalid business id %s",
				ps.RequestCtx.Elements[5])
			return ps
		}
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.MainlineInstanceTopology,
					Action: meta.Find,
				},
			},
		}

		return ps
	}

	// find business instance topology operation.
	// also is find mainline instance topology operation.
	if ps.hitRegexp(findBusinessInstanceTopologyLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find business instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse biz id from url failed, but got invalid business id %s",
				ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// find business instance topology operation.
	// also is find mainline instance topology operation.
	if ps.hitRegexp(findBusinessInstanceTopologyPathRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find business instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse biz id from url failed, but got invalid business id %s",
				ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// find business instance topology operation.
	// also is find mainline instance topology operation.
	if ps.hitRegexp(findBusinessInstanceTopologyWithStatisticsLatestRegexp, http.MethodPost) {
		if len(ps.RequestCtx.Elements) != 6 {
			ps.err = errors.New("find business instance topology, but got invalid url")
			return ps
		}

		bizID, err := strconv.ParseInt(ps.RequestCtx.Elements[5], 10, 64)
		if err != nil {
			ps.err = fmt.Errorf("parse biz id from url failed, but got invalid business id %s",
				ps.RequestCtx.Elements[5])
			return ps
		}

		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				BusinessID: bizID,
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.Find,
				},
			},
		}
		return ps
	}

	// get toponode host and service instance count
	if ps.hitRegexp(findTopoNodeHostAndServiceInstCountLatestRegexp, http.MethodPost) {
		ps.Attribute.Resources = []meta.ResourceAttribute{
			{
				Basic: meta.Basic{
					Type:   meta.ModelInstanceTopology,
					Action: meta.SkipAction,
				},
			},
		}
		return ps
	}

	return ps
}
