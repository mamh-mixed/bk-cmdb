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

// Package bsrelation TODO
package bsrelation

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/source_controller/cacheservice/event"
	mixevent "configcenter/src/source_controller/cacheservice/event/mix-event"
	"configcenter/src/storage/stream/task"
)

const (
	bizSetRelationLockKey = common.BKCacheKeyV3Prefix + "biz_set_relation:event_lock"
	bizSetRelationLockTTL = 1 * time.Minute
)

// BizSetRelation is the biz set relation event flow struct
type BizSetRelation struct {
	tasks []*task.Task
}

// GetWatchTasks returns the event watch tasks
func (b *BizSetRelation) GetWatchTasks() []*task.Task {
	return b.tasks
}

// NewBizSetRelation init and run biz set relation event watch
func NewBizSetRelation() (*BizSetRelation, error) {
	bizSetRel := &BizSetRelation{tasks: make([]*task.Task, 0)}

	base := mixevent.MixEventFlowOptions{
		MixKey:       event.BizSetRelationKey,
		EventLockKey: bizSetRelationLockKey,
		EventLockTTL: bizSetRelationLockTTL,
	}

	// watch biz set event
	bizSet := base
	bizSet.Key = event.BizSetKey
	bizSet.WatchFields = []string{common.BKBizSetIDField, common.BKBizSetScopeField}
	if err := bizSetRel.addWatchTask(bizSet); err != nil {
		blog.Errorf("watch biz set event for biz set relation failed, err: %v", err)
		return nil, err
	}
	blog.Info("watch biz set relation events, watch biz set success")

	// watch biz event
	biz := base
	biz.Key = event.BizKey
	if err := bizSetRel.addWatchTask(biz); err != nil {
		blog.Errorf("watch biz event for biz set relation failed, err: %v", err)
		return nil, err
	}
	blog.Info("watch biz set relation events, watch biz success")

	return bizSetRel, nil
}
