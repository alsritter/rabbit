package oss_utils

import (
	"fmt"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type ossInstance struct {
	client *oss.Client
}

var singleton *ossInstance
var once sync.Once

func Instance() *ossInstance {
	once.Do(func() {
		endpoint := ""
		accessKeyID := ""
		accessKeySecret := ""

		// 创建OSSClient实例。
		client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
		if err != nil {
			fmt.Println(err)
		}
		singleton = &ossInstance{
			client: client,
		}
	})
	return singleton
}

func (instance *ossInstance) Client() *oss.Client {
	return instance.client
}

func (instance *ossInstance) DefaultBucket() *oss.Bucket {
	// 获取存储空间。
	bucket, err := instance.client.Bucket("")
	if err != nil {
	}
	return bucket
}
