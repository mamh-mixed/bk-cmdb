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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

var needTranslateObjMap = map[string]bool{
	common.BKInnerObjIDApp:      true,
	common.BKInnerObjIDSet:      true,
	common.BKInnerObjIDModule:   true,
	common.BKInnerObjIDProc:     true,
	common.BKInnerObjIDHost:     true,
	common.BKInnerObjIDPlat:     true,
	common.BKInnerObjIDSwitch:   true,
	common.BKInnerObjIDRouter:   true,
	common.BKInnerObjIDBlance:   true,
	common.BKInnerObjIDFirewall: true,
	common.BKInnerObjIDBizSet:   true,
	common.BKInnerObjIDProject:  true,
}

// TranslateObjectName TODO
func (s *coreService) TranslateObjectName(defLang language.DefaultCCLanguageIf, obj *metadata.Object) string {
	return util.FirstNotEmptyString(defLang.Language("object_"+obj.ObjectID), obj.ObjectName, obj.ObjectID)
}

// TranslateInstName TODO
func (s *coreService) TranslateInstName(defLang language.DefaultCCLanguageIf, obj *metadata.Object) string {
	return util.FirstNotEmptyString(defLang.Language("inst_"+obj.ObjectID), obj.ObjectName, obj.ObjectID)
}

// TranslatePropertyName translate property name
func (s *coreService) TranslatePropertyName(defLang language.DefaultCCLanguageIf, att *metadata.Attribute) string {
	var prefix string
	switch att.PropertyID {
	case common.BKInstNameField, common.BKCreatedAt, common.BKCreatedBy, common.BKUpdatedAt, common.BKUpdatedBy:
		prefix = "common"
	default:
		prefix = att.ObjectID
	}

	return util.FirstNotEmptyString(defLang.Language(prefix+"_property_"+att.PropertyID), att.PropertyName,
		att.PropertyID)
}

// TranslatePlaceholder TODO
func (s *coreService) TranslatePlaceholder(defLang language.DefaultCCLanguageIf, att *metadata.Attribute) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_placeholder_"+att.PropertyID), att.Placeholder)
}

// TranslateEnumName TODO
func (s *coreService) TranslateEnumName(ctx context.Context, defLang language.DefaultCCLanguageIf,
	att *metadata.Attribute, val interface{}) interface{} {
	rid := util.ExtractRequestIDFromContext(ctx)
	options, err := metadata.ParseEnumOption(val)
	if err != nil {
		blog.Warnf("TranslateEnumName failed: %v, rid: %s", err, rid)
		return val
	}
	for index := range options {
		options[index].Name = util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_"+att.PropertyID+"_enum_"+options[index].ID),
			options[index].Name, options[index].ID)
	}
	return options
}

// TranslateTablePlaceholder table类型属性字段中英文转化
func (s *coreService) TranslateTablePlaceholder(defLang language.DefaultCCLanguageIf, objID string, propertyID string,
	optPropertyID string) string {
	return util.FirstNotEmptyString(defLang.Language(objID + "_" + propertyID + "_placeholder_" + optPropertyID))
}

// TranslatePropertyGroupName TODO
func (s *coreService) TranslatePropertyGroupName(defLang language.DefaultCCLanguageIf, att *metadata.Group) string {
	return util.FirstNotEmptyString(defLang.Language(att.ObjectID+"_property_group_"+att.GroupID), att.GroupName,
		att.GroupID)
}

// TranslateClassificationName TODO
func (s *coreService) TranslateClassificationName(defLang language.DefaultCCLanguageIf,
	att *metadata.Classification) string {
	return util.FirstNotEmptyString(defLang.Language("classification_"+att.ClassificationID), att.ClassificationName,
		att.ClassificationID)
}

// TranslateOperationChartName TODO
func (s *coreService) TranslateOperationChartName(defLang language.DefaultCCLanguageIf,
	att metadata.ChartConfig) string {
	return util.FirstNotEmptyString(defLang.Language("operation_chart_"+att.ReportType), att.Name, att.ReportType)
}

// TranslateAssociationType translate association type
func (s *coreService) TranslateAssociationType(defLang language.DefaultCCLanguageIf,
	assKind *metadata.AssociationKind) {

	assKind.AssociationKindName = util.FirstNotEmptyString(defLang.Language("unique_kind_name_"+
		assKind.AssociationKindID), assKind.AssociationKindName)
	assKind.SourceToDestinationNote = util.FirstNotEmptyString(defLang.Language("unique_kind_src_to_dest_"+
		assKind.AssociationKindID), assKind.SourceToDestinationNote)
	assKind.DestinationToSourceNote = util.FirstNotEmptyString(defLang.Language("unique_kind_dest_to_src_"+
		assKind.AssociationKindID), assKind.DestinationToSourceNote)
}

// TranslateServiceCategory TODO
func (s *coreService) TranslateServiceCategory(defLang language.DefaultCCLanguageIf,
	att *metadata.ServiceCategory) string {
	return util.FirstNotEmptyString(defLang.Language("service_category_"+strings.Replace(att.Name, " ", "_", -1)),
		att.Name)
}
