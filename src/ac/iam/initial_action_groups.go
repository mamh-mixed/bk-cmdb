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

// GenerateActionGroups generate all the resource action groups registered to IAM.
func GenerateActionGroups(objects []metadata.Object) []ActionGroup {
	ActionGroups := GenerateStaticActionGroups()

	// generate model instance manage action groups, contains model instance related actions which are dynamic
	ActionGroups = append(ActionGroups, GenModelInstanceManageActionGroups(objects)...)

	return ActionGroups
}

// GenerateStaticActionGroups generate all the static resource action groups.
func GenerateStaticActionGroups() []ActionGroup {
	ActionGroups := make([]ActionGroup, 0)

	// generate business set manage action groups, contains fulltext search related actions
	ActionGroups = append(ActionGroups, genFulltextSearchServiceActionGroups()...)

	// generate business set manage action groups, contains business set related actions
	ActionGroups = append(ActionGroups, genBizSetManageActionGroups()...)

	// generate business manage action groups, contains business related actions
	ActionGroups = append(ActionGroups, genBusinessManageActionGroups()...)

	// generate resource manage action groups, contains resource related actions
	ActionGroups = append(ActionGroups, genResourceManageActionGroups()...)

	// generate model manage action groups, contains model related actions
	ActionGroups = append(ActionGroups, genModelManageActionGroups()...)

	// generate operation statistic action groups, contains operation statistic and audit related actions
	ActionGroups = append(ActionGroups, genOperationStatisticActionGroups()...)

	// generate global settings action groups, contains global settings related actions
	ActionGroups = append(ActionGroups, genGlobalSettingsActionGroups()...)

	// generate container management action groups, contains container related actions
	ActionGroups = append(ActionGroups, genContainerManagementActionGroups()...)

	return ActionGroups
}

func genFulltextSearchServiceActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "检索服务",
			NameEn: "Fulltext Search Service",
			Actions: []ActionWithID{
				{
					ID: UseFulltextSearch,
				},
			},
		},
	}
}

func genBizSetManageActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "业务集管理",
			NameEn: "Business Set Manage",
			Actions: []ActionWithID{
				{
					ID: AccessBizSet,
				},
			},
		},
	}
}

func genBusinessManageActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "业务管理",
			NameEn: "Business Manage",
			Actions: []ActionWithID{
				{
					ID: ViewBusinessResource,
				},
			},
			SubGroups: []ActionGroup{
				{
					Name:   "业务主机",
					NameEn: "Business Host",
					Actions: []ActionWithID{
						{
							ID: EditBusinessHost,
						},
						{
							ID: BusinessHostTransferToResourcePool,
						},
						{
							ID: HostTransferAcrossBusiness,
						},
					},
				},
				{
					Name:   "业务拓扑",
					NameEn: "Business Topology",
					Actions: []ActionWithID{
						{
							ID: CreateBusinessTopology,
						},
						{
							ID: EditBusinessTopology,
						},
						{
							ID: DeleteBusinessTopology,
						},
					},
				},
				{
					Name:   "服务实例",
					NameEn: "Service Instance",
					Actions: []ActionWithID{
						{
							ID: CreateBusinessServiceInstance,
						},
						{
							ID: EditBusinessServiceInstance,
						},
						{
							ID: DeleteBusinessServiceInstance,
						},
					},
				},
				{
					Name:   "服务模版",
					NameEn: "Service Template",
					Actions: []ActionWithID{
						{
							ID: CreateBusinessServiceTemplate,
						},
						{
							ID: EditBusinessServiceTemplate,
						},
						{
							ID: DeleteBusinessServiceTemplate,
						},
					},
				},
				{
					Name:   "集群模版",
					NameEn: "Set Template",
					Actions: []ActionWithID{
						{
							ID: CreateBusinessSetTemplate,
						},
						{
							ID: EditBusinessSetTemplate,
						},
						{
							ID: DeleteBusinessSetTemplate,
						},
					},
				},
				{
					Name:   "服务分类",
					NameEn: "Service Category",
					Actions: []ActionWithID{
						{
							ID: CreateBusinessServiceCategory,
						},
						{
							ID: EditBusinessServiceCategory,
						},
						{
							ID: DeleteBusinessServiceCategory,
						},
					},
				},
				{
					Name:   "动态分组",
					NameEn: "Dynamic Grouping",
					Actions: []ActionWithID{
						{
							ID: CreateBusinessCustomQuery,
						},
						{
							ID: EditBusinessCustomQuery,
						},
						{
							ID: DeleteBusinessCustomQuery,
						},
					},
				},
				{
					Name:   "业务自定义字段",
					NameEn: "Business Custom Field",
					Actions: []ActionWithID{
						{
							ID: EditBusinessCustomField,
						},
					},
				},
				{
					Name:   "主机自动应用",
					NameEn: "Business Host Apply",
					Actions: []ActionWithID{
						{
							ID: EditBusinessHostApply,
						},
					},
				},
			},
		},
	}
}

func genResourceManageActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "资源管理",
			NameEn: "Resource Manage",
			SubGroups: []ActionGroup{
				{
					Name:   "主机池",
					NameEn: "Host Pool",
					Actions: []ActionWithID{
						{
							ID: ViewResourcePoolHost,
						},
						{
							ID: CreateResourcePoolHost,
						},
						{
							ID: EditResourcePoolHost,
						},
						{
							ID: DeleteResourcePoolHost,
						},
						{
							ID: ResourcePoolHostTransferToBusiness,
						},
						{
							ID: ResourcePoolHostTransferToDirectory,
						},
						{
							ID: CreateResourcePoolDirectory,
						},
						{
							ID: EditResourcePoolDirectory,
						},
						{
							ID: DeleteResourcePoolDirectory,
						},
						{
							ID: ManageHostAgentID,
						},
					},
				},
				{
					Name:   "业务",
					NameEn: "Business",
					Actions: []ActionWithID{
						{
							ID: CreateBusiness,
						},
						{
							ID: EditBusiness,
						},
						{
							ID: ArchiveBusiness,
						},
						{
							ID: FindBusiness,
						},
					},
				},
				{
					Name:   "项目",
					NameEn: "Project",
					Actions: []ActionWithID{
						{
							ID: CreateProject,
						},
						{
							ID: EditProject,
						},
						{
							ID: DeleteProject,
						},
						{
							ID: ViewProject,
						},
					},
				},
				{
					Name:   "业务集",
					NameEn: "BizSet",
					Actions: []ActionWithID{
						{
							ID: CreateBizSet,
						},
						{
							ID: EditBizSet,
						},
						{
							ID: DeleteBizSet,
						},
						{
							ID: ViewBizSet,
						},
					},
				},
				{
					Name:   "云账户",
					NameEn: "Cloud Account",
					Actions: []ActionWithID{
						{
							ID: CreateCloudAccount,
						},
						{
							ID: EditCloudAccount,
						},
						{
							ID: DeleteCloudAccount,
						},
						{
							ID: FindCloudAccount,
						},
					},
				},
				{
					Name:   "云资源任务",
					NameEn: "Cloud Resource Task",
					Actions: []ActionWithID{
						{
							ID: CreateCloudResourceTask,
						},
						{
							ID: EditCloudResourceTask,
						},
						{
							ID: DeleteCloudResourceTask,
						},
						{
							ID: FindCloudResourceTask,
						},
					},
				},
				{
					Name:   "管控区域",
					NameEn: "Cloud Area",
					Actions: []ActionWithID{
						{
							ID: ViewCloudArea,
						},
						{
							ID: CreateCloudArea,
						},
						{
							ID: EditCloudArea,
						},
						{
							ID: DeleteCloudArea,
						},
					},
				},
				{
					Name:   "事件监听",
					NameEn: "Event Watch",
					Actions: []ActionWithID{
						{
							ID: WatchHostEvent,
						},
						{
							ID: WatchHostRelationEvent,
						},
						{
							ID: WatchBizEvent,
						},
						{
							ID: WatchSetEvent,
						},
						{
							ID: WatchModuleEvent,
						},
						{
							ID: WatchProcessEvent,
						},
						{
							ID: WatchCommonInstanceEvent,
						},
						{
							ID: WatchMainlineInstanceEvent,
						},
						{
							ID: WatchInstAsstEvent,
						},
						{
							ID: WatchBizSetEvent,
						},
						{
							ID: WatchPlatEvent,
						},
						{
							ID: WatchKubeClusterEvent,
						},
						{
							ID: WatchKubeNodeEvent,
						},
						{
							ID: WatchKubeNamespaceEvent,
						},
						{
							ID: WatchKubeWorkloadEvent,
						},
						{
							ID: WatchKubePodEvent,
						},
						{
							ID: WatchProjectEvent,
						},
					},
				},
				{
					Name:   "全量同步缓存条件",
					NameEn: "Full Sync Condition",
					Actions: []ActionWithID{
						{ID: CreateFullSyncCond},
						{ID: ViewFullSyncCond},
						{ID: EditFullSyncCond},
						{ID: DeleteFullSyncCond},
					},
				},
				{
					Name:   "缓存",
					NameEn: "Cache",
					Actions: []ActionWithID{
						{ID: ViewGeneralCache},
					},
				},
			},
		},
	}
}

func genModelManageActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "模型管理",
			NameEn: "Model Manage",
			SubGroups: []ActionGroup{
				{
					Name:   "模型分组",
					NameEn: "Model Group",
					Actions: []ActionWithID{
						{ID: CreateModelGroup},
						{ID: EditModelGroup},
						{ID: DeleteModelGroup},
					},
				},
				{
					Name:   "模型关系",
					NameEn: "Model Relation",
					Actions: []ActionWithID{
						{
							ID: ViewModelTopo,
						},
						{
							ID: EditBusinessLayer,
						},
						{
							ID: EditModelTopologyView,
						},
					},
				},
				{
					Name:   "模型",
					NameEn: "Model",
					Actions: []ActionWithID{
						{
							ID: ViewSysModel,
						},
						{
							ID: CreateSysModel,
						},
						{
							ID: EditSysModel,
						},
						{
							ID: DeleteSysModel,
						},
					},
				},
				{
					Name:   "关联类型",
					NameEn: "Association Type",
					Actions: []ActionWithID{
						{ID: CreateAssociationType},
						{ID: EditAssociationType},
						{ID: DeleteAssociationType},
					},
				},
				{
					Name:   "字段组合模板",
					NameEn: "Field Grouping Template",
					Actions: []ActionWithID{
						{ID: CreateFieldGroupingTemplate},
						{ID: ViewFieldGroupingTemplate},
						{ID: EditFieldGroupingTemplate},
						{ID: DeleteFieldGroupingTemplate},
					},
				},
				{
					Name:   "ID规则自增ID",
					NameEn: "ID Rule Self-increasing ID",
					Actions: []ActionWithID{
						{ID: EditIDRuleIncrID},
					},
				},
			},
		},
	}
}

// GenModelInstanceManageActionGroups TODO
func GenModelInstanceManageActionGroups(objects []metadata.Object) []ActionGroup {
	if len(objects) == 0 {
		return make([]ActionGroup, 0)
	}

	subGroups := []ActionGroup{}
	for _, obj := range objects {
		subGroups = append(subGroups, genDynamicActionSubGroup(obj))
	}
	return []ActionGroup{
		{
			Name:      "模型实例管理",
			NameEn:    "Model instance Manage",
			SubGroups: subGroups,
		},
	}
}

func genContainerManagementActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "容器资源管理",
			NameEn: "Container Management",
			SubGroups: []ActionGroup{
				{
					Name:   "容器 Cluster",
					NameEn: "Container Cluster",
					Actions: []ActionWithID{
						{
							ID: CreateContainerCluster,
						},
						{
							ID: EditContainerCluster,
						},
						{
							ID: DeleteContainerCluster,
						},
					},
				}, {
					Name:   "容器 Node",
					NameEn: "Container Node",
					Actions: []ActionWithID{
						{
							ID: CreateContainerNode,
						},
						{
							ID: EditContainerNode,
						},
						{
							ID: DeleteContainerNode,
						},
					},
				}, {
					Name:   "容器命名空间",
					NameEn: "Container Namespace",
					Actions: []ActionWithID{
						{
							ID: CreateContainerNamespace,
						},
						{
							ID: EditContainerNamespace,
						},
						{
							ID: DeleteContainerNamespace,
						},
					},
				}, {
					Name:   "容器工作负载",
					NameEn: "Container Workload",
					Actions: []ActionWithID{
						{
							ID: CreateContainerWorkload,
						},
						{
							ID: EditContainerWorkload,
						},
						{
							ID: DeleteContainerWorkload,
						},
					},
				}, {
					Name:   "容器 Pod",
					NameEn: "Container Pod",
					Actions: []ActionWithID{
						{
							ID: CreateContainerPod,
						},
						{
							ID: DeleteContainerPod,
						},
					},
				},
			},
		},
	}
}

func genOperationStatisticActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "运营统计",
			NameEn: "Operation Statistic",
			SubGroups: []ActionGroup{
				{
					Name:   "运营统计",
					NameEn: "Operation Statistic",
					Actions: []ActionWithID{
						{
							ID: FindOperationStatistic,
						},
						{
							ID: EditOperationStatistic,
						},
					},
				},
				{
					Name:   "操作审计",
					NameEn: "Operation Audit",
					Actions: []ActionWithID{
						{
							ID: FindAuditLog,
						},
					},
				},
			},
		},
	}
}

func genGlobalSettingsActionGroups() []ActionGroup {
	return []ActionGroup{
		{
			Name:   "全局设置",
			NameEn: "Global Settings",
			SubGroups: []ActionGroup{
				{
					Name:   "全局设置",
					NameEn: "Global Settings",
					Actions: []ActionWithID{
						{
							ID: GlobalSettings,
						},
					},
				},
			},
		},
	}
}
