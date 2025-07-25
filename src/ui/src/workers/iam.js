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

import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
import { postData } from './utils.js'
import has from 'has'

const state = {
  config: {}
}

const actions = {
  config(state, payload) {
    state.config = payload
  },
  async preverify({ config }, { authList }) {
    try {
      const resources = TRANSFORM_TO_INTERNAL(authList)
      const authedList = []

      if (config.IAM_ENABLED) {
        const result = await postData(`${config.API_PREFIX}/auth/verify`, { resources })
        result.forEach((item, index) => {
          // 第0个元素是一个auth的key，第1个元素为鉴权是否通过，用于在使用数据时根据key来查找是否鉴权通过
          authedList.push([resources[index], item.is_pass])
        })
      } else {
        resources.forEach((item) => {
          authedList.push([item, true])
        })
      }

      // 将结果传递回主线程
      self.postMessage({ type: 'preverify', payload: authedList })
    } catch (error) {
      self.postMessage({ type: 'error', error })
    }
  }
}

self.onmessage = ({ data: { type, payload } }) => {
  if (has(actions, type)) {
    actions[type](state, payload)
  }
}

self.onerror = (error) => {
  self.postMessage({ type: 'error', error })
}
