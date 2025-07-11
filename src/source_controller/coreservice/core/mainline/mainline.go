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

// Package mainline TODO
package mainline

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common/language"
	"configcenter/src/source_controller/coreservice/core"
)

var _ core.TopoOperation = (*topoManager)(nil)

type topoManager struct {
	lang      language.CCLanguageIf
	clientSet apimachinery.ClientSetInterface
}

// New create a new model manager instance
func New(lang language.CCLanguageIf, clientSet apimachinery.ClientSetInterface) core.TopoOperation {

	return &topoManager{
		lang:      lang,
		clientSet: clientSet,
	}
}
