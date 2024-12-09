package base62

import (
	"fmt"
	"math"
	"strings"
)

const (
	base        = 62
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Encode 将整数编码为 Base62
func Encode(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var result strings.Builder
	for num > 0 {
		remainder := num % base
		result.WriteByte(base62Chars[remainder])
		num /= base
	}

	// Base62 编码结果是倒序的，需要反转
	encoded := result.String()
	return reverseString(encoded)
}

// Decode 将 Base62 字符串解码为整数
func Decode(encoded string) (int64, error) {
	var num int64
	for i, char := range encoded {
		index := strings.IndexRune(base62Chars, char)
		if index == -1 {
			return 0, fmt.Errorf("invalid character '%c' in Base62 string", char)
		}
		num += int64(index) * int64(math.Pow(base, float64(len(encoded)-1-i)))
	}
	return num, nil
}

// reverseString 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
