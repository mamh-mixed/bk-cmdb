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
	"net/http"

	"configcenter/src/common/http/rest"

	"github.com/emicklei/go-restful/v3"
)

func (s *cacheService) initCache(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/host/with_inner_ip",
		Handler: s.SearchHostWithInnerIPInCache})
	// Note: only for data-collection api!!!
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/host/with_agent_id",
		Handler: s.SearchHostWithAgentIDInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/host/with_host_id",
		Handler: s.SearchHostWithHostIDInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cache/host/with_host_id",
		Handler: s.ListHostWithHostIDInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cache/host/with_page",
		Handler: s.ListHostWithPageInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/biz/{bk_biz_id}",
		Handler: s.SearchBusinessInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cache/biz", Handler: s.ListBusinessInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/set/{bk_set_id}",
		Handler: s.SearchSetInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cache/set", Handler: s.ListSetsInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/module/{bk_module_id}",
		Handler: s.SearchModuleInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/cache/module",
		Handler: s.ListModulesInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/cache/{bk_obj_id}/{bk_inst_id}",
		Handler: s.SearchCustomLayerInCache})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "find/cache/topo/node_path/biz/{bk_biz_id}",
		Handler: s.SearchBizTopologyNodePath})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/cache/topo/brief/biz/{biz}",
		Handler: s.SearchBusinessBriefTopology})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/biz/{type}/topo", Handler: s.SearchBizTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/refresh/biz/{type}/topo", Handler: s.RefreshBizTopo})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/pod/label/key",
		Handler: s.ListPodLabelKey})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/kube/pod/label/value",
		Handler: s.ListPodLabelValue})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/refresh/kube/pod/label", Handler: s.RefreshPodLabel})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/watch/cache/event", Handler: s.WatchEvent})
	// Note: only for inner api!!!
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/inner/watch/cache/event", Handler: s.InnerWatchEvent})

	// conditional full sync scene resource cache related api
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/full/sync/cond",
		Handler: s.CreateFullSyncCond})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/full/sync/cond", Handler: s.UpdateFullSyncCond})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/full/sync/cond",
		Handler: s.DeleteFullSyncCond})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/full/sync/cond",
		Handler: s.ListFullSyncCond})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/resource/by_full_sync_cond",
		Handler: s.ListCacheByFullSyncCond})

	// general resource cache api
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/resource/by_ids",
		Handler: s.ListGeneralCacheByIDs})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/resource/by_unique_keys",
		Handler: s.ListGeneralCacheByUniqueKey})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/refresh/general/resource/id_list",
		Handler: s.RefreshGeneralResIDList})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/refresh/general/resource/detail/by_ids",
		Handler: s.RefreshGeneralResDetailByIDs})

	utility.AddToRestfulWebService(web)
}

func (s *cacheService) initService(web *restful.WebService) {

	s.initCache(web)
}
