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
  <div :class="['confirm-wrapper', { 'has-invalid': hasInvalid }]">
    <h2 class="title">{{$t('确认归还主机池')}}</h2>
    <i18n tag="p" path="确认归还主机池忽略主机数量" class="content" v-if="hasInvalid">
      <template #count><span class="count">{{count}}</span></template>
      <template #invalid><span class="invalid">{{invalidList.length}}</span></template>
      <template #idleSet><span>{{idleSet}}</span></template>
    </i18n>
    <i18n tag="p" path="确认归还主机池主机数量" class="content" v-else>
      <template #count><span class="count">{{count}}</span></template>
      <template #idleSet><span>{{idleSet}}</span></template>
    </i18n>
    <p class="content">{{$t('归还主机池提示')}}</p>
    <div class="directory">
      {{$t('归还至目录')}}
      <bk-select class="directory-selector ml10"
        v-model="target"
        searchable
        :clearable="false"
        :loading="$loading(request.list)"
        :popover-options="{
          boundary: 'window'
        }">
        <cmdb-auth-option v-for="directory in directories"
          :key="directory.bk_module_id"
          :id="directory.bk_module_id"
          :name="directory.bk_module_name"
          :auth="{ type: $OPERATION.HOST_TO_RESOURCE, relation: [[[bizId], [directory.bk_module_id]]] }">
        </cmdb-auth-option>
      </bk-select>
    </div>
    <invalid-list :title="$t('以下主机不能移除')" :list="invalidList"></invalid-list>
    <div class="options">
      <bk-button class="mr10" theme="primary"
        :disabled="!target"
        @click="handleConfirm">{{$t('确定')}}</bk-button>
      <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
    </div>
  </div>
</template>

<script>
  import InvalidList from './invalid-list'
  export default {
    name: 'cmdb-move-to-resource-confirm',
    components: {
      InvalidList
    },
    props: {
      count: {
        type: Number,
        default: 0
      },
      bizId: Number,
      invalidList: {
        type: Array,
        default: () => ([])
      }
    },
    data() {
      return {
        target: '',
        directories: [],
        request: {
          list: Symbol('list')
        }
      }
    },
    computed: {
      hasInvalid() {
        return !!this.invalidList.length
      },
      idleSet() {
        return this.$store.state.globalConfig.config.set
      }
    },
    created() {
      this.getDirectories()
    },
    methods: {
      async getDirectories() {
        try {
          const { info } = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
            params: {
              page: {
                sort: 'bk_module_name'
              }
            },
            config: {
              requestId: this.request.list
            }
          })
          this.directories = info
        } catch (error) {
          console.error(error)
        }
      },
      handleConfirm() {
        this.$emit('confirm', this.target)
      },
      handleCancel() {
        this.$emit('cancel')
      }
    }
  }
</script>

<style lang="scss" scoped>
    .confirm-wrapper {
        text-align: center;
        &.has-invalid {
            .content {
                padding: 0 26px;
                text-align: left;
            }
            .directory {
                padding: 0 0 0 26px;
                justify-content: flex-start;
                .directory-selector {
                    width: 514px;
                }
            }
        }
    }
    .title {
        margin: 45px 0 17px;
        line-height: 32px;
        font-size:24px;
        font-weight: normal;
        color: #313238;
    }
    .content {
        line-height:20px;
        font-size:14px;
        color: $textColor;
        .count {
            font-weight: bold;
            color: $successColor;
            padding: 0 4px;
        }
        .invalid {
            font-weight: bold;
            color: $dangerColor;
            padding: 0 4px;
        }
    }
    .directory {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        margin-top: 10px;
        .directory-selector {
            width: 305px;
            margin-left: 10px;
            text-align: left;
        }
    }
    .options {
        margin: 20px 0;
        font-size: 0;
    }
</style>
