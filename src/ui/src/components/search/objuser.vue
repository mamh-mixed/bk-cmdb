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
  <cmdb-form-objuser
    :exclude="false"
    v-model="localValue"
    v-bind="$attrs"
    :render-tag="renderTag"
    ref="objUserRef"
    @clear="() => $emit('clear')"
    @focus="handleToggle(true)"
    @blur="handleToggle(false)">
    <template #prepend>
      <slot name="prepend" />
    </template>
  </cmdb-form-objuser>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-objuser',
    mixins: [activeMixin],
    props: {
      value: {
        type: Array,
        default: () => ([])
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value ? this.value.join(',') : ''
        },
        set(value) {
          const values = value.split(',')
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    },
    methods: {
      renderTag(h, { _username, _index, user }) {
        const userSelector = this.$refs.objUserRef?.$refs?.userSelector
        return h('span', {
          class: ['user-selector-selected-value', { 'non-existent': !user.id }],
          directives: [
            {
              name: 'bkTooltips',
              value: {
                content: this.$t('该人员不存在'),
                disabled: user.id
              }
            }
          ],
        }, userSelector?.getDisplayText?.(user))
      }
    }
  }
</script>
