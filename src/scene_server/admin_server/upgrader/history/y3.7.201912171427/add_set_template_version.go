/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_7_201912171427

import (
	"context"
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader/history"
	"configcenter/src/storage/dal"
)

func addSetTemplateDefaultVersion(ctx context.Context, db dal.RDB, conf *history.Config) error {
	filter := map[string]interface{}{
		"version": map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	doc := map[string]interface{}{
		"version": 0,
	}
	if err := db.Table(common.BKTableNameSetTemplate).Update(ctx, filter, doc); err != nil {
		return fmt.Errorf("addSetTemplateDefaultVersion failed, err: %+v", err)
	}
	return nil
}

func addSetDefaultVersion(ctx context.Context, db dal.RDB, conf *history.Config) error {
	filter := map[string]interface{}{
		"set_template_version": map[string]interface{}{
			common.BKDBExists: false,
		},
	}
	doc := map[string]interface{}{
		"set_template_version": 0,
	}
	if err := db.Table(common.BKTableNameBaseSet).Update(ctx, filter, doc); err != nil {
		return fmt.Errorf("addSetDefaultVersion failed, err: %+v", err)
	}
	return nil
}

func addSetVersionField(ctx context.Context, db dal.RDB, conf *history.Config) error {
	filter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDSet,
		common.BKPropertyIDField: "set_template_version",
	}
	count, err := db.Table(common.BKTableNameObjAttDes).Find(filter).Count(ctx)
	if err != nil {
		return fmt.Errorf("check whether set_template_version attribute exist failed, err: %+v", err)
	}
	if count != 0 {
		return nil
	}

	id, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		return fmt.Errorf("generate attribute id failed, err: %+v", err)
	}

	now := metadata.Now()
	attribute := Attribute{
		ID:                int64(id),
		OwnerID:           conf.TenantID,
		ObjectID:          common.BKInnerObjIDSet,
		PropertyID:        "set_template_version",
		PropertyName:      "集群模板",
		PropertyGroup:     "default",
		PropertyGroupName: "default",
		PropertyIndex:     0,
		Unit:              "",
		Placeholder:       "",
		IsEditable:        true,
		IsPre:             true,
		IsRequired:        false,
		IsReadOnly:        true,
		IsOnly:            false,
		// IsSystem = true 时，字段标记系统内部使用的字段，不会返回到前端
		IsSystem: true,
		// IsAPI = true 时，字段对页面不可见
		IsAPI:        true,
		PropertyType: "int",
		Option:       "",
		Description:  "集群版本，从通集群模板同步",
		Creator:      conf.User,
		CreateTime:   &now,
		LastTime:     &now,
	}
	uniqueFields := []string{common.BKObjIDField, common.BKPropertyIDField, "bk_supplier_account"}
	if _, _, err := history.Upsert(ctx, db, common.BKTableNameObjAttDes, attribute, "id", uniqueFields,
		[]string{}); err != nil {
		blog.Errorf("addSetVersionField failed, add set_template_version attribute failed, err: %+v", err)
		return fmt.Errorf("add set_template_version attribute failed, err: %+v", err)
	}
	return nil
}

// Attribute attribute metadata definition
type Attribute struct {
	BizID             int64          `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" mapstructure:"bk_biz_id"`
	ID                int64          `field:"id" json:"id" bson:"id" mapstructure:"id"`
	OwnerID           string         `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	ObjectID          string         `field:"bk_obj_id" json:"bk_obj_id" bson:"bk_obj_id" mapstructure:"bk_obj_id"`
	PropertyID        string         `field:"bk_property_id" json:"bk_property_id" bson:"bk_property_id" mapstructure:"bk_property_id"`
	PropertyName      string         `field:"bk_property_name" json:"bk_property_name" bson:"bk_property_name" mapstructure:"bk_property_name"`
	PropertyGroup     string         `field:"bk_property_group" json:"bk_property_group" bson:"bk_property_group" mapstructure:"bk_property_group"`
	PropertyGroupName string         `field:"bk_property_group_name,ignoretomap" json:"bk_property_group_name" bson:"-" mapstructure:"bk_property_group_name"`
	PropertyIndex     int64          `field:"bk_property_index" json:"bk_property_index" bson:"bk_property_index" mapstructure:"bk_property_index"`
	Unit              string         `field:"unit" json:"unit" bson:"unit" mapstructure:"unit"`
	Placeholder       string         `field:"placeholder" json:"placeholder" bson:"placeholder" mapstructure:"placeholder"`
	IsEditable        bool           `field:"editable" json:"editable" bson:"editable" mapstructure:"editable"`
	IsPre             bool           `field:"ispre" json:"ispre" bson:"ispre" mapstructure:"ispre"`
	IsRequired        bool           `field:"isrequired" json:"isrequired" bson:"isrequired" mapstructure:"isrequired"`
	IsReadOnly        bool           `field:"isreadonly" json:"isreadonly" bson:"isreadonly" mapstructure:"isreadonly"`
	IsOnly            bool           `field:"isonly" json:"isonly" bson:"isonly" mapstructure:"isonly"`
	IsSystem          bool           `field:"bk_issystem" json:"bk_issystem" bson:"bk_issystem" mapstructure:"bk_issystem"`
	IsAPI             bool           `field:"bk_isapi" json:"bk_isapi" bson:"bk_isapi" mapstructure:"bk_isapi"`
	PropertyType      string         `field:"bk_property_type" json:"bk_property_type" bson:"bk_property_type" mapstructure:"bk_property_type"`
	Option            interface{}    `field:"option" json:"option" bson:"option" mapstructure:"option"`
	Default           interface{}    `field:"default" json:"default,omitempty" bson:"default" mapstructure:"default"`
	IsMultiple        *bool          `field:"ismultiple" json:"ismultiple,omitempty" bson:"ismultiple" mapstructure:"ismultiple"`
	Description       string         `field:"description" json:"description" bson:"description" mapstructure:"description"`
	TemplateID        int64          `field:"bk_template_id" json:"bk_template_id" bson:"bk_template_id" mapstructure:"bk_template_id"`
	Creator           string         `field:"creator" json:"creator" bson:"creator" mapstructure:"creator"`
	CreateTime        *metadata.Time `json:"create_time" bson:"create_time" mapstructure:"create_time"`
	LastTime          *metadata.Time `json:"last_time" bson:"last_time" mapstructure:"last_time"`
}
