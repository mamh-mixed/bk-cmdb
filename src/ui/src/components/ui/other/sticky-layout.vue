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
  <div class="sticky-layout"
    :class="{ 'has-scrollbar': hasScrollbar }">
    <div class="sticky-header" ref="header" v-if="$slots.header || $scopedSlots.header">
      <slot name="header" v-if="$slots.header"></slot>
      <slot name="header" v-if="$scopedSlots.header" v-bind:sticky="hasScrollbar"></slot>
    </div>
    <div class="sticky-content" ref="content">
      <slot v-bind:sticky="hasScrollbar"></slot>
    </div>
    <div class="sticky-footer" ref="footer" v-if="$slots.footer || $scopedSlots.footer">
      <slot name="footer" v-if="$slots.footer"></slot>
      <slot name="footer" v-if="$scopedSlots.footer" v-bind:sticky="hasScrollbar"></slot>
    </div>
  </div>
</template>

<script>
  import {
    addResizeListener,
    removeResizeListener
  } from '@/utils/resize-events'
  import throttle from 'lodash.throttle'
  export default {
    name: 'cmdb-sticky-layout',
    data() {
      return {
        hasScrollbar: false,
        scheduleResize: throttle(this.handleResize, 300)
      }
    },
    mounted() {
      addResizeListener(this.$refs.content, this.scheduleResize)
      addResizeListener(this.$el, this.scheduleResize)
    },
    beforeDestroy() {
      removeResizeListener(this.$refs.content, this.scheduleResize)
      removeResizeListener(this.$el, this.scheduleResize)
    },
    methods: {
      handleResize() {
        this.hasScrollbar = this.$el.clientHeight < this.$el.scrollHeight
      }
    }
  }
</script>

<style lang="scss" scoped>
    .sticky-layout {
        position: relative;
        &.has-scrollbar {
          .sticky-header,
          .sticky-footer {
            z-index: 2;
          }
          .sticky-content {
            z-index: 1;
          }
        }
        .sticky-header {
            position: sticky;
            top: 0;
            left: 0;
            width: 100%;
            z-index: 1;
        }
        .sticky-content {
            position: relative;
            white-space: normal;
            word-break: break-all;
            z-index: 2;
        }
        .sticky-footer {
            position: sticky;
            bottom: 0;
            left: 0;
            width: 100%;
            z-index: 1;
        }
    }
</style>
