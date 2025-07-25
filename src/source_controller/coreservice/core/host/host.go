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

// Package host TODO
package host

import (
	"configcenter/src/apimachinery/cacheservice/cache/host"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/host/searcher"
	"configcenter/src/source_controller/coreservice/core/host/transfer"
)

var _ core.HostOperation = (*hostManager)(nil)

type hostManager struct {
	hostTransfer *transfer.TransferManager
	dependent    transfer.OperationDependence
	hostSearcher searcher.Searcher
}

// New create a new model manager instance
func New(dependent transfer.OperationDependence, hostApplyDependence transfer.HostApplyRuleDependence,
	cacheHost host.Interface) core.HostOperation {

	coreMgr := &hostManager{
		dependent: dependent,
	}
	coreMgr.hostTransfer = transfer.New(dependent, hostApplyDependence)
	coreMgr.hostSearcher = searcher.New(cacheHost)
	return coreMgr
}
