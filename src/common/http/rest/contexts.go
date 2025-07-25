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

package rest

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"configcenter/src/ac"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	headerutil "configcenter/src/common/http/header/util"
	"configcenter/src/common/json"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful/v3"
)

// Kit TODO
type Kit struct {
	Rid             string
	Header          http.Header
	Ctx             context.Context
	CCError         errors.DefaultCCErrorIf
	User            string
	SupplierAccount string
}

// Contexts TODO
type Contexts struct {
	Kit            *Kit
	Request        *restful.Request
	resp           *restful.Response
	respStatusCode int
	uri            string
}

// DecodeInto TODO
func (c *Contexts) DecodeInto(to interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Request.Body)
	if err != nil {
		blog.ErrorfDepthf(1, "rid: %s, read request body failed, err: %v", c.Kit.Rid, err)
		return c.Kit.CCError.Error(common.CCErrCommHTTPReadBodyFailed)
	}

	if len(body) == 0 {
		blog.V(9).InfoDepthf(1, "request body is empty, rid: %s", c.Kit.Rid)
		return nil
	}

	if err := json.Unmarshal(body, to); err != nil {
		blog.ErrorfDepthf(1, "rid: %s, unmarshal request body failed, err: %v, body: %s", c.Kit.Rid, err, string(body))
		return c.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	return nil
}

// RespEntity TODO
func (c *Contexts) RespEntity(data interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	c.Response(&metadata.Response{
		BaseResp: metadata.BaseResp{Result: true, Code: common.CCSuccess, ErrMsg: common.CCSuccessStr},
		Data:     data,
	})
}

// RespString response the data format to a json string.
// the data is a string, and do not need marshal, can return directly.
func (c *Contexts) RespString(data *string) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	jsonBuffer := bytes.Buffer{}
	jsonBuffer.WriteString("{\"result\": true, \"bk_error_code\": 0, \"bk_error_msg\": \"success\", \"data\": ")
	if data == nil {
		jsonBuffer.WriteString("")
	} else {
		jsonBuffer.WriteString(*data)
	}
	jsonBuffer.WriteByte('}')
	c.resp.Write(jsonBuffer.Bytes())
}

// RespStringArray response the data format to a json string.
// the data is a string, and do not need marshal, can return directly.
func (c *Contexts) RespStringArray(jsonArray []string) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)

	if len(jsonArray) == 0 {
		jsonBuffer := bytes.Buffer{}
		jsonBuffer.WriteString("{\"result\": true, \"bk_error_code\": 0, \"bk_error_msg\": \"success\", \"data\": []")
		c.resp.Write(jsonBuffer.Bytes())
		return
	}

	last := len(jsonArray) - 1
	jsonBuffer := bytes.Buffer{}
	jsonBuffer.WriteString("{\"result\": true, \"bk_error_code\": 0, \"bk_error_msg\": \"success\", \"data\": ")
	// convert json string to json array format.
	jsonBuffer.WriteByte('[')
	for idx, val := range jsonArray {
		jsonBuffer.WriteString(val)
		if idx != last {
			jsonBuffer.WriteByte(',')
		}
	}
	jsonBuffer.WriteByte(']')
	// end of json
	jsonBuffer.WriteByte('}')
	_, err := c.resp.Write(jsonBuffer.Bytes())
	if err != nil {
		blog.ErrorfDepthf(1, "write response failed, err: %v, rid :%s", err, c.Kit.Rid)
		return
	}
}

// RespCountInfoString TODO
func (c *Contexts) RespCountInfoString(count int64, infoArray []string) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)

	// format data field
	last := len(infoArray) - 1
	dataBuffer := bytes.Buffer{}
	dataBuffer.WriteByte('{')
	dataBuffer.WriteString("\"count\":")
	dataBuffer.WriteString(strconv.FormatInt(count, 10))
	dataBuffer.WriteString(",\"info\":[")
	for idx, val := range infoArray {
		dataBuffer.WriteString(val)
		if idx != last {
			dataBuffer.WriteByte(',')
		}
	}
	dataBuffer.WriteString("]}")

	jsonBuffer := bytes.Buffer{}
	jsonBuffer.WriteString("{\"result\": true, \"bk_error_code\": 0, \"bk_error_msg\": \"success\", \"data\": ")
	jsonBuffer.Write(dataBuffer.Bytes())
	jsonBuffer.WriteByte('}')
	_, err := c.resp.Write(jsonBuffer.Bytes())
	if err != nil {
		blog.ErrorfDepthf(1, "write response failed, err: %v, rid :%s", err, c.Kit.Rid)
		return
	}
}

// RespEntityWithError TODO
func (c *Contexts) RespEntityWithError(data interface{}, err error) {
	c.collectErrorMetric()

	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	resp := metadata.Response{
		Data: data,
	}
	if err != nil {
		blog.ErrorfDepthf(1, "code: %s, user: %s, rid: %s, err: %v", httpheader.GetAppCode(c.Kit.Header), c.Kit.User,
			c.Kit.Rid, err)

		if err == ac.NoAuthorizeError {
			body, err := json.Marshal(data)
			if err != nil {
				blog.ErrorfDepthf(2, "rid: %s, marshal json response failed, err: %v", c.Kit.Rid, err)
				return
			}
			if _, err := c.resp.Write(body); err != nil {
				blog.ErrorfDepthf(2, "rid: %s, response http request failed, err: %v", c.Kit.Rid, err)
				return
			}
			return
		}
		t, yes := err.(errors.CCErrorCoder)
		var code int
		var errMsg string
		if yes {
			code = t.GetCode()
			errMsg = t.Error()
		} else {
			code = common.CCErrorUnknownOrUnrecognizedError
			errMsg = c.Kit.CCError.Error(code).Error()
		}
		resp.BaseResp = metadata.BaseResp{
			Result: false,
			ErrMsg: errMsg,
			Code:   code,
		}
	} else {
		resp.BaseResp = metadata.SuccessBaseResp
	}
	c.Response(&resp)
}

// CountInfo TODO
type CountInfo struct {
	Count int64       `json:"count"`
	Info  interface{} `json:"info"`
}

// RespEntityWithCount TODO
func (c *Contexts) RespEntityWithCount(count int64, info interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	resp := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data: CountInfo{
			Count: count,
			Info:  info,
		},
	}
	c.Response(&resp)
}

// WithStatusCode TODO
func (c *Contexts) WithStatusCode(statusCode int) *Contexts {
	c.respStatusCode = statusCode
	return c
}

// RespWithError TODO
// WriteError is used to write a error response to the http client, which means the request occur an error.
// It receive an err and an optional error code parameter.
// It will testify the err, if the err is a CCErrorCoder, then the error code inside it will be used.
// Otherwise, if errCode is set and > 0, then errCode value is used.
// Finally, if error code is not set and err is not CCErrorCoder, then it will be set with a default
// CCSystemBusy code.
// This function will also write a log when it's called which contains the request id field.
func (c *Contexts) RespWithError(err error, errCode int, format string, args ...interface{}) {
	c.collectErrorMetric()

	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	blog.ErrorfDepthf(1, "code: %s, user: %s, rid: %s, %s, err: %v", httpheader.GetAppCode(c.Kit.Header), c.Kit.User,
		c.Kit.Rid, fmt.Sprintf(format, args), err)

	var code int
	var errMsg string
	if err != nil {
		t, yes := err.(errors.CCErrorCoder)
		if yes {
			code = t.GetCode()
			errMsg = t.Error()
		} else {
			if errCode > 0 {
				code = errCode
				errMsg = c.Kit.CCError.Error(code).Error()
			} else {
				code = common.CCErrorUnknownOrUnrecognizedError
				errMsg = c.Kit.CCError.Error(code).Error()
			}
		}
		// log the error

	} else {
		code = common.CCErrorUnknownOrUnrecognizedError
		errMsg = c.Kit.CCError.Error(code).Error()
	}

	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: errMsg,
			Code:   code,
		},
		Data: nil,
	}

	c.Response(&body)
}

// RespAutoError TODO
func (c *Contexts) RespAutoError(err error) {
	c.collectErrorMetric()

	blog.ErrorfDepthf(1, "code: %s, user: %s, rid: %s, err: %v", httpheader.GetAppCode(c.Kit.Header), c.Kit.User,
		c.Kit.Rid, err)
	var code int
	var errMsg string
	if err != nil {
		t, yes := err.(errors.CCErrorCoder)
		if yes {
			code = t.GetCode()
			errMsg = t.Error()
		} else {
			code = common.CCErrorUnknownOrUnrecognizedError
			errMsg = err.Error()
		}
	} else {
		code = common.CCErrorUnknownOrUnrecognizedError
		errMsg = c.Kit.CCError.Error(code).Error()
	}

	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: errMsg,
			Code:   code,
		},
		Data: nil,
	}

	c.Response(&body)
}

// RespErrorCodeF TODO
// WriteErrorf used to write a error response to the request client.
// it will wrapper the error with error code and other errorf args.
// errorf is used to format multiple-language error message.
// it will also will log the error at the same time with logMsg.
func (c *Contexts) RespErrorCodeF(errCode int, logMsg string, errorf ...interface{}) {
	c.collectErrorMetric()

	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	blog.ErrorfDepthf(1, "code: %s, user: %s, rid: %s, %s", httpheader.GetAppCode(c.Kit.Header), c.Kit.User, c.Kit.Rid,
		logMsg)

	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: c.Kit.CCError.CCErrorf(errCode, errorf).Error(),
			Code:   errCode,
		},
		Data: nil,
	}
	c.Response(&body)
}

// RespErrorCodeOnly TODO
func (c *Contexts) RespErrorCodeOnly(errCode int, format string, args ...interface{}) {
	c.collectErrorMetric()

	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	blog.ErrorfDepthf(1, "code: %s, user: %s, %s, rid: %s", httpheader.GetAppCode(c.Kit.Header), c.Kit.User,
		fmt.Sprintf(format, args), c.Kit.Rid)

	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)
	body := metadata.Response{
		BaseResp: metadata.BaseResp{
			Result: false,
			ErrMsg: c.Kit.CCError.Error(errCode).Error(),
			Code:   errCode,
		},
		Data: nil,
	}

	c.Response(&body)
}

// RespBkEntity TODO
func (c *Contexts) RespBkEntity(data interface{}) {
	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)

	body := &metadata.BKResponse{
		BkBaseResp: metadata.BkBaseResp{
			Code:    common.CCSuccess,
			Message: common.CCSuccessStr,
		},
		Data: data,
	}
	if err := c.resp.WriteAsJson(body); err != nil {
		blog.ErrorfDepthf(1, fmt.Sprintf("rid: %s, response http request failed, err: %v", c.Kit.Rid, err))
		return
	}
}

// RespBkError TODO
func (c *Contexts) RespBkError(errCode int, errMsg string) {
	c.collectErrorMetric()

	blog.ErrorfDepthf(1, "code: %s, user: %s, rid: %s, errCode: %d, errMsg: %s", httpheader.GetAppCode(c.Kit.Header),
		c.Kit.User, c.Kit.Rid, errCode, errMsg)

	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)

	body := &metadata.BkBaseResp{
		Code:    errCode,
		Message: errMsg,
	}
	if err := c.resp.WriteAsJson(body); err != nil {
		blog.ErrorfDepthf(1, fmt.Sprintf("rid: %s, response http request failed, err: %v", c.Kit.Rid, err))
		return
	}
}

// Response TODO
func (c *Contexts) Response(resp *metadata.Response) {
	body, err := json.Marshal(resp)
	if err != nil {
		blog.ErrorfDepthf(2, "marshal json response failed, err: %v, rid: %s", err, c.Kit.Rid)
		return
	}
	if _, err := c.resp.Write(body); err != nil {
		blog.ErrorfDepthf(2, "response http request failed, err: %v, rid: %s", err, c.Kit.Rid)
		return
	}
}

// RespNoAuth respond no auth response
func (c *Contexts) RespNoAuth(resp *metadata.BaseResp) {
	c.collectNoAuthMetric()

	if c.respStatusCode != 0 {
		c.resp.WriteHeader(c.respStatusCode)
	}
	c.resp.Header().Set("Content-Type", "application/json")
	httpheader.AddRid(c.resp.Header(), c.Kit.Rid)

	if err := c.resp.WriteAsJson(resp); err != nil {
		blog.ErrorfDepthf(1, fmt.Sprintf("rid: %s, response http request failed, err: %v", c.Kit.Rid, err))
		return
	}
}

// NewContexts 产生一个新的contexts， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (c *Contexts) NewContexts() *Contexts {
	newHeader := headerutil.CCHeader(c.Kit.Header)
	c.Kit.Header = newHeader
	return &Contexts{
		Kit:            c.Kit,
		Request:        c.Request,
		resp:           c.resp,
		respStatusCode: 0,
	}
}

// NewHeader 产生一个新的header， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (c *Contexts) NewHeader() http.Header {
	return headerutil.CCHeader(c.Kit.Header)
}

// SetReadPreference TODO
func (c *Contexts) SetReadPreference(mode common.ReadPreferenceMode) {
	c.Kit.Ctx, c.Kit.Header = util.SetReadPreference(c.Kit.Ctx, c.Kit.Header, mode)
}

// NewKit 产生一个新的kit， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (kit *Kit) NewKit() *Kit {
	newHeader := headerutil.CCHeader(kit.Header)
	newKit := *kit
	newKit.Header = newHeader
	return &newKit
}

// NewHeader 产生一个新的header， 一般用于在创建新的协程的时候，这个时候会对header 做处理，删除不必要的http header。
func (kit *Kit) NewHeader() http.Header {
	return headerutil.CCHeader(kit.Header)
}

// NewKitFromHeader generate a new kit from http header.
func NewKitFromHeader(header http.Header, errorIf errors.CCErrorIf) *Kit {
	return &Kit{
		Rid:             httpheader.GetRid(header),
		Header:          header,
		Ctx:             util.NewContextFromHTTPHeader(header),
		CCError:         errorIf.CreateDefaultCCErrorIf(httpheader.GetLanguage(header)),
		User:            httpheader.GetUser(header),
		SupplierAccount: httpheader.GetSupplierAccount(header),
	}
}
