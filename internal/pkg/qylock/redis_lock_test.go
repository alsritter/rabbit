package qylock

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestRedisLock_SetWithContext(t *testing.T) {
	type fields struct {
		redisPool *redis.Pool
	}
	type args struct {
		ctx    context.Context
		key    string
		expire time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				redisPool: _test_redis_client,
			},
			args: args{
				ctx:    context.TODO(),
				key:    "merchant:store:check_vip_state_lock",
				expire: 30 * time.Second,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisLock{
				redisPool: tt.fields.redisPool,
			}
			got, _, err := r.SetWithContext(tt.args.ctx, tt.args.key, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisLock.SetWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisLock.SetWithContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisLock_Set(t *testing.T) {
	type fields struct {
		redisPool *redis.Pool
	}
	type args struct {
		key    string
		expire time.Duration
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantIsLock bool
		wantErr    bool
	}{
		{
			name: "",
			fields: fields{
				redisPool: _test_redis_client,
			},
			args: args{
				key:    "merchant:store:check_vip_state_lock",
				expire: 30 * time.Second,
			},
			wantIsLock: true,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisLock{
				redisPool: tt.fields.redisPool,
			}
			gotIsLock, _, err := r.Set(tt.args.key, tt.args.expire)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisLock.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsLock != tt.wantIsLock {
				t.Errorf("RedisLock.Set() gotIsLock = %v, want %v", gotIsLock, tt.wantIsLock)
			}
		})
	}
}
