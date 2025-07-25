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

// Package user TODO
package user

import (
	"configcenter/src/apimachinery/apiserver"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/web_server/app/options"

	"github.com/gin-gonic/gin"
)

// LoginInterface TODO
// User 登录系统抽象出来接口
type LoginInterface interface {
	// LoginUser 判断用户是否登录
	LoginUser(c *gin.Context) (isLogin bool)
	// GetLoginUrl 获取登录系统的URL
	GetLoginUrl(c *gin.Context) string
	// GetUserList 获取不同登录方式下对应的用户列表
	GetUserList(c *gin.Context) ([]*metadata.LoginSystemUserInfo, *errors.RawErrorInfo)
}

// NewUser return user instance by type
func NewUser(config options.Config, engine *backbone.Engine, cacheCli redis.Client,
	apiCli apiserver.ApiServerClientInterface) LoginInterface {

	return &publicUser{config, engine, cacheCli, apiCli}
}
