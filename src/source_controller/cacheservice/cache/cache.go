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

// Package cache TODO
package cache

import (
	"fmt"

	"configcenter/src/apimachinery/discovery"
	biztopo "configcenter/src/source_controller/cacheservice/cache/biz-topo"
	"configcenter/src/source_controller/cacheservice/cache/custom"
	"configcenter/src/source_controller/cacheservice/cache/general"
	"configcenter/src/source_controller/cacheservice/cache/mainline"
	"configcenter/src/source_controller/cacheservice/cache/topotree"
	"configcenter/src/source_controller/cacheservice/event/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/task"
)

// NewCache new cache service
func NewCache(isMaster discovery.ServiceManageInterface) (*ClientSet, error) {
	if err := mainline.NewMainlineCache(isMaster); err != nil {
		return nil, fmt.Errorf("new mainline cache failed, err: %v", err)
	}

	customCache, err := custom.New(isMaster)
	if err != nil {
		return nil, fmt.Errorf("new custom resource cache failed, err: %v", err)
	}

	watchCli := watch.NewClient(mongodb.Dal("watch"), mongodb.Dal(), redis.Client())

	generalCache, err := general.New(isMaster, watchCli)
	if err != nil {
		return nil, fmt.Errorf("new general resource cache failed, err: %v", err)
	}

	topoTreeClient, err := biztopo.New(isMaster, customCache.CacheSet(), watchCli)
	if err != nil {
		return nil, fmt.Errorf("new common topo cache failed, err: %v", err)
	}

	mainlineClient := mainline.NewMainlineClient(generalCache)

	cache := &ClientSet{
		Tree:     topotree.NewTopologyTree(mainlineClient),
		Business: mainlineClient,
		Topo:     topoTreeClient,
		Event:    watchCli,
		Custom:   customCache,
		General:  generalCache,
	}
	return cache, nil
}

// ClientSet is the cache client set
type ClientSet struct {
	Tree     *topotree.TopologyTree
	Topo     *biztopo.Topo
	Business *mainline.Client
	Event    *watch.Client
	Custom   *custom.Cache
	General  *general.Cache
}

// GetWatchTasks returns the event watch tasks
func (c *ClientSet) GetWatchTasks() []*task.Task {
	tasks := c.Topo.Watcher().GetWatchTasks()
	tasks = append(tasks, c.Custom.Watcher().GetWatchTasks()...)
	tasks = append(tasks, c.General.FullSyncCond().GetWatchTasks()...)
	return tasks
}
