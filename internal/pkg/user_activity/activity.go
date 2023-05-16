package user_activity

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

/**
 * @Doc: 这个包主要用于记录用户的活跃度
 * @Author: zhujl04
 */

type RedisUseLoginMap struct {
	redisPool *redis.Pool
}

func NewRedisUseLoginMap(redisPool *redis.Pool) *RedisUseLoginMap {
	return &RedisUseLoginMap{redisPool: redisPool}
}

func (r *RedisUseLoginMap) Add(domainId, userId int32) error {
	conn := r.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("PFADD", fmt.Sprintf("login:%s:%d", r.getDateKey(time.Now()), domainId), userId)
	return err
}

func (r *RedisUseLoginMap) getDateKey(date time.Time) string {
	return date.Format("20060102")
}

// GetThatDayCount 获取当天活跃用户数
func (r *RedisUseLoginMap) GetThatDayCount(domainId int64) (int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("PFCOUNT", fmt.Sprintf("login:%s:%d", r.getDateKey(time.Now()), domainId)))
}

// GetCountByDate 获取某天活跃用户数
func (r *RedisUseLoginMap) GetCountByDate(domainId int64, date time.Time) (int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("PFCOUNT", fmt.Sprintf("login:%s:%d", r.getDateKey(date), domainId)))
}

// GetCountByDateRange 获取某段时间内活跃用户数
func (r *RedisUseLoginMap) GetCountByDateRange(domainId int64, startDate, endDate time.Time) (int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	var keys []interface{}
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		keys = append(keys, fmt.Sprintf("login:%s:%d", r.getDateKey(date), domainId))
	}
	return redis.Int64(conn.Do("PFCOUNT", keys...))
}

// GetCountByDateRangeWithDate 获取某段时间内活跃用户数，返回map，key为日期，value为活跃用户数
func (r *RedisUseLoginMap) GetCountByDateRangeWithDate(domainId int64, startDate, endDate time.Time) (map[string]int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	var keys []interface{}
	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		keys = append(keys, fmt.Sprintf("login:%s:%d", r.getDateKey(date), domainId))
	}
	counts, err := redis.Int64s(conn.Do("PFCOUNT", keys...))
	if err != nil {
		return nil, err
	}
	result := make(map[string]int64)
	for i, date := 0, startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		result[r.getDateKey(date)] = counts[i]
		i++
	}
	return result, nil
}

// GetCountByMonth 获取某月活跃用户数
func (r *RedisUseLoginMap) GetCountByMonth(domainId int64, month time.Time) (int64, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	var keys []interface{}
	for date := month; date.Month() == month.Month(); date = date.AddDate(0, 0, 1) {
		keys = append(keys, fmt.Sprintf("login:%s:%d", date.Format("20060102"), domainId))
	}
	return redis.Int64(conn.Do("PFCOUNT", keys...))
}

// GetThatMonthCount 获取当月活跃用户数
func (r *RedisUseLoginMap) GetThatMonthCount(domainId int64) (int64, error) {
	return r.GetCountByMonth(domainId, time.Now())
}
