package shorturl

import (
	"context"

	"demo/redis-demo/2.07.short-url/base62"
	"github.com/redis/go-redis/v9"
)

// URLShortener 短链接服务接口

type URLShortener interface {
	// Shorten 长链接转短链接
	Shorten(url string) (string, error)
	// Restore 短链接还原为长链接
	Restore(shortID string) (string, error)
}

const (
	URLIDCounter   = "UrlShorty:id_counter"
	URLMappingHash = "UrlShorty:mapping_hash"
)

type RedisShortURL struct {
	client *redis.Client
}

// NewRedisShortURL 创建一个新的短链接服务
func NewRedisShortURL(client *redis.Client) *RedisShortURL {
	return &RedisShortURL{client: client}
}

// Shorten 为给定的网址创建并记录一个对应的短网址ID，然后将其返回
func (rsu RedisShortURL) Shorten(url string) (string, error) {
	originID, err := rsu.client.Incr(context.Background(), URLIDCounter).Result()
	if err != nil {
		return "", err
	}

	// 将原始ID转换为62进制
	shortID := base62.Encode(originID)
	err = rsu.client.HSet(context.Background(), URLMappingHash, shortID, url).Err()
	if err != nil {
		return "", err
	}

	return shortID, nil
}

// Restore 根据给定的短网址ID找出与之对应的原网址
func (rsu RedisShortURL) Restore(shortID string) (string, error) {
	url, err := rsu.client.HGet(context.Background(), URLMappingHash, shortID).Result()
	if err != nil {
		return "", err
	}

	return url, nil
}
