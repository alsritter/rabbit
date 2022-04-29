package util

import (
	"context"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestConnet(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.16.238.101:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Errorf("connect to etcd failed, err:%v\n", err)
		return
	}

	t.Log("connect to etcd success")
	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "hello", "world")
	cancel()
	if err != nil {
		t.Errorf("put to etcd failed, err:%v\n", err)
		return
	}

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "hello")
	cancel()
	if err != nil {
		t.Errorf("get from etcd failed, err:%v\n", err)
		return
	}

	for _, ev := range resp.Kvs {
		t.Logf("%s:%s\n", ev.Key, ev.Value)
	}
}
