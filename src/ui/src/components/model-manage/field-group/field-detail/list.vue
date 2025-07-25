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
  <div class="form-list-layout">
    <div class="toolbar">
      <p class="title">{{$t('列表值')}}</p>
      <i
        v-bk-tooltips.top-start="$t('通过列表项的值按照0-9，a-z排序')"
        :class="['sort-icon', `icon-cc-sort-${order > 0 ? 'up' : 'down'}`]"
        @click="handleSort">
      </i>
    </div>
    <vue-draggable
      class="form-list-wrapper"
      tag="ul"
      v-model="list"
      :options="dragOptions"
      @end="handleDragEnd">
      <li class="form-item clearfix" v-for="(item, index) in list" :key="index">
        <div class="list-name">
          <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`name${index}`) }">
            <bk-input type="text"
              class="cmdb-form-input"
              :placeholder="$t('请输入值')"
              v-model.trim="item.name"
              v-validate="`required|enumName|repeat:${getOtherName(index)}`"
              @input="handleInput"
              :disabled="isReadOnly"
              :name="`name${index}`"
              :ref="`name${index}`">
            </bk-input>
            <p class="form-error">{{errors.first(`name${index}`)}}</p>
          </div>
        </div>
        <bk-button text class="list-btn" @click="deleteList(index)" :disabled="list.length === 1 || isReadOnly">
          <i class="bk-icon icon-minus-circle-shape"></i>
        </bk-button>
        <bk-button text class="list-btn" @click="addList(index)"
          :disabled="isReadOnly" v-if="index === list.length - 1">
          <i class="bk-icon icon-plus-circle-shape"></i>
        </bk-button>
      </li>
    </vue-draggable>
    <div class="default-setting">
      <p class="title mb10">{{$t('默认值设置')}}</p>
      <div class="cmdb-form-item" :class="{ 'is-error': errors.has('defaultValueSelect') }">
        <bk-select style="width: 100%;"
          :clearable="true"
          :searchable="true"
          :disabled="isReadOnly"
          name="defaultValueSelect"
          data-vv-validate-on="change"
          v-model="defaultListValue"
          @change="handleSettingDefault">
          <bk-option v-for="option in settingList"
            :key="option.id"
            :id="option.id"
            :name="option.name">
          </bk-option>
        </bk-select>
        <p class="form-error">{{errors.first('defaultValueSelect')}}</p>
      </div>
    </div>
  </div>
</template>

<script>
  import vueDraggable from 'vuedraggable'
  import { v4 as uuidv4 } from 'uuid'
  export default {
    components: {
      vueDraggable
    },
    props: {
      value: {
        type: [Array, String],
        default: ''
      },
      isReadOnly: {
        type: Boolean,
        default: false
      },
      defaultValue: {
        type: String,
        default: ''
      }
    },
    data() {
      return {
        list: [{ name: '' }],
        settingList: [],
        defaultListValue: '',
        dragOptions: {
          animation: 300,
          disabled: false,
          filter: '.list-btn, .list-name',
          preventOnFilter: false,
          ghostClass: 'ghost'
        },
        order: 1
      }
    },
    watch: {
      value() {
        this.initValue()
      },
      list: {
        deep: true,
        handler(value) {
          this.settingList = (value || []).filter(val => val.name).map(item => ({
            id: uuidv4(),
            name: item.name
          }))
          if (this.defaultValue) {
            this.defaultListValue = this.settingList.find(item => item.name === this.defaultValue).id
          }
        }
      }
    },
    created() {
      this.initValue()
    },
    methods: {
      getOtherName(index) {
        const nameList = []
        this.list.forEach((item, _index) => {
          if (index !== _index) {
            nameList.push(item.name)
          }
        })
        return nameList.join(',')
      },
      initValue() {
        if (this.value === '') {
          this.list = [{ name: '' }]
        } else {
          this.list = this.value.map(name => ({ name }))
        }
      },
      handleInput() {
        this.$nextTick(async () => {
          const res = await this.$validator.validateAll()
          if (res) {
            const list = this.list.map(item => item.name)
            this.$emit('input', list)
          }
        })
      },
      addList(index) {
        this.list.push({ name: '' })
        this.$nextTick(() => {
          this.$refs[`name${index + 1}`] && this.$refs[`name${index + 1}`][0].focus()
        })
      },
      deleteList(index) {
        this.list.splice(index, 1)
        this.defaultListValue = ''
        this.handleInput()
      },
      validate() {
        return this.$validator.validateAll()
      },
      handleDragEnd() {
        const list = this.list.map(item => item.name)
        this.$emit('input', list)
      },
      handleSort() {
        this.order = this.order * -1
        this.list.sort((A, B) => A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' }) * this.order)

        const list = this.list.map(item => item.name)
        this.$emit('input', list)
      },
      handleSettingDefault(id) {
        const defaultValue =  id ? this.settingList.find(item => item.id === id).name : ''
        this.$emit('update:defaultValue', defaultValue)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .title {
        font-size: 14px;
    }
    .form-list-wrapper {
        .form-item {
            display: flex;
            align-items: center;
            position: relative;
            margin-bottom: 16px;
            padding: 2px 2px 2px 28px;
            font-size: 0;
            cursor: move;

            &:not(:first-child) {
                margin-top: 16px;
            }

            &::before {
                content: '';
                position: absolute;
                top: 12px;
                left: 8px;
                width: 3px;
                height: 3px;
                border-radius: 50%;
                background-color: #979ba5;
                box-shadow: 0 5px 0 0 #979ba5,
                    0 10px 0 0 #979ba5,
                    5px 0 0 0 #979ba5,
                    5px 5px 0 0 #979ba5,
                    5px 10px 0 0 #979ba5;
            }

            .list-name {
                float: left;
                width: 200px;
                input {
                    width: 100%;
                }
            }
            .list-btn {
                font-size: 0;
                color: #c4c6cc;
                margin: -2px 0 0 6px;
                .bk-icon {
                    width: 18px;
                    height: 18px;
                    line-height: 18px;
                    font-size: 18px;
                    text-align: center;
                }

                &.is-disabled {
                    color: #dcdee5;
                }
                &:not(.is-disabled):hover {
                    color: #979ba5;
                }
            }
        }
    }

    .toolbar {
        display: flex;
        margin-bottom: 10px;
        align-items: center;
        line-height: 20px;

        .sort-icon {
            width: 20px;
            height: 20px;
            margin-left: 10px;
            border: 1px solid #c4c6cc;
            background: #fff;
            border-radius: 2px;
            font-size: 16px;
            line-height: 18px;
            text-align: center;
            color: #c4c6cc;
            cursor: pointer;

            &:hover {
                color: #979ba5;
            }
        }
    }

    .ghost {
        border: 1px dashed $cmdbBorderFocusColor;
    }
</style>
