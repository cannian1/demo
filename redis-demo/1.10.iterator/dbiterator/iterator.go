package dbiterator

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 如果一个程序需要访问 redis 中的大量数据，最好用迭代器方式
// 迭代器是一个非常重要的数据库特性，它可以让用户分批获取数据库中的数据，而不是一次性获取整个数据库的全部数据，因为后一种做法在数据量巨大的时候很容易造成服务器阻塞。

type Iterator interface {
	// Next 返回迭代器的下一个元素
	Next() ([]string, error)
	// Rewind 重置迭代器
	Rewind()
}

const DefaultCount = 10

type DbIterator struct {
	client  *redis.Client
	count   int64
	cursor  uint64
	hasMore bool
}

// NewDbIterator 创建一个新的数据库迭代器。
// count 参数用于指定每次返回的键数量，默认为 DefaultCount。
func NewDbIterator(client *redis.Client, count int64) *DbIterator {
	if count <= 0 {
		count = DefaultCount
	}
	return &DbIterator{
		client:  client,
		count:   count,
		cursor:  0,
		hasMore: true,
	}
}

// Next 继续迭代数据库，返回本次迭代的键列表。
// 如果迭代已完成，返回空列表。
func (it *DbIterator) Next() ([]string, error) {
	ctx := context.Background()

	if !it.hasMore {
		// 迭代已完成，返回空列表
		return []string{}, nil
	}

	// 调用 SCAN 命令进行下一次迭代
	keys, newCursor, err := it.client.Scan(ctx, it.cursor, "", it.count).Result()
	if err != nil {
		return nil, err
	}

	// 判断迭代是否结束
	it.cursor = newCursor
	if it.cursor == 0 {
		it.hasMore = false
	}

	return keys, nil
}

// Rewind 重置迭代器，以便从头开始对数据库进行迭代。
func (it *DbIterator) Rewind() {
	it.cursor = 0
	it.hasMore = true
}
