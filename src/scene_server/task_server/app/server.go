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

package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/task_server/app/options"
	"configcenter/src/scene_server/task_server/logics"
	tasksvc "configcenter/src/scene_server/task_server/service"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful/v3"
)

// Run TODO
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		blog.Errorf("wrap server info failed, err: %v", err)
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(tasksvc.Service)
	taskSrv := new(TaskServer)

	input := &backbone.BackboneParameter{
		SrvRegdiscv:  backbone.SrvRegdiscv{Regdiscv: op.ServConf.RegDiscover, TLSConfig: op.ServConf.GetTLSClientConf()},
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: taskSrv.onHostConfigUpdate,
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		blog.Errorf("new backbone failed, err: %v", err)
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil != taskSrv.Config {
			configReady = true
			break
		}
		blog.Infof("waiting for config ready ...")
		time.Sleep(time.Second)
	}
	if false == configReady {
		blog.Infof("waiting config timeout.")
		return errors.New("configuration item not found")
	}
	taskSrv.Config.Mongo, err = engine.WithMongo()
	if err != nil {
		return err
	}
	taskSrv.Config.Redis, err = engine.WithRedis()
	if err != nil {
		return err
	}
	cacheDB, err := redis.NewFromConfig(taskSrv.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
		return fmt.Errorf("new redis client failed, err: %s", err.Error())
	}
	db, err := taskSrv.Config.Mongo.GetMongoClient()
	if err != nil {
		blog.Errorf("new mongo client failed, err: %s", err.Error())
		return fmt.Errorf("new mongo client failed, err: %s", err.Error())
	}

	initErr := db.InitTxnManager(cacheDB)
	if initErr != nil {
		blog.Errorf("init txn manager failed, err: %v", initErr)
		return initErr
	}

	service.Engine = engine
	service.Config = taskSrv.Config
	service.CacheDB = cacheDB
	service.DB = db
	taskSrv.Core = engine
	service.Logics = logics.NewLogics(engine.CoreAPI, db)
	taskSrv.Service = service

	// cron job delete history task
	go taskSrv.Service.TimerDeleteHistoryTask(ctx)

	if err := backbone.StartServer(ctx, cancel, engine, service.WebService(), true); err != nil {
		blog.Errorf("start backbone failed, err: %+v", err)
		return err
	}

	queue := service.NewQueue(taskSrv.taskQueue)
	queue.Start()
	select {
	case <-ctx.Done():
	}
	return nil
}

// TaskServer TODO
type TaskServer struct {
	Core      *backbone.Engine
	Config    *options.Config
	Service   *tasksvc.Service
	taskQueue map[string]tasksvc.TaskInfo
}

// WebService TODO
func (h *TaskServer) WebService() *restful.Container {
	return h.Service.WebService()
}

func (h *TaskServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	if h.Config == nil {
		h.Config = new(options.Config)
	}
}
