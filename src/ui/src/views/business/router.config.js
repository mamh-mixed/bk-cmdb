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

import Meta from '@/router/meta'
import {
  MENU_RESOURCE_BUSINESS,
  MENU_RESOURCE_BUSINESS_DETAILS,
  MENU_RESOURCE_BUSINESS_HISTORY,
  MENU_RESOURCE_MANAGEMENT
} from '@/dictionary/menu-symbol'
import { OPERATION } from '@/dictionary/iam-auth'

export default [{
  name: MENU_RESOURCE_BUSINESS,
  path: 'business',
  component: () => import('./index.vue'),
  meta: new Meta({
    menu: {
      i18n: '业务',
      relative: MENU_RESOURCE_MANAGEMENT
    },
    layout: {}
  })
}, {
  name: MENU_RESOURCE_BUSINESS_DETAILS,
  path: 'business/details/:bizId',
  component: () => import('./details.vue'),
  meta: new Meta({
    menu: {
      i18n: '业务详情',
      relative: MENU_RESOURCE_BUSINESS
    },
    layout: {}
  })
}, {
  name: MENU_RESOURCE_BUSINESS_HISTORY,
  path: 'business/history',
  component: () => import('./archived.vue'),
  meta: new Meta({
    menu: {
      i18n: '已归档业务',
      relative: MENU_RESOURCE_BUSINESS
    },
    auth: {
      view: { type: OPERATION.BUSINESS_ARCHIVE }
    }
  })
}]
