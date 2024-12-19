package location

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// 地理位置

const maxSearchSize = 10 // 10km 以内的搜索范围

type Location interface {
	// Pin 记录用户的地理位置
	Pin(user string, longitude, latitude float64) error
	// Locate 获取用户的地理位置
	Locate(user string) (longitude, latitude float64, err error)
	// Distance 以公里为单位返回两个用户之间的直线距禂
	Distance(user1, user2 string) (float64, error)
	// Size 获取已存储的用户位置数量
	Size() (int64, error)
	// Search 以指定用户为中心，搜索指定范围内的其他用户
	Search(user string, radius float64, limit ...int) ([]string, error)
}

type RedisLocation struct {
	client redis.Cmdable
	key    string
}

func NewRedisLocation(client redis.Cmdable, key string) *RedisLocation {
	return &RedisLocation{client: client, key: key}
}

// Pin 记录用户的地理位置
func (r *RedisLocation) Pin(user string, longitude, latitude float64) error {
	return r.client.GeoAdd(context.Background(), r.key, &redis.GeoLocation{
		Name:      user,
		Longitude: longitude,
		Latitude:  latitude,
	}).Err()
}

// Locate 获取用户的地理位置
func (r *RedisLocation) Locate(user string) (longitude, latitude float64, err error) {
	res, err := r.client.GeoPos(context.Background(), r.key, user).Result()
	if err != nil {
		return 0, 0, err
	}
	return res[0].Longitude, res[0].Latitude, nil
}

// Distance 以公里为单位返回两个用户之间的直线距离
func (r *RedisLocation) Distance(user1, user2 string) (float64, error) {
	return r.client.GeoDist(context.Background(), r.key, user1, user2, "km").Result()
}

// Size 获取已存储的用户位置数量
func (r *RedisLocation) Size() (int64, error) {
	return r.client.ZCard(context.Background(), r.key).Result()
}

// Search 以指定用户为中心，搜索指定范围内的其他用户
func (r *RedisLocation) Search(user string, radius float64, limit ...int) ([]string, error) {
	count := maxSearchSize
	if len(limit) > 0 {
		count = limit[0]
	}

	res, err := r.client.GeoSearch(context.Background(), r.key, &redis.GeoSearchQuery{
		Member:     user,
		Radius:     radius,
		RadiusUnit: "km",
		Sort:       "",
		Count:      count,
	}).Result()
	if err != nil {
		return nil, err
	}

	users := make([]string, 0, len(res)-1)
	// 排除自己
	for _, name := range res {
		if name == user {
			continue
		}
		users = append(users, name)
	}
	return users, nil
}
