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

package identifier

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	mixevent "configcenter/src/source_controller/cacheservice/event/mix-event"
	"configcenter/src/storage/stream/types"
)

const (
	hostIdentityLockKey = common.BKCacheKeyV3Prefix + "host_identity:event_lock"
	hostIdentityLockTTL = 1 * time.Minute
)

type identityOptions struct {
	key         event.Key
	watchFields []string
}

func (i *Identity) addWatchTask(opts identityOptions) error {
	identity := hostIdentity{
		identityOptions: opts,
		metrics:         event.InitialMetrics(opts.key.Collection(), "host_identifier"),
	}

	flowOpt := mixevent.MixEventFlowOptions{
		MixKey:       event.HostIdentityKey,
		Key:          opts.key,
		WatchFields:  opts.watchFields,
		EventLockTTL: hostIdentityLockTTL,
		EventLockKey: hostIdentityLockKey,
	}

	flow, err := mixevent.NewMixEventFlow(flowOpt, identity.rearrangeEvents, identity.parseEvent)
	if err != nil {
		return err
	}

	flowTask, err := flow.GenWatchTask()
	if err != nil {
		return err
	}

	i.tasks = append(i.tasks, flowTask)
	return nil
}

type hostIdentity struct {
	identityOptions
	metrics *event.EventMetrics
}

func (f *hostIdentity) rearrangeEvents(rid string, es []*types.Event) ([]*types.Event, error) {
	switch f.key.Collection() {
	case event.HostKey.Collection():
		return f.rearrangeHostEvents(es, rid), nil
	case event.ModuleHostRelationKey.Collection():
		return f.rearrangeHostRelationEvents(es, rid)
	case event.ProcessKey.Collection():
		return f.rearrangeProcessEvents(es, rid)
	case event.ProcessInstanceRelationKey.Collection():
		return f.rearrangeHostRelationEvents(es, rid)
	default:
		blog.ErrorJSON("received unsupported host identity event, skip, es: %s, rid :%s", es, rid)
		return es[:0], nil
	}
}

// parseEvent parse event into chain nodes, host identifier detail is formed when watched, do not store in redis
func (f *hostIdentity) parseEvent(e *types.Event, id uint64, rid string) (string, *watch.ChainNode, []byte, bool,
	error) {

	switch e.OperationType {
	case types.Insert, types.Update, types.Replace, types.Delete:
	case types.Invalidate:
		blog.Errorf("host identify event, received invalid event operation type, doc: %s, rid: %s", e.DocBytes, rid)
		return "", nil, nil, false, nil
	default:
		blog.Errorf("host identify event, received unsupported event operation type: %s, doc: %s, rid: %s",
			e.OperationType, e.DocBytes, rid)
		return "", nil, nil, false, nil
	}

	name := f.key.Name(e.DocBytes)
	cursor, err := genHostIdentifyCursor(f.key.Collection(), e, rid)
	if err != nil {
		blog.Errorf("get %s event cursor failed, name: %s, err: %v, oid: %s, rid: %s", f.key.Collection(), name,
			err, e.ID(), rid)
		return "", nil, nil, false, err
	}

	chainNode := &watch.ChainNode{
		ID:          id,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		// redirect all the event type to update.
		EventType: watch.ConvertOperateType(types.Update),
		Token:     e.Token.Data,
		Cursor:    cursor,
	}

	if instanceID := event.HostIdentityKey.InstanceID(e.DocBytes); instanceID > 0 {
		chainNode.InstanceID = instanceID
	}

	return e.TenantID, chainNode, nil, false, nil
}

func genHostIdentifyCursor(coll string, e *types.Event, rid string) (string, error) {
	curType := watch.UnknownType
	switch coll {
	case common.BKTableNameBaseHost:
		curType = watch.Host
	case common.BKTableNameModuleHostConfig:
		curType = watch.ModuleHostRelation
	case common.BKTableNameBaseProcess:
		curType = watch.Process
	case common.BKTableNameProcessInstanceRelation:
		curType = watch.ProcessInstanceRelation
	default:
		blog.ErrorJSON("unsupported host identity cursor type collection: %s, event: %s, oid: %s", coll, e, rid)
		return "", fmt.Errorf("unsupported host identity cursor type collection: %s", coll)
	}

	hCursor := &watch.Cursor{
		Type:        curType,
		ClusterTime: e.ClusterTime,
		Oid:         e.Oid,
		Oper:        e.OperationType,
		// UniqKey:     coll,
	}

	hCursorEncode, err := hCursor.Encode()
	if err != nil {
		blog.ErrorJSON("encode head node cursor failed, cursor: %s, err: %s, rid: %s", hCursor, err, rid)
		return "", err
	}

	return hCursorEncode, nil
}
