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
  <div class="relation-wrapper">
    <cmdb-tips
      class="mb10"
      tips-key="associationTips"
      :more-link="`${$helpDocUrlPrefix}/UserGuide/Feature/ModelRelationType.md`">
      {{$t('关联关系提示')}}
    </cmdb-tips>
    <p class="operation-box clearfix">
      <cmdb-auth class="inline-block-middle"
        :auth="{ type: $OPERATION.C_RELATION }">
        <bk-button slot-scope="{ disabled }"
          theme="primary"
          class="create-btn"
          :disabled="disabled"
          @click="createRelation">
          {{$t('新建')}}
        </bk-button>
      </cmdb-auth>
      <label class="search-input fr">
        <!-- <i class="bk-icon icon-search" @click="searchRelation(true)"></i> -->
        <bk-input type="text" class="cmdb-form-input"
          clearable
          v-model.trim="searchText"
          :right-icon="'bk-icon icon-search'"
          :placeholder="$t('请输入关联类型名称')"
          font-size="medium"
          @clear="searchRelation(true)"
          @enter="searchRelation(true)"
          @right-icon-click="searchRelation(true)">
        </bk-input>
      </label>
    </p>
    <bk-table
      v-bkloading="{ isLoading: $loading('searchAssociationType') }"
      :max-height="$APP.height - 229"
      :data="table.list"
      :pagination="table.pagination"
      :row-style="{ cursor: 'pointer' }"
      @row-click="handleShowDetails"
      @page-change="handlePageChange"
      @page-limit-change="handleSizeChange"
      @sort-change="handleSortChange">
      <bk-table-column prop="bk_asst_id" :label="$t('唯一标识')" sortable="custom"
        class-name="is-highlight" show-overflow-tooltip>
      </bk-table-column>
      <bk-table-column prop="bk_asst_name" :label="$t('名称')" sortable="custom" show-overflow-tooltip>
        <template slot-scope="{ row }">
          {{row['bk_asst_name'] || '--'}}
        </template>
      </bk-table-column>
      <bk-table-column prop="src_des" :label="$t('源->目标描述')" sortable="custom" show-overflow-tooltip></bk-table-column>
      <bk-table-column prop="dest_des" :label="$t('目标->源描述')" sortable="custom" show-overflow-tooltip></bk-table-column>
      <bk-table-column prop="count" :label="$t('使用数')"></bk-table-column>
      <bk-table-column
        fixed="right"
        prop="operation"
        :label="$t('操作')">
        <template slot-scope="{ row }">
          <cmdb-auth class="mr10"
            :auth="{ type: $OPERATION.U_RELATION, relation: [row.id] }"
            :ignore="row.ispre"
            v-bk-tooltips="{
              content: $t('禁止操作内置关联类型'),
              disabled: !row.ispre
            }">
            <bk-button slot-scope="{ disabled }"
              text
              theme="primary"
              :disabled="row.ispre || disabled"
              @click.stop="editRelation(row)">
              {{$t('编辑')}}
            </bk-button>
          </cmdb-auth>
          <cmdb-auth
            :auth="{ type: $OPERATION.D_RELATION, relation: [row.id] }"
            :ignore="row.ispre"
            v-bk-tooltips="{
              content: $t('禁止操作内置关联类型'),
              disabled: !row.ispre
            }">
            <bk-button slot-scope="{ disabled }"
              text
              theme="primary"
              :disabled="row.ispre || disabled"
              @click.stop="deleteRelation(row)">
              {{$t('删除')}}
            </bk-button>
          </cmdb-auth>
        </template>
      </bk-table-column>
      <cmdb-table-empty
        slot="empty"
        :stuff="table.stuff"
        :auth="{ type: $OPERATION.C_RELATION }"
        @create="createRelation"
        @clear="handleClearFilter"
      ></cmdb-table-empty>
    </bk-table>
    <bk-sideslider
      v-transfer-dom
      class="relation-slider"
      :width="450"
      :title="slider.title"
      :is-show.sync="slider.isShow"
      :before-close="handleSliderBeforeClose">
      <the-relation
        ref="relationForm"
        slot="content"
        class="model-slider-content"
        v-if="slider.isShow"
        :is-edit="slider.isEdit"
        :is-read-only="slider.isReadOnly"
        :relation="slider.relation"
        :save-btn-text="slider.isEdit ? $t('保存') : $t('提交')"
        @saved="saveRelation"
        @cancel="handleSliderBeforeClose">
      </the-relation>
    </bk-sideslider>
  </div>
</template>

<script>
  import theRelation from './_detail'
  import { mapActions } from 'vuex'
  import associationService from '@/service/association'
  import { escapeRegexChar } from '@/utils/util'
  export default {
    components: {
      theRelation
    },
    data() {
      return {
        slider: {
          isShow: false,
          isEdit: false,
          title: this.$t('新建关联类型'),
          relation: {},
          isReadOnly: false
        },
        searchText: '',
        table: {
          list: [],
          pagination: {
            count: 0,
            current: 1,
            ...this.$tools.getDefaultPaginationConfig()
          },
          defaultSort: '-ispre',
          sort: '-ispre',
          stuff: {
            type: 'default',
            payload: {
              resource: this.$t('关联类型')
            }
          }
        },
        sendSearchText: ''
      }
    },
    computed: {
      searchParams() {
        const params = {
          page: {
            start: (this.table.pagination.current - 1) * this.table.pagination.limit,
            limit: this.table.pagination.limit,
            sort: this.table.sort
          }
        }
        if (this.sendSearchText.length) {
          Object.assign(params, {
            condition: {
              bk_asst_name: {
                $regex: this.sendSearchText
              }
            }
          })
        }
        return params
      }
    },
    created() {
      this.searchRelation()
    },
    methods: {
      ...mapActions('objectAssociation', [
        'searchAssociationType',
        'deleteAssociationType'
      ]),
      searchRelation(fromClick) {
        if (fromClick) {
          this.sendSearchText = escapeRegexChar(this.searchText)
          this.table.pagination.current = 1
          this.table.stuff.type = 'search'
        }
        this.searchAssociationType({
          params: this.searchParams,
          config: {
            requestId: 'searchAssociationType'
          }
        }).then((data) => {
          if (data.count && !data.info.length) {
            this.table.pagination.current -= 1
            this.searchRelation()
          }
          this.table.list = data.info
          this.searchUsageCount()
          this.table.pagination.count = data.count
          this.$http.cancel('post_searchAssociationType')
        })
      },
      async searchUsageCount() {
        const asstIds = []
        this.table.list.forEach(({ bk_asst_id: asstId }) => asstIds.push(asstId))
        const res = await associationService.getAssociationCount({
          params: {
            asst_ids: asstIds
          }
        })
        this.table.list.forEach((item) => {
          const asst = res?.associations?.find(({ bk_asst_id: asstId }) => asstId === item.bk_asst_id)
          this.$set(item, 'count', asst ? asst.count : '--')
        })
        this.table.list.splice()
      },
      createRelation() {
        this.slider.title = this.$t('新建关联类型')
        this.slider.isReadOnly = false
        this.slider.isEdit = false
        this.slider.isShow = true
      },
      editRelation(relation) {
        this.slider.title = this.$t('编辑关联类型')
        this.slider.isReadOnly = false
        this.slider.relation = relation
        this.slider.isEdit = true
        this.slider.isShow = true
      },
      deleteRelation(relation) {
        this.$bkInfo({
          title: this.$tc('确定删除关联类型？', relation.bk_asst_name, { name: relation.bk_asst_name }),
          confirmFn: async () => {
            await this.deleteAssociationType({
              id: relation.id,
              config: {
                requestId: 'deleteAssociationType'
              }
            })
            this.$success(this.$t('删除成功'))
            this.searchRelation()
          }
        })
      },
      saveRelation() {
        this.slider.isShow = false
        const { isEdit } = this.slider
        const text = isEdit ? '编辑成功' : '创建成功'
        this.$success(this.$t(text))
        this.searchRelation()
      },
      handlePageChange(current) {
        this.table.pagination.current = current
        this.searchRelation()
      },
      handleSizeChange(size) {
        this.table.pagination.limit = size
        this.handlePageChange(1)
      },
      handleSortChange(sort) {
        this.table.sort = this.$tools.getSort(sort)
        this.searchRelation()
      },
      handleSliderBeforeClose() {
        const hasChanged = Object.keys(this.$refs.relationForm.changedValues).length
        if (hasChanged) {
          return this.$refs.relationForm.beforeClose(() => {
            this.slider.isShow = false
          })
        }
        this.slider.isShow = false
        return true
      },
      handleShowDetails(relation, event, column = {}) {
        if (column.property === 'operation') return
        this.slider.title = this.$t('关联类型详情')
        this.slider.relation = relation
        this.slider.isReadOnly = true
        this.slider.isEdit = true
        this.slider.isShow = true
      },
      handleClearFilter() {
        this.searchText = ''
        this.table.stuff.type = 'default'
        this.searchRelation(true)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .relation-wrapper {
        padding: 15px 20px 0;
    }
    .operation-box {
        margin: 0 0 14px 0;
        font-size: 0;
        .create-btn {
            margin: 0 10px 0 0;
        }
        .search-input {
            position: relative;
            display: inline-block;
            width: 300px;
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
