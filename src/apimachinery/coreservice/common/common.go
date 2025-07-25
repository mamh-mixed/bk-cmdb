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

// Package common TODO
package common

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// CommonInterface TODO
type CommonInterface interface {
	GetDistinctField(ctx context.Context, h http.Header, option *metadata.DistinctFieldOption) ([]interface{},
		errors.CCErrorCoder)
	GetDistinctCount(ctx context.Context, h http.Header, option *metadata.DistinctFieldOption) (int64,
		errors.CCErrorCoder)
}

// NewCommonInterfaceClient TODO
func NewCommonInterfaceClient(client rest.ClientInterface) CommonInterface {
	return &common{client: client}
}

type common struct {
	client rest.ClientInterface
}
