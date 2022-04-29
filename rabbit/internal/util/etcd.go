package util

import (
	"fmt"
	"strings"
	"time"

	"alsritter.icu/rabbit/internal/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcd(urls string) (*clientv3.Client, error) {
	cfg := clientv3.Config{
		Endpoints:   strings.Split(urls, ","),
		DialTimeout: 10 * time.Second,
	}

	username, password := config.GetEtcdUsername(), config.GetEtcdPassword()
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("client.New err: %v", err)
	}

	return cli, nil
}
