<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="module-selector-layout"
    v-bkloading="{ isLoading: loading }">
    <div class="wrapper">
      <cmdb-resize-layout class="topo-resize-layout"
        direction="right"
        :handler-offset="3"
        :min="480"
        :max="508">
        <div :class="['wrapper-column wrapper-left', { 'has-title': hasTitle }]">
          <h2 class="title" v-if="hasTitle">{{title}}</h2>
          <bk-select class="business-selector"
            v-model="targetBizId"
            :clearable="false"
            :searchable="true"
            :placeholder="$t('请选择xx', { name: $t('业务') })"
            @change="getModules">
            <cmdb-auth-option v-for="item in businessList"
              :key="item.bk_biz_id"
              :id="item.bk_biz_id"
              :name="`[${item.bk_biz_id}] ${item.bk_biz_name}`"
              :auth="{ type: authType, relation: item.relation }">
            </cmdb-auth-option>
          </bk-select>
          <template v-if="targetBizId">
            <bk-big-tree ref="tree" class="topology-tree"
              default-expand-all
              :options="{
                idKey: getNodeId,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
              }"
              :height="450"
              :node-height="36"
              :show-checkbox="false"
              @node-click="handleNodeClick"
              @check-change="handleNodeCheck">
              <template slot-scope="{ node, data }">
                <span :class="['node-checkbox fl', { 'is-checked': checked.includes(node) }]"
                  v-if="data.bk_obj_id === 'module'">
                </span>
                <i class="internal-node-icon fl"
                  v-if="data.default !== 0"
                  :class="getInternalNodeClass(node, data)">
                </i>
                <i v-else :class="['node-icon fl']">{{data.bk_obj_name[0]}}</i>
                <span class="node-name" :title="node.name">{{node.name}}</span>
              </template>
            </bk-big-tree>
          </template>
          <template v-else>
            <bk-exception class="business-tips" type="empty" scene="part">
              <span>{{$t('请先选择业务')}}</span>
            </bk-exception>
          </template>
        </div>
      </cmdb-resize-layout>
      <div class="wrapper-column wrapper-right">
        <module-checked-list :checked="checked" @delete="handleDeleteModule" @clear="handleClearModule" />
      </div>
    </div>
    <div class="layout-footer">
      <span v-bk-tooltips="{ content: $t('请先选择业务模块'), disabled: checked.length > 0 }">
        <bk-button class="mr10" theme="primary"
          :disabled="!checked.length"
          :loading="confirmLoading"
          @click="handleNextStep">
          {{$t('确定')}}
        </bk-button>
      </span>
      <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import { AuthRequestId, afterVerify } from '@/components/ui/auth/auth-queue.js'
  import ModuleCheckedList from './module-checked-list.vue'
  import { ONE_TO_ONE, MULTI_TO_ONE } from '@/dictionary/host-transfer-type.js'

  export default {
    name: 'cmdb-across-business-module-selector',
    components: {
      ModuleCheckedList
    },
    props: {
      confirmLoading: {
        type: Boolean,
        default: false
      },
      title: {
        type: String,
        default: ''
      },
      /**
       * 选择的机器所在业务，支持单个业务和多个业务。
       */
      business: {
        type: [Object, Array],
        required: true
      },
      /**
       * 跨业务转移类型，支持两种类型
       */
      type: {
        type: String,
        required: true,
        validator: type => [ONE_TO_ONE, MULTI_TO_ONE].includes(type)
      }
    },
    data() {
      return {
        loading: true,
        businessList: [],
        targetBizId: '',
        checked: [],
        request: {
          auth: AuthRequestId,
          internal: Symbol('internal'),
          list: Symbol('list')
        },
        nodeIconMap: {
          1: 'icon-cc-host-free-pool',
          2: 'icon-cc-host-breakdown',
          default: 'icon-cc-host-free-pool'
        }
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      bizIds() {
        if (this.business?.bk_biz_id) {
          return [this.business.bk_biz_id]
        }
        return this.business
      },
      targetBiz() {
        return this.businessList.find(biz => biz.bk_biz_id === this.targetBizId)
      },
      hasTitle() {
        return this.title && this.title.length
      },
      authType() {
        const { HOST_TRANSFER_ACROSS_BIZ, IDLE_HOST_TRANSFER_ACROSS_BIZ } = this.$OPERATION
        if (this.type === ONE_TO_ONE) {
          return HOST_TRANSFER_ACROSS_BIZ
        }
        if (this.type === MULTI_TO_ONE) {
          return IDLE_HOST_TRANSFER_ACROSS_BIZ
        }
        return ''
      },
    },
    async created() {
      this.getFullAmountBusiness()
    },
    methods: {
      generateRelation(businesses) {
        businesses.forEach((item) => {
          let relation = null

          if (this.type === ONE_TO_ONE) {
            relation = [[[this.bizIds[0]], [item.bk_biz_id]]]
          }

          if (this.type === MULTI_TO_ONE) {
            relation = []
            this.bizIds.forEach((bizId) => {
              relation.push([[bizId], [item.bk_biz_id]])
            })
          }

          if (!relation?.length) {
            throw Error('relation is required')
          }

          item.relation = relation
        })
      },
      async getFullAmountBusiness() {
        try {
          const data = await this.$http.get('biz/simplify?sort=bk_biz_id', { requestId: this.request.list })
          const availableBusiness = (data.info || []).filter(business => business.bk_biz_id !== this.bizId)
          this.generateRelation(availableBusiness)
          this.businessList = Object.freeze(availableBusiness)
        } catch (e) {
          console.error(e)
          this.businessList = []
        } finally {
          setTimeout(() => {
            afterVerify((authData) => {
              this.loading = false
              this.sortBusinessByAuth(authData)
            })
          }, 0)
        }
      },
      sortBusinessByAuth(authData) {
        if (!authData) {
          return false
        }
        const list = this.businessList?.map((item, index) => ({ ...item, is_pass: authData[index]?.is_pass }))
        list.sort((itemA, itemB) => itemB?.is_pass - itemA?.is_pass)
        this.businessList = list
      },
      async getModules() {
        try {
          this.checked = []
          const internalTop = await this.getInternalModules()
          this.$refs.tree.setData(internalTop)
        } catch (e) {
          this.$refs.tree.setData([])
          console.error(e)
        }
      },
      getInternalModules() {
        return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
          bizId: this.targetBizId,
          config: {
            requestId: this.request.internal
          }
        }).then(data => [{
          bk_inst_id: this.targetBizId,
          bk_inst_name: this.targetBiz.bk_biz_name,
          bk_obj_id: 'biz',
          bk_obj_name: this.getModelById('biz').bk_obj_name,
          default: 0,
          child: [{
            bk_inst_id: data.bk_set_id,
            bk_inst_name: data.bk_set_name,
            bk_obj_id: 'set',
            bk_obj_name: this.getModelById('set').bk_obj_name,
            default: 0,
            child: data.module.map(module => ({
              bk_inst_id: module.bk_module_id,
              bk_inst_name: module.bk_module_name,
              bk_obj_id: 'module',
              bk_obj_name: this.getModelById('module').bk_obj_name,
              default: module.default
            }))
          }]
        }])
      },
      getNodeId(data) {
        return `${data.bk_obj_id}-${data.bk_inst_id}`
      },
      getInternalNodeClass(node, data) {
        const clazz = []
        clazz.push(this.nodeIconMap[data.default] || this.nodeIconMap.default)
        return clazz
      },
      handleNodeClick(node) {
        if (node.data.bk_obj_id !== 'module') {
          return false
        }
        this.checked = [node]
      },
      handleNodeCheck(checked, currentNode) {
        const currentChecked = []
        const removeChecked = []
        checked.forEach((id) => {
          const node = this.$refs.tree.getNodeById(id)
          // 用户自定义空闲模块的 default 是一样的，所以增加 bk_inst_id 是否一致的判断
          if (node.data.default === currentNode.data.default && node.data.bk_inst_id === currentNode.data.bk_inst_id) {
            currentChecked.push(id)
          } else {
            removeChecked.push(id)
          }
        })
        this.checked = currentChecked.map(id => this.$refs.tree.getNodeById(id))
        this.$refs.tree.setChecked(removeChecked, { checked: false })
      },
      handleDeleteModule(node) {
        this.checked = this.checked.filter(exist => exist !== node)
        this.$refs.tree.setChecked(node.id, { checked: false })
      },
      handleClearModule() {
        this.$refs.tree.setChecked(this.checked.map(node => node.id), { checked: false })
        this.checked = []
      },
      handleCancel() {
        this.$emit('cancel')
      },
      handleNextStep() {
        this.$emit('confirm', this.checked, this.targetBizId)
      }
    }
  }
</script>

<style lang="scss" scoped>
    $leftPadding: 0 12px 0 23px;
    .module-selector-layout {
        height: var(--height, 600px);
        min-height: 300px;
        padding: 0 0 50px;
        position: relative;
        .layout-footer {
            position: sticky;
            bottom: 0;
            left: 0;
            width: 100%;
            height: 50px;
            padding: 8px 20px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
            text-align: right;
            font-size: 0;
            z-index: 100;
            .footer-tips {
                display: inline-block;
                vertical-align: middle;
            }
        }
    }
    .wrapper {
        display: flex;
        height: 100%;
        .topo-resize-layout {
            width: 508px;
            height: 100%;
            border-right: 1px solid #dcdee5;
        }
        .wrapper-column {
            flex: 1;
        }
        .wrapper-left {
            height: calc(100% - 24px);
            margin-top: 24px;
            .title {
                padding: 0 12px 24px 23px;
                font-size: 16px;
                font-weight: normal;
                color: #313238;
                line-height:22px;
            }
            .business-selector {
                display: block;
                margin: 0 12px 0 23px;
            }
            .business-tips {
                height: calc(100% - 90px);
                display: flex;
                align-items: center;
                justify-content: center;
                font-size: 12px;
                color: $textColor;
            }

            &.has-title {
                height: calc(100% - 18px);
                margin-top: 18px;
                .topology-tree {
                    max-height: calc(100% - 102px);
                }
            }
        }
        .wrapper-right {
            padding: 12px 23px 0;
            background: #f5f6fa;
            overflow: hidden;
        }
    }
    .topology-tree {
        width: 100%;
        max-height: calc(100% - 32px - 24px);
        padding: 0 0 0 12px;
        margin: 12px 0;
        @include scrollbar;
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
            &.is-template {
                background-color: #97AED6;
            }
            &.is-selected {
                background-color: #3A84FF;
            }
        }
        .node-name {
            height: 36px;
            line-height: 36px;
            overflow: hidden;
            @include ellipsis;
        }
        .node-checkbox {
            width: 16px;
            height: 16px;
            margin: 10px 5px 0 10px;
            background: #FFF;
            border-radius: 50%;
            border: 1px solid #979BA5;
            &.is-checked {
                padding: 3px;
                border-color: $primaryColor;
                background-color: $primaryColor;
                background-clip: content-box;
            }
        }
        .internal-node-icon{
            width: 20px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            margin: 8px 4px 8px 0;
        }
    }
</style>
