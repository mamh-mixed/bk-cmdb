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

<script>
  import PropertySelector from './property-selector.vue'
  export default {
    extends: PropertySelector,
    props: {
      objId: {
        type: String
      },
      properties: {
        type: Array,
        default: () => ([])
      },
      propertyGroups: {
        type: Array,
        default: () => ([])
      },
      propertySelected: {
        type: Array,
        default: () => ([])
      },
      handler: {
        type: Function,
        default: () => {}
      }
    },
    data() {
      return {
        selected: [...this.propertySelected],
      }
    },
    computed: {
      propertyMap() {
        const modelPropertyMap = { [this.objId]: this.properties }
        const ignoreProperties = [] // 预留，需要忽略的属性
        // eslint-disable-next-line max-len
        modelPropertyMap[this.objId] = modelPropertyMap[this.objId].filter(property => !ignoreProperties.includes(property.bk_property_id))
        return modelPropertyMap
      }
    },
    methods: {
      async confirm() {
        this.handler(this.selected)
        this.close()
      }
    }
  }
</script>
