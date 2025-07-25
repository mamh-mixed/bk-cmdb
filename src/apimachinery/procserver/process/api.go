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

package process

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/rest"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// ProcessClientInterface TODO
type ProcessClientInterface interface {
	CreateProcessInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response,
		err error)
	DeleteProcessInstance(ctx context.Context, h http.Header,
		data *metadata.DeleteProcessInstanceInServiceInstanceInput) error
	SearchProcessInstance(ctx context.Context, h http.Header, data *metadata.ListProcessInstancesOption) (
		[]metadata.ProcessInstance, error)
	UpdateProcessInstance(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response,
		err error)

	CreateProcessTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response,
		err error)
	DeleteProcessTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response,
		err error)
	SearchProcessTemplate(ctx context.Context, h http.Header, i *metadata.ListProcessTemplateWithServiceTemplateInput) (
		*metadata.MultipleProcessTemplate, errors.CCErrorCoder)
	UpdateProcessTemplate(ctx context.Context, h http.Header, data map[string]interface{}) (resp *metadata.Response,
		err error)
	ListProcessRelatedInfo(ctx context.Context, h http.Header, bizID int64,
		data metadata.ListProcessRelatedInfoOption) (resp *metadata.ListProcessRelatedInfoResponse, err error)
	ListProcessInstancesNameIDsInModule(ctx context.Context, h http.Header,
		data map[string]interface{}) (resp *metadata.Response, err error)
	ListProcessInstancesDetailsByIDs(ctx context.Context, h http.Header,
		data map[string]interface{}) (resp *metadata.Response, err error)
	ListProcessInstancesDetails(ctx context.Context, h http.Header, bizID int64,
		data metadata.ListProcessInstancesDetailsOption) (resp *metadata.MapArrayResponse, err error)
	UpdateProcessInstancesByIDs(ctx context.Context, h http.Header,
		data map[string]interface{}) (resp *metadata.Response, err error)
}

// NewProcessClientInterface TODO
func NewProcessClientInterface(client rest.ClientInterface) ProcessClientInterface {
	return &process{client: client}
}

type process struct {
	client rest.ClientInterface
}
