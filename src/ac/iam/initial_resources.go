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

package iam

import "configcenter/src/common/metadata"

var (
	businessParent = Parent{
		SystemID:   SystemIDCMDB,
		ResourceID: Business,
	}
)

// ResourceTypeIDMap TODO
var ResourceTypeIDMap = map[TypeID]string{
	Business:                 "业务",
	BizSet:                   "业务集",
	Project:                  "项目",
	BusinessForHostTrans:     "业务主机",
	SysCloudArea:             "管控区域",
	SysResourcePoolDirectory: "主机池目录",
	SysHostRscPoolDirectory:  "主机池主机",
	SysModelGroup:            "模型分组",
	SysInstanceModel:         "实例模型",
	SysModel:                 "模型",
	SysModelEvent:            "模型列表",
	MainlineModelEvent:       "资源事件",
	InstAsstEvent:            "实例关联事件",
	KubeWorkloadEvent:        "容器工作负载事件",
	// SysInstance:               "实例",
	SysAssociationType:        "关联类型",
	SysOperationStatistic:     "运营统计",
	SysAuditLog:               "操作审计",
	SysCloudAccount:           "云账户",
	SysCloudResourceTask:      "云资源任务",
	SysEventWatch:             "事件监听",
	Host:                      "主机",
	BizHostApply:              "主机自动应用",
	BizCustomQuery:            "动态分组",
	BizCustomField:            "自定义字段",
	BizProcessServiceInstance: "服务实例",
	BizProcessServiceCategory: "服务分类",
	BizSetTemplate:            "集群模板",
	BizTopology:               "业务拓扑",
	BizProcessServiceTemplate: "服务模板",
	FieldGroupingTemplate:     "字段组合模板",
	GeneralCache:              "通用缓存",
	Set:                       "集群",
	Module:                    "模块",
}

// GenerateResourceTypes generate all the resource types registered to IAM.
func GenerateResourceTypes(models []metadata.Object) []ResourceType {
	resourceTypeList := make([]ResourceType, 0)

	// add public and business resources
	resourceTypeList = append(resourceTypeList, GenerateStaticResourceTypes()...)

	// add dynamic resources
	resourceTypeList = append(resourceTypeList, genDynamicResourceTypes(models)...)

	return resourceTypeList
}

// GenerateStaticResourceTypes TODO
func GenerateStaticResourceTypes() []ResourceType {
	resourceTypeList := make([]ResourceType, 0)

	// add public resources
	resourceTypeList = append(resourceTypeList, genPublicResources()...)

	// add business resources
	resourceTypeList = append(resourceTypeList, genBusinessResources()...)
	return resourceTypeList
}

// GetResourceParentMap generate resource types' mapping to parents.
func GetResourceParentMap() map[TypeID][]TypeID {
	resourceParentMap := make(map[TypeID][]TypeID, 0)
	for _, resourceType := range GenerateStaticResourceTypes() {
		for _, parent := range resourceType.Parents {
			resourceParentMap[resourceType.ID] = append(resourceParentMap[resourceType.ID], parent.ResourceID)
		}
	}
	return resourceParentMap
}

func genBusinessResources() []ResourceType {
	return []ResourceType{
		{
			ID:            Host,
			Name:          ResourceTypeIDMap[Host],
			NameEn:        "Host",
			Description:   "主机",
			DescriptionEn: "hosts under a business or in resource pool",
			Parents: []Parent{{
				SystemID: SystemIDCMDB,
				// ResourceID: Module,
				ResourceID: Business,
			}, {
				SystemID:   SystemIDCMDB,
				ResourceID: SysResourcePoolDirectory,
			}},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizHostApply,
			Name:          ResourceTypeIDMap[BizHostApply],
			NameEn:        "Host Apply",
			Description:   "自动应用业务主机的属性信息",
			DescriptionEn: "apply business host attribute automatically",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizCustomQuery,
			Name:          ResourceTypeIDMap[BizCustomQuery],
			NameEn:        "Dynamic Grouping",
			Description:   "根据条件查询主机信息",
			DescriptionEn: "custom query the host instances",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizCustomField,
			Name:          ResourceTypeIDMap[BizCustomField],
			NameEn:        "Custom Field",
			Description:   "模型在业务下的自定义字段",
			DescriptionEn: "model's custom field under a business",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizProcessServiceInstance,
			Name:          ResourceTypeIDMap[BizProcessServiceInstance],
			NameEn:        "Service Instance",
			Description:   "服务实例",
			DescriptionEn: "service instance",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizProcessServiceCategory,
			Name:          ResourceTypeIDMap[BizProcessServiceCategory],
			NameEn:        "Service Category",
			Description:   "服务分类用于分类服务实例",
			DescriptionEn: "service category is to classify service instances",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizSetTemplate,
			Name:          ResourceTypeIDMap[BizSetTemplate],
			NameEn:        "Set Template",
			Description:   "集群模板用于实例化集群",
			DescriptionEn: "set template is used to instantiate a set",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizTopology,
			Name:          ResourceTypeIDMap[BizTopology],
			NameEn:        "Business Topology",
			Description:   "业务拓扑包含了业务拓扑树中所有的相关元素",
			DescriptionEn: "business topology contains all elements related to the business topology tree",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BizProcessServiceTemplate,
			Name:          ResourceTypeIDMap[BizProcessServiceTemplate],
			NameEn:        "Service Template",
			Description:   "服务模板用于实例化服务实例",
			DescriptionEn: "service template is used to instantiate a service instance ",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		// only for biz topology usage, not related to actions
		{
			ID:            Set,
			Name:          ResourceTypeIDMap[Set],
			NameEn:        "Set",
			Description:   "业务拓扑集群",
			DescriptionEn: "business topology set",
			Parents:       []Parent{businessParent},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            Module,
			Name:          ResourceTypeIDMap[Module],
			NameEn:        "Module",
			Description:   "业务拓扑模块",
			DescriptionEn: "business topology module",
			Parents: []Parent{{
				SystemID:   SystemIDCMDB,
				ResourceID: Set,
			}},
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
	}
}

func genPublicResources() []ResourceType {
	return []ResourceType{
		{
			ID:            BizSet,
			Name:          ResourceTypeIDMap[BizSet],
			NameEn:        "Business Set",
			Description:   "业务集",
			DescriptionEn: "business set",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            Business,
			Name:          ResourceTypeIDMap[Business],
			NameEn:        "Business",
			Description:   "业务列表",
			DescriptionEn: "all the business in blueking cmdb.",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            Project,
			Name:          ResourceTypeIDMap[Project],
			NameEn:        "Project",
			Description:   "项目列表",
			DescriptionEn: "all the project in blueking cmdb.",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            BusinessForHostTrans,
			Name:          ResourceTypeIDMap[BusinessForHostTrans],
			NameEn:        "Host In Business",
			Description:   "业务主机",
			DescriptionEn: "host in business",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysCloudArea,
			Name:          ResourceTypeIDMap[SysCloudArea],
			NameEn:        "Cloud Area",
			Description:   "管控区域",
			DescriptionEn: "cloud area",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysResourcePoolDirectory,
			Name:          ResourceTypeIDMap[SysResourcePoolDirectory],
			NameEn:        "Host Pool Directory",
			Description:   "主机池目录",
			DescriptionEn: "host pool directory",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysHostRscPoolDirectory,
			Name:          ResourceTypeIDMap[SysHostRscPoolDirectory],
			NameEn:        "Host In Host Pool Directory",
			Description:   "主机池主机",
			DescriptionEn: "host in host pool directory",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysModelGroup,
			Name:          ResourceTypeIDMap[SysModelGroup],
			NameEn:        "Model Group",
			Description:   "模型分组用于对模型进行归类",
			DescriptionEn: "group models by model group",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysInstanceModel,
			Name:          ResourceTypeIDMap[SysInstanceModel],
			NameEn:        "InstanceModel",
			Description:   "实例模型",
			DescriptionEn: "instance model",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysModel,
			Name:          ResourceTypeIDMap[SysModel],
			NameEn:        "Model",
			Description:   "模型",
			DescriptionEn: "model",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysAssociationType,
			Name:          ResourceTypeIDMap[SysAssociationType],
			NameEn:        "Association Type",
			Description:   "关联类型是模型关联关系的分类",
			DescriptionEn: "association type is the classification of model association",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysOperationStatistic,
			Name:          ResourceTypeIDMap[SysOperationStatistic],
			NameEn:        "Operational Statistics",
			Description:   "运营统计",
			DescriptionEn: "operational statistics",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysAuditLog,
			Name:          ResourceTypeIDMap[SysAuditLog],
			NameEn:        "Operation Audit",
			Description:   "操作审计",
			DescriptionEn: "audit log",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysCloudAccount,
			Name:          ResourceTypeIDMap[SysCloudAccount],
			NameEn:        "Cloud Account",
			Description:   "云账户",
			DescriptionEn: "cloud account",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysCloudResourceTask,
			Name:          ResourceTypeIDMap[SysCloudResourceTask],
			NameEn:        "Cloud Resource Task",
			Description:   "云资源任务",
			DescriptionEn: "cloud resource task",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysEventWatch,
			Name:          ResourceTypeIDMap[SysEventWatch],
			NameEn:        "Event Listen",
			Description:   "事件监听",
			DescriptionEn: "event watch",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            SysModelEvent,
			Name:          ResourceTypeIDMap[SysModelEvent],
			NameEn:        "Model List",
			Description:   "模型列表",
			DescriptionEn: "model list",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            MainlineModelEvent,
			Name:          ResourceTypeIDMap[MainlineModelEvent],
			NameEn:        "Resource Event",
			Description:   "资源事件",
			DescriptionEn: "resource event",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            InstAsstEvent,
			Name:          ResourceTypeIDMap[InstAsstEvent],
			NameEn:        "Instance Association Event",
			Description:   "实例关联事件",
			DescriptionEn: "instance association event",
			Parents:       nil,
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            KubeWorkloadEvent,
			Name:          ResourceTypeIDMap[KubeWorkloadEvent],
			NameEn:        "Kube Workload Event",
			Description:   "容器工作负载事件",
			DescriptionEn: "kube workload event",
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            FieldGroupingTemplate,
			Name:          ResourceTypeIDMap[FieldGroupingTemplate],
			NameEn:        "Field Grouping Template",
			Description:   "字段组合模板",
			DescriptionEn: "Field Grouping Template",
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
		{
			ID:            GeneralCache,
			Name:          ResourceTypeIDMap[GeneralCache],
			NameEn:        "General Resource Cache",
			Description:   "通用缓存",
			DescriptionEn: "general resource cache",
			ProviderConfig: ResourceConfig{
				Path: "/auth/v3/find/resource",
			},
			Version: 1,
		},
	}
}
