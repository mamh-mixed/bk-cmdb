/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package transfer

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
)

// TransferManager TODO
type TransferManager struct {
	dependence          OperationDependence
	hostApplyDependence HostApplyRuleDependence
}

// OperationDependence TODO
type OperationDependence interface {
	AutoCreateServiceInstanceModuleHost(kit *rest.Kit, hostIDs []int64, moduleIDs []int64) errors.CCErrorCoder
	SelectObjectAttWithParams(kit *rest.Kit, objID string, bizIDs []int64) (attribute []metadata.Attribute, err error)
	UpdateModelInstance(kit *rest.Kit, objID string, param metadata.UpdateOption) (*metadata.UpdatedCount, error)
	CreateAuditLogDependence(kit *rest.Kit, logs ...metadata.AuditLog) error
}

// HostApplyRuleDependence TODO
type HostApplyRuleDependence interface {
	RunHostApplyOnHosts(kit *rest.Kit, bizID int64, relations []metadata.ModuleHost) (metadata.MultipleHostApplyResult,
		errors.CCErrorCoder)
}

// New TODO
func New(dependence OperationDependence, hostApplyDependence HostApplyRuleDependence) *TransferManager {
	return &TransferManager{
		dependence:          dependence,
		hostApplyDependence: hostApplyDependence,
	}
}

// NewHostModuleTransfer business normal module transfer
func (manager *TransferManager) NewHostModuleTransfer(kit *rest.Kit, bizID int64, moduleIDArr []int64, isIncr bool,
	needAutoCreateSvcInst bool) *genericTransfer {
	return &genericTransfer{
		dependent:             manager.dependence,
		hostApplyDependence:   manager.hostApplyDependence,
		moduleIDArr:           moduleIDArr,
		bizID:                 bizID,
		isIncrement:           isIncr,
		needAutoCreateSvcInst: needAutoCreateSvcInst,
	}
}

// TransferToInnerModule transfer host to inner module, default module contain(idle module, fault module)
func (manager *TransferManager) TransferToInnerModule(kit *rest.Kit, input *metadata.TransferHostToInnerModule) error {

	transfer := manager.NewHostModuleTransfer(kit, input.ApplicationID, []int64{input.ModuleID}, false, false)

	exit, err := transfer.HasInnerModule(kit)
	if err != nil {
		blog.Errorf("check if moduleID is inner module failed, input: %v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}
	if !exit {
		blog.Errorf("module ID is not default module, input: %v, rid: %s", input, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCoreServiceModuleNotDefaultModuleErr, input.ModuleID,
			input.ApplicationID)
	}
	err = transfer.ValidParameter(kit)
	if err != nil {
		blog.Errorf("valid parameter failed, input: %v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	transferErr := transfer.Transfer(kit, input.HostID, false)
	if transferErr != nil {
		blog.Errorf("transfer module host relation failed, input: %v, hostID: %s, err: %v, rid: %s", input,
			input.HostID, transferErr, kit.Rid)
		return transferErr
	}

	return nil
}

// TransferToNormalModule transfer host to use add module
// 目标模块不能为空闲机模块
func (manager *TransferManager) TransferToNormalModule(kit *rest.Kit, input *metadata.HostsModuleRelation) error {
	// 确保目标模块不能为空闲机模块
	defaultModuleFilter := map[string]interface{}{
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultFlagDefaultValue,
		},
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: input.ModuleID,
		},
	}
	defaultModuleCount, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).
		Find(defaultModuleFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("filter default module failed, filter: %v, err: %v, rid: %s", defaultModuleFilter,
			common.BKTableNameBaseModule, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if defaultModuleCount > 0 {
		blog.Errorf("target module shouldn't be default module, input: %v, defaultModuleCount: %d, rid: %s", input,
			defaultModuleCount, kit.Rid)
		return kit.CCError.CCError(common.CCErrCoreServiceTransferToDefaultModuleUseWrongMethod)
	}

	transfer := manager.NewHostModuleTransfer(kit, input.ApplicationID, input.ModuleID, input.IsIncrement,
		!input.DisableAutoCreateSvcInst)

	err = transfer.ValidParameter(kit)
	if err != nil {
		blog.Errorf("valid parameter failed, input: %v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}
	err = transfer.Transfer(kit, input.HostID, input.DisableTransferHostAutoApply)
	if err != nil {
		blog.Errorf("transfer module host relation failed, input: %v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	return nil
}

// RemoveFromModule 将主机从模块中移出
// 如果主机属于n+1个模块（n>0），操作之后，主机属于n个模块
// 如果主机属于1个模块, 且非空闲机模块，操作之后，主机属于空闲机模块
// 如果主机属于空闲机模块，操作失败
// 如果主机属于故障机模块，操作失败
// 如果主机不在参数指定的模块中，操作失败
func (manager *TransferManager) RemoveFromModule(kit *rest.Kit, input *metadata.RemoveHostsFromModuleOption) error {
	hostConfigFilter := map[string]interface{}{
		common.BKHostIDField: input.HostID,
		common.BKAppIDField:  input.ApplicationID,
	}
	hostConfigs := make([]metadata.ModuleHost, 0)
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Find(hostConfigFilter).All(kit.Ctx,
		&hostConfigs)
	if err != nil {
		blog.Errorf("find host module config failed, filter: %+v, table: %s, err: %v, rid: %s", hostConfigFilter,
			common.BKTableNameModuleHostConfig, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrHostModuleConfigFailed, err)
	}

	// 如果主机不在参数指定的模块中，操作失败
	if len(hostConfigs) == 0 {
		blog.Errorf("host module config not found, input: %+v, rid: %s", input, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrHostModuleNotExist)
	}

	moduleIDs := make([]int64, 0)
	for _, hostConfig := range hostConfigs {
		moduleIDs = append(moduleIDs, hostConfig.ModuleID)
	}

	// 检查 moduleIDs 是否有空闲机或故障机模块
	// 如果主机属于空闲机模块，操作失败
	// 如果主机属于故障机模块，操作失败
	defaultModuleFilter := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDs,
		},
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultFlagDefaultValue,
		},
	}
	defaultModuleCount, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(
		defaultModuleFilter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get default module failed, filter: %+v, table: %s, err: %v, rid: %s", defaultModuleFilter,
			common.BKTableNameBaseModule, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrHostGetModuleFail, err)
	}
	if defaultModuleCount > 0 {
		blog.Errorf("default module shouldn't in target modules, input: %+v, rid: %v", input, kit.Rid)
		return kit.CCError.CCError(common.CCErrHostRemoveFromDefaultModuleFailed)
	}

	targetModuleIDs := make([]int64, 0)
	for _, moduleID := range moduleIDs {
		if moduleID != input.ModuleID {
			targetModuleIDs = append(targetModuleIDs, moduleID)
		}
	}
	if len(targetModuleIDs) > 0 {
		option := metadata.HostsModuleRelation{
			ApplicationID:            input.ApplicationID,
			HostID:                   []int64{input.HostID},
			ModuleID:                 targetModuleIDs,
			IsIncrement:              false,
			DisableAutoCreateSvcInst: true,
		}
		err = manager.TransferToNormalModule(kit, &option)
		if err != nil {
			blog.Errorf("TransferToNormalModule failed, input: %+v, option: %+v, err: %v, rid: %s", input, option, err,
				kit.Rid)
			return err
		}
		return nil
	}

	// transfer host to idle module
	idleModuleFilter := map[string]interface{}{
		common.BKAppIDField:   input.ApplicationID,
		common.BKDefaultField: common.DefaultResModuleFlag,
	}
	idleModule := metadata.ModuleHost{}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(idleModuleFilter).One(kit.Ctx,
		&idleModule); err != nil {
		return kit.CCError.CCErrorf(common.CCErrHostGetModuleFail, err)
	}
	innerModuleOption := metadata.TransferHostToInnerModule{
		ApplicationID: input.ApplicationID,
		ModuleID:      idleModule.ModuleID,
		HostID:        []int64{input.HostID},
	}

	if err := manager.TransferToInnerModule(kit, &innerModuleOption); err != nil {
		blog.Errorf("TransferToInnerModule failed, filter: %+v, option: %+v, err: %v, rid: %s", input,
			innerModuleOption, err, kit.Rid)
		return err
	}
	return nil
}

// TransferToAnotherBusiness Host cross-business transfer
func (manager *TransferManager) TransferToAnotherBusiness(kit *rest.Kit,
	input *metadata.TransferHostsCrossBusinessRequest) error {
	transfer := manager.NewHostModuleTransfer(kit, input.DstApplicationID, input.DstModuleIDArr, false, true)
	transfer.SetCrossBusiness(kit, input.SrcApplicationIDs)
	var err error
	err = transfer.ValidParameter(kit)
	if err != nil {
		blog.Errorf("validParameter failed, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	// attributes in legacy business
	legacyAttributes, err := transfer.dependent.SelectObjectAttWithParams(kit, common.BKInnerObjIDHost,
		input.SrcApplicationIDs)
	if err != nil {
		blog.Errorf("select objectAtt with params failed, bizIDs: %v, err: %v, rid: %s", input.SrcApplicationIDs, err,
			kit.Rid)
		return err
	}

	err = transfer.Transfer(kit, input.HostIDArr, false)
	if err != nil {
		blog.Errorf("transfer to another business failed, transfer module host relation error, "+
			"err: %v, input: %+v, rid: %s", err, input, kit.Rid)
		return err
	}

	// reset private field in legacy business
	if err := manager.clearLegacyPrivateField(kit, legacyAttributes, input.HostIDArr...); err != nil {
		blog.Errorf("clearLegacyPrivateField failed, hostID: %v, attributes: %+v, err: %v, rid: %s", input.HostIDArr,
			legacyAttributes, err, kit.Rid)
		// we should go on setting default value for new private field
	}

	return nil
}

func (manager *TransferManager) clearLegacyPrivateField(kit *rest.Kit, attributes []metadata.Attribute,
	hostIDs ...int64) errors.CCErrorCoder {
	doc := make(map[string]interface{}, 0)
	for _, attribute := range attributes {
		if attribute.BizID == 0 {
			continue
		}
		doc[attribute.PropertyID] = nil
	}
	if len(doc) == 0 {
		return nil
	}
	reset := types.ModeUpdate{
		Op:  "unset",
		Doc: doc,
	}
	filter := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostIDs,
		},
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).UpdateMultiModel(kit.Ctx, filter,
		reset); err != nil {
		blog.Errorf("update host info failed, table: %s, filter: %+v, doc: %+v, err: %v, rid: %s",
			common.BKTableNameBaseHost, filter, doc, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

func (manager *TransferManager) setDefaultPrivateField(kit *rest.Kit, attributes []metadata.Attribute,
	hostID ...int64) errors.CCErrorCoder {
	doc := make(map[string]interface{})
	for _, attribute := range attributes {
		if attribute.BizID == 0 {
			continue
		}
		doc[attribute.PropertyID] = nil
	}
	if len(doc) == 0 {
		return nil
	}
	updateOption := metadata.UpdateOption{
		Data: doc,
		Condition: map[string]interface{}{
			common.BKHostIDField: map[string]interface{}{
				common.BKDBIN: hostID,
			},
		},
	}
	_, err := manager.dependence.UpdateModelInstance(kit, common.BKInnerObjIDHost, updateOption)
	if err != nil {
		blog.Errorf("UpdateModelInstance failed, option: %+v, err: %v, rid: %s", common.BKTableNameBaseHost,
			updateOption, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

// GetHostModuleRelation get host module relation
func (manager *TransferManager) GetHostModuleRelation(kit *rest.Kit, input *metadata.HostModuleRelationRequest) (
	*metadata.HostConfigData, error) {

	if input.Empty() {
		blog.Errorf("input empty, input: %#v, rid: %s", input, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "GetHostModuleRelation input")
	}
	moduleHostCond := condition.CreateCondition()
	if input.ApplicationID > 0 {
		moduleHostCond.Field(common.BKAppIDField).Eq(input.ApplicationID)
	}
	if len(input.HostIDArr) > 0 {
		moduleHostCond.Field(common.BKHostIDField).In(input.HostIDArr)
	}
	if len(input.ModuleIDArr) > 0 {
		moduleHostCond.Field(common.BKModuleIDField).In(input.ModuleIDArr)
	}
	if len(input.SetIDArr) > 0 {
		moduleHostCond.Field(common.BKSetIDField).In(input.SetIDArr)
	}
	cond := moduleHostCond.ToMapStr()
	if len(cond) == 0 {
		return nil, nil
	}

	cnt, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Find(cond).
		Fields(input.Fields...).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get module host config count failed, err: %v, cond: %#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	hostModuleArr := make([]metadata.ModuleHost, 0)
	db := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).
		Find(cond).
		Fields(input.Fields...).
		Start(uint64(input.Page.Start)).
		Sort(input.Page.Sort)

	if input.Page.Limit > 0 {
		db = db.Limit(uint64(input.Page.Limit))
	}

	err = db.All(kit.Ctx, &hostModuleArr)
	if err != nil {
		blog.Errorf("get module host config failed, err: %v, cond: %#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	return &metadata.HostConfigData{
		Count: int64(cnt),
		Info:  hostModuleArr,
		Page:  input.Page,
	}, nil
}

// DeleteFromSystem delete host module relation and host info
func (manager *TransferManager) DeleteFromSystem(kit *rest.Kit, input *metadata.DeleteHostRequest) error {
	transfer := manager.NewHostModuleTransfer(kit, input.ApplicationID, nil, false, false)
	return transfer.DeleteHosts(kit, input.HostIDArr)
}

func (manager *TransferManager) getHostIDModuleMapByHostID(kit *rest.Kit, appID int64,
	hostIDArr []int64) (map[int64][]metadata.ModuleHost, errors.CCErrorCoder) {
	moduleHostCond := condition.CreateCondition()
	moduleHostCond.Field(common.BKAppIDField).Eq(appID)
	moduleHostCond.Field(common.BKHostIDField).In(hostIDArr)
	cond := moduleHostCond.ToMapStr()

	var dataArr []metadata.ModuleHost
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Find(cond).All(kit.Ctx, &dataArr)
	if err != nil {
		blog.Errorf("get module host config error, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	result := make(map[int64][]metadata.ModuleHost, 0)
	for _, item := range dataArr {
		result[item.HostID] = append(result[item.HostID], item)
	}
	return result, nil
}

// GetDistinctHostIDsByTopoRelation get all  host ids by topology relation condition
func (manager *TransferManager) GetDistinctHostIDsByTopoRelation(kit *rest.Kit,
	input *metadata.DistinctHostIDByTopoRelationRequest) ([]int64, error) {
	if input.Empty() {
		blog.Errorf("input empty, input: %#v, rid: %s", input, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)
	}
	moduleHostCond := condition.CreateCondition()
	if len(input.ApplicationIDArr) > 0 {
		moduleHostCond.Field(common.BKAppIDField).In(input.ApplicationIDArr)
	}
	if len(input.HostIDArr) > 0 {
		moduleHostCond.Field(common.BKHostIDField).In(input.HostIDArr)
	}
	if len(input.ModuleIDArr) > 0 {
		moduleHostCond.Field(common.BKModuleIDField).In(input.ModuleIDArr)
	}
	if len(input.SetIDArr) > 0 {
		moduleHostCond.Field(common.BKSetIDField).In(input.SetIDArr)
	}
	cond := moduleHostCond.ToMapStr()
	if len(cond) == 0 {
		return nil, nil
	}
	cond = moduleHostCond.ToMapStr()

	// 根据约束cond,获得去重后的主机id.
	ret, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Distinct(kit.Ctx,
		common.BKHostIDField, cond)
	if err != nil {
		blog.Errorf("get module host config failed, err: %v, cond: %#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	// 将ret转化为[]int64
	var hostIDArr []int64
	if hostIDArr, err = util.SliceInterfaceToInt64(ret); err != nil {
		blog.Errorf("slice interface to int64 failed, err: %v, cond: %#v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	return hostIDArr, nil
}

// TransferResourceDirectory transfer host in resource directory
func (manager *TransferManager) TransferResourceDirectory(kit *rest.Kit,
	input *metadata.TransferHostResourceDirectory) errors.CCErrorCoder {

	// validate input bk_module_id
	module, err := manager.validTransferResourceDirParams(kit, input)
	if err != nil {
		blog.Errorf("validate input(%+v) failed, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	cond := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: input.HostID,
		},
	}
	deleteErr := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Delete(kit.Ctx, cond)
	if deleteErr != nil {
		blog.Errorf("delete module host config failed, err: %v, cond: %v, rid: %v", deleteErr, cond, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}

	data := make([]metadata.ModuleHost, 0, len(input.HostID))
	for _, hostID := range input.HostID {
		data = append(data, metadata.ModuleHost{
			SetID:    module.SetID,
			ModuleID: module.ModuleID,
			HostID:   hostID,
			AppID:    module.BizID,
			TenantID: kit.TenantID,
		})
	}
	insertErr := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Insert(kit.Ctx, data)
	if insertErr != nil {
		blog.Errorf("insert host relation failed, err: %v, data: %v, rid: %v", insertErr, data, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return nil
}

func (manager *TransferManager) validTransferResourceDirParams(kit *rest.Kit,
	input *metadata.TransferHostResourceDirectory) (*metadata.ModuleInst, errors.CCErrorCoder) {

	biz := new(metadata.BizInst)
	filter := mapstr.MapStr{common.BKDefaultField: 1}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Find(filter).Fields(common.BKAppIDField).
		One(kit.Ctx, biz)
	if err != nil {
		blog.Errorf("get resource pool biz failed, err: %v, cond: %v, rid: %s", err, filter, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	// valid bk_module_id,资源池目录default=4,空闲机default=1
	cond := mapstr.MapStr{
		common.BKAppIDField:    biz.BizID,
		common.BKModuleIDField: input.ModuleID,
		common.BKDefaultField:  mapstr.MapStr{common.BKDBIN: []int{1, 4}},
	}
	module := new(metadata.ModuleInst)
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(cond).One(kit.Ctx, module)
	if err != nil {
		if mongodb.IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCoreServiceResourceDirectoryNotExistErr)
		}
		blog.Errorf("get module failed, err: %v, cond: %v, rid: %s", err, cond, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	// 确保主机在资源池目录下(default=1的业务)
	opt := mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: input.HostID},
		common.BKAppIDField: biz.BizID}
	existHostIDs, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameModuleHostConfig).Distinct(kit.Ctx,
		common.BKHostIDField, opt)
	if err != nil {
		blog.Errorf("get host ids in resource pool failed, err: %v, cond: %v, rid: %s", err, opt, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(existHostIDs) < len(input.HostID) {
		parsed, _ := util.SliceInterfaceToInt64(existHostIDs)
		wrongHostIDs := util.IntArrDiff(input.HostID, parsed)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotUnderAnyResourceDirectory, wrongHostIDs)
	}

	return module, nil
}
