package relation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// 社交关系

type Relation interface {
	// Follow 关注某个用户
	Follow(target string) (bool, error)
	// UnFollow 取消关注某个用户
	UnFollow(target string) (bool, error)
	// IsFollowing 检测当前用户是否正在关注目标用户
	IsFollowing(target string) (bool, error)
	// IsFollower 判断当前用户是否被某个用户关注
	IsFollower(target string) (bool, error)
	// IsFriend 判断是否是好友(互相关注)
	IsFriend(target string) (bool, error)
	// FollowingsCount 获取关注数量
	FollowingsCount() (int64, error)
	// FollowersCount 获取粉丝数量
	FollowersCount() (int64, error)
}

type RedisRelation struct {
	client *redis.Client
	user   string
}

func makeFollowingKey(user string) string {
	return fmt.Sprintf("Relation:%s:following", user)
}

func makeFollowerKey(user string) string {
	return fmt.Sprintf("Relation:%s:follower", user)
}

func NewRedisRelation(client *redis.Client, user string) *RedisRelation {
	return &RedisRelation{client: client, user: user}
}

func (rr *RedisRelation) Follow(target string) (bool, error) {
	ctx := context.Background()
	tx := rr.client.TxPipeline()

	currentTime := time.Now().Unix()
	// 关注操作
	tx.ZAdd(ctx, makeFollowingKey(rr.user), redis.Z{Score: cast.ToFloat64(currentTime), Member: target})
	tx.ZAdd(ctx, makeFollowerKey(target), redis.Z{Score: cast.ToFloat64(currentTime), Member: rr.user})

	cmds, err := tx.Exec(ctx)
	if err != nil {
		return false, err
	}
	cmdVal1 := cmds[0].(*redis.IntCmd).Val()
	cmdVal2 := cmds[1].(*redis.IntCmd).Val()

	if cmdVal1 == 0 || cmdVal2 == 0 {
		return false, errors.New("follow failed")
	}
	return true, nil
}

func (rr *RedisRelation) UnFollow(target string) (bool, error) {
	ctx := context.Background()
	tx := rr.client.TxPipeline()

	// 取消关注操作
	tx.ZRem(ctx, makeFollowingKey(rr.user), target)
	tx.ZRem(ctx, makeFollowerKey(target), rr.user)

	cmds, err := tx.Exec(ctx)
	if err != nil {
		return false, err
	}
	cmdVal1 := cmds[0].(*redis.IntCmd).Val()
	cmdVal2 := cmds[1].(*redis.IntCmd).Val()

	if cmdVal1 == 0 || cmdVal2 == 0 {
		return false, errors.New("unfollow failed")
	}
	return true, nil
}

func (rr *RedisRelation) IsFollowing(target string) (bool, error) {
	ctx := context.Background()
	// 注意：ZRank 是从 0 开始的，所以返回 0 说明是第一个，也就是存在，redis-cli 中不存在返回 nil
	// ZScore 时间复杂度更低
	err := rr.client.ZScore(ctx, makeFollowingKey(rr.user), target).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rr *RedisRelation) IsFollower(target string) (bool, error) {
	ctx := context.Background()
	err := rr.client.ZScore(ctx, makeFollowerKey(rr.user), target).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rr *RedisRelation) IsFriend(target string) (bool, error) {
	// 这个应该不需要事务，如果不一致就说明某个时刻不是好友
	isFollowing, err := rr.IsFollowing(target)
	if err != nil {
		return false, err
	}
	isFollower, err := rr.IsFollower(target)
	if err != nil {
		return false, err
	}
	return isFollowing && isFollower, nil
}

func (rr *RedisRelation) FollowingsCount() (int64, error) {
	ctx := context.Background()
	return rr.client.ZCard(ctx, makeFollowingKey(rr.user)).Result()
}

func (rr *RedisRelation) FollowersCount() (int64, error) {
	ctx := context.Background()
	return rr.client.ZCard(ctx, makeFollowerKey(rr.user)).Result()
}
