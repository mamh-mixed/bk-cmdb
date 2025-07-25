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

import Vue from 'vue'
import i18n from '@/i18n'
import magicbox from 'bk-magic-vue'
import 'bk-magic-vue/dist/bk-magic-vue.min.css'
import './magicbox.scss'

const magicboxLanguageMap = {
  zh_CN: magicbox.locale.lang.zhCN,
  en: magicbox.locale.lang.enUS
}

export const setLocale = (targetLocale) => {
  const locale = targetLocale || i18n.locale
  i18n.mergeLocaleMessage(locale, magicboxLanguageMap[locale])
  magicbox.locale.use(magicboxLanguageMap[locale])
}
setLocale()

Vue.use(magicbox, {
  'bk-sideslider': {
    quickClose: true,
    width: 800
  },
  // 'bk-input': {
  //   fontSize: 'medium'
  // },
  // 'bk-select': {
  //   fontSize: 'medium'
  // },
  'bk-big-tree': {
    useDefaultEmpty: true
  },
  i18n: (key, value) => i18n.t(key, value)
})

export const $error = (message, delay = 3000) => magicbox.bkMessage({
  message,
  delay,
  theme: 'error',
  ellipsisLine: 0
})

export const $success = (message, delay = 3000) => magicbox.bkMessage({
  message,
  delay,
  theme: 'success',
  ellipsisLine: 0
})

export const $info = (message, delay = 3000) => magicbox.bkMessage({
  message,
  delay,
  theme: 'primary',
  ellipsisLine: 0
})

export const $warn = (message, delay = 3000) => magicbox.bkMessage({
  message,
  delay,
  theme: 'warning',
  hasCloseIcon: true,
  ellipsisLine: 0
})

export const $bkInfo = magicbox.bkInfoBox

export const { $bkPopover } = Vue.prototype

Vue.prototype.$error = $error
Vue.prototype.$success = $success
Vue.prototype.$info = $info
Vue.prototype.$warn = $warn
Vue.prototype.$bkInfo = magicbox.bkInfoBox
