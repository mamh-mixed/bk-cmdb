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

package cloudsync

import (
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

const (
	// SyncPeriodMinutesMin 同步周期最小值
	SyncPeriodMinutesMin = 5
)

// SyncPeriodMinutes 同步周期，单位为分钟
var SyncPeriodMinutes int

// 任务处理器
type taskProcessor struct {
	scheduler *taskScheduler
}

// NewTaskProcessor TODO
func NewTaskProcessor(scheduler *taskScheduler) *taskProcessor {
	return &taskProcessor{scheduler}
}

// TaskChanLoop 不断给taskChan提供任务数据
func (t *taskProcessor) TaskChanLoop(taskChan chan *metadata.CloudSyncTask) {
	go func() {
		for {
			tasks, err := t.scheduler.GetTaskList()
			if err != nil {
				blog.Errorf("scheduler GetTaskList err:%s", err.Error())
				time.Sleep(time.Duration(SyncPeriodMinutes) * time.Minute)
				continue
			}
			for i := range tasks {
				taskChan <- tasks[i]
			}
			time.Sleep(time.Duration(SyncPeriodMinutes) * time.Minute)
		}
	}()
}
