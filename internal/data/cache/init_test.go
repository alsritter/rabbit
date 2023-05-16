package cache

import (
	"fmt"

	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/gomodule/redigo/redis"
)

var (
	_test_db           *data.Data
	_test_redis_client *redis.Pool
	_test_log          = log.NewHelper(log.DefaultLogger)
)

func init() {
	fmt.Println("================================ 初始化测试 DB ================================")

	confData := &conf.Data{
		Database: &conf.Data_Database{
			Dsn: "alsritter:xxxxxxxxxxxxxxx@(dev.alsritter.icu:32131)/rabbit_dev?charset=utf8mb4&parseTime=True&loc=Local&timeout=60s",
		},
	}

	_test_redis_client = data.NewRedisConn(&conf.Data{Redis: &conf.Data_Redis{
		Passwd: "xxxxxxxxxxxxxxx",
		Db:     0,
		Addr:   "dev.alsritter.icu:31129",
	}})

	db := data.NewDB(confData, log.DefaultLogger)
	var err error
	_test_db, _, err = data.NewData(confData, log.DefaultLogger, db, _test_redis_client)
	if err != nil {
		fmt.Printf("init db error = %v", err)
	}

}
