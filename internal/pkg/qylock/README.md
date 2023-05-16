## lock模块
### 说明
锁模块基于单实例版redis实现分布式锁，参考实现：https://redis.io/topics/distlock

### 代码示例（支持单测的写法）

```go
type ServiceA struct {
    RedisLock lock.RedisLockIface
}

// 创建对象
s := ServiceA {
    RedisLock: lock.RedisLock{},
}

// 加锁
// key命名推荐：服务名称:模块名称:其他信息
key := "project:goods:collect:1"
isLock, randVal, err := s.RedisLock.Set(stark.RedisConn, key, 5)
if err != nil {
    // 错误逻辑判断
}
if isLock {
    // 加锁成功逻辑处理
}

// 释放锁
err := s.RedisLock.Release(stark.RedisConn, key, randVal)
if err != nil {
    // 错误逻辑判断
}
```
