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

package common

import (
	"context"
	"net/http"

	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
)

// GetDistinctField TODO
func (p *common) GetDistinctField(ctx context.Context, h http.Header,
	option *metadata.DistinctFieldOption) ([]interface{}, errors.CCErrorCoder) {
	ret := new(metadata.ArrayResponse)
	subPath := "/findmany/common/distinct_field"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateServiceCategory failed, http request failed, err: %+v", err)
		return nil, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return nil, ret.CCError()
	}

	return ret.Data, nil
}

// GetDistinctCount 根据条件获取指定表中满足条件数据的数量
func (p *common) GetDistinctCount(ctx context.Context, h http.Header, option *metadata.DistinctFieldOption) (int64,
	errors.CCErrorCoder) {
	ret := new(metadata.HostCountResponse)
	subPath := "/findmany/common/distinct_count"

	err := p.client.Post().
		WithContext(ctx).
		Body(option).
		SubResourcef(subPath).
		WithHeaders(h).
		Do().
		Into(ret)

	if err != nil {
		blog.Errorf("CreateServiceCategory failed, http request failed, err: %+v", err)
		return ret.Data, errors.CCHttpError
	}
	if ret.CCError() != nil {
		return ret.Data, ret.CCError()
	}

	return ret.Data, nil
}
