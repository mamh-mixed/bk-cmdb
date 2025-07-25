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
  <bk-sideslider
    :transfer="true"
    :width="800"
    :title="$t('批量编辑')"
    :is-show.sync="isShow"
    :before-close="beforeClose"
    @hidden="handleHidden">
    <cmdb-form-multiple slot="content"
      ref="form"
      v-bkloading="{ isLoading: pending }"
      :properties="properties"
      :property-groups="propertyGroups"
      :uneditable-properties="bindedProperties"
      @on-submit="handleSaveProcess"
      @on-cancel="beforeClose">
    </cmdb-form-multiple>
  </bk-sideslider>
</template>

<script>
  import { mapGetters } from 'vuex'
  import {
    processPropertyRequestId,
    processPropertyGroupsRequestId
  } from './symbol'
  export default {
    props: {
      serviceTemplateId: Number,
      processTemplateId: Number,
      submitHandler: Function
    },
    data() {
      return {
        isShow: false,
        properties: [],
        propertyGroups: [],
        // 绑定信息暂不支持批量编辑
        bindedProperties: ['bind_info'],
        pending: true
      }
    },
    computed: {
      ...mapGetters(['supplierAccount'])
    },
    async created() {
      try {
        const request = [
          this.getProperties(),
          this.getPropertyGroups()
        ]
        if (this.processTemplateId) {
          request.push(this.getProcessTemplate())
        }
        await Promise.all(request)
      } catch (error) {
        console.error(error)
      } finally {
        this.pending = false
      }
    },
    methods: {
      show() {
        this.isShow = true
      },
      async getProperties() {
        try {
          this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
            params: {
              bk_obj_id: 'process',
              bk_supplier_account: this.supplierAccount
            },
            config: {
              requestId: processPropertyRequestId,
              fromCache: true
            }
          })
        } catch (error) {
          console.error(error)
          this.properties = []
        }
      },
      async getPropertyGroups() {
        try {
          this.propertyGroups = await this.$store.dispatch('objectModelFieldGroup/searchGroup', {
            objId: 'process',
            params: {},
            config: {
              requestId: processPropertyGroupsRequestId,
              fromCache: true
            }
          })
        } catch (error) {
          console.error(error)
          this.propertyGroups = []
        }
      },
      async getProcessTemplate() {
        try {
          const { property } = await this.$store.dispatch('processTemplate/getProcessTemplate', {
            params: {
              processTemplateId: this.processTemplateId
            },
            config: {
              cancelPrevious: true
            }
          })
          const bindedProperties = []
          Object.keys(property).forEach((key) => {
            if (property[key].as_default_value) {
              bindedProperties.push(key)
            }
          })
          this.bindedProperties = bindedProperties
        } catch (error) {
          console.error(error)
        }
      },
      handleHidden() {
        this.$emit('close')
      },
      async handleSaveProcess(values, changedValues, instance) {
        try {
          this.pending = true
          await this.submitHandler(values, changedValues, instance)
          this.isShow = false
        } catch (error) {
          console.error(error)
        } finally {
          this.pending = false
        }
      },
      beforeClose() {
        if (this.$refs.form.hasChange) {
          this.$refs.form.setChanged(true)
          return this.$refs.form.beforeClose(() => {
            this.isShow = false
          })
        }
        this.isShow = false
        return Promise.resolve(true)
      }
    }
  }
</script>
