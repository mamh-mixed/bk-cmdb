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

package plugin

import (
	"sync/atomic"

	"configcenter/src/common/metric"
)

// NewIntegerCounter TODO
// integer
func NewIntegerCounter(name, help string) *IntegerCounter {
	return &IntegerCounter{
		name:  name,
		help:  help,
		value: 0,
	}
}

// IntegerCounter counter for integer
type IntegerCounter struct {
	name  string
	help  string
	value int64
}

// GetMeta TODO
func (c *IntegerCounter) GetMeta() metric.MetricMeta {
	return metric.MetricMeta{
		Name: c.name,
		Help: c.help,
	}
}

// GetValue TODO
func (c *IntegerCounter) GetValue() (*metric.FloatOrString, error) {
	return metric.FormFloatOrString(c.value)
}

// GetExtension TODO
func (c *IntegerCounter) GetExtension() (*metric.MetricExtension, error) {
	return nil, nil
}

// Reset TODO
func (c *IntegerCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

// Inc TODO
func (c *IntegerCounter) Inc(i int64) {
	atomic.AddInt64(&c.value, i)
}

// Dec TODO
func (c *IntegerCounter) Dec(i int64) {
	atomic.AddInt64(&c.value, -i)
}
