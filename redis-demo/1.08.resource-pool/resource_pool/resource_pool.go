package resource_pool

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 资源池

//go:embed srandmenber_and_smove.lua
var luaScriptSRandMemberAndSMove string

type ResourcePool interface {
	// Associate 关联资源
	Associate(resource string) (bool, error)
	// Disassociate 解除资源关联
	Disassociate(resource string) (bool, error)
	// Acquire 获取资源
	Acquire() (string, error)
	// Release 释放资源
	Release(resource string) (bool, error)
}

type RedisResourcePool struct {
	client          *redis.Client
	availableSetKey string
	OccupiedSetKey  string
}

func NewRedisResourcePool[T any](client *redis.Client, availableSetPoolName, occupiedSetPoolName T) *RedisResourcePool {
	return &RedisResourcePool{
		client:          client,
		availableSetKey: fmt.Sprintf("ResourcePool:%v:available", availableSetPoolName),
		OccupiedSetKey:  fmt.Sprintf("ResourcePool:%v:occupied", occupiedSetPoolName),
	}
}

func (rp *RedisResourcePool) Associate(resource string) (bool, error) {
	ctx := context.Background()
	txf := func(tx *redis.Tx) error {
		// 检查资源是否已存在
		available, err := tx.SIsMember(ctx, rp.availableSetKey, resource).Result()
		if err != nil {
			return err
		}
		occupied, err := tx.SIsMember(ctx, rp.OccupiedSetKey, resource).Result()
		if err != nil {
			return err
		}
		if available || occupied {
			// 资源已存在
			return nil
		}
		// 添加资源到可用集合
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.SAdd(ctx, rp.availableSetKey, resource)
			return nil
		})
		return err
	}

	err := rp.client.Watch(ctx, txf, rp.availableSetKey, rp.OccupiedSetKey)
	if errors.Is(err, redis.TxFailedErr) {
		return false, nil // 事务失败，可能需要重试
	}
	return err == nil, err
}

func (rp *RedisResourcePool) Disassociate(resource string) (bool, error) {
	ctx := context.Background()

	pipe := rp.client.TxPipeline()
	pipe.SRem(ctx, rp.availableSetKey, resource)
	pipe.SRem(ctx, rp.OccupiedSetKey, resource)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	for _, cmd := range cmds {
		if cmd.(*redis.IntCmd).Val() == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (rp *RedisResourcePool) Acquire() (string, error) {
	ctx := context.Background()

	result, err := rp.client.Eval(ctx, luaScriptSRandMemberAndSMove, []string{rp.availableSetKey, rp.OccupiedSetKey}).Result()
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return result.(string), nil
}

// Release 将占用资源释放
func (rp *RedisResourcePool) Release(resource string) (bool, error) {
	ctx := context.Background()
	return rp.client.SMove(ctx, rp.OccupiedSetKey, rp.availableSetKey, resource).Result()
}

// AvailableCount 获取可用资源数量
func (rp *RedisResourcePool) AvailableCount() (int64, error) {
	ctx := context.Background()
	return rp.client.SCard(ctx, rp.availableSetKey).Result()
}

// OccupiedCount 获取占用资源数量
func (rp *RedisResourcePool) OccupiedCount() (int64, error) {
	ctx := context.Background()
	return rp.client.SCard(ctx, rp.OccupiedSetKey).Result()
}

// TotalCount 获取资源池总资源数量
func (rp *RedisResourcePool) TotalCount() (int64, error) {
	ctx := context.Background()

	pipe := rp.client.TxPipeline()
	availableCmd := pipe.SCard(ctx, rp.availableSetKey)
	occupiedCmd := pipe.SCard(ctx, rp.OccupiedSetKey)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return availableCmd.Val() + occupiedCmd.Val(), nil
}

// IsAvailable 检查资源是否可用
func (rp *RedisResourcePool) IsAvailable(resource string) (bool, error) {
	ctx := context.Background()
	return rp.client.SIsMember(ctx, rp.availableSetKey, resource).Result()
}

// IsOccupied 检查资源是否占用
func (rp *RedisResourcePool) IsOccupied(resource string) (bool, error) {
	ctx := context.Background()
	return rp.client.SIsMember(ctx, rp.OccupiedSetKey, resource).Result()
}

// Has 检查资源是否存在
func (rp *RedisResourcePool) Has(resource string) (bool, error) {
	ctx := context.Background()

	pipe := rp.client.TxPipeline()
	availableCmd := pipe.SIsMember(ctx, rp.availableSetKey, resource)
	occupiedCmd := pipe.SIsMember(ctx, rp.OccupiedSetKey, resource)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return availableCmd.Val() || occupiedCmd.Val(), nil
}
