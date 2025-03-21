package copier

import "github.com/jinzhu/copier"

func Copy(dst interface{}, src interface{}) error {
	return copier.Copy(dst, src)
}

func CopyWithOption(dst interface{}, src interface{}, opts ...func(*copier.Option)) error {
	return copier.CopyWithOption(dst, src, MergeOptions(opts...))
}

func MergeOptions(opts ...func(*copier.Option)) copier.Option {
	option := copier.Option{}
	for _, opt := range opts {
		opt(&option)
	}
	return option
}
