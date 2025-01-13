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

package adminserver

import (
	"context"
	"net/http"

	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// ClearDatabase TODO
func (a *adminServer) ClearDatabase(ctx context.Context, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/clear"

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// Set TODO
func (a *adminServer) Set(ctx context.Context, ownerID string, h http.Header) (resp *metadata.Response, err error) {
	resp = new(metadata.Response)
	subPath := "/migrate/system/hostcrossbiz/%s"

	err = a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath, ownerID).
		WithHeaders(h).
		Do().
		Into(resp)
	return
}

// Migrate TODO
func (a *adminServer) Migrate(ctx context.Context, h http.Header) error {
	resp := new(metadata.Response)
	subPath := "/migrate/database"

	err := a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return err
	}

	return nil
}

// RunSyncDBIndex TODO
func (a *adminServer) RunSyncDBIndex(ctx context.Context, h http.Header) error {
	resp := new(metadata.Response)
	subPath := "/migrate/sync/db/index"

	err := a.client.Post().
		WithContext(ctx).
		Body(nil).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(resp)

	if err != nil {
		return errors.CCHttpError
	}

	if err = resp.CCError(); err != nil {
		return err
	}

	return nil
}
