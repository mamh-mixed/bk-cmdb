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

package metadata

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIssue1720(t *testing.T) {
	testData := ReadModelWithAttributeResult{
		Data: QueryModelWithAttributeDataResult{
			Info: []SearchModelInfo{
				{
					Spec: Object{
						CreateTime: &Time{Time: time.Now()},
					},
					Attributes: []Attribute{
						{
							CreateTime: &Time{Time: time.Now()},
						},
					},
				},
			},
		},
	}

	inputData, err := json.Marshal(testData)

	require.NoError(t, err)
	t.Logf("input data:%s", inputData)

	out := &ReadModelWithAttributeResult{}
	err = json.Unmarshal(inputData, out)
	require.NoError(t, err)
	t.Logf("out:%s %s", out.Data.Info[0].Spec.CreateTime.String(), out.Data.Info[0].Attributes[0].CreateTime.String())

}
