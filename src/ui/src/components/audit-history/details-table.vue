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
  <div class="details-layout">
    <div class="details-info">
      <div class="info-group" v-if="details.bk_biz_id">
        <label class="info-label">{{$t('所属业务')}}</label>
        <span class="info-content" v-bk-overflow-tips>
          <audit-business-selector type="info" :value="details.bk_biz_id"></audit-business-selector>
        </span>
      </div>
      <div class="info-group">
        <label class="info-label">{{$t('操作对象')}}</label>
        <span class="info-content" v-bk-overflow-tips>
          {{resourceType ? resourceType.name : details.resource_type}}
        </span>
      </div>
      <div class="info-group">
        <label class="info-label">{{$t('动作')}}</label>
        <span class="info-content" v-bk-overflow-tips>{{action ? action.name : details.action}}</span>
      </div>
      <div class="info-group" v-if="modelId">
        <label class="info-label">{{$t('模型类型')}}</label>
        <span class="info-content" v-bk-overflow-tips>{{modelName}}</span>
      </div>
      <div class="info-group">
        <label class="info-label">{{$t('操作实例')}}</label>
        <span class="info-content" v-bk-overflow-tips>{{details.resource_name}}</span>
      </div>
      <div class="info-group">
        <label class="info-label">{{$t('操作描述')}}</label>
        <span class="info-content" v-bk-overflow-tips>{{description}}</span>
      </div>
      <div class="info-group">
        <label class="info-label">{{$t('操作时间')}}</label>
        <span class="info-content" v-bk-overflow-tips>{{$tools.formatTime(details.operation_time)}}</span>
      </div>
      <div class="info-group">
        <label class="info-label">{{$t('操作账号')}}</label>
        <span class="info-content" v-bk-overflow-tips>
          <cmdb-form-objuser type="info" :value="details.user"></cmdb-form-objuser>
        </span>
      </div>
    </div>
    <div class="details-table">
      <bk-table
        row-border
        col-border
        :data="tableList"
        :max-height="$APP.height - 300"
        :cell-style="getCellStyle">
        <bk-table-column :label="$t('属性')" prop="field" align="right"></bk-table-column>
        <bk-table-column :label="$t('变更前')" prop="before" show-overflow-tooltip
          v-if="showBefore">
          <template slot-scope="{ row }">
            <cmdb-property-value
              v-if="row.type === 'property'"
              :property="row.property"
              :value="row.before">
            </cmdb-property-value>
            <div class="details-custom-content" v-else v-html="row.before"></div>
          </template>
        </bk-table-column>
        <bk-table-column :label="$t('变更后')" prop="after" show-overflow-tooltip
          v-if="showAfter">
          <template slot-scope="{ row }">
            <cmdb-property-value
              v-if="row.type === 'property'"
              :property="row.property"
              :value="row.after">
            </cmdb-property-value>
            <div class="details-custom-content" v-else v-html="row.after"></div>
          </template>
        </bk-table-column>
      </bk-table>
      <div class="details-toggle" v-if="showToggle">
        <bk-link
          theme="primary"
          @click="toggleDetails">
          {{ isShowAllFields ? $t('收起') : $t('展开') }}
        </bk-link>
      </div>
    </div>
  </div>
</template>

<script>
  import AuditBusinessSelector from '@/components/audit-history/audit-business-selector'
  export default {
    name: 'details-table',
    components: {
      AuditBusinessSelector
    },
    props: {
      details: Object
    },
    data() {
      return {
        list: [],
        diffList: [],
        dictionary: [],
        properties: [],
        isShowAllFields: false
      }
    },
    computed: {
      showBefore() {
        return !['create'].includes(this.details.action)
      },
      showAfter() {
        return !['delete'].includes(this.details.action)
      },
      resourceType() {
        return this.dictionary.find(type => type.id === this.details.resource_type)
      },
      action() {
        if (!this.resourceType) {
          return null
        }
        return this.resourceType.operations.find(action => action.id === this.details.action)
      },
      description() {
        const actionName = this.action ? this.action.name : this.details.action
        return `${actionName}${this.details.resource_name}`
      },
      modelId() {
        return this.details?.operation_detail?.bk_obj_id
      },
      modelName() {
        const model = this.$store.getters['objectModelClassify/getModelById'](this.modelId)
        return model.bk_obj_name
      },
      tableList() {
        if (this.showToggle) {
          return this.isShowAllFields ? this.list : this.diffList
        }
        return this.list
      },
      showToggle() {
        return this.diffList.length && this.diffList.length !== this.list.length
      }
    },
    async created() {
      await Promise.all([
        this.getAuditDictionary(),
        this.getModelProperty()
      ])
      this.setList()
    },
    methods: {
      async getAuditDictionary() {
        try {
          this.dictionary = await this.$store.dispatch('audit/getDictionary', {
            fromCache: true
          })
        } catch (error) {
          this.dictionary = []
        }
      },
      async getModelProperty() {
        try {
          if (!this.modelId) {
            return false
          }
          const params = { bk_obj_id: this.modelId }
          if (this.details.bk_biz_id) {
            params.bk_biz_id = this.details.bk_biz_id
          }
          this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', { params })
        } catch (error) {
          this.properties = []
          console.error(error)
        }
      },
      setList() {
        const isTopologyPathChange = ['unassign_host', 'assign_host', 'transfer_host_module'].includes(this.details.action)
        if (isTopologyPathChange) {
          this.setTopologyPathList()
        } else {
          this.setNormalList()
        }
      },
      setTopologyPathList() {
        this.list = [{
          field: this.$t('关联关系'),
          type: 'topology',
          before: this.getTopoPath(this.details?.operation_detail?.pre_data),
          after: this.getTopoPath(this.details?.operation_detail?.cur_data)
        }]
        this.diffList = [...this.list]
      },
      getTopoPath(data) {
        const paths = []
        data.set.forEach((set) => {
          const path = [data.bk_biz_name, set.bk_set_name]
          set.module.forEach((module) => {
            paths.push([...path, module.bk_module_name].join('→'))
          })
        })
        return paths.join('<br>')
      },
      setNormalList() {
        const operationDetails = this.details?.operation_detail?.details || {}
        const before = operationDetails.pre_data || {}
        const update = operationDetails.update_fields || {}
        const after = Object.assign(operationDetails.cur_data || {}, before, update)
        this.list = this.properties.map((property) => {
          const field = property.bk_property_id
          return {
            field: property.bk_property_name,
            type: 'property',
            property,
            before: before[field],
            after: after[field]
          }
        })
        this.diffList = this.list.filter(row => row.before !== row.after)
      },
      toggleDetails() {
        this.isShowAllFields = !this.isShowAllFields
      },
      getCellStyle({ row, columnIndex }) {
        if (!this.showBefore || !this.showAfter) {
          return {}
        }
        if (columnIndex > 0 && row.before !== row.after) {
          return {
            backgroundColor: '#e9faf0'
          }
        }
        return {}
      }
    }
  }
</script>

<style lang="scss" scoped>
    .details-layout{
        height: 100%;
    }
    .details-info {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        .info-group {
            display: flex;
            width: 50%;
            padding-right: 25px;
            font-size: 14px;
            line-height: 34px;
            .info-label {
                display: inline-block;
                width: 100px;
                padding-right: 10px;
                color: #63656e;
                &:after {
                    content: ":";
                    padding: 0 2px;
                }
            }
            .info-content {
                display: inline-block;
                width: calc(100% - 100px);
                color: #313237;
                @include ellipsis;
            }
        }
    }
    .details-custom-content {
        line-height: 24px;
        overflow: auto;
    }
    .details-toggle {
        text-align: right;
    }

    [lang="en"] {
      .details-info .info-group .info-label {
        width: 160px;
      }
    }
</style>
