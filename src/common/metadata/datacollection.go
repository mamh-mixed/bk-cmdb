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

package metadata

const (
	// HostSnapDataSourcesDelayQueue host snap data comes from the delay queue.
	HostSnapDataSourcesDelayQueue = "delay_queue"
	// HostSnapDataSourcesChannel the source of host snap data is channels such as redis or kafka.
	HostSnapDataSourcesChannel = "channel"

	// IPv4LoopBackIpPrefix ipv4 loopback address
	IPv4LoopBackIpPrefix = "127.0.0.1"

	// IPv6LoopBackIp ipv6 loopback address
	IPv6LoopBackIp = "::1"

	// IPv6LinkLocalAddressPrefix link local address
	IPv6LinkLocalAddressPrefix = "fe80"
)
