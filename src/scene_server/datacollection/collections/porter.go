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

package collections

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/redis"

	goredis "github.com/go-redis/redis/v7"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/tidwall/gjson"
)

const (
	// metricsNamespacePrefix is prefix of metrics namespace.
	metricsNamespacePrefix = "cmdb_collector"

	// defaultMessageChanTimeout is default timeout of analyzer message channel.
	defaultMessageChanTimeout = time.Second

	// defaultMessageChanSize is default size of analyzer message channel, 1MB * 5000.
	defaultMessageChanSize = 5000

	// defaultFusingCheckInterval is default internal for fusing.
	defaultFusingCheckInterval = 30 * time.Second

	// defaultFusingThresholdPercent is default fusing threshold percent.
	defaultFusingThresholdPercent = 90

	// defaultFusingPercent is default fusing percent.
	defaultFusingPercent = 50

	// defaultReSubscribeWaitDuration is default wait duration before re-subscribe.
	defaultReSubscribeWaitDuration = time.Second

	// defaultDebugInterval is default internal for debuging.
	defaultDebugInterval = 10 * time.Second

	// needInternalDebug is flag for internal debug.
	needInternalDebug = false

	// defaultHostSnapDelayQueueBatchNum by default, the maximum number of members fetched from the delayed queue at
	// one time.
	defaultHostSnapDelayQueueBatchNum = 200
)

// SimplePorter is simple porter handles message from collectors.
// You can impls your own Porter base on Porter and Analyzer interfaces.
type SimplePorter struct {
	engine *backbone.Engine

	// name of this porter.
	name string

	// hash collections hash object, that updates target nodes in dynamic mode,
	// and calculates node base on hash key of data.
	hash *Hash

	// analyzer analyzes message from collectors.
	analyzer Analyzer

	// msgChan is message channel that analyzer consumes from.
	msgChan chan *string

	// message channel redis, read collector data
	// from it by subscribe target topics.
	redisCli redis.Client

	// collectors message channels redis topics.
	topics []string

	// metrics.
	// receiveTotal is message received total stat.
	receiveTotal prometheus.Counter

	// receiveInvalidTotal is message received but it's not valid data total stat.
	receiveInvalidTotal prometheus.Counter

	// receiveShardingTotal is message received and with suitable sharding total stat.
	receiveShardingTotal prometheus.Counter

	// receiveTimeoutTotal is message received but timeout to add to analyze channel total stat.
	receiveTimeoutTotal prometheus.Counter

	// analyzeTotal is message analyzed total stat with target labels.
	analyzeTotal *prometheus.CounterVec

	// analyzeDuration is analyze cost duration stat.
	analyzeDuration prometheus.Histogram

	// fusingTotal is message channel fused count total stat.
	fusingTotal prometheus.Counter

	// registry is prometheus register.
	registry prometheus.Registerer

	needDebug bool
}

// NewSimplePorter creates a new SimplePorter object.
func NewSimplePorter(name string, engine *backbone.Engine, hash *Hash, analyzer Analyzer,
	redisCli redis.Client, topics []string, registry prometheus.Registerer) *SimplePorter {

	return &SimplePorter{
		name:      name,
		engine:    engine,
		hash:      hash,
		analyzer:  analyzer,
		msgChan:   make(chan *string, defaultMessageChanSize),
		redisCli:  redisCli,
		topics:    topics,
		registry:  registry,
		needDebug: needInternalDebug,
	}
}

// init inits a new simple porter.
func (p *SimplePorter) init() {
	// register metrics.
	p.registerMetrics()
	blog.Infof("SimplePorter[%s]| init metrics success!", p.name)
}

// registerMetrics registers prometheus metrics.
func (p *SimplePorter) registerMetrics() {
	p.registry.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_%s_analyze_queue", metricsNamespacePrefix, p.name),
			Help: "current number of analyze message queue.",
		},
		func() float64 { return float64(len(p.msgChan)) },
	))

	p.receiveTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_receive_total", metricsNamespacePrefix, p.name),
			Help: "total number of received message.",
		},
	)
	p.registry.MustRegister(p.receiveTotal)

	p.receiveInvalidTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_receive_invalid_total", metricsNamespacePrefix, p.name),
			Help: "total number of received invalid message.",
		},
	)
	p.registry.MustRegister(p.receiveInvalidTotal)

	p.receiveShardingTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_receive_sharding_total", metricsNamespacePrefix, p.name),
			Help: "total number of received sharding message.",
		},
	)
	p.registry.MustRegister(p.receiveShardingTotal)

	p.receiveTimeoutTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_receive_timeout_total", metricsNamespacePrefix, p.name),
			Help: "total number of received but sending to analyze queue timeout message.",
		},
	)
	p.registry.MustRegister(p.receiveTimeoutTotal)

	p.analyzeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_analyze_total", metricsNamespacePrefix, p.name),
			Help: "total number of analyzed message.",
		},
		[]string{"status"},
	)
	p.registry.MustRegister(p.analyzeTotal)

	p.analyzeDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_%s_analyze_duration", metricsNamespacePrefix, p.name),
			Help: "analyze duration of each message.",
		},
	)

	p.fusingTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_fusing_total", metricsNamespacePrefix, p.name),
			Help: "total number of received but fusing from analyze queue.",
		},
	)

	p.registry.MustRegister(p.analyzeDuration)
}

// Name returns name of this porter.
func (p *SimplePorter) Name() string {
	return p.name
}

// Mock handles a mock message.
func (p *SimplePorter) Mock() error {
	mock := p.analyzer.Mock()
	if err := p.AddMessage(&mock); err != nil {
		return fmt.Errorf("mock failed, %+v", err)
	}
	return nil
}

// AddMessage TODO
// Mock handles a mock message.
func (p *SimplePorter) AddMessage(message *string) error {
	if message == nil {
		return fmt.Errorf("message nil")
	}

	select {
	case p.msgChan <- message:

	case <-time.After(defaultMessageChanTimeout):
		return fmt.Errorf("channel timeout")
	}
	return nil
}

// handleCollectorMsg processing collected and reported messages.
func (p *SimplePorter) handleCollectorMsg(msg *string, sourceType string) {

	cost := time.Now()
	if _, err := p.analyzer.Analyze(msg, sourceType); err != nil {
		blog.Errorf("SimplePorter[%s], analyze message failed, %+v", p.name, err)
		// metrics stats for analyze failed.
		p.analyzeTotal.WithLabelValues("failed").Inc()
	} else {
		// metrics stats for analyze success.
		p.analyzeTotal.WithLabelValues("success").Inc()
	}

	// metrics stats for analyze duration.
	p.analyzeDuration.Observe(time.Since(cost).Seconds())

}

func (p *SimplePorter) getMonitorMsgFromDelayQueue() {

	blog.Infof("SimplePorter[%s], loop to get host snapshot data in delayed queue now!", p.name)

	for {
		// 1、get 200 message data from the delay queue by default each time for compensation processing.
		msgs, err := p.redisCli.ZRangeWithScores(context.Background(), common.RedisHostSnapMsgDelayQueue, 0,
			defaultHostSnapDelayQueueBatchNum).Result()
		if err != nil {
			blog.Errorf("get members failed in delay queue, err: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var wg sync.WaitGroup

		pipeline := make(chan bool, 10)

		for _, msg := range msgs {
			pipeline <- true
			wg.Add(1)
			go func(msg goredis.Z) {
				defer func() {
					wg.Done()
					<-pipeline
				}()
				// if the timestamp of the data is 10 minutes ago, drop it directly。
				if msg.Score < float64(time.Now().Add(-time.Minute*10).Unix()) {
					err := p.redisCli.ZRem(context.Background(), common.RedisHostSnapMsgDelayQueue, msg.Member).Err()
					if err != nil {
						blog.Errorf("remove member failed, msg: %v,err: %v", msg, err)
						return
					}
					return
				}
				// incorrect data type, discard directly.
				val, ok := msg.Member.(string)
				if !ok {
					blog.Errorf("convert data msg: %v, failed, type: %v", msg, reflect.TypeOf(msg))
					// if the type of element taken from the queue is wrong, delete it directly from the queue
					err := p.redisCli.ZRem(context.Background(), common.RedisHostSnapMsgDelayQueue, msg.Member).Err()
					if err != nil {
						blog.Errorf("remove delay queue member failed, msg: %v, err: %v", msg, err)
					}
					return
				}
				// process legitimate messages in the delay queue
				p.handleCollectorMsg(&val, metadata.HostSnapDataSourcesDelayQueue)
			}(msg)
		}

		wg.Wait()

		time.Sleep(10 * time.Second)
	}
}

// analyzeLoop keeps analyzing message from collectors.
func (p *SimplePorter) analyzeLoop() {
	blog.Infof("SimplePorter[%s]| start a new analyze loop now!", p.name)

	for msg := range p.msgChan {
		p.handleCollectorMsg(msg, metadata.HostSnapDataSourcesChannel)
	}
}

// collectLoop keeps subscribe redis topic and collecting messages from collectors.
func (p *SimplePorter) collectLoop() error {
	for {
		// subscribe target topics and handle message base on the redis pubsun channel.
		subChan := p.redisCli.Subscribe(context.Background(), p.topics...)
		blog.V(4).Infof("SimplePorter[%s]| subscribe topics[%+v] success, receiving message now!", p.name, p.topics)

		// receiving message.
		for {
			// ReceiveMessage returns a Message or error ignoring Subscription or Pong
			// messages. It automatically reconnects to Redis Server and resubscribes
			// to topics in case of network errors.
			newMsg, err := subChan.ReceiveMessage()
			if err != nil {
				blog.Errorf("SimplePorter[%s]| receive topics[%+v] message failed, %+v", p.name, p.topics, err)

				// internal errors, unsubscribe and try to sub-recv again.
				subChan.Unsubscribe(p.topics...)
				subChan.Close()
				break
			}
			// metrics stats for message receiving.
			p.receiveTotal.Inc()

			// ignoring invalid payloads.
			if len(newMsg.Payload) == 0 {
				blog.Errorf("SimplePorter[%s]| recved a message with empty payload!", p.name)

				// metrics stats for invalid message.
				p.receiveInvalidTotal.Inc()
				continue
			}
			// the outermost "ip" field is GSE's control ip. In the gse2.0 scenario, there is no concept of control ip.
			// at this time, the "ip" field is empty, so the hashkey needs to be combined in order to be compatible with
			//different scenarios ip and bk_agent_id
			ip := gjson.Get(newMsg.Payload, "ip").String()
			agentID := gjson.Get(newMsg.Payload, "bk_agent_id").String()
			// message data sharding hashring check.
			hashKey, err := p.analyzer.Hash(gjson.Get(newMsg.Payload, "cloudid").String(), ip+agentID)
			if err != nil {
				blog.Errorf("SimplePorter[%s]| calculates message hash key failed, %+v", p.name, err)

				// metrics stats for invalid message.
				p.receiveInvalidTotal.Inc()
				continue
			}

			if !p.hash.IsMatch(hashKey) {
				// ignore message.
				continue
			}

			// metrics stats for suitable sharding message.
			p.receiveShardingTotal.Inc()

			if err := p.AddMessage(&newMsg.Payload); err != nil {
				blog.Errorf("SimplePorter[%s]| add message to analyze, %+v", p.name, err)

				// metrics stats for message sending timeout.
				p.receiveTimeoutTotal.Inc()
			}

			// end of once message receiving loop.
			continue
		}
	}

	// should no-reach.
	return nil
}

// fusing is fuse controller, it would weed out the stacked message in channel,
// in order to keep the newest message could be analyzed in time.
func (p *SimplePorter) fusing() {
	blog.Infof("SimplePorter[%s]| fusing running now!", p.name)

	ticker := time.NewTicker(defaultFusingCheckInterval)
	defer ticker.Stop()

	// is message channel stacked in last interval check.
	isStackedInLastCheck := false

	// keep checking message channel status, and ctrl the fusing.
	for now := range ticker.C {
		stackedN := len(p.msgChan)
		percent := (stackedN * 100) / defaultMessageChanSize

		// check threshold percent.
		if percent < defaultFusingThresholdPercent {
			isStackedInLastCheck = false

			blog.V(4).Infof("SimplePorter[%s]| no-need fusing now, percent[%d] < threshold percent[%d]",
				p.name, percent, defaultFusingThresholdPercent)
			continue
		}

		// stacked in current check.
		if !isStackedInLastCheck {
			isStackedInLastCheck = true

			blog.V(4).Infof("SimplePorter[%s]| no-need fusing now, percent[%d] > threshold percent[%d] first time!",
				p.name, percent, defaultFusingThresholdPercent)
			continue
		}

		// need to fuse stacked message channel.
		fuseCount := 0
		fuseMaxCount := (defaultMessageChanSize * defaultFusingPercent) / 100
		needEarlyStop := false

		blog.Warnf("SimplePorter[%s]| time[%+v] fusing now! stackedNum[%d] chanSize[%d] threshold[%d] interval[%+v] "+
			"fuseMaxCount[%d]", p.name, now, stackedN, defaultMessageChanSize, defaultFusingThresholdPercent,
			defaultFusingCheckInterval, fuseMaxCount)

		for (fuseCount < fuseMaxCount) && !needEarlyStop {
			select {
			case <-p.msgChan:
				fuseCount++
				p.fusingTotal.Inc()

			case <-time.After(defaultMessageChanTimeout):
				// no more stacked data, mark early stop flag, DO NOT use 'break' here.
				needEarlyStop = true
			}
		}
		blog.Warnf("SimplePorter[%s]| time[%+v] fusing done! fuseCount[%d] needEarlyStop[%+v] cost[%+v]",
			p.name, now, fuseCount, needEarlyStop, time.Since(now).Seconds())

		// stacked message cleaned in this round.
		isStackedInLastCheck = false
	}
}

// Run runs the porter.
func (p *SimplePorter) Run() error {
	// init new simple porter.
	p.init()

	// setups analyze goroutines.
	for i := 0; i < runtime.NumCPU(); i++ {
		go p.analyzeLoop()
	}

	// periodically obtain the data reported by the host from the delay queue for compensation update.
	go p.getMonitorMsgFromDelayQueue()

	// fuse controller.
	go p.fusing()

	// internal debug infos.
	go p.debug()

	// NOTE: keep collecting message here.
	p.collectLoop()

	return nil
}

// debug stats and prints internal debug infos in duration.
func (p *SimplePorter) debug() {
	if !p.needDebug {
		return
	}

	ticker := time.NewTicker(defaultDebugInterval)
	defer ticker.Stop()

	for now := range ticker.C {
		msgChanStage := len(p.msgChan)

		// debug recvTotalNum.
		pbmetric := &dto.Metric{}
		p.receiveTotal.Write(pbmetric)
		recvTotalNum := pbmetric.GetCounter().GetValue()

		// debug receiveInvalidTotalNum.
		pbmetric = &dto.Metric{}
		p.receiveInvalidTotal.Write(pbmetric)
		receiveInvalidTotalNum := pbmetric.GetCounter().GetValue()

		// debug receiveShardingTotalNum.
		pbmetric = &dto.Metric{}
		p.receiveShardingTotal.Write(pbmetric)
		receiveShardingTotalNum := pbmetric.GetCounter().GetValue()

		// debug receiveTimeoutTotalNum.
		pbmetric = &dto.Metric{}
		p.receiveTimeoutTotal.Write(pbmetric)
		receiveTimeoutTotalNum := pbmetric.GetCounter().GetValue()

		// debug analyzeTotalFailedNum.
		pbmetric = &dto.Metric{}
		analyzeTotalFailedNum := float64(0)
		if counter, err := p.analyzeTotal.GetMetricWithLabelValues("failed"); err == nil {
			counter.Write(pbmetric)
			analyzeTotalFailedNum = pbmetric.GetCounter().GetValue()
		}

		// debug analyzeTotalSuccNum.
		pbmetric = &dto.Metric{}
		analyzeTotalSuccNum := float64(0)
		if counter, err := p.analyzeTotal.GetMetricWithLabelValues("success"); err == nil {
			counter.Write(pbmetric)
			analyzeTotalSuccNum = pbmetric.GetCounter().GetValue()
		}

		// debug analyzeDuration.
		pbmetric = &dto.Metric{}
		p.analyzeDuration.Write(pbmetric)
		analyzeDuration := float64(0)

		// debug analyzeDuration
		simpleSum := pbmetric.GetHistogram().GetSampleSum()
		simpleCount := pbmetric.GetHistogram().GetSampleCount()
		if simpleCount != 0 {
			analyzeDuration = simpleSum / float64(simpleCount)
		}

		// debug fusing.
		pbmetric = &dto.Metric{}
		p.fusingTotal.Write(pbmetric)
		fusingTotalNum := pbmetric.GetCounter().GetValue()

		// debug infos.
		blog.Infof("SimplePorter[%s]| DEBUG[%+v], msgChanStage[%d], recvTotal[%f] recvInvalid[%f]"+
			"recvSharding[%f] recvTimeout[%f] analyzeFailedTotal[%f] analyzeSuccTotal[%f] analyzeDuration[%f] "+
			"fuseTotal[%f]", p.name, now, msgChanStage, recvTotalNum, receiveInvalidTotalNum, receiveShardingTotalNum,
			receiveTimeoutTotalNum, analyzeTotalFailedNum, analyzeTotalSuccNum, analyzeDuration, fusingTotalNum)
	}
}
