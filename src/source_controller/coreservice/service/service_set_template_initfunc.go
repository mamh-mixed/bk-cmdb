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

func (s *coreService) initSetTemplate(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/topo/set_template/bk_biz_id/{bk_biz_id}/", Handler: s.CreateSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/", Handler: s.UpdateSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/topo/set_template/bk_biz_id/{bk_biz_id}/", Handler: s.DeleteSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/", Handler: s.GetSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/topo/set_template/bk_biz_id/{bk_biz_id}/", Handler: s.ListSetTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/topo/set_template/count_instances/bk_biz_id/{bk_biz_id}/", Handler: s.CountSetTplInstances})
	utility.AddHandler(rest.Action{
		Verb:    http.MethodGet,
		Path:    "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/service_templates_relations",
		Handler: s.ListSetServiceTemplateRelations,
	})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/findmany/topo/set_template/{set_template_id}/bk_biz_id/{bk_biz_id}/service_templates", Handler: s.ListSetTplRelatedSvcTpl})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/set_template/attribute",
		Handler: s.CreateSetTemplateAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/set_template/attribute",
		Handler: s.UpdateSetTemplateAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/set_template/attribute",
		Handler: s.DeleteSetTemplateAttribute})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/set_template/attribute",
		Handler: s.ListSetTemplateAttribute})

	utility.AddToRestfulWebService(web)
}
