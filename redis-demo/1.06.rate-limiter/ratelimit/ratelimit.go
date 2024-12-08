package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"math"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

//go:embed incr_and_ex.lua
var luaScriptIncrAndExpire string

//go:embed key_and_ttl.lua
var luaScriptKeyAndTTL string

type RateLimiter interface {
	// IsAllowed 判断给定用户当前是否可以执行指定行为
	IsAllowed(uid int64) (bool, error)
	// Remaining 返回给定用户当前还可以执行指定行为的次数
	Remaining(uid int64) (int64, error)
	// Duration 返回给定用户下一次可以执行指定行为的等待时间，单位秒
	Duration(uid int64) (int64, error)
	// Revoke 清除给定用户的限流器状态
	Revoke(uid int64) error
}

type RedisRateLimiter struct {
	client   *redis.Client
	action   string // 行为名称
	interval int64  // 限流时间间隔
	maximum  int64  // 限流次数
}

func NewRedisRateLimiter(client *redis.Client, action string, interval, maximum int64) *RedisRateLimiter {
	return &RedisRateLimiter{client: client, action: action, interval: interval, maximum: maximum}
}

func (r *RedisRateLimiter) IsAllowed(uid int64) (bool, error) {
	key := fmt.Sprintf("%s:%d", r.action, uid)

	// 使用 Lua 脚本原子性地增加计数器并设置过期时间
	// 可以使用 ScriptLoad 方法预先加载脚本，然后使用 EvalSha 方法执行脚本，以减少网络开销
	res, err := r.client.Eval(context.Background(), luaScriptIncrAndExpire, []string{key}, r.maximum).Result()
	if err != nil {
		return false, err
	}

	count := cast.ToInt64(res)
	return count <= r.maximum, nil
}

func (r *RedisRateLimiter) Remaining(uid int64) (int64, error) {
	key := fmt.Sprintf("%s:%d", r.action, uid)

	// 获取当前计数器的值
	currentTimesStr := r.client.Get(context.Background(), key).Val()
	currentTimes := cast.ToInt64(currentTimesStr)

	// 如果当前计数器的值大于限流次数，则返回 0
	if currentTimes > r.maximum {
		return 0, nil
	}
	// 否则返回剩余次数
	return r.maximum - currentTimes, nil
}

func (r *RedisRateLimiter) Duration(uid int64) (int64, error) {
	key := fmt.Sprintf("%s:%d", r.action, uid)

	// 使用 Lua 脚本获取当前计数器的值和过期时间
	res, err := r.client.Eval(context.Background(), luaScriptKeyAndTTL, []string{key}, r.maximum, r.interval).Result()
	if err != nil {
		return math.MaxInt64, err
	}

	resSlice := res.([]interface{})
	currentTimes := cast.ToInt64(resSlice[0])
	ttl := cast.ToInt64(resSlice[1])

	// 如果当前计数器的值大于限流次数，则返回过期时间（多久之后可以执行）
	if currentTimes != 0 && currentTimes > r.maximum {
		return ttl, nil
	}
	// 否则返回 0
	return 0, nil
}

func (r *RedisRateLimiter) Revoke(uid int64) error {
	// 删除给定用户的限流器状态
	key := fmt.Sprintf("%s:%d", r.action, uid)
	return r.client.Del(context.Background(), key).Err()
}
