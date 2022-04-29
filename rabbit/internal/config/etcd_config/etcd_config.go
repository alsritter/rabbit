package etcdconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"alsritter.icu/rabbit/internal/util"
)

const (
	ROOT            = "/"
	SERVICE         = "service"
	DEFAULT_CLUSTER = "default"
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

func (c *ServiceConfig) GetConfig() (*Config, error) {
	cli, err := util.NewEtcd(c.EtcdServerUrl)
	if err != nil {
		return nil, fmt.Errorf("generating etcd client failed to %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	key := c.GetKeyName(c.ServerName)
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

func (c *ServiceConfig) WriteConfig(cf Config) error {
	cli, err := util.NewEtcd(c.EtcdServerUrl)
	if err != nil {
		return fmt.Errorf("generating etcd client failed to %v", err)
	}
	key := c.GetKeyName(c.ServerName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfJson, err := marshalToString(&cf)
	if err != nil {
		return fmt.Errorf("json.MarshalToString err: %v", err)
	}

	_, err = cli.Put(ctx, key, cfJson)
	if err != nil {
		return fmt.Errorf("cli.Put err: %v", err)
	}

	return nil
}

func (c *ServiceConfig) GetKeyName(serverName string) string {
	return ROOT + SERVICE + "." + serverName + "." + DEFAULT_CLUSTER
}

func marshalToString(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}
