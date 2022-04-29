package etcdconfig

import (
	"context"
	"reflect"
	"testing"

	"alsritter.icu/rabbit/internal/util"
)

func TestServiceConfig_GetConfig(t *testing.T) {
	type fields struct {
		EtcdServerUrl string
		ServerName    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Config
		wantErr bool
	}{
		{
			"test read service config info from etcd",
			fields{
				EtcdServerUrl: "172.16.238.101:2379",
				ServerName:    "rabbit-test-read-server",
			},
			&Config{
				ServiceVersion: "1.0",
				ServicePort:    "6060",
				HttpPort:       "7070",
				IsSsl:          true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ServiceConfig{
				EtcdServerUrl: tt.fields.EtcdServerUrl,
				ServerName:    tt.fields.ServerName,
			}

			cli, _ := util.NewEtcd(c.EtcdServerUrl)
			key := c.GetKeyName(c.ServerName)
			cli.Put(context.Background(), key, `{"service_version":"1.0","service_port":"6060","http_port":"7070","is_ssl":true}`)

			got, err := c.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceConfig.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceConfig.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceConfig_WriteConfig(t *testing.T) {
	type fields struct {
		EtcdServerUrl string
		ServerName    string
	}
	tests := []struct {
		name    string
		fields  fields
		config  Config
		wantErr bool
	}{
		{
			"test write service config info to etcd",
			fields{
				EtcdServerUrl: "172.16.238.101:2379",
				ServerName:    "rabbit-test-write-server",
			},
			Config{
				ServiceVersion: "1.0",
				ServicePort:    "8080",
				HttpPort:       "9090",
				IsSsl:          true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ServiceConfig{
				EtcdServerUrl: tt.fields.EtcdServerUrl,
				ServerName:    tt.fields.ServerName,
			}
			cli, _ := util.NewEtcd(c.EtcdServerUrl)
			key := c.GetKeyName(c.ServerName)
			cli.Delete(context.Background(), key)

			if err := c.WriteConfig(tt.config); (err != nil) != tt.wantErr {
				t.Errorf("ServiceConfig.WriteConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			serviceInfo, _ := cli.Get(context.Background(), key)
			if string(serviceInfo.Kvs[0].Value) != `{"service_version":"1.0","service_port":"8080","http_port":"9090","is_ssl":true}` {
				t.Errorf("write failed")
			}
		})
	}
}
