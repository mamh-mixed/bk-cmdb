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
  <table class="audit-business-options">
    <colgroup>
      <col width="3%">
      <col width="25%">
      <col width="6%">
      <col width="25%">
      <col width="6%">
      <col width="25%">
      <col width="10%">
    </colgroup>
    <tr>
      <td align="right"><label class="option-label">{{$t('业务')}}</label></td>
      <td>
        <audit-business-selector class="option-value"
          searchable
          :placeholder="$t('请选择xx', { name: $t('业务') })"
          v-model="condition.bk_biz_id">
        </audit-business-selector>
      </td>
      <td align="right"><label class="option-label">{{$t('操作对象')}}</label></td>
      <td>
        <audit-target-selector class="option-value"
          searchable
          :placeholder="$t('请选择xx', { name: $t('操作对象') })"
          v-model="condition.resource_type">
        </audit-target-selector>
      </td>
      <td align="right"><label class="option-label">{{$t('动作')}}</label></td>
      <td>
        <audit-action-selector class="option-value"
          :target="condition.resource_type"
          :placeholder="$t('请选择xx', { name: $t('动作') })"
          :empty-text="$t('请先选择操作对象')"
          v-model="condition.action">
        </audit-action-selector>
      </td>
      <td></td>
    </tr>
    <tr>
      <td align="right"><label class="option-label">{{$t('时间')}}</label></td>
      <td>
        <cmdb-form-date-range class="option-value"
          font-size="medium"
          :placeholder="$t('请选择xx', { name: $t('时间') })"
          :clearable="false"
          v-model="condition.operation_time">
        </cmdb-form-date-range>
      </td>
      <td align="right"><label class="option-label">{{$t('账号')}}</label></td>
      <td><audit-user-selector class="option-value" v-model="condition.condition.user" /></td>
      <td align="right"><label class="option-label">{{$t('实例')}}</label></td>
      <td>
        <bk-input class="option-value"
          v-model.trim="instanceFilter"
          :placeholder="$t('请输入xx', { name: instanceType === 'resource_id' ? 'ID' : $t('名称') })">
          <bk-select class="option-type" slot="prepend"
            :clearable="false"
            v-model="instanceType">
            <bk-option id="resource_name" :name="$t('名称')"></bk-option>
            <bk-option id="resource_id" name="ID"></bk-option>
          </bk-select>
          <bk-checkbox class="option-exact" slot="append" size="small" v-model="condition.fuzzy_query">
            {{$t('模糊')}}
          </bk-checkbox>
        </bk-input>
      </td>
      <td>
        <div class="options-button">
          <bk-button class="mr10" theme="primary" @click="handleSearch(1)">{{$t('查询')}}</bk-button>
          <bk-button theme="default" @click="handleReset">{{$t('清空')}}</bk-button>
        </div>
      </td>
    </tr>
  </table>
</template>

<script>
  import AuditTargetSelector from './audit-target-selector'
  import AuditActionSelector from './audit-action-selector'
  import AuditUserSelector from './audit-user-selector'
  import AuditBusinessSelector from '@/components/audit-history/audit-business-selector'
  import RouterQuery from '@/router/query'
  import tools from '@/utils/tools'


  const today = tools.formatTime(new Date(), 'YYYY-MM-DD')
  const defaultCondition = () => ({
    bk_biz_id: '',
    resource_type: '',
    action: [],
    operation_time: [today, today],
    resource_id: '',
    resource_name: '',
    category: 'business',
    fuzzy_query: false,
    condition: {
      user: ['in', '']
    }
  })

  export default {
    name: 'audit-business-options',
    components: {
      AuditTargetSelector,
      AuditActionSelector,
      AuditUserSelector,
      AuditBusinessSelector
    },
    data() {
      return {
        instanceType: 'resource_name',
        condition: {
          ...defaultCondition()
        }
      }
    },
    computed: {
      instanceFilter: {
        get() {
          return this.condition[this.instanceType]
        },
        set(value) {
          this.condition[this.instanceType] = value
        }
      }
    },
    watch: {
      instanceType() {
        this.condition.resource_id = ''
        this.condition.resource_name = ''
      }
    },
    created() {
      this.handleSearch()
    },
    methods: {
      handleSearch(isEvent) {
        this.$emit('condition-change', this.condition)
        RouterQuery.set({
          tab: 'business',
          page: 1,
          _t: Date.now(),
          _e: isEvent
        })
      },
      handleReset() {
        this.condition = defaultCondition()
        this.handleSearch()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .audit-business-options {
        width: 100%;
        padding: 5px 0;
        tr {
            td {
                padding: 5px 0;
            }
        }
        .option-label {
            font-size: 14px;
            padding: 0 10px;
            @include ellipsis;
        }
        .option-value {
            width: 100%;
            min-width: 230px;
            max-width: 400px;
        }
        .option-type {
            width: 80px;
            margin-top: -1px;
            border-color: #c4c6cc transparent;
            box-shadow: none;
        }
        .option-exact {
          white-space: nowrap;
          padding: 0 4px;
          margin: 7px 0;
        }
        .options-button {
            display: flex;
            align-items: center;
            margin-left: 10px;
        }
    }
</style>
