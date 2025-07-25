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

package topo_server_test

import (
	"context"
	"encoding/json"

	"configcenter/pkg/filter"
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	commonutil "configcenter/src/common/util"
	"configcenter/src/kube/types"
	"configcenter/src/test"
	"configcenter/src/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("namespace test", func() {
	ctx := context.Background()

	var bizID int64
	var clusterID int64
	var namespaceID int64
	nsName := "nsName"

	It("prepare environment, create business, cluster", func() {
		test.DeleteAllBizs()

		biz := map[string]interface{}{
			common.BKMaintainersField: "kube",
			common.BKAppNameField:     "biz_for_kube",
			"time_zone":               "Africa/Accra",
			"language":                "1",
		}
		bizResp, err := apiServerClient.CreateBiz(ctx, "0", header, biz)
		util.RegisterResponseWithRid(bizResp, header)
		Expect(err).NotTo(HaveOccurred())
		Expect(bizResp.Result).To(Equal(true))
		bizID, err = commonutil.GetInt64ByInterface(bizResp.Data[common.BKAppIDField])
		Expect(err).NotTo(HaveOccurred())

		clusterName := "cluster"
		schedulingEngine := "k8s"
		uid := "BCS-K8S-25001"
		xid := "cls-hox2lkf2"
		version := "0.1"
		networkType := "underlay"
		region := "shenzhen"
		vpc := "vpc-q6awe02n"
		environment := "prod"
		network := []string{"1.1.1.0/21"}
		clusterType := types.IndependentClusterType
		createCluster := &types.Cluster{
			BizID:            bizID,
			Name:             &clusterName,
			SchedulingEngine: &schedulingEngine,
			Uid:              &uid,
			Xid:              &xid,
			Version:          &version,
			NetworkType:      &networkType,
			Environment:      &environment,
			Region:           &region,
			Vpc:              &vpc,
			NetWork:          &network,
			Type:             &clusterType,
		}

		id, err := kubeClient.CreateCluster(ctx, header, createCluster)
		util.RegisterResponseWithRid(id, header)
		Expect(err).NotTo(HaveOccurred())
		clusterID = id

	})

	It("create namespace", func() {
		label := map[string]string{
			"test":  "test",
			"test2": "test2",
		}
		resourceQuotas := []types.ResourceQuota{
			{
				Hard: map[string]string{
					"memory": "200Gi",
					"pods":   "100",
					"cpu":    "10k",
				},
				ScopeSelector: &types.ScopeSelector{MatchExpressions: []types.ScopedResourceSelectorRequirement{
					{
						Values:    []string{"high"},
						Operator:  "In",
						ScopeName: "PriorityClass",
					},
				}},
			},
		}
		ns := types.Namespace{
			ClusterSpec: types.ClusterSpec{
				ClusterID: clusterID,
			},
			Name:           nsName,
			Labels:         &label,
			ResourceQuotas: &resourceQuotas,
		}
		createOpt := types.NsCreateOption{
			BizID: bizID,
			Data:  []types.Namespace{ns},
		}

		result, err := kubeClient.CreateNamespace(ctx, header, &createOpt)
		util.RegisterResponseWithRid(result, header)
		Expect(err).NotTo(HaveOccurred())
		namespaceID = result.IDs[0]
	})

	It("update namespace", func() {
		nsName = "nsName"
		label := map[string]string{
			"test": "test2",
		}
		resourceQuotas := []types.ResourceQuota{
			{
				Hard: map[string]string{
					"memory": "200Gi",
					"pods":   "1000",
					"cpu":    "15k",
				},
				ScopeSelector: &types.ScopeSelector{MatchExpressions: []types.ScopedResourceSelectorRequirement{
					{
						Values:    []string{"high"},
						Operator:  "In",
						ScopeName: "PriorityClass",
					},
				}},
			},
		}
		ns := &types.Namespace{
			ClusterSpec: types.ClusterSpec{
				ClusterID: clusterID,
			},
			Name:           nsName,
			Labels:         &label,
			ResourceQuotas: &resourceQuotas,
		}
		updateOpt := types.NsUpdateOption{
			BizID: bizID,
			NsUpdateByIDsOption: types.NsUpdateByIDsOption{
				IDs:  []int64{namespaceID},
				Data: ns,
			},
		}

		err := kubeClient.UpdateNamespace(ctx, header, &updateOpt)
		util.RegisterResponseWithRid(err, header)
		Expect(err).NotTo(HaveOccurred())
	})

	It("find namespace", func() {
		filter := &filter.Expression{
			RuleFactory: &filter.CombinedRule{
				Condition: filter.And,
				Rules: []filter.RuleFactory{
					&filter.AtomRule{
						Field:    types.BKClusterIDFiled,
						Operator: filter.Equal.Factory(),
						Value:    clusterID,
					},
					&filter.AtomRule{
						Field:    common.BKFieldName,
						Operator: filter.Equal.Factory(),
						Value:    nsName,
					},
					&filter.AtomRule{
						Field:    "labels",
						Operator: filter.Object.Factory(),
						Value: &filter.AtomRule{
							Field:    "test",
							Operator: filter.Equal.Factory(),
							Value:    "test2",
						},
					},
				},
			},
		}

		// get namespace data
		page := metadata.BasePage{
			Start: 0,
			Limit: 10,
		}
		fields := []string{common.BKFieldID}
		queryOpt := types.NsQueryOption{
			BizID:  bizID,
			Filter: filter,
			Page:   page,
			Fields: fields,
		}
		result, err := kubeClient.ListNamespace(ctx, header, &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(result.Info)).To(Equal(1))
		Expect(result.Info[0][common.BKFieldID].(json.Number).Int64()).To(Equal(namespaceID))

		// get namespace count
		page = metadata.BasePage{
			EnableCount: true,
		}
		queryOpt = types.NsQueryOption{
			BizID:  bizID,
			Filter: filter,
			Page:   page,
		}
		queryOpt.Page.EnableCount = true
		result, err = kubeClient.ListNamespace(ctx, header, &queryOpt)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Count).To(Equal(1))
	})

	It("delete namespace", func() {
		deleteOpt := types.NsDeleteOption{
			BizID: bizID,
			NsDeleteByIDsOption: types.NsDeleteByIDsOption{
				IDs: []int64{
					namespaceID,
				},
			},
		}

		err := kubeClient.DeleteNamespace(ctx, header, &deleteOpt)
		Expect(err).NotTo(HaveOccurred())
	})
})
