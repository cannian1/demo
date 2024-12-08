package binary_record

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 用位图实现二元操作记录器
// 以签到记录为例，每天签到一次，签到记录可以用一个二进制位表示
// 选择一天作为基准日期，签到记录的 key 为 user:1:sign_in
// 二进制位索引为 0 表示基准日期，1 表示第一天，以此类推
// 二进制位值为 1 表示签到，0 表示未签到
// 在这基础上，如果想要检查用户的全勤情况，只需要统计指定索引范围二进制位值为 1 的个数即可

type BinaryRecorder interface {
	// SetBit 将制定索引的二进制位设置为 1
	SetBit(index int64) error
	// ClearBit 将指定索引的二进制位设置为 0
	ClearBit(index int64) error
	// GetBit 获取指定索引的二进制位的值
	GetBit(index int64) (int64, error)
	// CountBit 统计指定范围内二进制位值为 1 的个数
	CountBit(start, end int64) (int64, error)
}

type BinaryRecord struct {
	client *redis.Client
	key    string
}

func NewBinaryRecord(client *redis.Client, key string) *BinaryRecord {
	return &BinaryRecord{client: client, key: key}
}

func (b *BinaryRecord) SetBit(index int64) error {
	return b.client.SetBit(context.Background(), b.key, index, 1).Err()
}

func (b *BinaryRecord) ClearBit(index int64) error {
	return b.client.SetBit(context.Background(), b.key, index, 0).Err()
}

func (b *BinaryRecord) GetBit(index int64) (int64, error) {
	return b.client.GetBit(context.Background(), b.key, index).Result()
}

func (b *BinaryRecord) CountBit(start, end int64) (int64, error) {
	return b.client.BitCount(context.Background(), b.key, &redis.BitCount{Start: start, End: end, Unit: redis.BitCountIndexBit}).Result()
}
