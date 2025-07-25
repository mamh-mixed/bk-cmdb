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

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"configcenter/src/ac/extensions"
	bkc "configcenter/src/common"
	"configcenter/src/common/backbone"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/storage/dal/redis"
)

// Discover TODO
type Discover struct {
	ctx        context.Context
	httpHeader http.Header

	redisCli redis.Client
	*backbone.Engine
	authManager *extensions.AuthManager
}

var msgHandlerCnt = int64(0)

// NewDiscover new discover
func NewDiscover(ctx context.Context, redisCli redis.Client, backbone *backbone.Engine,
	authManager *extensions.AuthManager) *Discover {
	header := headerutil.BuildHeader(bkc.CCSystemCollectorUserName, bkc.BKDefaultOwnerID)

	discover := &Discover{
		redisCli:    redisCli,
		ctx:         ctx,
		httpHeader:  header,
		authManager: authManager,
	}
	discover.Engine = backbone
	return discover
}

// Hash returns hash value base on message.
func (d *Discover) Hash(cloudid, ip string) (string, error) {
	if len(cloudid) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, cloudid empty")
	}
	if len(ip) == 0 {
		return "", fmt.Errorf("can't make hash from invalid message format, ip empty")
	}

	hash := fmt.Sprintf("%s:%s", cloudid, ip)

	return hash, nil
}

// Mock returns local mock message for testing.
func (d *Discover) Mock() string {
	return MockMessage
}

// Analyze analyze discover data
func (d *Discover) Analyze(msg *string, sourceType string) (bool, error) {
	err := d.UpdateOrCreateInst(msg)
	if err != nil {
		return false, fmt.Errorf("create inst err: %v, raw: %s", err, msg)
	}
	return false, nil
}

// MockMessage TODO
var MockMessage = `{
    "meta": {
        "model": {
            "bk_obj_id": "bk_apache",
            "bk_supplier_account": "0"
        }
    },
    "data": {
        "bk_inst_name": "apache",
        "bk_ip": "127.0.0.1"
    }
}`
