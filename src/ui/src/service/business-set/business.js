/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import $http from '@/api'

/**
 * 获取业务信息
 * @param {Object} pathParams
 * @param {number} pathParams.bizSetId 业务集 ID
 * @param {number} pathParams.bizId 业务 ID
 * @param {Object} params 查询参数
 * @param {Object} config 请求配置
 * @returns {Promise}
 */
export const findOne = ({
  bizSetId,
  bizId,
}, config) => $http.post('find/biz_set/biz_list', {
  bk_biz_set_id: bizSetId,
  filter: {
    condition: 'AND',
    rules: [
      {
        field: 'bk_biz_id',
        operator: 'equal',
        value: bizId
      }
    ]
  },
  page: {
    start: 0,
    limit: 1
  }
}, config)

export const BusinessService = {
  findOne
}
