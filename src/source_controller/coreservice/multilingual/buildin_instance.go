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

package multilingual

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
)

// BuildInInstanceNamePkg TODO
var BuildInInstanceNamePkg = map[string]map[string][]string{
	common.BKInnerObjIDModule: {
		"1": {"inst_module_idle", common.BKModuleNameField},
		"2": {"inst_module_fault", common.BKModuleNameField},
		"3": {"inst_module_recycle", common.BKModuleNameField},
	},
	common.BKInnerObjIDApp: {
		"1": {"inst_biz_default", common.BKAppNameField},
	},
	common.BKInnerObjIDSet: {
		"1": {"inst_set_default", common.BKSetNameField},
	},
}

// TranslateInstanceName is used to translate build-in model(module/set/biz/plat) instance's name to the
// corresponding language.
// Note: these instances's name is related it's default field's value, different value have different name.
// such as the module's instance, the different meaning of default value is as follows:
// 0: a common module
// 1: a idle module
// 2: a fault module
// 3: a recycle module
func TranslateInstanceName(defLang language.DefaultCCLanguageIf, objectID string, instances []mapstr.MapStr) {

	if m, ok := BuildInInstanceNamePkg[objectID]; ok {
		for idx, inst := range instances {
			// 如果用户将默认的内置模块名或者内置"空闲机池"名修改了，就不需要国际化，直接跳过
			if v, ok := inst[common.BKModuleNameField]; ok &&
				!(v == common.DefaultResModuleName || v == common.DefaultFaultModuleName ||
					v == common.DefaultRecycleModuleName) {
				continue
			}
			if v, ok := inst[common.BKSetNameField]; ok && v != common.DefaultResSetName {
				continue
			}
			// get the default's value and it's corresponding infos from defaultNameLanguagePkg
			subResult := m[fmt.Sprint(instances[idx][common.BKDefaultField])]

			if len(subResult) >= 2 {
				instances[idx][subResult[1]] = util.FirstNotEmptyString(defLang.Language(subResult[0]),
					fmt.Sprint(instances[idx][subResult[1]]))
			}
		}
		return
	}

	// translate unassigned cloud area name
	if objectID == common.BKInnerObjIDPlat {
		for idx, inst := range instances {
			cloudAreaName := util.GetStrByInterface(inst[common.BKCloudNameField])
			if cloudAreaName == common.UnassignedCloudAreaName {
				instances[idx][common.BKCloudNameField] = util.FirstNotEmptyString(
					defLang.Language("inst_plat_unassigned"), cloudAreaName)
			}
		}
		return

	}
}
