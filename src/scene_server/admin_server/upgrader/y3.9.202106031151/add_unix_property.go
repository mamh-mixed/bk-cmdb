// Package y3_9_202106031151 TODO
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
package y3_9_202106031151

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

type attribute struct {
	ID                int64       `json:"id" bson:"id"`
	OwnerID           string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	ObjectID          string      `json:"bk_obj_id" bson:"bk_obj_id"`
	PropertyID        string      `json:"bk_property_id" bson:"bk_property_id"`
	PropertyName      string      `json:"bk_property_name" bson:"bk_property_name"`
	PropertyGroup     string      `json:"bk_property_group" bson:"bk_property_group"`
	PropertyGroupName string      `json:"bk_property_group_name" bson:"-"`
	PropertyIndex     int64       `json:"bk_property_index" bson:"bk_property_index"`
	Unit              string      `json:"unit" bson:"unit"`
	Placeholder       string      `json:"placeholder" bson:"placeholder"`
	IsEditable        bool        `json:"editable" bson:"editable"`
	IsPre             bool        `json:"ispre" bson:"ispre"`
	IsRequired        bool        `json:"isrequired" bson:"isrequired"`
	IsReadOnly        bool        `json:"isreadonly" bson:"isreadonly"`
	IsOnly            bool        `json:"isonly" bson:"isonly"`
	IsSystem          bool        `json:"bk_issystem" bson:"bk_issystem"`
	IsAPI             bool        `json:"bk_isapi" bson:"bk_isapi"`
	PropertyType      string      `json:"bk_property_type" bson:"bk_property_type"`
	Option            interface{} `json:"option" bson:"option"`
	Description       string      `json:"description" bson:"description"`
	Creator           string      `json:"creator" bson:"creator"`
	CreateTime        *time.Time  `json:"create_time" bson:"create_time"`
	LastTime          *time.Time  `json:"last_time" bson:"last_time"`
}

func addUnixProperty(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(common.BKDefaultOwnerID)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)
	cond.Field(common.BKPropertyIDField).Eq(common.BKOSTypeField)
	cond.Field(common.BKAppIDField).Eq(0)

	ostypeProperty := attribute{}
	err := db.Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).One(ctx, &ostypeProperty)
	if err != nil {
		return err
	}

	enumOpts, err := metadata.ParseEnumOption(ostypeProperty.Option)
	if err != nil {
		return err
	}
	for _, enum := range enumOpts {
		if enum.ID == common.HostOSTypeEnumUNIX {
			return nil
		}
	}

	aixEnum := metadata.EnumVal{
		ID:   common.HostOSTypeEnumUNIX,
		Name: "Unix",
		Type: "text",
	}
	enumOpts = append(enumOpts, aixEnum)

	data := mapstr.MapStr{
		common.BKOptionField: enumOpts,
	}

	err = db.Table(common.BKTableNameObjAttDes).Update(ctx, cond.ToMapStr(), data)
	if err != nil {
		return err
	}
	return nil
}
