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
  <div class="business-layout">
    <div class="business-options clearfix">
      <cmdb-auth class="fl" :auth="{ type: $OPERATION.C_BUSINESS }">
        <bk-button slot-scope="{ disabled }"
          class="fl"
          theme="primary"
          :disabled="disabled"
          @click="handleCreate">
          {{$t('新建')}}
        </bk-button>
      </cmdb-auth>
      <cmdb-auth :auth="{ type: $OPERATION.U_BUSINESS }">
        <template #default="{ disabled }">
          <bk-button
            class="ml10"
            :disabled="selectedRows.length === 0 || disabled"
            @click="handleBatchEdit"
          >
            {{ $t("批量编辑") }}
          </bk-button>
        </template>
      </cmdb-auth>
      <cmdb-button-group
        class="mr10"
        :buttons="buttons"
        :expand="false">
      </cmdb-button-group>
      <div class="options-button fr">
        <cmdb-auth class="inline-block-middle" :auth="{ type: $OPERATION.BUSINESS_ARCHIVE }">
          <icon-button slot-scope="{ disabled }"
            class="mr10"
            icon="icon-cc-history"
            v-bk-tooltips.top="$t('查看已归档业务')"
            :disabled="disabled"
            @click="routeToHistory">
          </icon-button>
        </cmdb-auth>
        <icon-button
          icon="icon-cc-setting"
          v-bk-tooltips.top="$t('列表显示属性配置')"
          @click="columnsConfig.show = true">
        </icon-button>
      </div>
      <div class="options-filter clearfix fr">
        <cmdb-property-selector
          class="filter-selector fl"
          v-model="filter.field"
          :properties="fastSearchProperties">
        </cmdb-property-selector>
        <component class="filter-value fl r0"
          :is="`cmdb-search-${filterType}`"
          :placeholder="filterPlaceholder"
          :property="filterProperty"
          :class="filterType"
          :fuzzy="true"
          v-bind="filterComponentProps"
          v-model="filter.value"
          @change="handleFilterValueChange"
          @enter="handleFilterValueEnter"
          @clear="handleFilterValueEnter">
        </component>
      </div>
    </div>
    <bk-table class="business-table"
      v-bkloading="{ isLoading: $loading('post_searchBusiness_list') }"
      :data="table.visibleList"
      :pagination="table.pagination"
      :max-height="$APP.height - 200"
      @sort-change="handleSortChange"
      @page-limit-change="handleSizeChange"
      @page-change="handlePageChange">
      <batch-selection-column
        width="60px"
        row-key="bk_biz_id"
        ref="batchSelectionColumn"
        indeterminate
        :cross-page="table.visibleList.length >= table.pagination.limit"
        :selected-rows.sync="selectedRows"
        :data="table.visibleList"
        :full-data="table.list"
      >
      </batch-selection-column>
      <bk-table-column v-for="column in table.header"
        sortable="custom"
        :key="column.id"
        :prop="column.id"
        :label="column.name"
        :min-width="$tools.getHeaderPropertyMinWidth(column.property, { hasSort: true })"
        :show-overflow-tooltip="$tools.isShowOverflowTips(column.property)">
        <template slot-scope="{ row }">
          <cmdb-property-value
            :theme="column.id === 'bk_biz_id' ? 'primary' : 'default'"
            :value="row[column.id]"
            :show-unit="false"
            :property="column.property"
            :instance="row"
            show-on="cell"
            @click.native.stop="handleValueClick(row, column)">
          </cmdb-property-value>
        </template>
      </bk-table-column>
      <bk-table-column :label="$t('操作')" fixed="right">
        <template slot-scope="{ row }">
          <cmdb-auth @click.native.stop :auth="{ type: $OPERATION.BUSINESS_ARCHIVE, relation: [row.bk_biz_id] }">
            <template slot-scope="{ disabled }">
              <span class="text-primary"
                style="color: #dcdee5 !important; cursor: not-allowed;"
                v-if="row['bk_biz_name'] === '蓝鲸' && !disabled"
                v-bk-tooltips.top="$t('内置业务不可归档')"
                @click.stop>
                {{$t('归档')}}
              </span>
              <bk-button v-else
                theme="primary"
                :disabled="disabled"
                :text="true"
                @click.stop="handleDelete(row)">
                {{$t('归档')}}
              </bk-button>
            </template>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff"
        :auth="{ type: $OPERATION.C_BUSINESS }"
        @clear="handleClearFilter">
        <bk-exception type="403" scene="part">
          <i18n path="业务列表提示语" class="table-empty-tips">
            <template #auth>
              <bk-link theme="primary" @click="handleApplyPermission">{{$t('申请查看权限')}}</bk-link>
            </template>
            <template #create>
              <cmdb-auth :auth="{ type: $OPERATION.C_BUSINESS }">
                <bk-button slot-scope="{ disabled }" text
                  theme="primary"
                  class="text-btn"
                  :disabled="disabled"
                  @click="handleCreate">
                  {{$t('立即创建')}}
                </bk-button>
              </cmdb-auth>
            </template>
          </i18n>
        </bk-exception>
      </cmdb-table-empty>
    </bk-table>
    <bk-sideslider
      v-transfer-dom
      :is-show.sync="slider.show"
      :title="slider.title"
      :width="800"
      :before-close="handleSliderBeforeClose">
      <bk-tab :active.sync="tab.active" type="unborder-card" slot="content" v-if="slider.show">
        <bk-tab-panel name="attribute" :label="$t('属性')" style="width: calc(100% + 40px);margin: 0 -20px;">
          <cmdb-form
            ref="form"
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="attribute.inst.edit"
            :is-main-line="true"
            :type="attribute.type"
            :save-auth="saveAuth"
            @on-submit="handleSave"
            @on-cancel="handleSliderBeforeClose">
          </cmdb-form>
        </bk-tab-panel>
      </bk-tab>
    </bk-sideslider>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="batchUpdateSlider.show"
      :title="$t('批量修改业务')"
      :width="800"
      :before-close="handleBatchUpdateSliderBeforeClose"
    >
      <bk-tab
        :active.sync="tab.active"
        type="unborder-card"
        slot="content"
        v-if="batchUpdateSlider.show"
      >
        <bk-tab-panel
          name="attribute"
          :label="$t('属性')"
          style="width: calc(100% + 40px); margin: 0 -20px"
        >
          <cmdb-form-multiple
            ref="batchUpdateForm"
            :properties="properties"
            :property-groups="propertyGroups"
            :save-auth="saveAuth"
            :show-default-value="true"
            @on-submit="handleMultipleSave"
            :loading="batchUpdateSlider.loading"
            @on-cancel="handleBatchUpdateSliderBeforeClose"
          >
          </cmdb-form-multiple>
        </bk-tab-panel>
      </bk-tab>
    </bk-sideslider>

    <bk-sideslider
      v-transfer-dom
      :is-show.sync="columnsConfig.show"
      :width="600"
      :title="$t('列表显示属性配置')"
      :before-close="handleColumnsConfigSliderBeforeClose"
    >
      <cmdb-columns-config
        slot="content"
        v-if="columnsConfig.show"
        ref="cmdbColumnsConfig"
        :properties="properties"
        :selected="columnsConfig.selected"
        :disabled-columns="columnsConfig.disabledColumns"
        @on-apply="handleApplayColumnsConfig"
        @on-cancel="handleColumnsConfigSliderBeforeClose"
        @on-reset="handleResetColumnsConfig">
      </cmdb-columns-config>
    </bk-sideslider>
    <cmdb-model-fast-link :obj-id="objId"></cmdb-model-fast-link>
  </div>
</template>

<script>
  import { translateAuth } from '@/setup/permission'
  import { mapState, mapGetters, mapActions } from 'vuex'
  import { MENU_RESOURCE_BUSINESS_HISTORY, MENU_RESOURCE_BUSINESS_DETAILS } from '@/dictionary/menu-symbol'
  import cmdbColumnsConfig from '@/components/columns-config/columns-config.vue'
  import cmdbPropertySelector from '@/components/property-selector'
  import RouterQuery from '@/router/query'
  import Utils from '@/components/filters/utils'
  import throttle from 'lodash.throttle'
  import BatchSelectionColumn from '@/components/batch-selection-column'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants.js'
  import cmdbModelFastLink from '@/components/model-fast-link'
  import cmdbButtonGroup from '@/components/ui/other/button-group'
  import { QUERY_OPERATOR } from '@/utils/query-builder-operator'

  export default {
    components: {
      cmdbColumnsConfig,
      cmdbPropertySelector,
      BatchSelectionColumn,
      cmdbModelFastLink,
      cmdbButtonGroup
    },
    data() {
      return {
        properties: [],
        propertyGroups: [],
        table: {
          header: [],
          list: [],
          visibleList: [],
          pagination: {
            count: 0,
            current: 1,
            ...this.$tools.getDefaultPaginationConfig()
          },
          defaultSort: 'bk_biz_id',
          sort: 'bk_biz_id',
          stuff: {
            type: 'default',
            payload: {
              resource: this.$t('业务')
            }
          }
        },
        selectedRows: [],
        editAll: false,
        filter: {
          field: 'bk_biz_name',
          value: '',
          operator: ''
        },
        slider: {
          show: false,
          title: ''
        },
        batchUpdateSlider: {
          show: false,
          loading: false,
        },
        tab: {
          active: 'attribute'
        },
        attribute: {
          type: null,
          inst: {
            edit: {},
            details: {}
          }
        },
        columnsConfig: {
          show: false,
          selected: [],
          disabledColumns: ['bk_biz_id', 'bk_biz_name']
        },
        columnsConfigKey: 'biz_custom_table_columns'
      }
    },
    computed: {
      ...mapState('userCustom', ['globalUsercustom']),
      ...mapGetters(['supplierAccount', 'userName']),
      ...mapGetters('userCustom', ['usercustom']),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('objectModelClassify', ['getModelById']),
      objId() {
        return BUILTIN_MODELS.BUSINESS
      },
      customBusinessColumns() {
        return this.usercustom[this.columnsConfigKey] || []
      },
      globalCustomColumns() {
        return this.globalUsercustom.biz_global_custom_table_columns || []
      },
      saveAuth() {
        const { type } = this.attribute
        if (type === 'create') {
          return { type: this.$OPERATION.C_BUSINESS }
        } if (type === 'update') {
          return { type: this.$OPERATION.U_BUSINESS }
        }
        return null
      },
      model() {
        return this.getModelById('biz') || {}
      },
      filterProperty() {
        const property = this.properties.find(property => property.bk_property_id === this.filter.field)
        return property || null
      },
      filterType() {
        if (this.filterProperty) {
          return this.filterProperty.bk_property_type
        }
        return 'singlechar'
      },
      filterPlaceholder() {
        return Utils.getPlaceholder(this.filterProperty)
      },
      filterComponentProps() {
        return Utils.getBindProps(this.filterProperty)
      },
      fastSearchProperties() {
        return this.properties.filter(item => item.bk_property_type !== PROPERTY_TYPES.INNER_TABLE)
      },
      buttons() {
        const buttonConfig = [{
          id: 'export',
          text: this.$t('导出选中'),
          handler: this.exportField,
          disabled: !this.selectedRows.length
        }, {
          id: 'batchExport',
          text: this.$t('导出全部'),
          handler: () => this.exportField('all'),
          disabled: !this.table.pagination.count
        }]
        return buttonConfig
      },
    },
    watch: {
      'filter.field'() {
        this.genFilterCondition()
      },
      'slider.show'(show) {
        if (!show) {
          this.tab.active = 'attribute'
        }
      },
      customBusinessColumns() {
        this.setTableHeader()
      }
    },
    async created() {
      try {
        this.properties = await this.searchObjectAttribute({
          injectId: 'biz',
          params: {
            bk_obj_id: 'biz',
            bk_supplier_account: this.supplierAccount
          },
          config: {
            requestId: 'post_searchObjectAttribute_biz',
            fromCache: true
          }
        })
        await Promise.all([
          this.getPropertyGroups(),
          this.setTableHeader()
        ])
        this.throttleGetTableData = throttle(this.getTableData, 300, { leading: false, trailing: true })
        this.unwatch = RouterQuery.watch('*', async ({
          page = 1,
          limit = this.table.pagination.limit,
          filter = '',
          operator = '',
          field = 'bk_biz_name'
        }) => {
          this.filter.field = field
          this.table.pagination.current = parseInt(page, 10)
          this.table.pagination.limit = parseInt(limit, 10)

          await this.$nextTick()

          this.genFilterCondition(filter, operator)

          this.throttleGetTableData()
        }, { immediate: true })
      } catch (e) {
        // ignore
      }

      if (this.$route.query.create) {
        this.handleCreate()
      }
    },
    beforeDestroy() {
      this.unwatch()
    },
    methods: {
      ...mapActions('objectModelFieldGroup', ['searchGroup']),
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),
      ...mapActions('objectBiz', [
        'searchBusiness',
        'archiveBusiness',
        'updateBusiness',
        'batchUpdateBusiness',
        'createBusiness',
        'searchBusinessById'
      ]),
      /**
       * 通过 bk_property_id 或自定义过滤项和过滤操作符来生成过滤条件
       * @param {string} filter 过滤项
       * @param {string} operator 过滤操作符
       */
      genFilterCondition(filter = '', operator = '') {
        const defaultData = Utils.getDefaultData(this.filterProperty)
        const isBizName = this.filterProperty.bk_property_id === 'bk_biz_name'
        const { LIKE } = QUERY_OPERATOR

        // 这里的业务名称因只支持模糊搜索，且 getDefaultData 返回的过滤方式没有命中，所以只能单独加个判断
        if (isBizName) {
          this.filter.operator = LIKE
          this.filter.value = filter || ''
        } else {
          this.filter.operator = operator || defaultData.operator
          this.filter.value = this.formatFilterValue(
            { value: filter, operator: this.filter.operator },
            defaultData.value
          )
        }
      },
      getCondition() {
        // 这里先直接复用转换通用模型实例查询条件的方法
        const condition = {
          [this.filterProperty.id]: {
            value: this.filter.value,
            operator: this.filter.operator
          }
        }
        const {
          conditions,
          time_condition
        } = Utils.transformGeneralModelCondition(condition, this.properties)
        return { conditions, time_condition }
      },
      async exportField(type = 'select') {
        const useExport = await import('@/components/export-file')
        const title = type === 'select' ? '导出选中' : '导出全部'
        const count = type === 'select' ? this.selectedRows.length : this.table.pagination.count

        useExport.default({
          title: this.$t(title),
          bk_obj_id: 'biz',
          defaultSelectedFields: this.table.header.map(item => item.id),
          count,
          steps: [{ title: this.$t('选择字段'), icon: 1 }],
          confirmBtnText: this.$t('导出'),
          submit: (state, task) => {
            const { fields } = state
            const params = {
              export_custom_fields: fields.value.map(property => property.bk_property_id),
            }
            if (type === 'select') {
              const selected = this.selectedRows.map(e => e.bk_biz_id)
              params.bk_biz_ids = selected
              params.export_condition = {
                page: {
                  start: 0,
                  limit: selected.length,
                  sort: this.table.sort
                }
              }
            }
            if (type === 'all') {
              const {
                conditions,
                time_condition: timeCondition
              } = this.getCondition()
              params.export_condition = {
                filter: conditions,
                time_condition: timeCondition,
                page: {
                  ...task.current.value.page,
                  sort: this.table.sort
                }
              }
            }

            return this.$http.download({
              url: `${window.API_HOST}biz/export`,
              method: 'post',
              name: task.current.value.name,
              data: params
            })
          }
        }).show()
      },
      handleBatchEdit() {
        this.batchUpdateSlider.show = true
      },
      async handleMultipleSave(changedValues) {
        const includeBizIds = this.selectedRows.map(r => r.bk_biz_id)
        const condition = {
          ...this.getSearchParams()?.condition
        }

        if (includeBizIds?.length > 0) {
          condition.bk_biz_id = { $in: includeBizIds }
        }

        this.batchUpdateSlider.loading = true
        this.batchUpdateBusiness({
          params: {
            properties: changedValues,
            condition
          },
        })
          .then(() => {
            this.$refs.batchSelectionColumn.clearSelection()
            this.batchUpdateSlider.show = false
            RouterQuery.set({
              _t: Date.now(),
            })
          })
          .catch((err) => {
            console.log(err)
          })
          .finally(() => {
            this.batchUpdateSlider.loading = false
          })
      },
      getPropertyGroups() {
        return this.searchGroup({
          objId: 'biz',
          params: {},
          config: {
            fromCache: true,
            requestId: 'post_searchGroup_biz'
          }
        }).then((groups) => {
          this.propertyGroups = groups
          return groups
        })
      },
      setTableHeader() {
        return new Promise((resolve) => {
          // eslint-disable-next-line max-len
          const customColumns = this.customBusinessColumns.length ? this.customBusinessColumns : this.globalCustomColumns
          // eslint-disable-next-line max-len
          const headerProperties = this.$tools.getHeaderProperties(this.properties, customColumns, this.columnsConfig.disabledColumns)
          resolve(headerProperties)
        }).then((properties) => {
          this.updateTableHeader(properties)
          this.columnsConfig.selected = properties.map(property => property.bk_property_id)
        })
      },
      updateTableHeader(properties) {
        // 数组length在没有变化时候，需要先清空数组在赋值。否则表头无法实时更新
        this.table.header = []
        this.$nextTick(() => {
          this.table.header = properties.map(property => ({
            id: property.bk_property_id,
            name: this.$tools.getHeaderPropertyName(property),
            property
          }))
        })
      },
      handleValueClick(row, column) {
        if (column.id !== 'bk_biz_id') {
          return false
        }
        this.$routerActions.redirect({
          name: MENU_RESOURCE_BUSINESS_DETAILS,
          params: {
            bizId: row.bk_biz_id
          },
          history: true
        })
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort)
        this.handlePageChange(1)
        this.getTableData()
      },
      handleSizeChange(size) {
        this.table.pagination.limit = size
        this.handlePageChange(1)
        this.renderVisibleList()
      },
      handlePageChange(page) {
        this.table.pagination.current = page
        this.renderVisibleList()
      },
      renderVisibleList() {
        const { limit, current } = this.table.pagination
        this.table.visibleList = this.table.list.slice((current - 1) * limit, current * limit)
      },
      getBusinessList(config = { cancelPrevious: true }) {
        return this.searchBusiness({
          params: this.getSearchParams(),
          config: Object.assign({ requestId: 'post_searchBusiness_list' }, config)
        })
      },
      formatFilterValue({ value: currentValue, operator }, defaultValue) {
        let value = currentValue.toString().length ? currentValue : defaultValue
        const isNumber = ['int', 'float'].includes(this.filterType)
        if (isNumber && value) {
          value = parseFloat(value, 10)
        } else if (operator === '$in') {
          // eslint-disable-next-line no-nested-ternary
          value = Array.isArray(value) ? value : !!value ? [value] : []
        } else if (Array.isArray(value)) {
          value = value.filter(value => !!value)
        }
        return value
      },
      handleFilterValueChange() {
        const hasEnterEvnet = ['float', 'int', 'longchar', 'singlechar']
        if (hasEnterEvnet.includes(this.filterType)) return
        this.handleFilterValueEnter()
      },
      handleFilterValueEnter() {
        this.$refs.batchSelectionColumn.clearSelection()
        RouterQuery.set({
          _t: Date.now(),
          page: 1,
          field: this.filter.field,
          filter: this.filter.value,
          operator: this.filter.operator
        })
      },
      getTableData() {
        this.getBusinessList({ cancelPrevious: true, globalPermission: false }).then((data) => {
          const { current, limit } = this.table.pagination
          const { count } = data

          this.table.pagination.current = Math.min(current, Math.ceil(count / limit))
          this.table.list = data.info
          this.table.pagination.count = data.count

          this.table.stuff.type = this.filter.value.toString().length ? 'search' : 'default'

          this.renderVisibleList()

          return data
        })
          .catch(({ permission }) => {
            if (!permission) return
            this.table.stuff = {
              type: 'permission',
              payload: { permission }
            }
          })
      },
      getSearchParams() {
        const params = {
          condition: {
            bk_data_status: { $ne: 'disabled' }
          },
          fields: [],
          page: {
            start: 0,
            limit: 10000,
            sort: this.table.sort
          }
        }

        if (!this.filter.value.toString()) {
          return params
        }

        if (this.filterType === 'time') {
          const condition = {
            [this.filter.field]: {
              value: this.filter.value,
              operator: this.filter.operator
            }
          }
          const { time_condition: timeCondition } = Utils.transformGeneralModelCondition(condition, this.properties)
          if (timeCondition) {
            params.time_condition = timeCondition
          }
          return params
        }

        if (this.filter.operator === '$range') {
          const [start, end] = this.filter.value
          params.condition[this.filter.field] = {
            $gte: start,
            $lte: end
          }
          return params
        }
        if (this.filterType === 'objuser') {
          const filterValue  = this.filter.value
          const multiple = filterValue.length > 1
          params.condition[this.filter.field] = multiple ? { $in: filterValue } : filterValue.toString()
          return params
        }

        if (this.filter.field === 'bk_biz_name') {
          params.condition[this.filter.field] = this.filter.value.toString()
          return params
        }

        params.condition[this.filter.field] = { [this.filter.operator]: this.filter.value }

        return params
      },
      async handleEdit(inst) {
        const bizNameProperty = this.$tools.getProperty(this.properties, 'bk_biz_name')
        bizNameProperty.isreadonly = inst.bk_biz_name === '蓝鲸'
        this.attribute.inst.edit = inst
        this.attribute.type = 'update'
      },
      handleCreate() {
        this.attribute.type = 'create'
        this.attribute.inst.edit = {}
        this.slider.show = true
        this.slider.title = `${this.$t('创建')} ${this.model.bk_obj_name}`
      },
      handleDelete(inst) {
        this.$bkInfo({
          title: this.$t('确认要归档', { name: inst.bk_biz_name }),
          subTitle: this.$t('归档确认信息'),
          confirmFn: () => {
            this.archiveBusiness(inst.bk_biz_id).then(() => {
              this.slider.show = false
              this.$success(this.$t('归档成功'))
              this.getTableData()
              this.$http.cancel('post_searchBusiness_$ne_disabled')
            })
          }
        })
      },
      handleSave(values, changedValues, originalValues, type) {
        if (type === 'update') {
          this.updateBusiness({
            bizId: originalValues.bk_biz_id,
            params: values
          }).then(() => {
            this.attribute.inst.details = Object.assign({}, originalValues, values)
            RouterQuery.refresh()
            this.closeCreateSlider()
            this.$success(this.$t('修改成功'))
            this.$http.cancel('post_searchBusiness_$ne_disabled')
          })
        } else {
          delete values.bk_biz_id // properties中注入了前端自定义的bk_biz_id属性
          this.createBusiness({
            params: values
          }).then(() => {
            RouterQuery.refresh()
            this.closeCreateSlider()
            this.$success(this.$t('创建成功'))
            this.$http.cancel('post_searchBusiness_$ne_disabled')
          })
        }
      },
      closeCreateSlider() {
        if (this.attribute.type === 'create') {
          this.slider.show = false
        }
      },
      handleApplayColumnsConfig(properties) {
        this.$store.dispatch('userCustom/saveUsercustom', {
          [this.columnsConfigKey]: properties.map(property => property.bk_property_id)
        })
        this.columnsConfig.show = false
      },
      handleResetColumnsConfig() {
        this.$store.dispatch('userCustom/saveUsercustom', {
          [this.columnsConfigKey]: []
        })
        this.columnsConfig.show = false
      },
      routeToHistory() {
        this.$routerActions.redirect({
          name: MENU_RESOURCE_BUSINESS_HISTORY,
          history: true
        })
      },
      handleSliderBeforeClose() {
        this.addDoubleConfirm(this.$refs.form, this.closeCreateSlider)
      },
      handleBatchUpdateSliderBeforeClose() {
        this.addDoubleConfirm(this.$refs.batchUpdateForm, () => {
          this.batchUpdateSlider.show = false
        })
      },
      addDoubleConfirm(componentRef, confirmCallback) {
        const { changedValues } = componentRef
        if (this.tab.active === 'attribute') {
          if (Object.keys(changedValues).length) {
            componentRef.setChanged(true)
            return componentRef.beforeClose(confirmCallback)
          }

          confirmCallback && confirmCallback()

          return true
        }

        confirmCallback && confirmCallback()

        return true
      },
      handleColumnsConfigSliderBeforeClose() {
        const refColumns = this.$refs.cmdbColumnsConfig
        if (!refColumns) {
          return
        }
        const { columnsChangedValues } = refColumns
        if (columnsChangedValues?.()) {
          refColumns.setChanged(true)
          return refColumns.beforeClose(() => {
            this.columnsConfig.show = false
          })
        }
        this.columnsConfig.show = false
      },
      async handleApplyPermission() {
        try {
          const permission = translateAuth({
            type: this.$OPERATION.R_BUSINESS,
            relation: []
          })
          const url = await this.$store.dispatch('auth/getSkipUrl', { params: permission })
          window.open(url)
        } catch (e) {
          console.error(e)
        }
      },
      handleClearFilter() {
        RouterQuery.clear()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .business-layout {
        padding: 15px 20px 0;
    }
    .options-filter{
        position: relative;
        margin-right: 10px;
        width: 439px;
        .filter-selector{
            width: 120px;
            border-radius: 2px 0 0 2px;
            margin-right: -1px;
        }
        .filter-value{
            width: 320px;
            border-radius: 0 2px 2px 0;
            /deep/ .bk-form-input {
                border-radius: 0 2px 2px 0;
            }
        }
        .filter-search{
            position: absolute;
            right: 10px;
            top: 9px;
            cursor: pointer;
        }
    }
    .options-button{
        font-size: 0;
        .bk-button {
            width: 32px;
            padding: 0;
            /deep/ .bk-icon {
                line-height: 14px;
            }
        }
    }
    .business-table{
        margin-top: 14px;
    }
    .table-empty-tips {
        display: flex;
        align-items: center;
        justify-content: center;
        .text-btn {
            font-size: 14px;
        }
    }
</style>
