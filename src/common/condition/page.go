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

package condition

import (
	"fmt"
	"strconv"

	"configcenter/src/common"
)

const (
	// PageName TODO
	PageName = "page"
	// PageSort TODO
	PageSort = "sort"
	// PageStart TODO
	PageStart = "start"
	// DBFields TODO
	DBFields = "fields"
	// DBQueryCondition TODO
	DBQueryCondition = "condition"
)

// BasePage for paging query
type BasePage struct {
	Sort  string `json:"sort,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Start int    `json:"start,omitempty"`
}

// ParsePage TODO
func ParsePage(origin interface{}) BasePage {
	if origin == nil {
		return BasePage{Limit: common.BKNoLimit}
	}
	page, ok := origin.(map[string]interface{})
	if !ok {
		return BasePage{Limit: common.BKNoLimit}
	}
	result := BasePage{}
	if sort, ok := page["sort"]; ok && sort != nil {
		result.Sort = fmt.Sprint(sort)
	}
	if start, ok := page["start"]; ok {
		result.Start, _ = strconv.Atoi(fmt.Sprint(start))
	}
	if limit, ok := page["limit"]; ok {
		result.Limit, _ = strconv.Atoi(fmt.Sprint(limit))
		if result.Limit <= 0 {
			result.Limit = common.BKNoLimit
		}
	}
	return result
}
