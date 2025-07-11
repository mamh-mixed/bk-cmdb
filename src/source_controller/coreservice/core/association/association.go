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

// Package association TODO
package association

import (
	"configcenter/src/apimachinery"
	"configcenter/src/source_controller/coreservice/core"
)

var _ core.AssociationOperation = (*associationManager)(nil)

type associationManager struct {
	*associationKind
	*associationInstance
	*associationModel
	clientSet apimachinery.ClientSetInterface
}

// New create a new association manager instance
func New(dependent OperationDependencies, clientSet apimachinery.ClientSetInterface) core.AssociationOperation {
	asstModel := &associationModel{}
	asstKind := &associationKind{
		associationModel: asstModel,
	}
	return &associationManager{
		associationKind: asstKind,
		associationInstance: &associationInstance{
			associationKind:  asstKind,
			associationModel: asstModel,
			dependent:        dependent,
			clientSet:        clientSet,
		},
		associationModel: &associationModel{},
		clientSet:        clientSet,
	}
}
