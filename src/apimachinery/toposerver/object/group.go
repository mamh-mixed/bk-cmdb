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

package object

import (
	"context"
	"net/http"

	"configcenter/src/common/metadata"
)

// CreatePropertyGroup TODO
func (t *object) CreatePropertyGroup(ctx context.Context, h http.Header, dat metadata.Group) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/create/objectattgroup"

	err = t.client.Post().
		WithContext(ctx).
		Body(dat).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdatePropertyGroup TODO
func (t *object) UpdatePropertyGroup(ctx context.Context, h http.Header,
	cond *metadata.PropertyGroupCondition) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objectattgroup"

	err = t.client.Put().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// ExchangePropertyGroupIndex TODO
func (t *object) ExchangePropertyGroupIndex(ctx context.Context, h http.Header, cond *metadata.ExchangeGroupIndex) (
	*metadata.Response, error) {

	resp := new(metadata.Response)
	subPath := "/update/objectattgroup/groupindex"

	err := t.client.Put().
		WithContext(ctx).
		Body(cond).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	if ccErr := resp.CCError(); ccErr != nil {
		return nil, ccErr
	}

	return resp, nil
}

// DeletePropertyGroup TODO
func (t *object) DeletePropertyGroup(ctx context.Context, groupID string, h http.Header) (resp *metadata.Response,
	err error) {
	resp = new(metadata.Response)
	subPath := "/delete/objectattgroup/%s"

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, groupID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// UpdatePropertyGroupObjectAtt TODO
func (t *object) UpdatePropertyGroupObjectAtt(ctx context.Context, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/update/objectattgroupproperty"

	err = t.client.Put().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// DeletePropertyGroupObjectAtt TODO
func (t *object) DeletePropertyGroupObjectAtt(ctx context.Context, ownerID string, objID string, propertyID string,
	groupID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/objectatt/group/owner/%s/object/%s/propertyids/%s/groupids/%s"

	err = t.client.Delete().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID, objID, propertyID, groupID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// SelectPropertyGroupByObjectID TODO
func (t *object) SelectPropertyGroupByObjectID(ctx context.Context, objID string, h http.Header,
	data map[string]interface{}) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/find/objectattgroup/object/%s"

	err = t.client.Post().
		WithContext(ctx).
		Body(data).
		SubResourcef(subPath, objID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}
