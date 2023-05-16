package data

import (
	"context"
	"time"

	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/pkg/qylogger"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/google/wire"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-kratos/kratos/v2/log"
	grom_tracing "gorm.io/plugin/opentelemetry/tracing"

	"github.com/gomodule/redigo/redis"
	redigotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/garyburd/redigo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewRedisConn, NewDB)

type Data struct {
	db          *gorm.DB
	redisClient *redis.Pool
}

func (d *Data) GetDB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *Data) GetRedisWithCtx(ctx context.Context) redis.Conn {
	conn := d.redisClient.Get()
	return conn
}

func (d *Data) GetRedis(ctx context.Context) *redis.Pool {
	return d.redisClient
}

func (d *Data) CloseRedis() error {
	// d.redisSpan.End()
	return d.redisClient.Close()
}

func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, redisClient *redis.Pool) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, redisClient: redisClient}, cleanup, nil
}

func NewDB(c *conf.Data, logger log.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		// 这里再注册 SQL 直方图
		Logger: qylogger.New(log.NewHelper(logger),
			qylogger.WithSeconds(prom.NewHistogram(_sqlHistogramTracing))),
		PrepareStmt: false,
	})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.Use(grom_tracing.NewPlugin()); err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetConnMaxLifetime(time.Duration(60) * time.Second)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(30)
	return db
}

func NewRedisConn(c *conf.Data) *redis.Pool {
	redisSetting := c.Redis
	if redisSetting == nil {
		panic("redis config is nil")
	}
	if redisSetting.Addr == "" {
		panic("lack of redisSetting.Addr")
	}
	maxIdle := 10
	maxActive := 15
	idleTimeout := 240
	if redisSetting.MaxActive > 0 && redisSetting.MaxIdle > 0 {
		maxIdle = int(redisSetting.MaxIdle)
		maxActive = int(redisSetting.MaxActive)
	}
	if redisSetting.IdleTimeout > 0 {
		idleTimeout = int(redisSetting.IdleTimeout)
	}
	// tracing.WithTracerProvider(trace.NewNoopTracerProvider())
	redisPool := &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redigotrace.Dial("tcp", redisSetting.Addr,
				redigotrace.WithServiceName("my-redis-backend"),
				redis.DialKeepAlive(time.Minute),
			)
			if err != nil {
				return nil, err
			}
			if redisSetting.Passwd != "" {
				if _, err := c.Do("AUTH", redisSetting.Passwd); err != nil {
					c.Close()
					return nil, err
				}
			}
			if redisSetting.Db > 0 {
				if _, err := c.Do("SELECT", redisSetting.Db); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return redisPool
}

// 创建 OSS 客户端实例
// oss 文档: https://help.aliyun.com/document_detail/32143.html
func NewOssClient(c *conf.Data) *oss.Client {
	client, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKeyId, c.Oss.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	return client
}
