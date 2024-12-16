package bloom_filter

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spaolacci/murmur3"
)

type BloomFilter interface {
	Add(data string) error
	Exists(data string) (bool, error)
}

type BloomFilterService struct {
	client    *redis.Client
	key       string
	size      uint // 位数组大小
	hashCount uint // 哈希函数数量
	lock      sync.Mutex
}

func NewBloomFilterService(client *redis.Client, key string, expectedItems uint, falsePositiveRate float64) *BloomFilterService {
	size, hashCount := calculateOptimalParams(expectedItems, falsePositiveRate)

	return &BloomFilterService{
		client:    client,
		key:       key,
		size:      size,
		hashCount: hashCount,
	}
}

func (bf *BloomFilterService) Add(data string) error {
	bf.lock.Lock()
	defer bf.lock.Unlock()

	for i := uint(0); i < bf.hashCount; i++ {
		hashValue := hashWithSeed(uint32(i), []byte(data))
		offset := hashValue % uint64(bf.size) // 计算偏移量

		_, err := bf.client.SetBit(context.Background(), bf.key, int64(offset), 1).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists 检查元素是否可能存在
func (bf *BloomFilterService) Exists(data string) (bool, error) {
	bf.lock.Lock()
	defer bf.lock.Unlock()

	for i := uint(0); i < bf.hashCount; i++ {
		hashValue := hashWithSeed(uint32(i), []byte(data))
		offset := hashValue % uint64(bf.size) // 计算偏移量
		bit, err := bf.client.GetBit(context.Background(), bf.key, int64(offset)).Result()
		if err != nil {
			return false, fmt.Errorf("error getting bit from Redis: %w", err)
		}
		if bit == 0 {
			return false, nil // 任何一个位为 0，元素一定不存在
		}
	}
	return true, nil
}

func hashWithSeed(seed uint32, data []byte) uint64 {
	hasher := murmur3.New64WithSeed(seed)
	hasher.Write(data)
	return hasher.Sum64()
}

// 动态计算布隆过滤器的最佳大小和哈希函数数量
func calculateOptimalParams(expectedItems uint, falsePositiveRate float64) (uint, uint) {
	if expectedItems == 0 {
		panic("expectedItems must be greater than 0")
	}
	if falsePositiveRate <= 0.0 || falsePositiveRate >= 1.0 {
		panic("falsePositiveRate must be between 0 and 1 (exclusive)")
	}

	size := uint(-float64(expectedItems) * math.Log(falsePositiveRate) / (math.Ln2 * math.Ln2)) // 最优位数组大小
	hashCount := uint(math.Ceil(math.Ln2 * float64(size) / float64(expectedItems)))             // 最优哈希函数数量
	return size, hashCount
}
