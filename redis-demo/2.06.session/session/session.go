package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// 登录会话

const (
	// 会话超时时间
	defaultTimeout = time.Hour * 24
	// 会话状态
	TokenValid    Status = "TokenValid"
	TokenInvalid  Status = "TokenInvalid"
	TokenNotFound Status = "TokenNotFound"
)

type Status string

type Session interface {
	// Create 创建一个新的会话
	Create(timeout ...time.Duration) (string, error)
	// Validate 验证会话是否有效
	Validate(token string) (Status, error)
	// Destroy 销毁会话
	Destroy() error
}

type RedisSession struct {
	client redis.Cmdable
	uid    string
}

func NewRedisSession(client redis.Cmdable, uid string) *RedisSession {
	return &RedisSession{client: client, uid: uid}
}

func (rs *RedisSession) Create(timeout ...time.Duration) (string, error) {
	t := defaultTimeout
	if len(timeout) > 0 {
		t = timeout[0]
	}

	tokenKey := makeTokenKey(rs.uid)
	token := generateToken()

	err := rs.client.Set(context.Background(), tokenKey, token, t).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateToken() string {
	return uuid.New().String()
}

func makeTokenKey(uid string) string {
	return fmt.Sprintf("User:%s:token", uid)
}

func (rs *RedisSession) Validate(token string) (Status, error) {
	ctx := context.Background()
	tokenKey := makeTokenKey(rs.uid)

	val, err := rs.client.Get(ctx, tokenKey).Result()
	if errors.Is(err, redis.Nil) {
		return TokenNotFound, nil
	}
	if err != nil {
		return TokenInvalid, err
	}

	if val != token {
		return TokenInvalid, nil
	}

	return TokenValid, nil
}

func (rs *RedisSession) Destroy() error {
	ctx := context.Background()
	tokenKey := makeTokenKey(rs.uid)

	return rs.client.Del(ctx, tokenKey).Err()
}
