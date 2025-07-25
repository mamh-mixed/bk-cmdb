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

export const jsonp = (url, data) => {
  if (!url) throw new Error('invalid URL')
  const callback = `CALLBACK${Math.random().toString()
    .slice(9, 18)}`
  const JSONP = document.createElement('script')
  JSONP.setAttribute('type', 'text/javascript')

  const headEle = document.getElementsByTagName('head')[0]

  let query = ''
  if (data) {
    if (typeof data === 'string') {
      query = `&${data}`
    } else if (typeof data === 'object') {
      for (const [key, value] of Object.entries(data)) {
        query += `&${key}=${encodeURIComponent(value)}`
      }
    }
    query += `&_time=${Date.now()}`
  }

  let promiseRejecter = null

  JSONP.src = `${url}?callback=${callback}${query}`
  JSONP.onerror = function (event) {
    promiseRejecter?.(event)
  }

  return new Promise((resolve, reject) => {
    promiseRejecter = reject
    try {
      window[callback] = (result) => {
        resolve(result)
        headEle.removeChild(JSONP)
        delete window[callback]
      }
      headEle.appendChild(JSONP)
    } catch (err) {
      reject(err)
    }
  })
}
