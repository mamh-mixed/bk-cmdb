// Package x18_12_12_06 TODO
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
package x18_12_12_06

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func keyfunc(a, b string) string { return a + ":" + b }

func reconcilUnique(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	attrs := []metadata.Attribute{}

	attrCond := condition.CreateCondition()
	attrCond.Field(common.BKObjIDField).In([]string{
		common.BKInnerObjIDSwitch,
		common.BKInnerObjIDRouter,
		common.BKInnerObjIDBlance,
		common.BKInnerObjIDFirewall,
	})

	objStrMapArr := make([]map[string]string, 0)
	err := db.Table(common.BKTableNameObjDes).Find(attrCond.ToMapStr()).Fields(common.BKObjIDField).All(ctx, &objStrMapArr)
	if err != nil {
		return err
	}

	err = db.Table(common.BKTableNameObjAttDes).Find(attrCond.ToMapStr()).All(ctx, &attrs)
	if err != nil {
		return err
	}
	var propertyIDToProperty = map[string]metadata.Attribute{}

	for _, attr := range attrs {
		propertyIDToProperty[keyfunc(attr.ObjectID, attr.PropertyID)] = attr
	}

	var shouldCheck []string
	var uniques []objectUnique
	for _, objMap := range objStrMapArr {
		objID := objMap[common.BKObjIDField]
		switch objID {
		case common.BKInnerObjIDSwitch:
			shouldCheck = append(shouldCheck, keyfunc(common.BKInnerObjIDSwitch, common.BKAssetIDField))
			uniques = append(uniques, buildUnique(propertyIDToProperty, common.BKInnerObjIDSwitch, common.BKAssetIDField))

		case common.BKInnerObjIDRouter:
			shouldCheck = append(shouldCheck, keyfunc(common.BKInnerObjIDRouter, common.BKAssetIDField))
			uniques = append(uniques, buildUnique(propertyIDToProperty, common.BKInnerObjIDRouter, common.BKAssetIDField))

		case common.BKInnerObjIDBlance:
			shouldCheck = append(shouldCheck, keyfunc(common.BKInnerObjIDBlance, common.BKAssetIDField))
			uniques = append(uniques, buildUnique(propertyIDToProperty, common.BKInnerObjIDBlance, common.BKAssetIDField))

		case common.BKInnerObjIDFirewall:
			shouldCheck = append(shouldCheck, keyfunc(common.BKInnerObjIDFirewall, common.BKAssetIDField))
			uniques = append(uniques, buildUnique(propertyIDToProperty, common.BKInnerObjIDFirewall, common.BKAssetIDField))
		}
	}

	if notExistFields := checkKeysShouldExists(propertyIDToProperty, shouldCheck); len(notExistFields) > 0 {
		return fmt.Errorf("expected field not exists: %v", notExistFields)
	}

	uniqueIDs, err := db.NextSequences(ctx, common.BKTableNameObjUnique, len(uniques))
	if err != nil {
		return err
	}

	for index, unique := range uniques {
		exists, err := isUniqueExists(ctx, db, conf, unique)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		unique.ID = uniqueIDs[index]
		if err := db.Table(common.BKTableNameObjUnique).Insert(ctx, unique); err != nil {
			return err
		}
	}

	return nil
}

func buildUnique(propertyIDToProperty map[string]metadata.Attribute, model, field string) objectUnique {
	return objectUnique{
		ObjID:     model,
		MustCheck: true,
		Keys: []UniqueKey{
			{
				Kind: UniqueKeyKindProperty,
				ID:   uint64(propertyIDToProperty[keyfunc(model, field)].ID),
			},
		},
		Ispre:    false,
		OwnerID:  common.BKDefaultOwnerID,
		LastTime: Now(),
	}
}

func checkKeysShouldExists(m map[string]metadata.Attribute, shouldExistKeys []string) []string {
	notValidKeys := []string{}
	for _, k := range shouldExistKeys {
		if _, ok := m[k]; !ok {
			notValidKeys = append(notValidKeys, k)
		}
	}
	return notValidKeys
}

func isUniqueExists(ctx context.Context, db dal.RDB, conf *upgrader.Config, unique objectUnique) (bool, error) {
	keyhash := unique.KeysHash()
	uniqueCond := condition.CreateCondition()
	uniqueCond.Field(common.BKObjIDField).Eq(unique.ObjID)
	uniqueCond.Field(common.BKOwnerIDField).Eq(conf.OwnerID)
	existUniques := []objectUnique{}

	err := db.Table(common.BKTableNameObjUnique).Find(uniqueCond.ToMapStr()).All(ctx, &existUniques)
	if err != nil {
		return false, err
	}

	for _, uni := range existUniques {
		if uni.KeysHash() == keyhash {
			return true, nil
		}
	}
	return false, nil

}

type objectUnique struct {
	ID        uint64      `json:"id" bson:"id"`
	ObjID     string      `json:"bk_obj_id" bson:"bk_obj_id"`
	MustCheck bool        `json:"must_check" bson:"must_check"`
	Keys      []UniqueKey `json:"keys" bson:"keys"`
	Ispre     bool        `json:"ispre" bson:"ispre"`
	OwnerID   string      `json:"bk_supplier_account" bson:"bk_supplier_account"`
	LastTime  time.Time   `json:"last_time" bson:"last_time"`
}

// Now TODO
func Now() time.Time {
	return time.Now().UTC()
}

// UniqueKey TODO
type UniqueKey struct {
	Kind string `json:"key_kind" bson:"key_kind"`
	ID   uint64 `json:"key_id" bson:"key_id"`
}

// KeysHash TODO
func (o objectUnique) KeysHash() string {
	keys := []string{}
	for _, key := range o.Keys {
		keys = append(keys, fmt.Sprintf("%s:%d", key.Kind, key.ID))
	}
	sort.Strings(keys)
	return strings.Join(keys, "#")
}

const (
	// UniqueKeyKindProperty TODO
	UniqueKeyKindProperty = "property"
)
