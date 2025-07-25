/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package logics

import (
	"net/http"

	"configcenter/src/ac/extensions"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
	"configcenter/src/thirdparty/esbserver"
)

// Logics TODO
type Logics struct {
	*backbone.Engine
	esbServ     esbserver.EsbClientInterface
	ErrHandle   errors.DefaultCCErrorIf
	header      http.Header
	rid         string
	ownerID     string
	user        string
	ccErr       errors.DefaultCCErrorIf
	ccLang      language.DefaultCCLanguageIf
	AuthManager *extensions.AuthManager
	timerSpec   string
}

// NewLogics get logics handle
func NewLogics(b *backbone.Engine, header http.Header, authManager *extensions.AuthManager, spec string) *Logics {
	lang := httpheader.GetLanguage(header)
	return &Logics{
		Engine:      b,
		header:      header,
		rid:         httpheader.GetRid(header),
		ccErr:       b.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:      b.Language.CreateDefaultCCLanguageIf(lang),
		user:        httpheader.GetUser(header),
		ownerID:     httpheader.GetSupplierAccount(header),
		AuthManager: authManager,
		timerSpec:   spec,
	}
}

// NewFromHeader new Logic from header
func (lgc *Logics) NewFromHeader(header http.Header) *Logics {
	lang := httpheader.GetLanguage(header)
	rid := httpheader.GetRid(header)
	if rid == "" {
		if lgc.rid == "" {
			rid = util.GenerateRID()
		} else {
			rid = lgc.rid
		}
		httpheader.SetRid(header, rid)
	}
	newLgc := &Logics{
		header:    header,
		Engine:    lgc.Engine,
		rid:       rid,
		esbServ:   lgc.esbServ,
		user:      httpheader.GetUser(header),
		ownerID:   httpheader.GetSupplierAccount(header),
		timerSpec: lgc.timerSpec,
	}
	// if language not exist, use old language
	if lang == "" {
		newLgc.ccErr = lgc.ccErr
		newLgc.ccLang = lgc.ccLang
	} else {
		newLgc.ccErr = lgc.CCErr.CreateDefaultCCErrorIf(lang)
		newLgc.ccLang = lgc.Language.CreateDefaultCCLanguageIf(lang)
	}
	return newLgc
}

// Logic TODO
type Logic struct {
	*backbone.Engine
}
