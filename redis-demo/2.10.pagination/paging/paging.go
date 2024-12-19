package paging

import (
	"context"
	"math"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// 分页器

const defaultSize = 10 // 默认每页返回10条数据

type Pagination interface {
	// Add 向分页器中添加若干项
	Add(items ...any) (int64, error)
	// Get 获取指定页数的数据
	Get(page int, pageSize ...int) ([]string, error)
	// Length 获取分页器中的数据总数
	Length() (int64, error)
	// Number 返回在获取指定数量的元素时，分页列表包含的页数量，如果分页列表为空则返回0
	Number(pageSize ...int) (int64, error)
}

type Paging struct {
	client redis.Cmdable
	key    string
}

// NewPaging 创建一个新的分页器
func NewPaging(client redis.Cmdable, key string) *Paging {
	return &Paging{client: client, key: key}
}

// Add 向分页器中添加若干项
func (p Paging) Add(items ...any) (int64, error) {
	return p.client.LPush(context.Background(), p.key, items...).Result()
}

// Get 获取指定页数的数据
func (p Paging) Get(page int, pageSize ...int) ([]string, error) {
	size := defaultSize
	if len(pageSize) != 0 {
		size = pageSize[0]
	}

	start := int64((page - 1) * size)
	stop := int64(page * size)
	return p.client.LRange(context.Background(), p.key, start, stop).Result()
}

// Length 获取分页器中的数据总数
func (p Paging) Length() (int64, error) {
	return p.client.LLen(context.Background(), p.key).Result()
}

// Number 返回在获取指定数量的元素时，分页列表包含的页数量，如果分页列表为空则返回0
func (p Paging) Number(pageSize ...int) (int64, error) {
	size := defaultSize
	if len(pageSize) != 0 {
		size = pageSize[0]
	}
	length, err := p.Length()
	if err != nil {
		return 0, err
	}

	return cast.ToInt64(math.Ceil(float64(length) / float64(size))), nil
}
