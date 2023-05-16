package qylogger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metrics"
	"go.elastic.co/apm/module/apmsql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Logger struct {
	Klogger                   *log.Helper
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	op                        *options
}

type options struct {
	// SQL 时间计量
	seconds metrics.Observer
}

// Option is metrics option.
type Option func(*options)

// WithSeconds with seconds histogram.
func WithSeconds(c metrics.Observer) Option {
	return func(o *options) {
		o.seconds = c
	}
}

func New(logger *log.Helper, opts ...Option) Logger {
	op := options{}
	for _, o := range opts {
		o(&op)
	}

	return Logger{
		Klogger:                   logger,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		op:                        &op,
	}
}

func (l Logger) SetAsDefault() {
	gormlogger.Default = l
}

func (l Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Logger{
		Klogger:                   l.Klogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger(ctx).Debugf(str, args...)
}

func (l Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger(ctx).Warnf(str, args...)
}

func (l Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger(ctx).Errorf(str, args...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	sql, rows := fc()
	sql = strings.Replace(sql, "\n", " ", -1)
	sql = strings.Replace(sql, "\t", " ", -1)

	elapsed := time.Since(begin)

	if l.op.seconds != nil {
		l.op.seconds.With(apmsql.QuerySignature(sql)).Observe(elapsed.Seconds())
	}

	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel == gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		logger.Log(log.LevelError,
			"err", fmt.Sprintf("%+v", err),
			// time.Duration 可以转成字符类型，别管代码检查错误
			"elapsed", fmt.Sprintf("%s", elapsed), //nolint:all
			"rows", rows,
			"sql", sql,
		)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel == gormlogger.Warn:
		logger.Log(log.LevelWarn,
			"elapsed", fmt.Sprintf("%s", elapsed), //nolint:all
			"rows", rows,
			"sql", sql,
		)
	case l.LogLevel != gormlogger.Silent:
		logger.Log(log.LevelInfo,
			"elapsed", fmt.Sprintf("%s", elapsed), //nolint:all
			"rows", rows,
			"sql", sql,
		)
	default:
		logger.Log(log.LevelDebug,
			"elapsed", fmt.Sprintf("%s", elapsed), //nolint:all
			"rows", rows,
			"sql", sql,
		)
	}
}

func (l Logger) logger(ctx context.Context) *log.Helper {
	return l.Klogger.WithContext(ctx)
}
