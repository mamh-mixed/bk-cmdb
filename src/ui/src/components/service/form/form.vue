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
    :title="internalTitle"
    :is-show.sync="isShow"
    :before-close="beforeClose"
    @hidden="handleHidden">
    <template slot="content">
      <cmdb-details v-if="internalType === 'view'"
        :properties="properties"
        :property-groups="propertyGroups"
        :inst="instance"
        :show-options="showOptions"
        :show-delete="false"
        :edit-auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
        :invisible-name-properties="invisibleNameProperties"
        :flex-properties="flexProperties"
        @on-edit="handleChangeInternalType">
      </cmdb-details>
      <cmdb-form v-else
        ref="form"
        v-bkloading="{ isLoading: pending }"
        :type="internalType"
        :inst="instance"
        :properties="visibleProperties"
        :property-groups="propertyGroups"
        :disabled-properties="bindedProperties"
        :invisible-name-properties="invisibleNameProperties"
        :flex-properties="flexProperties"
        :render-append="renderAppend"
        :custom-validator="validateCustomComponent"
        @on-submit="handleSaveProcess"
        @on-cancel="handleCancel">
        <template slot="bind_info">
          <process-form-property-table
            ref="bindInfo"
            v-model="bindInfo"
            :options="bindInfoProperty.option || []">
          </process-form-property-table>
        </template>
      </cmdb-form>
    </template>
  </bk-sideslider>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import {
    processPropertyRequestId,
    processPropertyGroupsRequestId
  } from './symbol'
  import RenderAppend from './process-form-append-render'
  import ProcessFormPropertyTable from './process-form-property-table'
  import { MENU_BUSINESS_SET_TOPOLOGY } from '@/dictionary/menu-symbol'
  import { ProcessTemplateService } from '@/service/business-set/process-template.js'
  import router from '@/router'

  export default {
    components: {
      ProcessFormPropertyTable
    },
    props: {
      type: String,
      serviceTemplateId: Number,
      processTemplateId: Number,
      instance: {
        type: Object,
        default: () => ({})
      },
      title: String,
      hostId: Number,
      bizId: Number,
      submitHandler: Function,
      invisibleProperties: {
        type: Array,
        default: () => ([])
      },
      /**
       * 是否展示操作按钮
       */
      showOptions: {
        type: Boolean,
        default: true
      }
    },
    provide() {
      return {
        form: this
      }
    },
    data() {
      return {
        isShow: false,
        internalType: this.type,
        internalTitle: this.title,
        properties: [],
        propertyGroups: [],
        bindedProperties: [],
        processTemplate: {},
        pending: true,
        invisibleNameProperties: ['bind_info'],
        flexProperties: ['bind_info'],
        formValuesReflect: {}
      }
    },
    computed: {
      ...mapState('bizSet', ['bizSetId']),
      ...mapGetters(['supplierAccount']),
      isBizSet() {
        return router.currentRoute.name === MENU_BUSINESS_SET_TOPOLOGY
      },
      bindInfoProperty() {
        return this.properties.find(property => property.bk_property_id === 'bind_info') || {}
      },
      bindInfo: {
        get() {
          return this.formValuesReflect.bind_info || []
        },
        set(values) {
          this.formValuesReflect.bind_info = values
        }
      },
      visibleProperties() {
        return this.properties.filter(property => !this.invisibleProperties.includes(property.bk_property_id))
      }
    },
    watch: {
      internalType() {
        this.updateFormWatcher()
      }
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
    mounted() {
      this.updateFormWatcher()
    },
    beforeDestroy() {
      this.teardownWatcher()
    },
    methods: {
      show() {
        this.isShow = true
      },
      teardownWatcher() {
        this.unwatchName && this.unwatchName()
        this.unwatchFormValues && this.unwatchFormValues()
      },
      updateFormWatcher() {
        if (this.internalType === 'view') {
          this.teardownWatcher()
        } else {
          this.$nextTick(() => {
            const { form } = this.$refs
            if (!form) {
              return this.updateFormWatcher() // 递归nextTick等待form创建完成
            }
            // watch form组件表单值，用于获取bind_info字段给进程表格字段组件使用
            this.unwatchFormValues = this.$watch(() => form.values, (values) => {
              this.formValuesReflect = values
            }, { immediate: true })
            // watch 名称，在用户未修改进程别名时，自动同步进程名称到进程别名
            this.unwatchName = this.$watch(() => form.values.bk_func_name, (newVal, oldValue) => {
              if (form.values.bk_process_name === oldValue) {
                form.values.bk_process_name = newVal
              }
            })
          })
        }
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
          const reqParams = {
            processTemplateId: this.processTemplateId
          }
          const reqConfig = {
            cancelPrevious: true
          }

          let processTemplate = null

          if (this.isBizSet) {
            processTemplate = await ProcessTemplateService.findOne({
              bizSetId: this.bizSetId,
              processTemplateId: this.processTemplateId
            }, reqConfig)
          } else {
            processTemplate = await this.$store.dispatch('processTemplate/getProcessTemplate', {
              params: reqParams,
              config: reqConfig
            })
          }

          this.processTemplate = processTemplate

          const { property } = this.processTemplate
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
      async validateCustomComponent() {
        const customComponents = []
        const { bindInfo } = this.$refs
        if (bindInfo) {
          customComponents.push(bindInfo)
        }
        const validatePromise = []
        customComponents.forEach((component) => {
          validatePromise.push(component.$validator.validateAll())
          validatePromise.push(component.$validator.validateScopes())
        })
        const results = await Promise.all(validatePromise)
        return results.every(result => result)
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
      async handleCancel() {
        const userConfirm = await this.beforeClose()
        if (!userConfirm) {
          return false
        }
        if (this.type === 'view') {
          this.internalType = this.type
          this.internalTitle = this.title
        } else {
          this.isShow = false
        }
      },
      beforeClose() {
        if (this.internalType === 'view') return Promise.resolve(true)
        const formChanged = !!Object.values(this.$refs.form.changedValues).length
        if (formChanged) {
          this.$refs.form.setChanged(true)
          return this.$refs.form.beforeClose()
        }
        return Promise.resolve(true)
      },
      renderAppend(h, { property }) {
        if (this.bindedProperties.includes(property.bk_property_id)) {
          // eslint-disable-next-line new-cap
          return RenderAppend(h, {
            serviceTemplateId: this.serviceTemplateId,
            property,
            bizId: this.bizId
          })
        }
        return ''
      },
      handleChangeInternalType() {
        this.internalType = 'update'
        this.internalTitle = this.$t('编辑进程')
      }
    }
  }
</script>
