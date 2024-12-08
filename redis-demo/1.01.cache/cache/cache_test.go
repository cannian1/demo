package cache

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Set(t *testing.T) {
	// 创建 Redis mock
	client, mock := redismock.NewClientMock()
	cache := &RedisCache{client: client}

	key := "test-key"
	value := "test-value"
	ttl := 5 * time.Second

	// 模拟 Redis `SET` 操作的成功
	mock.ExpectSet(key, value, ttl).SetVal("OK")

	// 测试 Set 方法
	err := cache.Set(key, value, ttl)
	assert.NoError(t, err)                        // 确保没有错误
	assert.NoError(t, mock.ExpectationsWereMet()) // 确保所有期望都已匹配
}

// TestRedisCache_Get tests the Get method
func TestRedisCache_Get(t *testing.T) {
	// 创建 Redis mock
	client, mock := redismock.NewClientMock()
	cache := &RedisCache{client: client}

	key := "test-key"
	expectedValue := "test-value"

	// 模拟 Redis `GET` 操作的成功
	mock.ExpectGet(key).SetVal(expectedValue)

	// 测试 Get 方法
	value, err := cache.Get(key)
	assert.NoError(t, err)                        // 确保没有错误
	assert.Equal(t, expectedValue, value)         // 确保返回值正确
	assert.NoError(t, mock.ExpectationsWereMet()) // 确保所有期望都已匹配
}

// TestRedisCache_GetNotFound tests the Get method for missing keys
func TestRedisCache_GetNotFound(t *testing.T) {
	// 创建 Redis mock
	client, mock := redismock.NewClientMock()
	cache := &RedisCache{client: client}

	key := "nonexistent-key"

	// 模拟 Redis `GET` 操作返回错误
	mock.ExpectGet(key).RedisNil()

	// 测试 Get 方法
	value, err := cache.Get(key)
	assert.Error(t, err)                          // 确保返回错误
	assert.Empty(t, value)                        // 确保返回值为空
	assert.Equal(t, "redis: nil", err.Error())    // 确保错误类型正确
	assert.NoError(t, mock.ExpectationsWereMet()) // 确保所有期望都已匹配
}
