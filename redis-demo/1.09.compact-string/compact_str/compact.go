package compact_str

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
)

// 紧凑字符串

// CompactString 用于在一个 key 里存储多个字符串的数据结构
// 通过指定的分隔符将多个字符串连接在一起存储

const defaultSeparator = "\n"

type CompactString interface {
	// Append 将给定的字符串添加至已有字符串值的末尾
	Append(str string) (int64, error)
	// GetBytes 这个方法将返回一个列表，其中可以包含零个或任意多个字符串。
	GetBytes(...OptBytesRange) ([]string, error)
}

type OptBytesRange struct {
	Start int64
	End   int64
}

type CompactStr struct {
	client    redis.Cmdable
	key       string
	separator string
}

func NewCompactStr(client redis.Cmdable, key string, separator string) *CompactStr {
	if separator == "" {
		separator = defaultSeparator
	}
	return &CompactStr{client: client, key: key, separator: separator}
}

func (cs *CompactStr) Append(str string) (int64, error) {
	str += cs.separator
	return cs.client.Append(context.Background(), cs.key, str).Result()
}

func (cs *CompactStr) GetBytes(bytesRange ...OptBytesRange) ([]string, error) {
	ctx := context.Background()
	start, end := int64(0), int64(-1)

	if len(bytesRange) > 0 {
		start = bytesRange[0].Start
		end = bytesRange[0].End
	}

	// 获取指定范围的字符串数据
	content, err := cs.client.GetRange(ctx, cs.key, start, end).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []string{}, nil // 键不存在
		}
		return nil, err
	}

	// 处理内容为空的情况
	if content == "" {
		return []string{}, nil
	}

	// 根据分隔符分割内容
	listOfStrings := strings.Split(content, cs.separator)

	// 移除列表中可能的空字符串
	cleanedStrings := make([]string, 0, len(listOfStrings))
	for _, s := range listOfStrings {
		if s != "" {
			cleanedStrings = append(cleanedStrings, s)
		}
	}

	return cleanedStrings, nil
}
