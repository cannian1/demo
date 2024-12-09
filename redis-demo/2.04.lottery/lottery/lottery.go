package lottery

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 抽奖

type Lottery interface {
	// Join 添加参与抽奖的玩家
	Join(player string) error
	// Draw 从抽奖池中抽取指定数量的玩家，如果 isRemove 为 true，则抽取到的玩家将从抽奖池中移除
	Draw(number int64, isRemove bool) ([]string, error)
	// Size 获取抽奖池内玩家的数量
	Size() (int64, error)
}

type RedisLottery struct {
	client *redis.Client
	key    string
}

func NewRedisLottery(client *redis.Client, key string) *RedisLottery {
	return &RedisLottery{client: client, key: key}
}

func (rl *RedisLottery) Join(player string) error {
	ctx := context.Background()
	return rl.client.SAdd(ctx, rl.key, player).Err()
}

func (rl *RedisLottery) Draw(number int64, isRemove bool) ([]string, error) {
	ctx := context.Background()
	if isRemove {
		// 移除并返回集合中的一个或多个随机元素
		// SPOP key [count] 通常用于多个级别的抽奖
		return rl.client.SPopN(ctx, rl.key, number).Result()
	}

	// 仅返回集合中的一个或多个随机元素，不会移除元素，通常用于 m 个参与者中抽 n个 简单单次抽奖
	return rl.client.SRandMemberN(ctx, rl.key, number).Result()
}

func (rl *RedisLottery) Size() (int64, error) {
	ctx := context.Background()
	return rl.client.SCard(ctx, rl.key).Result()
}
