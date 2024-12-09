package shorturl

import (
	"context"

	"demo/redis-demo/2.07.short-url/base62"
	"github.com/redis/go-redis/v9"
)

const URLMappingCache = "UrlShorty:mapping_cache"

type RedisShortURLWithCache struct {
	client *redis.Client
}

// NewRedisShortURLWithCache 创建一个新的短链接服务
func NewRedisShortURLWithCache(client *redis.Client) *RedisShortURLWithCache {
	return &RedisShortURLWithCache{client: client}
}

// Shorten 为给定的网址创建并记录一个对应的短网址ID，然后将其返回。
// 如果该网址之前已经创建过相应的短网址ID，那么直接返回之前创建的ID。
func (rsu RedisShortURLWithCache) Shorten(url string) (string, error) {
	ctx := context.Background()
	cachedShortID := rsu.client.HGet(ctx, URLMappingCache, url).Val()
	if len(cachedShortID) > 0 {
		return cachedShortID, nil
	} else {
		originID, err := rsu.client.Incr(context.Background(), URLIDCounter).Result()
		if err != nil {
			return "", err
		}

		// 将原始ID转换为62进制
		shortID := base62.Encode(originID)

		pipeline := rsu.client.TxPipeline()
		pipeline.HSet(context.Background(), URLMappingHash, shortID, url)
		pipeline.HSet(context.Background(), URLMappingCache, url, shortID)
		_, err = pipeline.Exec(ctx)
		if err != nil {
			return "", err
		}

		return shortID, nil
	}
}

// Restore 根据给定的短网址ID找出与之对应的原网址
func (rsu RedisShortURLWithCache) Restore(shortID string) (string, error) {
	url, err := rsu.client.HGet(context.Background(), URLMappingHash, shortID).Result()
	if err != nil {
		return "", err
	}

	return url, nil
}
