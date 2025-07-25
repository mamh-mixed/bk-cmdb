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

// Package instance TODO
package instance

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// InstanceClientInterface TODO
type InstanceClientInterface interface {
	CreateInstance(ctx context.Context, h http.Header, objID string, input *metadata.CreateModelInstance) (
		*metadata.CreateOneDataResult, error)
	CreateManyInstance(ctx context.Context, h http.Header, objID string, input *metadata.CreateManyModelInstance) (
		*metadata.CreateManyDataResult, errors.CCErrorCoder)
	// BatchCreateInstance batch create instance, if one of instances fails to create, an error is returned.
	BatchCreateInstance(ctx context.Context, h http.Header, objID string, input *metadata.BatchCreateModelInstOption) (
		*metadata.BatchCreateInstRespData, errors.CCErrorCoder)
	SetManyInstance(ctx context.Context, h http.Header, objID string, input *metadata.SetManyModelInstance) (
		resp *metadata.SetOptionResult, err error)
	UpdateInstance(ctx context.Context, h http.Header, objID string, input *metadata.UpdateOption) (
		*metadata.UpdatedCount, errors.CCErrorCoder)
	ReadInstance(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition) (
		*metadata.InstDataInfo, error)
	DeleteInstance(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (
		*metadata.DeletedCount, error)
	DeleteInstanceCascade(ctx context.Context, h http.Header, objID string, input *metadata.DeleteOption) (
		resp *metadata.DeletedOptionResult, err error)
	// ReadInstanceStruct 按照结构体返回实例数据
	ReadInstanceStruct(ctx context.Context, h http.Header, objID string, input *metadata.QueryCondition,
		result interface{}) (err errors.CCErrorCoder)

	// CountInstances counts model instances num.
	CountInstances(ctx context.Context, header http.Header, objID string, input *metadata.Condition) (
		*metadata.CountResponseContent, error)
	GetInstanceObjectMapping(ctx context.Context, h http.Header, ids []int64) ([]metadata.ObjectMapping,
		errors.CCErrorCoder)
}

// NewInstanceClientInterface TODO
func NewInstanceClientInterface(client rest.ClientInterface) InstanceClientInterface {
	return &instance{client: client}
}

type instance struct {
	client rest.ClientInterface
}
