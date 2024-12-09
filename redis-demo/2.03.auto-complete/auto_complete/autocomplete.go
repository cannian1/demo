package auto_complete

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 自动补全

type AutoComplete interface {
	// Feed 向自动补全系统中添加一个新的内容，并分词添加补全建议
	Feed(content string, weight float64) error
	// FeedEx 自动移除冷门建议
	FeedEx(content string, weight float64, ttl time.Duration) error
	// Hint 根据输入的前缀获取补全建议，各个建议按权重从高到低排序
	Hint(segment string, limit int64) ([]string, error)
	// Set 设置补全建议的权重，可以用于从 ES 或其他数据源中同步埋点权重
	Set(content string, weight float64) error
}

type RedisAutoComplete struct {
	client  *redis.Client
	subject string
}

func makeACKey(subject, segment string) string {
	return fmt.Sprintf("AutoComplete:%s:%s", subject, segment)
}

func NewRedisAutoComplete(client *redis.Client, subject string) *RedisAutoComplete {
	return &RedisAutoComplete{client: client, subject: subject}
}

func (r *RedisAutoComplete) Feed(content string, weight float64) error {
	ctx := context.Background()
	tx := r.client.TxPipeline()

	// 为输入的内容创建分词, 并为每个分词添加权重
	for _, segment := range createSegments(content) {
		key := makeACKey(r.subject, segment)
		tx.ZIncrBy(ctx, key, weight, content)
	}
	_, err := tx.Exec(ctx)
	return err
}

func createSegments(content string) []string {
	var segments []string
	runes := []rune(content)

	for i := 1; i <= len(runes); i++ {
		segments = append(segments, string(runes[:i]))
	}

	return segments
}

func (r *RedisAutoComplete) FeedEx(content string, weight float64, ttl time.Duration) error {
	ctx := context.Background()
	tx := r.client.TxPipeline()

	for _, segment := range createSegments(content) {
		key := makeACKey(r.subject, segment)
		tx.ZIncrBy(ctx, key, weight, content)
		tx.Expire(ctx, key, ttl)
	}
	_, err := tx.Exec(ctx)
	return err
}

func (r *RedisAutoComplete) Hint(segment string, limit int64) ([]string, error) {
	ctx := context.Background()
	key := makeACKey(r.subject, segment)
	// 从高到低获取指定数量的建议
	return r.client.ZRevRange(ctx, key, 0, limit-1).Result()
}

func (r *RedisAutoComplete) Set(content string, weight float64) error {
	ctx := context.Background()
	tx := r.client.TxPipeline()

	// 为输入的内容设置提示和权重
	for _, segment := range createSegments(content) {
		key := makeACKey(r.subject, segment)
		tx.ZAdd(ctx, key, redis.Z{Score: weight, Member: content})
	}
	_, err := tx.Exec(ctx)
	return err
}
