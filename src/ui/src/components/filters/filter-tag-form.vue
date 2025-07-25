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
  <bk-form form-type="vertical">
    <bk-form-item>
      <label class="form-label">
        {{property.bk_property_name}}
        <span class="form-label-suffix">({{labelSuffix}})</span>
      </label>
      <div class="form-wrapper">
        <operator-selector class="form-operator"
          v-if="!withoutOperator.includes(property.bk_property_type)"
          :property="property"
          :custom-type-map="customOperatorTypeMap"
          :symbol-map="operatorSymbolMap"
          :desc-map="operatorDescMap"
          v-model="operator"
          @change="handleOperatorChange"
          @toggle="handleActiveChange">
        </operator-selector>
        <div class="form-value">
          <component class="form-el"
            :is="getComponentType()"
            :placeholder="getPlaceholder()"
            v-bind="getBindProps()"
            v-model.trim="value"
            @active-change="handleActiveChange">
          </component>
        </div>
      </div>
      <div class="form-options">
        <bk-button class="mr10" text @click="handleConfirm">{{$t('确定')}}</bk-button>
        <bk-button class="mr10" text @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </bk-form-item>
  </bk-form>
</template>

<script>
  import OperatorSelector from './operator-selector'
  import FilterStore from './store'
  import Utils from './utils'
  import { mapGetters } from 'vuex'
  import { isContainerObject } from '@/service/container/common'
  import { QUERY_OPERATOR, QUERY_OPERATOR_HOST_SYMBOL, QUERY_OPERATOR_HOST_DESC } from '@/utils/query-builder-operator'

  export default {
    components: {
      OperatorSelector
    },
    props: {
      property: {
        type: Object,
        required: true
      }
    },
    data() {
      const { IN, NIN, LIKE, CONTAINS, EQ, NE, GTE, LTE, RANGE } = QUERY_OPERATOR
      return {
        withoutOperator: ['date', 'time', 'bool', 'service-template'],
        localOperator: null,
        localValue: null,
        active: false,
        customOperatorTypeMap: {
          float: [EQ, NE, GTE, LTE, RANGE, IN],
          int: [EQ, NE, GTE, LTE, RANGE, IN],
          longchar: [IN, NIN, CONTAINS, LIKE],
          singlechar: [IN, NIN, CONTAINS, LIKE],
          array: [IN, NIN, CONTAINS, LIKE],
          object: [IN, NIN, CONTAINS, LIKE]
        },
        operatorSymbolMap: QUERY_OPERATOR_HOST_SYMBOL,
        operatorDescMap: QUERY_OPERATOR_HOST_DESC
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      labelSuffix() {
        const model = this.getModelById(this.property.bk_obj_id)
        return model ? model.bk_obj_name : this.property.bk_obj_id
      },
      operator: {
        get() {
          return this.localOperator || FilterStore.condition[this.property.id].operator
        },
        set(operator) {
          this.localOperator = operator
        }
      },
      value: {
        get() {
          if (this.localValue === null) {
            return FilterStore.condition[this.property.id].value
          }
          return this.localValue
        },
        set(value) {
          this.localValue = value
        }
      }
    },
    methods: {
      getPlaceholder() {
        return Utils.getPlaceholder(this.property)
      },
      getComponentType() {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = this.property
        const normal = `cmdb-search-${propertyType}`

        if (modelId === 'biz' && propertyId === 'bk_biz_name' && ![QUERY_OPERATOR.CONTAINS, QUERY_OPERATOR.LIKE].includes(this.operator)) {
          return `cmdb-search-${modelId}`
        }

        if (!FilterStore.bizId) {
          return normal
        }

        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'

        if ((isSetName || isModuleName) && ![QUERY_OPERATOR.CONTAINS, QUERY_OPERATOR.LIKE].includes(this.operator)) {
          return `cmdb-search-${modelId}`
        }

        // 数字类型int 和 float支持in操作符
        if (Utils.numberUseIn(this.property, this.operator)) {
          return 'cmdb-search-singlechar'
        }

        return normal
      },
      getBindProps() {
        const props = Utils.getBindProps(this.property)
        if (!FilterStore.bizId) {
          return props
        }
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = this.property
        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'
        if (isSetName || isModuleName) {
          return Object.assign(props, { bizId: FilterStore.bizId })
        }

        // 容器对象标签属性，需要注入标签kv数据作为选项
        if (isContainerObject(modelId) && propertyType === 'map') {
          return Object.assign(props, { options: FilterStore.containerPropertyMapValue?.[modelId]?.[propertyId] })
        }

        // 数字类型int 和 float支持in操作符
        if (Utils.numberUseIn(this.property, this.operator)) {
          props.onlyNumber = true
          props.fuzzy = false
        }

        return props
      },
      resetCondition() {
        this.operator = null
        this.value = null
      },
      handleOperatorChange(operator) {
        this.value = Utils.getOperatorSideEffect(this.property, operator, this.value)
      },
      // 当失去焦点时，激活状态的改变做一个延时，避免点击表单外部时直接隐藏了表单对应的tooltips
      handleActiveChange(active) {
        this.timer && clearTimeout(this.timer)
        if (active) {
          this.active = active
        } else {
          this.timer = setTimeout(() => {
            this.active = active
          }, 100)
        }
      },
      handleConfirm() {
        FilterStore.resetPage(true)
        FilterStore.updateCondition(this.property, this.operator, this.value)
        this.$emit('confirm')
      },
      handleCancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-label {
        display: block;
        font-size: 14px;
        font-weight: 400;
        line-height: 32px;
        @include ellipsis;
        .form-label-suffix {
            font-size: 12px;
            color: #979ba5;
        }
    }
    .form-wrapper {
        width: 380px;
        display: flex;
        .form-operator {
            flex: 128px 0 0;
            margin-right: 8px;
            align-self: baseline;
        }
        .form-value {
            flex: 1;
            position: relative;
        }
        .form-el {
          width: 100%;
          max-width: 100%;
        }
    }
    .form-options {
        display: flex;
        height: 32px;
        align-items: center;
        justify-content: flex-end;
    }
</style>
