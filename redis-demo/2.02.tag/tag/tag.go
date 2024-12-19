package tag

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

// 标签系统
// 需要为带标签的每个目标创建一个集合，用于记录目标带有的所有标签
// 与此同时还需要为每个标签创建一个集合，用于记录所有带有该标签的目标
// 标签系统中根据多个标签查找目标的操作涉及复杂的交集运算，因此可以为其设置缓存以复用运算结果并缩短操作的响应时间。

const CacheTTL = 60 // 缓存过期时间（秒）

type Service interface {
	// Add 为目标添加标签
	Add(target string, tags []string) (int64, error)
	// Remove 移除目标的标签
	Remove(target string, tags []string) (int64, error)
	// GetTagsByTarget 获取目标的所有标签
	GetTagsByTarget(target string) ([]string, error)
	// GetTargetByTags 根据多个标签获取目标集合（交集）
	GetTargetByTags(tags []string) ([]string, error)
	// GetCachedTargetByTags 带缓存的版本：根据多个标签获取目标集合（交集）
	GetCachedTargetByTags(tags []string) ([]string, error)
}

type Tag struct {
	client redis.Cmdable
}

func NewTag(client redis.Cmdable) *Tag {
	return &Tag{client: client}
}

// 生成目标的标签集合的键，用于记录目标关联的所有标签
func makeTargetKey(target string) string {
	return fmt.Sprintf("Tag:target:%s", target)
}

// 生成标签的目标集合的键，用于记录带有该标签的所有目标
func makeTagKey(tag string) string {
	return fmt.Sprintf("Tag:tag:%s", tag)
}

// 缓存多标签交集运算结果的集合
func makeCachedTargetsKey(tags []string) string {
	// 使用 Sort 确保多个集合输入无论如何排列都会产生相同的缓存
	slices.Sort(tags)
	return fmt.Sprintf("Tag:cached_targets:%v", tags)
}

// Add 为目标添加标签，并返回成功添加的标签数量
func (t *Tag) Add(target string, tags []string) (int64, error) {
	ctx := context.Background()
	pipe := t.client.TxPipeline()

	// 将 target 添加到每个 tag 对应的集合中
	for _, tag := range tags {
		pipe.SAdd(ctx, makeTagKey(tag), target)
	}

	// 将所有 tag 添加到 target 对应的集合中
	targetKey := makeTargetKey(target)
	cmd := pipe.SAdd(ctx, targetKey, tags)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// Remove 移除目标的标签，并返回成功移除的标签数量
func (t *Tag) Remove(target string, tags []string) (int64, error) {
	ctx := context.Background()
	pipe := t.client.TxPipeline()

	// 从每个 tag 对应的集合中移除 target
	for _, tag := range tags {
		pipe.SRem(ctx, makeTagKey(tag), target)
	}

	// 从 target 的标签集合中移除 tag
	targetKey := makeTargetKey(target)
	cmd := pipe.SRem(ctx, targetKey, tags)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return 0, err
	}
	return cmd.Val(), nil
}

// GetTagsByTarget 获取目标的所有标签
func (t *Tag) GetTagsByTarget(target string) ([]string, error) {
	ctx := context.Background()
	targetKey := makeTargetKey(target)
	cmd := t.client.SMembers(ctx, targetKey)
	return cmd.Result()
}

// GetTargetByTags 根据多个标签获取目标集合（交集）
func (t *Tag) GetTargetByTags(tags []string) ([]string, error) {
	ctx := context.Background()
	tagKeys := make([]string, len(tags))
	for i, tag := range tags {
		tagKeys[i] = makeTagKey(tag)
	}
	cmd := t.client.SInter(ctx, tagKeys...)
	return cmd.Result()
}

// GetCachedTargetByTags 带缓存的版本：根据多个标签获取目标集合（交集）
func (t *Tag) GetCachedTargetByTags(tags []string) ([]string, error) {
	ctx := context.Background()
	cacheKey := makeCachedTargetsKey(tags)

	// 尝试直接从缓存获取
	cachedTargets, err := t.client.SMembers(ctx, cacheKey).Result()
	if err == nil && len(cachedTargets) > 0 {
		return cachedTargets, nil
	}

	// 缓存不存在，计算交集并存储到缓存中
	tagKeys := make([]string, len(tags))
	for i, tag := range tags {
		tagKeys[i] = makeTagKey(tag)
	}

	pipe := t.client.TxPipeline()
	pipe.SInterStore(ctx, cacheKey, tagKeys...)
	// 然后再设置过期时间
	pipe.Expire(ctx, cacheKey, CacheTTL*time.Second)
	cmd := pipe.SMembers(ctx, cacheKey)
	_, err = pipe.Exec(ctx)

	if err != nil {
		return nil, err
	}
	// 返回交集元素
	return cmd.Val(), nil
}
