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

// Package backbone TODO
package backbone

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	crd "configcenter/src/common/confregdiscover"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metrics"
	"configcenter/src/common/ssl"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/logplatform/opentelemetry"
	"configcenter/src/thirdparty/monitor"

	"github.com/rs/xid"
)

// connect svcManager retry connect time
const maxRetry = 200

// BackboneParameter Used to constrain different services to ensure
// consistency of service startup capabilities
type BackboneParameter struct {
	// ConfigUpdate handle process config change
	ConfigUpdate cc.ProcHandlerFunc
	ExtraUpdate  cc.ProcHandlerFunc
	// config path
	ConfigPath string
	// http server parameter
	SrvInfo *types.ServerInfo
	SrvRegdiscv
}

func newSvcManagerClient(ctx context.Context, svcManagerAddr string,
	tlsConfig *ssl.TLSClientConfig) (*zk.ZkClient, error) {
	var err error
	for retry := 0; retry < maxRetry; retry++ {
		client := zk.NewZkClient(svcManagerAddr, 40*time.Second, tlsConfig)
		if err = client.Start(); err != nil {
			blog.Errorf("connect regdiscv [%s] failed: %v", svcManagerAddr, err)
			time.Sleep(time.Second * 2)
			continue
		}

		if err = client.Ping(); err != nil {
			blog.Errorf("connect regdiscv [%s] failed: %v", svcManagerAddr, err)
			time.Sleep(time.Second * 2)
			continue
		}

		return client, nil
	}

	return nil, err
}

func newConfig(ctx context.Context, srvInfo *types.ServerInfo, discovery discovery.DiscoveryInterface,
	apiMachineryConfig *util.APIMachineryConfig) (*Config, error) {

	machinery, err := apimachinery.NewApiMachinery(apiMachineryConfig, discovery)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}
	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, common.GetIdentification(), srvInfo.IP)

	bonC := &Config{
		RegisterPath: regPath,
		RegisterInfo: *srvInfo,
		CoreAPI:      machinery,
	}

	return bonC, nil
}

func newApiMachinery(disc discovery.DiscoveryInterface,
	config *util.APIMachineryConfig) (apimachinery.ClientSetInterface, error) {

	machinery, err := apimachinery.NewApiMachinery(config, disc)
	if err != nil {
		return nil, fmt.Errorf("new api machinery failed, err: %v", err)
	}

	return machinery, nil
}

func validateParameter(input *BackboneParameter) error {
	if !input.Disable && input.Regdiscv == "" {
		return fmt.Errorf("regdiscv can not be empty")
	}
	if input.SrvInfo.IP == "" {
		return fmt.Errorf("addrport ip can not be empty")
	}
	if input.SrvInfo.Port <= 0 || input.SrvInfo.Port > 65535 {
		return fmt.Errorf("addrport port must be 1-65535")
	}
	if input.ConfigUpdate == nil && input.ExtraUpdate == nil {
		return fmt.Errorf("service config change funcation can not be empty")
	}
	// to prevent other components which doesn't set it from failing
	if input.SrvInfo.RegisterIP == "" {
		input.SrvInfo.RegisterIP = input.SrvInfo.IP
	}
	if input.SrvInfo.UUID == "" {
		input.SrvInfo.UUID = xid.New().String()
	}
	return nil
}

// NewBackbone new backbone.
func NewBackbone(ctx context.Context, input *BackboneParameter) (*Engine, error) {
	if err := validateParameter(input); err != nil {
		return nil, err
	}

	metricService := metrics.NewService(metrics.Config{ProcessName: common.GetIdentification(),
		ProcessInstance: input.SrvInfo.Instance()})

	common.SetServerInfo(input.SrvInfo)

	engine, err := New()
	if err != nil {
		return nil, fmt.Errorf("new engine failed, err: %v", err)
	}
	engine.registerPath = getRegisterPath(input.SrvInfo.IP)
	engine.srvInfo = input.SrvInfo
	engine.metric = metricService
	engine.Disable = input.Disable

	handler := &cc.CCHandler{
		// 扩展这个函数， 新加传递错误
		OnProcessUpdate:  input.ConfigUpdate,
		OnExtraUpdate:    input.ExtraUpdate,
		OnLanguageUpdate: engine.onLanguageUpdate,
		OnErrorUpdate:    engine.onErrorUpdate,
		OnMongodbUpdate:  engine.onMongodbUpdate,
		OnRedisUpdate:    engine.onRedisUpdate,
	}

	if !input.Disable {
		client, err := newSvcManagerClient(ctx, input.Regdiscv, input.TLSConfig)
		if err != nil {
			return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
		}
		serviceDiscovery, err := discovery.NewServiceDiscovery(client, input.SrvInfo.Environment)
		if err != nil {
			return nil, fmt.Errorf("connect regdiscv [%s] failed: %v", input.Regdiscv, err)
		}
		disc, err := NewServiceRegister(client)
		if err != nil {
			return nil, fmt.Errorf("new service discover failed, err:%v", err)
		}

		engine.client = client
		engine.discovery = serviceDiscovery
		engine.ServiceManageInterface = serviceDiscovery
		engine.SvcDisc = disc

		// add default configcenter
		zkdisc := crd.NewZkRegDiscover(client)
		configCenter := &cc.ConfigCenter{
			Type:               common.BKDefaultConfigCenter,
			ConfigCenterDetail: zkdisc,
		}
		cc.AddConfigCenter(configCenter)

		tlsConf, err := getTLSConf()
		if err != nil {
			blog.Errorf("get tls config error, err: %v", err)
			return nil, err
		}
		engine.apiMachineryConfig = &util.APIMachineryConfig{
			QPS:       1000,
			Burst:     2000,
			TLSConfig: tlsConf,
		}

		machinery, err := newApiMachinery(serviceDiscovery, engine.apiMachineryConfig)
		if err != nil {
			return nil, err
		}
		engine.CoreAPI = machinery

		if err = handleNotice(ctx, client.Client(), input.SrvInfo.Instance()); err != nil {
			return nil, fmt.Errorf("handle notice failed, err: %v", err)
		}
	}

	// get the real configuration center.
	currentConfigCenter := cc.CurrentConfigCenter()

	if err = cc.NewConfigCenter(ctx, currentConfigCenter, input.ConfigPath, handler); err != nil {
		return nil, fmt.Errorf("new config center failed, err: %v", err)
	}

	if err := monitor.InitMonitor(); err != nil {
		return nil, fmt.Errorf("init monitor failed, err: %v", err)
	}

	if err := opentelemetry.InitOpenTelemetryConfig(); err != nil {
		return nil, fmt.Errorf("init openTelemetry config failed, err: %v", err)
	}

	if err := opentelemetry.InitTracer(ctx); err != nil {
		return nil, fmt.Errorf("init tracer failed, err: %v", err)
	}

	return engine, nil
}

// StartServer TODO
func StartServer(ctx context.Context, cancel context.CancelFunc, e *Engine, HTTPHandler http.Handler,
	pprofEnabled bool) error {
	tlsConf, err := getTLSConf()
	if err != nil {
		blog.Errorf("get tls config error, err: %v", err)
		return err
	}

	if isTLS(tlsConf) {
		e.srvInfo.Scheme = "https"
	}

	e.server = Server{
		ListenAddr:   e.srvInfo.IP,
		ListenPort:   e.srvInfo.Port,
		Handler:      e.Metric().HTTPMiddleware(HTTPHandler),
		TLS:          tlsConf,
		PProfEnabled: pprofEnabled,
	}

	if err := ListenAndServe(e.server, e.SvcDisc, cancel); err != nil {
		return err
	}

	// wait for a while to see if ListenAndServe in goroutine is successful
	// to avoid registering an invalid server address on zk
	time.Sleep(time.Second)

	if e.Disable {
		return nil
	}

	return e.SvcDisc.Register(e.registerPath, *e.srvInfo)
}

// New new engine
func New() (*Engine, error) {
	return &Engine{
		Language: language.NewFromCtx(language.EmptyLanguageSetting),
		CCErr:    errors.NewFromCtx(errors.EmptyErrorsSetting),
		CCCtx:    newCCContext(),
	}, nil
}

// SrvRegdiscv service registration discovery
type SrvRegdiscv struct {
	client                 *zk.ZkClient
	ServiceManageInterface discovery.ServiceManageInterface
	SvcDisc                ServiceRegisterInterface
	discovery              discovery.DiscoveryInterface
	// registerPath the path registered to the Service Discovery Center
	registerPath string
	// service component addr
	Regdiscv string
	// Disable disable service registration discovery
	Disable bool
	// TLS config
	TLSConfig *ssl.TLSClientConfig
}

// Discovery return discovery
func (s *SrvRegdiscv) Discovery() discovery.DiscoveryInterface {
	return s.discovery
}

// ServiceManageClient return service manage client
func (s *SrvRegdiscv) ServiceManageClient() *zk.ZkClient {
	return s.client
}

// Engine TODO
type Engine struct {
	CoreAPI            apimachinery.ClientSetInterface
	apiMachineryConfig *util.APIMachineryConfig
	metric             *metrics.Service
	sync.Mutex
	server   Server
	srvInfo  *types.ServerInfo
	Language language.CCLanguageIf
	CCErr    errors.CCErrorIf
	CCCtx    CCContextInterface
	SrvRegdiscv
}

// ApiMachineryConfig TODO
func (e *Engine) ApiMachineryConfig() *util.APIMachineryConfig {
	return e.apiMachineryConfig
}

// Metric TODO
func (e *Engine) Metric() *metrics.Service {
	return e.metric
}

func (e *Engine) onLanguageUpdate(previous, current map[string]language.LanguageMap) {
	e.Lock()
	defer e.Unlock()
	if e.Language == nil {
		e.Language = language.NewFromCtx(current)
		blog.Infof("load language config success.")
		return
	}
	e.Language.Load(current)
	blog.V(3).Infof("load new language config success.")
}

func (e *Engine) onErrorUpdate(previous, current map[string]errors.ErrorCode) {
	e.Lock()
	defer e.Unlock()
	if e.CCErr == nil {
		e.CCErr = errors.NewFromCtx(current)
		blog.Infof("load error code config success.")
		return
	}
	e.CCErr.Load(current)
	blog.V(3).Infof("load new error config success.")
}

func (e *Engine) onMongodbUpdate(previous, current cc.ProcessConfig) {
	e.Lock()
	defer e.Unlock()
	if err := cc.SetMongodbFromByte(current.ConfigData); err != nil {
		blog.Errorf("parse mongo config failed, err: %s, data: %s", err.Error(), string(current.ConfigData))
	}
}

func (e *Engine) onRedisUpdate(previous, current cc.ProcessConfig) {
	e.Lock()
	defer e.Unlock()
	if err := cc.SetRedisFromByte(current.ConfigData); err != nil {
		blog.Errorf("parse redis config failed, err: %s, data: %s", err.Error(), string(current.ConfigData))
	}
}

// Ping TODO
func (e *Engine) Ping() error {
	if e.SrvRegdiscv.Disable {
		return nil
	}
	return e.SvcDisc.Ping()
}

// WithRedis TODO
func (e *Engine) WithRedis(prefixes ...string) (redis.Config, error) {
	// use default prefix if no prefix is specified, or use the first prefix
	var prefix string
	if len(prefixes) == 0 {
		prefix = "redis"
	} else {
		prefix = prefixes[0]
	}

	return cc.Redis(prefix)
}

// WithMongo TODO
func (e *Engine) WithMongo(prefixes ...string) (mongo.Config, error) {
	var prefix string
	if len(prefixes) == 0 {
		prefix = "mongodb"
	} else {
		prefix = prefixes[0]
	}

	return cc.Mongo(prefix)
}

func getRegisterPath(ip string) string {
	return fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, common.GetIdentification(), ip)
}

// GetSrvInfo get service info
func (e *Engine) GetSrvInfo() *types.ServerInfo {
	return e.srvInfo
}

func getTLSConf() (*ssl.TLSClientConfig, error) {
	config, err := cc.NewTLSClientConfigFromConfig("tls")
	return &config, err
}

func isTLS(config *ssl.TLSClientConfig) bool {
	if config == nil || len(config.CertFile) == 0 || len(config.KeyFile) == 0 {
		return false
	}
	return true
}
