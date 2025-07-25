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

// Package reflector TODO
package reflector

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/common/blog"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
)

// Interface TODO
type Interface interface {
	Lister(ctx context.Context, opts *types.ListOptions) (ch chan *types.Event, err error)
	Watcher(ctx context.Context, opts *types.WatchOptions, cap *Capable) error
	ListWatcher(ctx context.Context, opts *types.ListWatchOptions, cap *Capable) error
}

// NewReflector TODO
func NewReflector(conf local.MongoConf) (Interface, error) {
	s, err := stream.NewStream(conf)
	if err != nil {
		return nil, fmt.Errorf("new stream faiiled, err: %v", err)
	}

	return &Reflector{Stream: s}, nil
}

// Reflector TODO
type Reflector struct {
	Stream stream.Interface
}

// Lister TODO
func (r *Reflector) Lister(ctx context.Context, opts *types.ListOptions) (ch chan *types.Event, err error) {
	return r.Stream.List(ctx, opts)
}

// Watcher TODO
func (r *Reflector) Watcher(ctx context.Context, opts *types.WatchOptions, cap *Capable) error {
	if cap == nil {
		return errors.New("invalid Capable value, must be a pointer and not nil")
	}
	if cap.OnChange.OnAdd == nil || cap.OnChange.OnUpdate == nil || cap.OnChange.OnDelete == nil {
		return errors.New("invalid Capable value")
	}

	if cap.OnChange.OnLister != nil || cap.OnChange.OnListerDone != nil {
		return errors.New("watch can not have OnLister or OnListerDone in Capable")
	}

	watch, err := r.Stream.Watch(ctx, opts)
	if err != nil {
		return err
	}

	go r.loopWatch(opts.Collection, watch, cap)
	return nil
}

// ListWatcher TODO
func (r *Reflector) ListWatcher(ctx context.Context, opts *types.ListWatchOptions, cap *Capable) error {
	if cap == nil {
		return errors.New("invalid Capable value, must be a pointer and not nil")
	}
	if cap.OnChange.OnLister == nil || cap.OnChange.OnAdd == nil ||
		cap.OnChange.OnUpdate == nil || cap.OnChange.OnDelete == nil {
		return errors.New("invalid Capable value")
	}

	watch, err := r.Stream.ListWatch(ctx, opts)
	if err != nil {
		return err
	}

	go r.loopWatch(opts.Collection, watch, cap)
	return nil
}

func (r *Reflector) loopWatch(coll string, w *types.Watcher, cap *Capable) {
	for event := range w.EventChan {
		switch event.OperationType {
		case types.Lister:
			if cap.OnChange.OnLister != nil {
				cap.OnChange.OnLister(event)
			}
		case types.Insert:
			if cap.OnChange.OnAdd != nil {
				cap.OnChange.OnAdd(event)
			}
		case types.Update, types.Replace:
			if cap.OnChange.OnUpdate != nil {
				cap.OnChange.OnUpdate(event)
			}
		case types.Delete:
			if cap.OnChange.OnDelete != nil {
				cap.OnChange.OnDelete(event)
			}
		case types.ListDone:
			// list operation has been done.
			if cap.OnChange.OnListerDone != nil {
				cap.OnChange.OnListerDone()
			}
		case types.Invalidate:
			blog.ErrorJSON("watch collection: %s, received a invalidate event, doc: %s.", coll, event.Document)
		default:
			blog.ErrorJSON("watch collection: %s, received a unsupported event type: %s, doc: %s.", coll, event.OperationType, event.Document)
		}
	}

}
