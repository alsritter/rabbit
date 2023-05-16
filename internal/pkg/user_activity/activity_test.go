package user_activity

import (
	"testing"
	"time"

	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/data"

	"github.com/gomodule/redigo/redis"
)

func TestRedisUseLoginMap_Add(t *testing.T) {
	type fields struct {
		redisPool *redis.Pool
	}
	type args struct {
		domainId int32
		userId   int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				redisPool: data.NewRedisConn(&conf.Data{
					Redis: &conf.Data_Redis{
						Addr: "dev.alsritter.icu:30379",
					},
				}),
			},
			args: args{
				domainId: 1,
				userId:   3344,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisUseLoginMap{
				redisPool: tt.fields.redisPool,
			}
			if err := r.Add(tt.args.domainId, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("RedisUseLoginMap.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisUseLoginMap_GetThatDayCount(t *testing.T) {
	type fields struct {
		redisPool *redis.Pool
	}
	type args struct {
		domainId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				redisPool: data.NewRedisConn(&conf.Data{
					Redis: &conf.Data_Redis{
						Addr: "dev.alsritter.icu:30379",
					},
				}),
			},
			args: args{
				domainId: 1,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisUseLoginMap{
				redisPool: tt.fields.redisPool,
			}
			got, err := r.GetThatDayCount(tt.args.domainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisUseLoginMap.GetThatDayCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisUseLoginMap.GetThatDayCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisUseLoginMap_GetCountByMonth(t *testing.T) {
	type fields struct {
		redisPool *redis.Pool
	}
	type args struct {
		domainId int64
		month    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				redisPool: data.NewRedisConn(&conf.Data{
					Redis: &conf.Data_Redis{
						Addr: "dev.alsritter.icu:30379",
					},
				}),
			},
			args: args{
				domainId: 1,
				month:    time.Now(),
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisUseLoginMap{
				redisPool: tt.fields.redisPool,
			}
			got, err := r.GetCountByMonth(tt.args.domainId, tt.args.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisUseLoginMap.GetCountByMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisUseLoginMap.GetCountByMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
