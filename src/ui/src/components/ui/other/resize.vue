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
  <div class="resize-layout" :class="localDirections">
    <slot></slot>
    <template v-if="!disabled">
      <i v-for="(_direction, index) in localDirections"
        :key="index"
        :class="['resize-handler', _direction]"
        :style="getHandlerStyle(_direction)"
        @mousedown.left="handleMousedown($event, _direction)">
      </i>
    </template>
    <i :class="['resize-proxy', state.direction]" ref="resizeProxy"></i>
    <div class="resize-mask" ref="resizeMask"></div>
  </div>
</template>

<script>
  /* eslint-disable max-len */
  export default {
    name: 'cmdb-resize-layout',
    props: {
      direction: {
        default() {
          return ['bottom', 'right']
        },
        validator(val) {
          const validDirections = ['bottom', 'right']
          if (typeof val === 'string') {
            return validDirections.includes(val)
          } if (val instanceof Array) {
            return !val.some(direction => !validDirections.includes(direction))
          }
          return false
        }
      },
      min: {
        default() {
          return {
            bottom: 0,
            right: 0
          }
        },
        validator(val) {
          return ['object', 'number'].includes(typeof val)
        }
      },
      max: {
        default() {
          return {
            bottom: Infinity,
            right: Infinity
          }
        },
        validator(val) {
          return ['object', 'number'].includes(typeof val)
        }
      },
      handlerWidth: {
        type: Number,
        default: 5
      },
      handlerOffset: {
        type: Number,
        default: 0
      },
      disabled: Boolean,
      /**
       * 将 resize 信息保存在本地的唯一标识，如果拥有此 id 将会记住上次 resize 的信息。出现重复的 id 将会互相覆盖信息，请保持此 id 的唯一性。
       */
      storeId: {
        type: String,
        default: ''
      },
      /**
       * resize layout 的宽度，用来手动设置宽度
       */
      width: {
        type: Number,
        default: undefined
      },
    },
    data() {
      return {
        state: {},
        storeState: {
          resizeWidth: '',
          resizeHeight: '',
        },
      }
    },
    computed: {
      localDirections() {
        if (typeof this.direction === 'string') {
          return [this.direction]
        }
        return this.direction
      },
      localMin() {
        const min = {
          bottom: 0,
          right: 0
        }
        if (typeof this.min === 'number') {
          min.bottom = this.min
          min.right = this.min
        } else {
          Object.assign(min, this.min)
        }
        return min
      },
      localMax() {
        const max = {
          bottom: Infinity,
          right: Infinity
        }
        if (typeof this.max === 'number') {
          max.bottom = this.max
          max.right = this.max
        } else {
          Object.assign(max, this.max)
        }
        return max
      }
    },
    watch: {
      width(width) {
        if (this.direction === 'right' && width !== undefined && typeof width === 'number') {
          this.$el.style.width = `${width}px`
        }
      }
    },
    mounted() {
      this.initStoreState()
    },
    methods: {
      getHandlerStyle(direction) {
        const style = {}
        if (direction === 'right') {
          style.width = `${this.handlerWidth}px`
          style.marginLeft = `${this.handlerOffset - this.handlerWidth}px`
        } else {
          style.height = `${this.handlerWidth}px`
          style.marginTop = `${this.handlerOffset - this.handlerWidth}px`
        }
        return style
      },
      handleMousedown(event, direction) {
        const $handler = event.currentTarget
        const handlerRect = $handler.getBoundingClientRect()
        const $container = this.$el
        const containerRect = $container.getBoundingClientRect()
        const $resizeProxy = this.$refs.resizeProxy
        const $resizeMask = this.$refs.resizeMask
        $resizeProxy.style.visibility = 'visible'
        $resizeMask.style.display = 'block'
        if (direction === 'right') {
          this.state = {
            direction,
            startMouseLeft: event.clientX,
            startLeft: handlerRect.right - containerRect.left
          }
          $resizeProxy.style.top = 0
          $resizeProxy.style.left = `${this.state.startLeft}px`
          $resizeMask.style.cursor = 'col-resize'
        } else {
          this.state = {
            direction,
            startMouseTop: event.clientY,
            startTop: handlerRect.bottom - containerRect.top
          }
          $resizeProxy.style.left = 0
          $resizeProxy.style.top = `${this.state.startTop}px`
          $resizeMask.style.cursor = 'row-resize'
        }
        document.onselectstart = () => false
        document.ondragstart = () => false
        const handleMouseMove = (event) => {
          if (direction === 'right') {
            const deltaLeft = event.clientX - this.state.startMouseLeft
            const proxyLeft = this.state.startLeft + deltaLeft
            const maxLeft = this.localMax.right
            const minLeft = this.localMin.right
            $resizeProxy.style.left = `${Math.min(maxLeft, Math.max(minLeft, proxyLeft)) + this.handlerOffset}px`
          } else {
            const deltaTop = event.clientY - this.state.startMouseTop
            const proxyTop = this.state.startTop + deltaTop
            const maxTop = this.localMax.bottom
            const minTop = this.localMin.bottom
            $resizeProxy.style.top = `${Math.min(maxTop, Math.max(minTop, proxyTop)) + this.handlerOffset}px`
          }
        }
        const handleMouseUp = () => {
          if (direction === 'right') {
            const finalLeft = parseInt($resizeProxy.style.left, 10)
            this.$el.style.width = `${finalLeft}px`
            this.storeState.resizeWidth = `${finalLeft}px`
          } else {
            const finalTop = parseInt($resizeProxy.style.top, 10)
            this.$el.style.height = `${finalTop}px`
            this.storeState.resizeHeight = `${finalTop}px`
          }
          $resizeProxy.style.visibility = 'hidden'
          this.$refs.resizeMask.style.display = 'none'
          document.removeEventListener('mousemove', handleMouseMove)
          document.removeEventListener('mouseup', handleMouseUp)
          document.onselectstart = null
          document.ondragstart = null
          this.setStoreState()
        }
        document.addEventListener('mousemove', handleMouseMove)
        document.addEventListener('mouseup', handleMouseUp)
      },
      setStoreState() {
        if (this.storeId) {
          localStorage.setItem(`${this.storeId}-resizeHeight`, this.storeState.resizeHeight)
          localStorage.setItem(`${this.storeId}-resizeWidth`, this.storeState.resizeWidth)
        }
      },
      initStoreState() {
        if (this.storeId) {
          const storedResizeHeight = localStorage.getItem(`${this.storeId}-resizeHeight`)
          const storedResizeWidth = localStorage.getItem(`${this.storeId}-resizeWidth`)


          if (storedResizeWidth && this.direction === 'right') {
            this.$el.style.width = storedResizeWidth
            this.$refs.resizeProxy.style.left = storedResizeHeight
          }

          if (storedResizeHeight && this.direction === 'bottom') {
            this.$el.style.height = storedResizeHeight
            this.$refs.resizeProxy.style.top = storedResizeHeight
          }
        }
      },
    }
  }
</script>

<style lang="scss" scoped>
    .resize-layout {
        position: relative;
    }
    .resize-handler {
        position: absolute;
        background-color: transparent;
        &.right {
            top: 0;
            left: 100%;
            width: 5px;
            height: 100%;
            cursor: col-resize;
            &:hover {
                background-image: linear-gradient(to right, transparent, transparent 2px, $primaryColor 2px, $primaryColor 3px, transparent 3px, transparent);
            }
        }
        &.bottom {
            top: 100%;
            left: 0;
            width: 100%;
            height: 5px;
            cursor: row-resize;
            &:hover {
                background-image: linear-gradient(to bottom, transparent, transparent 2px, $primaryColor 2px, $primaryColor 3px, transparent 3px, transparent);
            }
        }
    }
    .resize-proxy{
        visibility: hidden;
        position: absolute;
        pointer-events: none;
        z-index: 9999;
        &.right {
            top: 0;
            height: 100%;
            border-left: 1px dashed $primaryColor;
        }
        &.bottom {
            left: 0;
            width: 100%;
            border-top: 1px dashed $primaryColor;
        }
    }
    .resize-mask {
        display: none;
        position: fixed;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        z-index: 9999;
    }
</style>
