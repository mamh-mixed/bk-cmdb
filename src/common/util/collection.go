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

package util

import "strings"

// CalSliceDiff TODO
func CalSliceDiff(oldSlice, newSlice []string) (subs, plugs []string) {
	subs = make([]string, 0)
	plugs = make([]string, 0)
	for _, a := range oldSlice {
		if !Contains(newSlice, a) {
			subs = append(subs, a)
		}
	}
	for _, b := range newSlice {
		if !Contains(oldSlice, b) {
			plugs = append(plugs, b)
		}
	}
	return
}

// CaseInsensitiveContains TODO
func CaseInsensitiveContains(s string, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Contains if string target in array
func Contains(set []string, substr string) bool {
	for _, s := range set {
		if s == substr {
			return true
		}
	}
	return false
}

// ContainsInt64 if int64 target in array
func ContainsInt64(set []int64, sub int64) bool {
	for _, s := range set {
		if s == sub {
			return true
		}
	}
	return false
}

// ContainsInt if int target in array
func ContainsInt(set []int64, sub int64) bool {
	for _, s := range set {
		if s == sub {
			return true
		}
	}
	return false
}

// CalSliceInt64Diff TODO
func CalSliceInt64Diff(oldSlice, newSlice []int64) (subs, inter, plugs []int64) {
	subs = make([]int64, 0)
	inter = make([]int64, 0)
	plugs = make([]int64, 0)
	for _, a := range oldSlice {
		if !ContainsInt64(newSlice, a) {
			subs = append(subs, a)
		} else {
			inter = append(inter, a)
		}
	}
	for _, b := range newSlice {
		if !ContainsInt64(oldSlice, b) {
			plugs = append(plugs, b)
		}
	}
	return
}
