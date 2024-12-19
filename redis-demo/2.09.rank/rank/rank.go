package rank

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// Ranker 排行榜

type Ranker interface {
	// SetWeight 设置/更新某个项的权重，true 表示设置成功，false 表示更新成功
	SetWeight(item string, weight int) (bool, error)
	// GetWeight 获取某个项的权重
	GetWeight(item string) (int64, error)
	// UpdateWeight 更新某个项的权重，true 表示更新成功，false 表示不存在该项
	UpdateWeight(item string, weight int) (bool, error)
	// Remove 移除某个项，true 表示移除成功，false 表示不存在该项
	Remove(item string) (bool, error)
	// Length 获取排行榜的长度，即排行榜中的项数
	Length() (int64, error)
	// Top 以降序的方式获取排行榜中排名靠前的若干项
	Top(n int) ([]string, error)
	// Bottom 以升序的方式获取排行榜中排名靠后的若干项
	Bottom(n int) ([]string, error)
}

type Rank struct {
	client redis.Cmdable
	key    string
}

// NewRank 创建一个新的排行榜
func NewRank(client redis.Cmdable, key string) *Rank {
	return &Rank{client: client, key: key}
}

// SetWeight 设置/更新某个项的权重，true 表示设置成功，false 表示更新成功
func (r Rank) SetWeight(item string, weight int) (bool, error) {
	ctx := context.Background()
	res, err := r.client.ZAdd(ctx, r.key, redis.Z{Score: float64(weight), Member: item}).Result()
	return res > 0, err
}

// GetWeight 获取某个项的权重
func (r Rank) GetWeight(item string) (int64, error) {
	ctx := context.Background()
	res, err := r.client.ZScore(ctx, r.key, item).Result()
	return cast.ToInt64(res), err
}

// UpdateWeight 更新某个项的权重，true 表示更新成功，false 表示不存在该项
func (r Rank) UpdateWeight(item string, weight int) (bool, error) {
	ctx := context.Background()
	res, err := r.client.ZIncrBy(ctx, r.key, cast.ToFloat64(weight), item).Result()
	return res > 0, err
}

// Remove 移除某个项，true 表示移除成功，false 表示不存在该项
func (r Rank) Remove(item string) (bool, error) {
	ctx := context.Background()
	res, err := r.client.ZRem(ctx, r.key, item).Result()
	return res > 0, err
}

// Length 获取排行榜的长度，即排行榜中的项数
func (r Rank) Length() (int64, error) {
	ctx := context.Background()
	return r.client.ZCard(ctx, r.key).Result()
}

// Top 以降序的方式获取排行榜中排名靠前的若干项
func (r Rank) Top(n int) ([]string, error) {
	ctx := context.Background()
	return r.client.ZRevRange(ctx, r.key, 0, int64(n-1)).Result()
}

// Bottom 以升序的方式获取排行榜中排名靠后的若干项
func (r Rank) Bottom(n int) ([]string, error) {
	ctx := context.Background()
	return r.client.ZRange(ctx, r.key, 0, int64(n-1)).Result()
}
