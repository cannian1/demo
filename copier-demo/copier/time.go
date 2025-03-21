package copier

import (
	"fmt"
	"time"

	"github.com/jinzhu/copier"
)

// TimeToInt64 将time.Time转换为int64时间戳
func TimeToInt64() func(*copier.Option) {
	return func(opt *copier.Option) {
		opt.Converters = append(opt.Converters, copier.TypeConverter{
			SrcType: time.Time{},
			DstType: int64(0),
			Fn: func(src interface{}) (dst interface{}, err error) {
				if t, ok := src.(time.Time); ok {
					return t.Unix(), nil
				}
				return nil, fmt.Errorf("invalid type: %T", src)
			},
		})
	}
}

// Int64ToTime 将int64时间戳转换为time.Time
func Int64ToTime() func(*copier.Option) {
	return func(opt *copier.Option) {
		opt.Converters = append(opt.Converters, copier.TypeConverter{
			SrcType: int64(0),
			DstType: time.Time{},
			Fn: func(src interface{}) (dst interface{}, err error) {
				if i, ok := src.(int64); ok {
					return time.Unix(i, 0), nil
				}
				return nil, fmt.Errorf("invalid type: %T", src)
			},
		})
	}
}
