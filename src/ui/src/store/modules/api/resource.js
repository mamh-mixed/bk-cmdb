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

import $http from '@/api'
const directory = {
  namespaced: true,
  actions: {
    create(context, { params, config }) {
      return $http.post('create/resource/directory', params, config)
    },
    delete(context, { id, config }) {
      return $http.delete(`delete/resource/directory/${id}`, config)
    },
    update(context, { id, params, config }) {
      return $http.put(`update/resource/directory/${id}`, params, config)
    },
    findMany(context, { params, config }) {
      return $http.post('findmany/resource/directory', params, config)
    }
  }
}

const host = {
  namespaced: true,
  modules: {
    transfer: {
      namespaced: true,
      actions: {
        directory(context, { params, config }) {
          return $http.post('host/transfer/resource/directory', params, config)
        },
        idle(context, { params, config }) {
          return $http.post('hosts/modules/resource/idle', params, config)
        }
      }
    }
  }
}

export default {
  namespaced: true,
  modules: {
    directory,
    host
  }
}
