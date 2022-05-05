package etcdconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"alsritter.icu/rabbit/internal/util"
	"github.com/jpillora/backoff"
)

const (
	ROOT            = "/"
	SERVICE         = "service"
	DEFAULT_CLUSTER = "default"
)

var (
	svcConfigMutex sync.Mutex
	svcConfigMap   sync.Map
	tryTime        = time.Now()

	backOff = &backoff.Backoff{
		Min:    time.Millisecond * 50,
		Max:    time.Second * 60,
		Factor: 2,
	}
)

type Config struct {
	ServiceVersion string `json:"service_version"`
	ServicePort    string `json:"service_port"`
	HttpPort       string `json:"http_port"`
	IsSsl          bool   `json:"is_ssl"`
}

type ServiceConfig struct {
	EtcdServerUrl string
	ServerName    string
	Config
}

func NewServiceConfig(etcdServerUrl, serverName string) *ServiceConfig {
	return &ServiceConfig{
		EtcdServerUrl: etcdServerUrl,
		ServerName:    serverName,
	}
}

func (s *ServiceConfig) GetKeyName(serverName string) string {
	return ROOT + SERVICE + "." + serverName + "." + DEFAULT_CLUSTER
}

func (s *ServiceConfig) GetConfig() (*Config, error) {
	cli, err := util.NewEtcd(s.EtcdServerUrl)
	if err != nil {
		return nil, fmt.Errorf("generating etcd client failed to %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := s.GetKeyName(s.ServerName)
	serviceInfo, err := cli.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("etcd client Get key %s failed to %v", key, err)
	}

	var config Config
	if len(serviceInfo.Kvs) > 0 {
		err := json.Unmarshal(serviceInfo.Kvs[0].Value, &config)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal err: %v", err)
		}
	}

	if config.ServicePort == "" {
		return nil, fmt.Errorf("servicePort is empty, key: %s", key)
	}

	return &config, nil
}

func (s *ServiceConfig) WriteConfig(cf Config) error {
	cli, err := util.NewEtcd(s.EtcdServerUrl)
	if err != nil {
		return fmt.Errorf("generating etcd client failed to %v", err)
	}
	key := s.GetKeyName(s.ServerName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfJson, err := s.marshalToString(&cf)
	if err != nil {
		return fmt.Errorf("json.MarshalToString err: %v", err)
	}

	_, err = cli.Put(ctx, key, cfJson)
	if err != nil {
		return fmt.Errorf("cli.Put err: %v", err)
	}

	return nil
}

func (s *ServiceConfig) InitCache() {
	svcConfigMutex.Lock()
	defer svcConfigMutex.Unlock()

	configs, err := s.getConfigs()
	if err != nil {
		return
	}

	for key, cfg := range configs {
		svcConfigMap.Store(key, cfg)
	}
}

func (s *ServiceConfig) GetCacheConfig() (cfg *Config, err error) {
	key := s.GetKeyName(s.ServerName)
	if c, ok := svcConfigMap.Load(key); ok {
		return c.(*Config), nil
	}

	svcConfigMutex.Lock()
	defer svcConfigMutex.Unlock()

	// double check.
	if c, ok := svcConfigMap.Load(key); ok {
		return c.(*Config), nil
	}

	// in backoff time.
	if time.Now().Before(tryTime) {
		return nil, fmt.Errorf("get config err, backoff: %s", tryTime.Format("2006-01-02 15:04:05"))
	}

	cfg, err = s.GetConfig()
	if err != nil {
		tryTime = time.Now().Add(backOff.Duration())
		return nil, err
	}

	svcConfigMap.Store(key, cfg)
	backOff.Reset()

	return cfg, nil
}

func (s *ServiceConfig) getConfigs() (map[string]*Config, error) {
	cli, err := util.NewEtcd(s.EtcdServerUrl)
	if err != nil {
		return nil, fmt.Errorf("generating etcd client failed to %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serviceInfos, err := cli.Get(ctx, "/", nil)
	if err != nil {
		return nil, fmt.Errorf("cli.Get err: %v", err)
	}

	configs := make(map[string]*Config)
	for _, info := range serviceInfos.Kvs {
		if len(info.Value) > 0 {
			index := strings.Index(string(info.Key), ROOT+SERVICE)
			if index == 0 {
				config := &Config{}
				err := json.Unmarshal(info.Value, config)
				if err != nil {
					return nil, fmt.Errorf("json.UnmarshalByte err: %v", err)
				}

				configs[string(info.Key)] = config
			}
		}
	}

	return configs, nil
}

func (s *ServiceConfig) marshalToString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}
