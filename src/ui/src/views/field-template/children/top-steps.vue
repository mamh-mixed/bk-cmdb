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

<script setup>
  import { t } from '@/i18n'

  const props = defineProps({
    steps: {
      type: Array,
      default: () => ([
        { title: t('基础信息'), icon: 1 },
        { title: t('字段设置'), icon: 2 },
      ])
    },
    current: {
      type: Number
    },
    width: {
      type: String
    }
  })
</script>

<template>
  <cmdb-steps :steps="props.steps" :current="current" :width="props.width"></cmdb-steps>
</template>
