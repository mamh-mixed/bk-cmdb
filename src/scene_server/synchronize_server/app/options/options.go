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

// Package options TODO
package options

import (
	"github.com/spf13/pflag"

	"configcenter/src/common/auth"
	"configcenter/src/common/core/cc/config"
)

// ServerOption define option of server in flags
type ServerOption struct {
	ServConf *config.CCAPIConfig
}

// NewServerOption create a ServerOption object
func NewServerOption() *ServerOption {
	s := ServerOption{
		ServConf: config.NewCCAPIConfig(),
	}

	return &s
}

// AddFlags add flags
func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {
	s.ServConf.AddFlags(fs, "127.0.0.1:60006")
	fs.Var(auth.EnableAuthFlag, "enable-auth", "The auth center enable status, true for enabled, false for disabled")
}

// Config config file set
type Config struct {
	Names           []string
	exceptionDir    string
	ConifgItemArray []*ConfigItem
	Trigger         TriggerTime
}

const (
	// TriggerTimeTypeTiming TODO
	TriggerTimeTypeTiming = "timing"
	// TriggerTimeTypeInterval TODO
	TriggerTimeTypeInterval = "interval"
)

// TriggerTime  define synchronize task trigger style and role
type TriggerTime struct {
	// timing, interval , default value timing
	TriggerType string
	Role        string
}

// IsTiming judge is timing trigger
func (t TriggerTime) IsTiming() bool {
	if t.TriggerType != TriggerTimeTypeInterval {
		return true
	}
	return false
}

// ConfigItem config item info
type ConfigItem struct {
	Name string
	// White list, default false
	WhiteList bool
	// White list  = true, need synchronize app list
	// White list  = true,  out of synchronize business name,
	AppNames []string

	// ObjectID array source
	ObjectIDArr []string
	// resource pool sync config
	SyncResource bool

	// TargetHost target data logics
	TargetHost string

	// FieldSign source data fields
	FieldSign string

	// SynchronizeFlag current server data synchronize flag
	SynchronizeFlag string

	// SupplerAccount string
	SupplerAccount []string

	exceptionDirectory string

	// Retry error max retry count
	ExceptionFileCount int

	// Unsynchronized model related properties
	// 使用忽略模型属性变的模式。 用户在目标中新加对应的模型，模型的属性。
	// 满足同步的实例将会同步到目的cmdb。 在目标系统中新建相同的唯一标识模型或者模型的字段。内容会自动展示出来
	IgnoreModelAttr bool

	// EnableInstFilter  是否开启实例数据根据同步身份过滤
	EnableInstFilter bool
}
