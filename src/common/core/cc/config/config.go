/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package config TODO
package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

// CCAPIConfig define configuration of ccapi server
type CCAPIConfig struct {
	AddrPort    string
	RegDiscover string
	RegisterIP  string
	ExConfig    string
	Environment string
	Qps         int64
	Burst       int64
}

// NewCCAPIConfig create ccapi config object
func NewCCAPIConfig() *CCAPIConfig {
	return &CCAPIConfig{
		AddrPort:    "127.0.0.1:8081",
		RegDiscover: "",
		RegisterIP:  "",
		Qps:         1000,
		Burst:       2000,
	}
}

// GetAddress TODO
// IPV6 addr port, like ::1:80
// IPV4 addr port, like 127.0.0.1:80
// GetAddress get the address
func (conf *CCAPIConfig) GetAddress() (string, error) {
	addrPort := strings.TrimSpace(conf.AddrPort)
	if err := checkAddrPort(addrPort); err != nil {
		return "", err
	}
	if isIPV6(addrPort) {
		return getIPV6Adrr(addrPort)
	}
	return getIPV4Adrr(addrPort)
}

// GetPort get the port
func (conf *CCAPIConfig) GetPort() (uint, error) {
	addrPort := strings.TrimSpace(conf.AddrPort)
	if err := checkAddrPort(addrPort); err != nil {
		return uint(0), err
	}
	if isIPV6(addrPort) {
		return getIPV6Port(addrPort)
	}
	return getIPV4Port(addrPort)
}

func checkAddrPort(addrPort string) error {
	if strings.Count(addrPort, ":") == 0 {
		return fmt.Errorf("the value of flag[AddrPort: %s] is wrong", addrPort)
	}
	return nil
}

func isIPV6(addrPort string) bool {
	return strings.Count(addrPort, ":") > 1
}

func getIPV6Adrr(addrPort string) (string, error) {
	idx := strings.LastIndex(addrPort, ":")
	return addrPort[:idx], nil
}

func getIPV4Adrr(addrPort string) (string, error) {
	idx := strings.LastIndex(addrPort, ":")
	return addrPort[:idx], nil
}

func getIPV6Port(addrPort string) (uint, error) {
	return getPortFunc(addrPort)
}

func getIPV4Port(addrPort string) (uint, error) {
	return getPortFunc(addrPort)
}

func getPortFunc(addrPort string) (uint, error) {
	idx := strings.LastIndex(addrPort, ":")
	// the port can't be empty, len(":port") can't less than 2
	if len(addrPort[idx:]) < 2 {
		return 0, fmt.Errorf("the value of flag[AddrPort: %s] is wrong", addrPort)
	}
	port, err := strconv.ParseUint(addrPort[idx+1:], 10, 0)
	if err != nil {
		return uint(0), err
	}
	return uint(port), nil
}

// AddFlags add common flags for cc api config
func (conf *CCAPIConfig) AddFlags(fs *pflag.FlagSet, defaultAddrPort string) {
	fs.StringVar(&conf.AddrPort, "addrport", defaultAddrPort, "The ip address and port for the serve on")
	fs.StringVar(&conf.RegDiscover, "regdiscv", "", "hosts of register and discover server. e.g: 127.0.0.1:2181")
	fs.StringVar(&conf.ExConfig, "config", "", "The config path. e.g conf/api.conf")
	fs.StringVar(&conf.RegisterIP, "register-ip", "", "the ip address registered on zookeeper, it can be domain")
	fs.StringVar(&conf.Environment, "env", "", "the environment of the server, used for service discovery")
}
