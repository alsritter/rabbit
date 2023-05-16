package lock

import (
	"context"
	"errors"
	"time"

	"alsritter.icu/rabbit-template/internal/pkg/qylock"
)

var (
	ErrLockAcquire = errors.New("lock acquire failed")
)

type contextKey string

type base struct {
	ctx       context.Context
	redisLock qylock.RedisLockIface
	expire    time.Duration
}

func (b *base) Acquire(key string) error {
	isLock, randVal, err := b.redisLock.SetWithContext(b.ctx, b.genLockKey(key), b.expire)
	if err != nil {
		return err
	}
	b.ctx = context.WithValue(b.ctx, b.genValueKey(key), randVal)
	if !isLock {
		return ErrLockAcquire
	}
	return nil
}

func (b *base) Release(key string) error {
	val := b.ctx.Value(b.genValueKey(key))
	if val == nil {
		return nil
	}
	randVal := val.(string)
	err := b.redisLock.ReleaseWithContext(b.ctx, b.genLockKey(key), randVal)
	if err != nil {
		return err
	}
	return nil
}

func (b *base) genLockKey(key string) string {
	return key
}

func (b *base) genValueKey(key string) contextKey {
	return contextKey(b.genLockKey(key) + "_randVal")
}
