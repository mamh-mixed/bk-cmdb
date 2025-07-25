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
	"fmt"
	"time"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/operation_server/app/options"
	"configcenter/src/scene_server/operation_server/service"
)

// Run TODO
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		blog.Errorf("fail to new server information. err: %s", err.Error())
		return fmt.Errorf("make server information failed, err:%v", err)
	}

	operationSvr := new(service.OperationServer)

	input := &backbone.BackboneParameter{
		ConfigUpdate: operationSvr.OnOperationConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		SrvRegdiscv:  backbone.SrvRegdiscv{Regdiscv: op.ServConf.RegDiscover, TLSConfig: op.ServConf.GetTLSClientConf()},
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if operationSvr.Config.Ready() {
			configReady = true
			break
		}
		time.Sleep(time.Second)
	}
	if false == configReady {
		return fmt.Errorf("waiting for configuration timeout, maybe parse configuration failed")
	}
	operationSvr.Config.Mongo, err = engine.WithMongo()
	if err != nil {
		return err
	}

	operationSvr.Config.Auth, err = iam.ParseConfigFromKV("authServer", nil)
	if err != nil {
		blog.Warnf("parse auth center config failed: %v", err)
	}

	iamCli := new(iam.IAM)
	if auth.EnableAuthorize() {
		blog.Info("enable auth center access")
		iamCli, err = iam.NewIAM(operationSvr.Config.Auth, engine.Metric().Registry())
		if err != nil {
			return fmt.Errorf("new iam client failed: %v", err)
		}
	} else {
		blog.Infof("disable auth center access")
	}
	operationSvr.AuthManager = extensions.NewAuthManager(engine.CoreAPI, iamCli)

	operationSvr.Engine = engine

	go operationSvr.InitFunc()
	if err := backbone.StartServer(ctx, cancel, engine, operationSvr.WebService(), true); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("operation will exit!")
	}

	return nil
}
