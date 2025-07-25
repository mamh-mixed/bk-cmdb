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

package mainline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// SearchMainlineInstanceTopo get topo tree of mainline model
func (m *topoManager) SearchMainlineInstanceTopo(ctx context.Context, header http.Header, bkBizID int64,
	withDetail bool) (*metadata.TopoInstanceNode, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	bizTopoNode, err := m.SearchMainlineModelTopo(ctx, header, false)
	if err != nil {
		blog.Errorf("get mainline model topo info failed, %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get mainline model topo info failed, %+v", err)
	}
	blog.V(9).Infof("model mainline: %+v, rid: %s", bizTopoNode, rid)

	im, err := NewInstanceMainline(m.lang.CreateDefaultCCLanguageIf(httpheader.GetLanguage(header)), mongodb.Client(),
		bkBizID)
	if err != nil {
		blog.Errorf("SearchMainlineInstanceTopo failed, NewInstanceMainline failed, bizID: %d, err: %+v, rid: %s",
			bkBizID, err, rid)
		return nil, fmt.Errorf("new mainline instance by business:%d failed, %+v", bkBizID, err)
	}

	im.SetModelTree(ctx, bizTopoNode)
	im.LoadModelParentMap(ctx)

	if err := im.LoadSetInstances(ctx, header); err != nil {
		blog.Errorf("get set instances by business:%d failed, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("get set instances by business:%d failed, %+v", bkBizID, err)
	}

	if err := im.LoadModuleInstances(ctx, header); err != nil {
		blog.Errorf("get module instances by business:%d failed, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("get module instances by business:%d failed, %+v", bkBizID, err)
	}

	if err := im.LoadMainlineInstances(ctx, header); err != nil {
		blog.Errorf("get other mainline instances by business:%d failed, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("get other mainline instances by business:%d failed, %+v", bkBizID, err)
	}

	if err := im.ConstructBizTopoInstance(ctx, header, withDetail); err != nil {
		blog.Errorf("construct business:%d detail as topo instance failed, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("construct business:%d detail as topo instance failed, %+v", bkBizID, err)
	}

	if err := im.OrganizeSetInstance(ctx, withDetail); err != nil {
		blog.Errorf("organize set instance failed, businessID:%d, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("organize set instance failed, businessID:%d, %+v", bkBizID, err)
	}

	if err := im.OrganizeModuleInstance(ctx, withDetail); err != nil {
		blog.Errorf("organize module instance failed, businessID:%d, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("organize module instance failed, businessID:%d, %+v", bkBizID, err)
	}

	if err := im.OrganizeMainlineInstance(ctx, withDetail); err != nil {
		blog.Errorf("organize other mainline instance failed, businessID:%d, %+v, rid: %s", bkBizID, err, rid)

		return nil, fmt.Errorf("organize other mainline instance failed, businessID:%d, %+v", bkBizID, err)
	}

	instanceMap := im.GetInstanceMap(ctx)
	instanceMapStr, err := json.Marshal(instanceMap)
	if err != nil {
		blog.Errorf("json encode instanceMap:%+v failed, %+v, rid: %s", instanceMap, err, rid)
		return nil, fmt.Errorf("json encode instanceMap:%+v failed, %+v", instanceMap, err)
	}
	blog.V(5).Infof("instanceMap before check is: %s, rid: %s", instanceMapStr, rid)

	if err := im.CheckAndFillingMissingModels(ctx, header, withDetail); err != nil {
		blog.Errorf("check and filling missing models failed, business:%d %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("check and filling missing models failed, business:%d %+v", bkBizID, err)
	}

	instanceMapStr, err = json.Marshal(im.GetInstanceMap(ctx))
	if err != nil {
		blog.Errorf("json encode instanceMap failed, %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("json encode instanceMap failed, %+v", err)
	}
	blog.V(5).Infof("instanceMap after check: %s, rid: %s", instanceMapStr, rid)

	if err := im.ConstructInstanceTopoTree(ctx, header, withDetail); err != nil {
		blog.Errorf("get other mainline instances by business:%d failed, %+v, rid: %s", bkBizID, err, rid)
		return nil, fmt.Errorf("get other mainline instances by business:%d failed, %+v", bkBizID, err)
	}

	root := im.GetRoot(ctx)
	blog.V(9).Infof("topo instance tree root is: %+v, rid: %s", root, rid)
	treeData, err := json.Marshal(root)
	if err != nil {
		blog.Errorf("get other mainline instances by business:%d failed, %+v, rid: %s", bkBizID, err, rid)
		return root, nil
	}
	blog.V(9).Infof("topo instance tree root data is: %s, rid: %s", treeData, rid)
	return root, nil
}
