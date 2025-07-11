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

package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/json"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/tools/cmdb_ctl/app/config"
)

// WithRedColor TODO
func WithRedColor(str string) string {
	return fmt.Sprintf("%c[1;40;31m>> %s %c[0m\n", 0x1B, str, 0x1B)
}

// WithGreenColor TODO
func WithGreenColor(str string) string {
	return fmt.Sprintf("%c[1;40;32m>> %s %c[0m\n", 0x1B, str, 0x1B)
}

// WithBlueColor TODO
func WithBlueColor(str string) string {
	return fmt.Sprintf("%c[1;40;34m>> %s %c[0m\n", 0x1B, str, 0x1B)
}

func doCmdbHttpRequest(ccModule, path, tenantID string, body interface{}) (*http.Response, error) {
	// get server address from zk
	zk, err := config.NewZkService(config.Conf.ZkAddr, &config.Conf.ZkTLS)
	if err != nil {
		fmt.Printf("new zk client failed, err: %v\n", err)
		return nil, err
	}

	zkPath := types.CC_SERV_BASEPATH + "/" + ccModule
	children, err := zk.ZkCli.GetChildren(zkPath)
	if err != nil {
		fmt.Printf("get %s server failed, err: %v\n", ccModule, err)
		return nil, err
	}

	server := ""
	for _, child := range children {
		node, err := zk.ZkCli.Get(zkPath + "/" + child)
		if err != nil {
			return nil, err
		}
		svr := new(types.EventServInfo)
		if err := json.Unmarshal([]byte(node), svr); err != nil {
			return nil, err
		}
		server = fmt.Sprintf("%s:%d", svr.RegisterIP, svr.Port)
		break
	}

	if server == "" {
		return nil, fmt.Errorf("%s server not found", ccModule)
	}

	// do http request
	url := fmt.Sprintf("http://%s/%s", server, strings.TrimPrefix(path, "/"))

	data, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("marshal request body %+v failed, err: %v\n", body, err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		fmt.Printf("new http request failed, err: %v, url %s, body: %s\n", err, url, string(data))
		return nil, err
	}
	req.Header = headerutil.GenCommonHeader("cmdb_tool", tenantID, util.GenerateRID())

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("do request failed, err: %v, url: %s, body: %s\n", err, url, string(data))
		return nil, err
	}

	return resp, nil
}
