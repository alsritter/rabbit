package util

import (
	"fmt"
	"strings"
	"time"

	"alsritter.icu/rabbit/internal/config"
	"go.etcd.io/etcd/client"
)

func NewEtcd(urls string) (client.Client, error) {
	cfg := client.Config{
		Endpoints:               strings.Split(urls, ","),
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	}

	username, password := config.GetEtcdUsername(), config.GetEtcdPassword()
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	cli, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("client.New err: %v", err)
	}

	return cli, nil
}
