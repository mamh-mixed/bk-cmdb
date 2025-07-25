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

import http from '@/api'

const findAllRequsetId = Symbol('findAllRequsetId')

const find = async ({ params, config }) => {
  try {
    const { count = 0, info: list = [] } = await http.post(`biz/search/${window.Supplier.account}`, params, config)
    return { count, list }
  } catch (error) {
    console.error(error)
  }
}

const findOne = async ({ bk_biz_id: bizId, config }) => {
  try {
    const { info } = await http.post(`biz/search/${window.Supplier.account}`, {
      condition: { bk_biz_id: { $eq: bizId } },
      fields: [],
      page: { start: 0, limit: 1 }
    }, config)
    const [instance] = info || [null]
    return instance
  } catch (error) {
    console.error(error)
    return null
  }
}

const findByIds = async ({ ids, config }) => {
  try {
    const { count = 0, info: list = [] } = await http.post(`biz/search/${window.Supplier.account}`, {
      condition: { bk_biz_id: { $in: ids } },
      fields: [],
      page: { start: 0, limit: ids.length }
    }, config)
    return { count, list }
  } catch (error) {
    console.error(error)
  }
}

const findAll = async () => {
  const data = await http.get('biz/simplify?sort=bk_biz_id', {
    requestId: findAllRequsetId,
    fromCache: false
  })

  return Object.freeze(data.info || [])
}

export default {
  find,
  findOne,
  findAll,
  findByIds
}
