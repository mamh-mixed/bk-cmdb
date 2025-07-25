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

package service

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/zkclient"

	"github.com/emicklei/go-restful/v3"
)

// Limiter TODO
type Limiter struct {
	zkCli        *zkclient.ZkClient
	rules        map[string]*metadata.LimiterRule
	lock         sync.RWMutex
	syncDuration time.Duration
}

// NewLimiter TODO
func NewLimiter(zkCli *zkclient.ZkClient) *Limiter {
	return &Limiter{
		zkCli:        zkCli,
		syncDuration: 5 * time.Second,
	}
}

// SyncLimiterRules sync the api limiter rules from zk
func (l *Limiter) SyncLimiterRules() error {
	blog.Info("begin SyncLimiterRules")
	path := types.CC_SERVLIMITER_BASEPATH
	go func() {
		for {
			err := l.syncLimiterRules(path)
			if err != nil {
				blog.Errorf("fail to syncLimiterRules for path:%s, err:%s", path, err.Error())
			}
			time.Sleep(l.syncDuration)
		}
	}()
	return nil
}

func (l *Limiter) syncLimiterRules(path string) error {
	blog.V(5).Infof("syncing limiter rules for path:%s", path)
	children, err := l.zkCli.GetChildren(path)
	if err != nil {
		if err == zkclient.ErrNoNode {
			// if no rules, set rules to be empty
			l.setRules(make(map[string]*metadata.LimiterRule))
			return nil
		}
		blog.Errorf("fail to GetChildren for path:%s, err:%s", path, err.Error())
		return err
	}

	rules := make(map[string]*metadata.LimiterRule)
	for _, child := range children {
		data, err := l.zkCli.Get(path + "/" + child)
		if err != nil {
			blog.Errorf("fail to Get for path:%s, err:%s", path, err.Error())
			continue
		}

		rule := new(metadata.LimiterRule)
		err = json.Unmarshal([]byte(data), rule)
		if err != nil {
			blog.Errorf("fail to Unmarshal for child:%s, data:%s, err:%s", child, data, err.Error())
			continue
		}

		err = rule.Verify()
		if err != nil {
			blog.Errorf("fail to Verify for child:%s, rule:%v, err:%s", child, rule, err.Error())
			continue
		}

		rules[rule.RuleName] = rule
	}

	l.setRules(rules)
	return nil
}

func (l *Limiter) setRules(rules map[string]*metadata.LimiterRule) {
	l.lock.Lock()
	if reflect.DeepEqual(rules, l.rules) {
		blog.V(5).Info("setRules skip, nothing is changed")
		l.lock.Unlock()
		return
	}
	l.rules = rules
	l.lock.Unlock()
	blog.InfoJSON("setRules success, current rules is %s", rules)
}

// GetRules get all rules of limiter
func (l *Limiter) GetRules() map[string]*metadata.LimiterRule {
	l.lock.RLock()
	defer l.lock.RUnlock()
	rules := make(map[string]*metadata.LimiterRule)
	for k, v := range l.rules {
		rule := *v
		rules[k] = &rule
	}
	return rules
}

// LenOfRules get the count of limiter's rules
func (l *Limiter) LenOfRules() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.rules)
}

// GetMatchedRule get the matched limiter rule according request
func (l *Limiter) GetMatchedRule(req *restful.Request) *metadata.LimiterRule {
	header := req.Request.Header
	var matchedRule *metadata.LimiterRule
	var min int64 = 999999
	rules := l.GetRules()
	for _, r := range rules {
		if r.AppCode == "" && r.User == "" && r.IP == "" && r.Url == "" && r.Method == "" {
			blog.Errorf("wrong rule format, one of appcode, user, ip, url, method must be set, r:%#v", *r)
			continue
		}
		if r.AppCode != "" {
			if r.AppCode != httpheader.GetAppCode(header) {
				continue
			}
		}
		if r.User != "" {
			if r.User != httpheader.GetUser(header) {
				continue
			}
		}
		if r.IP != "" {
			hit := false
			ips := strings.Split(r.IP, ",")
			for _, ip := range ips {
				if strings.TrimSpace(ip) == strings.TrimSpace(httpheader.GetReqRealIP(header)) {
					hit = true
					break
				}
			}
			if hit == false {
				continue
			}
		}
		if r.Method != "" {
			if util.Normalize(r.Method) != util.Normalize(req.Request.Method) {
				continue
			}
		}
		if r.Url != "" {
			match, err := regexp.MatchString(r.Url, req.Request.RequestURI)
			if err != nil {
				blog.Errorf("MatchString failed, r.Url:%s, reqURI:%s, err:%s", r.Url, req.Request.RequestURI,
					err.Error())
				continue
			}
			if !match {
				continue
			}
		}

		if r.DenyAll == true {
			matchedRule = r
			break
		}
		if r.Limit < min {
			min = r.Limit
			matchedRule = r
		}
	}
	return matchedRule
}
