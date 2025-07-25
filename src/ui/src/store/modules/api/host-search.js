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

/* eslint-disable no-unused-vars */

import has from 'has'
import $http from '@/api'
import { transformHostSearchParams, localSort } from '@/utils/tools'

const state = {

}

const getters = {

}

const actions = {
  /**
     * 根据条件查询主机
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
  searchHost({ commit, state, dispatch }, { params, config }) {
    return $http.post('hosts/search', transformHostSearchParams(params), config).then((data) => {
      if (has(data, 'info')) {
        data.info.forEach((host) => {
          localSort(host.module, 'bk_module_name')
          localSort(host.set, 'bk_set_name')
        })
      }
      return data
    })
  },

  /**
     * 获取主机详情
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {String} bkSupplierAccount 开发商账号
     * @param {Number} bkHostId 主机id
     * @return {Promise} promise 对象
     */
  getHostBaseInfo({ commit, state, dispatch, rootGetters }, { hostId, config }) {
    return $http.get(`hosts/${rootGetters.supplierAccount}/${hostId}`)
  },

  /**
     * 根据主机id获取主机快照数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} bkHostId 主机id
     * @return {Promise} promise 对象
     */
  getHostSnapshot({ commit, state, dispatch }, { hostId, config }) {
    return $http.get(`hosts/snapshot/${hostId}`, config)
  },

  /**
     * 根据主机id获取主机快照数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {Promise} promise 对象
     */
  searchHostByCondition({ commit, state, dispatch }, { params }) {
    return $http.post('hosts/snapshot/asstdetail', params)
  }
}

const mutations = {

}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
