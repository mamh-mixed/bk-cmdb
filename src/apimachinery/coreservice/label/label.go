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

// Package label TODO
package label

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/selector"
)

// LabelInterface TODO
type LabelInterface interface {
	AddLabel(ctx context.Context, h http.Header, tableName string, option selector.LabelAddOption) errors.CCErrorCoder
	RemoveLabel(ctx context.Context, h http.Header, tableName string,
		option selector.LabelRemoveOption) errors.CCErrorCoder
	UpdateLabel(ctx context.Context, h http.Header, tableName string,
		option *selector.LabelUpdateOption) errors.CCErrorCoder
}

// NewLabelInterfaceClient TODO
func NewLabelInterfaceClient(client rest.ClientInterface) LabelInterface {
	return &label{client: client}
}

type label struct {
	client rest.ClientInterface
}
