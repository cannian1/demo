package timeline

import (
	"context"
	"math"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

const defaultSize int64 = 10 // 默认每页返回10条数据

type Timeline interface {
	// Add 向时间线中添加若干项
	Add(content ...Content) (int64, error)
	// Get 获取指定页数的数据
	Get(page int64, pageSize ...int64) ([]string, error)
	// GetWithTime 从指定分页中取出指定数量的项（包括元素和时间戳）
	GetWithTime(page int64, pageSize ...int64) ([]Content, error)
	// GetByTimeRange 对位于指定时间戳范围内的元素进行分页
	GetByTimeRange(min, max, page int64, pageSize ...int64) ([]Content, error)
	// Length 获取时间线包含的元素数量
	Length() (int64, error)
	// Number 返回在获取指定数量的元素时，时间线包含的页数量，如果时间线为空则返回0
	Number(pageSize ...int64) (int64, error)
}

// Paging 时间线分页器
type Paging struct {
	client *redis.Client
	key    string
}

// NewPaging 创建一个新的时间线分页器
func NewPaging(client *redis.Client, key string) *Paging {
	return &Paging{client: client, key: key}
}

type Content struct {
	Item   any
	Weight int
}

// Add 向分页器中添加若干项
func (p Paging) Add(content ...Content) (int64, error) {
	ctx := context.Background()

	tx := p.client.TxPipeline()

	z := make([]redis.Z, 0, len(content))
	for _, c := range content {
		z = append(z, redis.Z{Score: float64(c.Weight), Member: c.Item})
	}

	tx.ZAdd(ctx, p.key, z...) // 添加元素
	tx.ZCard(ctx, p.key)      // 获取时间线包含的元素数量

	cmds, err := tx.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return cmds[1].(*redis.IntCmd).Val(), nil
}

// Get 获取指定页数的数据
func (p Paging) Get(page int64, pageSize ...int64) ([]string, error) {
	size := defaultSize
	if len(pageSize) != 0 {
		size = pageSize[0]
	}

	start := (page - 1) * size
	stop := page * size
	return p.client.ZRevRange(context.Background(), p.key, start, stop).Result()
}

// GetWithTime 从指定分页中取出指定数量的项（包括元素和时间戳）
func (p Paging) GetWithTime(page int64, pageSize ...int64) ([]Content, error) {
	size := defaultSize
	if len(pageSize) != 0 {
		size = pageSize[0]
	}

	start := (page - 1) * size
	stop := page * size

	vals, err := p.client.ZRevRangeWithScores(context.Background(), p.key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	content := make([]Content, 0, len(vals))
	for _, v := range vals {
		content = append(content, Content{Item: v.Member, Weight: cast.ToInt(v.Score)})
	}

	return content, nil
}

// GetByTimeRange 对位于指定时间戳范围内的元素进行分页
func (p Paging) GetByTimeRange(min, max, page int64, pageSize ...int64) ([]Content, error) {
	size := defaultSize
	if len(pageSize) != 0 {
		size = pageSize[0]
	}

	offset := (page - 1) * size

	vals, err := p.client.ZRevRangeByScoreWithScores(context.Background(), p.key, &redis.ZRangeBy{
		Min:    cast.ToString(min),
		Max:    cast.ToString(max),
		Offset: offset,
		Count:  size,
	}).Result()
	if err != nil {
		return nil, err
	}

	content := make([]Content, 0, len(vals))
	for _, v := range vals {
		content = append(content, Content{Item: v.Member, Weight: cast.ToInt(v.Score)})
	}

	return content, nil
}

// Length 获取时间线包含的元素数量
func (p Paging) Length() (int64, error) {
	return p.client.ZCard(context.Background(), p.key).Result()
}

// Number 返回在获取指定数量的元素时，时间线包含的页数量，如果时间线为空则返回0
func (p Paging) Number(pageSize ...int64) (int64, error) {
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
