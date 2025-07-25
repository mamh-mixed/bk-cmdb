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
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
    </div>
    <div class="result-desc" @click="data.linkTo(data.source)">
      <div class="desc-item hl" v-html="`${$t('唯一标识')}：${getHighlightValue(data.source.bk_obj_id, data)}`"></div>
      <div class="desc-item hl" v-html="`${$t('模型名称')}：${getHighlightValue(data.source.bk_obj_name, data)}`"></div>
      <div class="desc-item">{{$t('所属模型分组')}}：{{classificationName}}</div>
      <dl class="model-group-list">
        <div class="group" v-for="(item, index) in groupedProperties" :key="index">
          <dt class="group-name">{{item.group.bk_group_name}}</dt>
          <dd class="property-list">
            <div class="property-item" v-for="(property, childIndex) in item.properties" :key="childIndex">
              （<span class="hl" v-html="`${$t('字段名称')}：${getHighlightValue(property.bk_property_name, data)}`"></span>
              <span class="hl" v-html="`${$t('唯一标识')}：${getHighlightValue(property.bk_property_id, data)}`"></span>
              <span>{{$t('字段类型')}}：{{fieldTypeMap[property.bk_property_type]}}</span>）
            </div>
          </dd>
        </div>
      </dl>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed } from 'vue'
  import store from '@/store'
  import { getText, getHighlightValue } from './use-item.js'
  import { PROPERTY_TYPE_NAMES } from '@/dictionary/property-constants'

  export default defineComponent({
    name: 'item-model',
    props: {
      data: {
        type: Object,
        default: () => ({})
      },
      propertyMap: {
        type: Object,
        default: () => ({})
      },
      propertyGroupMap: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props) {
      const { data, propertyMap, propertyGroupMap } = toRefs(props)

      const objId = computed(() => data.value.source.bk_obj_id)
      const properties = computed(() => propertyMap.value[objId.value])

      const classificationName = computed(() => {
        const classifications = store.getters['objectModelClassify/classifications']
        const id = data.value.source.bk_classification_id
        return (classifications.find(item => item.bk_classification_id === id) || {}).bk_classification_name
      })

      const propertyGroups = computed(() => propertyGroupMap.value[objId.value])

      // eslint-disable-next-line max-len
      const sortProperties = computed(() => (properties.value || []).sort((propertyA, propertyB) => propertyA.bk_property_index - propertyB.bk_property_index))

      // eslint-disable-next-line max-len
      const sortedPropertyGroups = computed(() => propertyGroups.value.sort((groupA, groupB) => groupA.bk_group_index - groupB.bk_group_index))

      const groupedProperties = computed(() => sortedPropertyGroups.value.map(group => ({
        group,
        properties: sortProperties.value.filter((property) => {
          if (['default', 'none'].includes(property.bk_property_group) && group.bk_group_id === 'default') {
            return true
          }
          return property.bk_property_group === group.bk_group_id
        })
      })))

      const fieldTypeMap = PROPERTY_TYPE_NAMES

      return {
        properties,
        groupedProperties,
        fieldTypeMap,
        classificationName,
        getText,
        getHighlightValue
      }
    }
  })
</script>

<style lang="scss" scoped>
  .model-group-list {
    display: flex;
    flex-wrap: wrap;
    .group {
      display: flex;
      margin-bottom: 6px;

      .group-name {
        flex: none;
        &::after {
          content: '：';
        }
      }

      .property-list {
        display: flex;
        flex-wrap: wrap;

        .property-item {
          flex: none;
          margin-right: 8px;
        }
      }
    }
  }
</style>
