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

// Package ac auth control
package ac

import (
	"context"
	"errors"
	"net/http"

	"configcenter/src/ac/meta"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/storage/dal/redis"
)

// NoAuthorizeError TODO
var NoAuthorizeError = errors.New("no authorize")

// Viewer is a interface to operate iam view
type Viewer interface {
	// CreateView create iam view for objects
	CreateView(ctx context.Context, h http.Header, objects []metadata.Object, redisCli redis.Client, rid string) error
	// DeleteView delete iam view for objects
	DeleteView(ctx context.Context, h http.Header, objects []metadata.Object, redisCli redis.Client, rid string) error
	// UpdateView update iam view for objects
	UpdateView(ctx context.Context, h http.Header, objects []metadata.Object, redisCli redis.Client, rid string) error
}

// AuthorizeInterface TODO
type AuthorizeInterface interface {
	AuthorizeBatch(ctx context.Context, h http.Header, user meta.UserInfo, resources ...meta.ResourceAttribute) (
		[]types.Decision, error)

	AuthorizeAnyBatch(ctx context.Context, h http.Header, user meta.UserInfo, resources ...meta.ResourceAttribute) (
		[]types.Decision, error)

	ListAuthorizedResources(ctx context.Context, h http.Header, input meta.ListAuthorizedResourcesParam) (
		*types.AuthorizeList, error)
	GetNoAuthSkipUrl(ctx context.Context, h http.Header, input *metadata.IamPermission) (string, error)
	GetPermissionToApply(ctx context.Context, h http.Header, input []meta.ResourceAttribute) (*metadata.IamPermission,
		error)
	RegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstanceWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)
	BatchRegisterResourceCreatorAction(ctx context.Context, h http.Header, input metadata.IamInstancesWithCreator) (
		[]metadata.IamCreatorActionPolicy, error)
}
