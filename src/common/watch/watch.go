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

// Package watch TODO
package watch

import (
	"encoding/json"
	"errors"
	"fmt"

	"configcenter/src/common/metadata"
)

// WatchEventOptions TODO
type WatchEventOptions struct {
	// event types you want to care, empty means all.
	EventTypes []EventType `json:"bk_event_types"`
	// the fields you only care, if nil, means all.
	Fields []string `json:"bk_fields"`
	// unix seconds time to where you want to watch from.
	// it's like Cursor, but StartFrom and Cursor can not use at the same time.
	StartFrom int64 `json:"bk_start_from"`
	// the cursor you hold previous, means you want to watch event form here.
	Cursor string `json:"bk_cursor"`
	// the resource kind you want to watch
	Resource CursorType       `json:"bk_resource"`
	Filter   WatchEventFilter `json:"bk_filter"`
}

// WatchEventFilter TODO
type WatchEventFilter struct {
	// SubResource the sub resource you want to watch, eg. object ID of the instance resource, watch all if not set
	SubResource string `json:"bk_sub_resource,omitempty"`
	// SubResources is the sub resources you want to watch, NOTE: this is a special parameter for internal use only
	SubResources []string `json:"-"`
}

// Validate watch event options
func (w *WatchEventOptions) Validate(isInner bool) error {
	if len(w.EventTypes) != 0 {
		for _, e := range w.EventTypes {
			if err := e.Validate(); err != nil {
				return err
			}
		}
	}

	if !isInner {
		switch w.Resource {
		case Host, Biz, Set, Module:
			if len(w.Fields) == 0 {
				return fmt.Errorf("%s event must have fields", w.Resource)
			}
		}
	}

	// use either StartFrom or Cursor.
	if w.StartFrom != 0 && len(w.Cursor) != 0 {
		return errors.New("bk_start_from and bk_cursor can not use at the same time")
	}

	if len(w.Filter.SubResource) > 0 || len(w.Filter.SubResources) > 0 {
		switch w.Resource {
		case ObjectBase, MainlineInstance, InstAsst, KubeWorkload:
		default:
			return fmt.Errorf("%s event cannot have sub resource", w.Resource)
		}
	}

	return nil
}

// WatchEventResp watch event response
type WatchEventResp struct {
	metadata.BaseResp `json:",inline"`
	Data              *WatchResp `json:"data"`
}

// WatchResp TODO
type WatchResp struct {
	// watched events or not
	Watched bool                `json:"bk_watched"`
	Events  []*WatchEventDetail `json:"bk_events"`
}

// WatchEventDetail TODO
type WatchEventDetail struct {
	Cursor    string     `json:"bk_cursor"`
	Resource  CursorType `json:"bk_resource"`
	EventType EventType  `json:"bk_event_type"`
	// Default instance is JsonString type
	Detail DetailInterface `json:"bk_detail"`

	// ChainNode is the chain node of this watch event
	// NOTE: this is a special return value for internal use only
	ChainNode *ChainNode `json:"-"`
}

type jsonWatchEventDetail struct {
	Cursor    string          `json:"bk_cursor"`
	Resource  CursorType      `json:"bk_resource"`
	EventType EventType       `json:"bk_event_type"`
	Detail    json.RawMessage `json:"bk_detail"`
}

// UnmarshalJSON TODO
func (w *WatchEventDetail) UnmarshalJSON(data []byte) error {
	watchEventDetail := jsonWatchEventDetail{}
	if err := json.Unmarshal(data, &watchEventDetail); err != nil {
		return err
	}

	w.Cursor = watchEventDetail.Cursor
	w.EventType = watchEventDetail.EventType
	w.Resource = watchEventDetail.Resource

	if watchEventDetail.Detail == nil {
		return nil
	}

	if w.Detail == nil {
		w.Detail = JsonString(watchEventDetail.Detail)
	}
	switch w.Detail.Name() {
	case "JsonString":
		w.Detail = JsonString(watchEventDetail.Detail)
	}
	return nil
}

// DetailInterface TODO
type DetailInterface interface {
	Name() string
}

// JsonString TODO
type JsonString string

// Name TODO
func (j JsonString) Name() string {
	return "JsonString"
}

// MarshalJSON TODO
func (j JsonString) MarshalJSON() ([]byte, error) {
	if j == "" {
		j = "{}"
	}
	return []byte(j), nil
}
