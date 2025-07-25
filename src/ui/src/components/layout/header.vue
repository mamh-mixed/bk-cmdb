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
  <header class="header-layout" v-test-id.global="'header'">
    <div class="logo">
      <router-link class="logo-link" to="/index">
        {{appName}}
      </router-link>
    </div>
    <nav class="header-nav" v-test-id.global="'headerNav'">
      <router-link class="header-link"
        v-for="nav in visibleMenu"
        :to="getHeaderLink(nav)"
        :key="nav.id"
        :class="{
          active: isLinkActive(nav)
        }">
        {{$t(nav.i18n)}}
      </router-link>
    </nav>
    <section class="header-info">
      <bk-popover class="info-item"
        theme="light header-info-popover"
        :arrow="false"
        :tippy-options="{
          animateFill: false,
          hideOnClick: false,
          offset: '0, 10'
        }">
        <i :class="`bk-icon icon-${currentSysLang.icon} lang-icon`"></i>
        <div slot="content">
          <div
            v-for="(lang, index) in sysLangs"
            :key="index"
            :class="['link-item', { active: $i18n.locale === lang.id }]"
            @click="handleToggleLang(lang.id)">
            <i :class="`bk-icon icon-${lang.icon} lang-icon`"></i>
            <span>{{lang.name}}</span>
          </div>
        </div>
      </bk-popover>
      <bk-popover class="info-item"
        theme="light header-info-popover"
        animation="fade"
        ref="popover"
        :arrow="false"
        :tippy-options="{
          animateFill: false,
          hideOnClick: false,
          offset: '0, 10'
        }">
        <i class="question-icon icon-cc-help"></i>
        <template slot="content">
          <a class="link-item" target="_blank" :href="helpDocUrl">{{$t('产品文档')}}</a>
          <a class="link-item" target="_blank" @click="handleChangeLog()" style="cursor:pointer">{{$t('版本日志')}}</a>
          <a class="link-item" target="_blank" href="https://bk.tencent.com/s-mart/community">{{$t('问题反馈')}}</a>
          <a class="link-item" target="_blank" href="https://github.com/TencentBlueKing/bk-cmdb">{{$t('开源社区')}}</a>
        </template>
      </bk-popover>
      <bk-popover class="info-item"
        theme="light header-info-popover"
        animation="fade"
        placement="bottom-end"
        :arrow="false"
        :tippy-options="{
          animateFill: false,
          hideOnClick: false,
          offset: '0, 10'
        }">
        <span class="info-user">
          <span class="user-name">{{userName}}</span>
          <i class="user-icon bk-icon icon-angle-down"></i>
        </span>
        <template slot="content">
          <a class="link-item" href="javascript:void(0)"
            @click="handleLogout">
            {{$t('退出登录')}}
          </a>
        </template>
      </bk-popover>
    </section>
    <versionLog
      :current-version="currentVersion"
      :version-list="versionList"
      :show.sync="isShowChangeLogs">
    </versionLog>
  </header>
</template>

<script>
  import has from 'has'
  import menu from '@/dictionary/menu'
  import {
    MENU_BUSINESS,
    MENU_BUSINESS_SET,
    MENU_BUSINESS_SET_TOPOLOGY,
    MENU_BUSINESS_HOST_AND_SERVICE
  } from '@/dictionary/menu-symbol'
  import { mapGetters, mapActions } from 'vuex'
  import {
    getBizSetIdFromStorage,
    getBizSetRecentlyUsed
  } from '@/utils/business-set-helper.js'
  import { changeLocale } from '@/i18n'
  import { LANG_SET } from '@/i18n/constants'
  import { gotoLoginPage } from '@/utils/login-helper'
  import versionLog from '../version-log'
  import logoSvg from '@/assets/images/logo.svg'

  export default {
    components: {
      versionLog
    },
    data() {
      return {
        isShowChangeLogs: false,
        versionList: [],
        currentVersion: '',
        sysLangs: LANG_SET
      }
    },
    computed: {
      ...mapGetters(['userName']),
      ...mapGetters('objectBiz', ['bizId']),
      ...mapGetters('globalConfig', ['config']),
      helpDocUrl() {
        return `${this.$helpDocUrlPrefix}/UserGuide/Introduce/Overview.md`
      },
      visibleMenu() {
        return menu.filter((menuItem) => {
          if (!has(menuItem, 'visibility')) {
            return true
          }

          if (typeof menuItem.visibility === 'function') {
            return menuItem.visibility(this)
          }
          return menuItem.visibility
        })
      },
      isBusinessView() {
        const { matched: [topRoute] } = this.$route
        return topRoute?.name === MENU_BUSINESS_SET || topRoute?.name === MENU_BUSINESS
      },
      currentSysLang() {
        return this.sysLangs.find(lang => lang.id === this.$i18n.locale) || {}
      },
      appName() {
        return this.config.publicConfig?.i18n?.productName ?? this.config.site.name ?? this.$t('蓝鲸配置平台')
      },
      appLogo() {
        const src = this.config.publicConfig.appLogo || logoSvg
        return `url(${src})`
      }
    },
    async mounted() {
      const oldCurrentVersion = localStorage.getItem('newVersion')
      const versionList = await this.getLogList()
      const formatData = versionList.map(item => ({
        title: item.version,
        date: item.time
      }))
      this.versionList = this.$tools.versionSort(formatData, 'title')
      this.currentVersion = versionList.find(item => item.is_current === true)?.version || ''
      if (oldCurrentVersion !== this.currentVersion) {
        this.isShowChangeLogs = true
        localStorage.setItem('newVersion', this.currentVersion)
      }
    },
    methods: {
      ...mapActions('versionLog', [
        'getLogList',
      ]),
      isLinkActive(nav) {
        const { matched: [topRoute] } = this.$route
        if (!topRoute) {
          return false
        }
        return topRoute.name === nav.id
          || (topRoute.name === MENU_BUSINESS_SET && nav.id === MENU_BUSINESS) // 业务集时高亮业务菜单
      },
      getHeaderLink(nav) {
        const link = { name: nav.id }
        if (nav.id === MENU_BUSINESS) {
          const isBizSetUsed = getBizSetRecentlyUsed()
          const id = isBizSetUsed ? getBizSetIdFromStorage() : this.bizId
          const paramName = isBizSetUsed ? 'bizSetId' : 'bizId'
          link.name = isBizSetUsed ? MENU_BUSINESS_SET_TOPOLOGY : MENU_BUSINESS_HOST_AND_SERVICE
          link.params = {
            [paramName]: id
          }
        }
        return link
      },
      handleToggleLang(locale) {
        changeLocale(locale)
      },
      handleLogout() {
        this.$http.post(`${window.API_HOST}logout`, {
          http_scheme: window.location.protocol.replace(':', '')
        }).then((data) => {
          gotoLoginPage(data.url, true)
        })
      },
      handleChangeLog() {
        this.isShowChangeLogs = true
        this.$refs.popover.instance.hide()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .header-layout {
    position: relative;
    display: flex;
    height: 58px;
    background-color: #182132;
    z-index: 1002;
  }
  .logo {
    flex: 292px 0 0;
    font-size: 0;
    .logo-link {
      display: inline-block;
      vertical-align: middle;
      height: 58px;
      line-height: 58px;
      margin-left: 24px;
      padding-left: 44px;
      color: #fff;
      font-size: 16px;
      background: v-bind(appLogo) no-repeat 0 center;
      background-size: 28px;
    }
  }
  .header-nav {
      flex: 3;
      font-size: 0;
      white-space: nowrap;
    .header-link {
        display: inline-block;
        vertical-align: middle;
        height: 58px;
        line-height: 58px;
        padding: 0 25px;
        color: #96A2B9;
        font-size: 14px;
      &:hover {
          background-color: rgba(49, 64, 94, 0.5);
          color: #C2CEE5;
      }
      &.router-link-active,
      &.active {
          background-color: rgba(49, 64, 94, 1);
          color: #fff;
      }
    }
  }
  .header-info {
      flex: 1;
      text-align: right;
      white-space: nowrap;
      font-size: 0;
      @include middleBlockHack;
  }
  .info-item {
      @include inlineBlock;
      margin: 0 12px 0 0;
      text-align: left;
      font-size: 0;
      cursor: pointer;
    &:last-child {
        margin-right: 24px;
    }
    .tippy-active {
        .bk-icon {
            color: #fff;
        }
        .user-icon {
            transform: rotate(-180deg);
        }
    }
    .question-icon,
    .lang-icon {
        font-size: 16px;
        color: #96A2B9;
        width: 32px;
        height: 32px;
        display: flex;
        justify-content: center;
        align-items: center;

        &:hover {
            color: #fff;
            background: linear-gradient(270deg,#253047,#263247);
            border-radius: 100%;
        }
    }
    .info-user {
        font-size: 14px;
        font-weight: bold;
        color: #96A2B9;
        line-height: 32px;
        margin-left: 6px;
        display: inline-block;
      .user-name {
          max-width: 150px;
          @include inlineBlock;
          @include ellipsis;
      }
      .user-icon {
          margin-left: -4px;
          transition: transform .2s linear;
          font-size: 20px;
          color: #96A2B9;
      }
      &:hover {
           color: #fff;
           .user-icon {
              color: #fff;
          }
        }
    }
    .lang-icon {
        font-size: 20px;
    }
  }

  .header-info-popover-theme {
    .link-item {
        display: flex;
        padding: 0 20px;
        height: 40px;
        font-size: 14px;
        white-space: nowrap;
        align-items: center;
        cursor: pointer;

        &:hover,
        &.active {
            background-color: #f1f7ff;
            color: #3a84ff;
        }

        .lang-icon {
            font-size: 16px;
            margin-right: 4px;
        }
    }
  }
</style>

<style>
  .tippy-tooltip.header-info-popover-theme {
      padding-left: 0 !important;
      padding-right: 0 !important;
      overflow: hidden;
      border-radius: 2px !important;
  }
</style>
