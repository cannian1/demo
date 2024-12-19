package vote

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 投票

type Voter interface {
	// Up 用户给某个对象投赞成票
	Up(user string) (bool, error)
	// Down 用户给某个对象投反对票
	Down(user string) (bool, error)
	// IsVoted 用户是否已经投过票
	IsVoted(user string) (bool, error)
	// UnVote 取消用户的投票
	UnVote(user string) (bool, error)
	// UpCount 获取目前投赞成票的用户数
	UpCount() (int64, error)
	// DownCount 获取目前投反对票的用户数
	DownCount() (int64, error)
	// Total 获取目前参与投票的总用户数
	Total() (int64, error)
}

type Vote struct {
	client      redis.Cmdable
	voteUpKey   string
	voteDownKey string
}

func NewVote(client redis.Cmdable, subject string) *Vote {
	return &Vote{
		client:      client,
		voteUpKey:   fmt.Sprintf("Vote:%s:up", subject),
		voteDownKey: fmt.Sprintf("Vote:%s:down", subject),
	}
}

// Up 用户给某个对象投赞成票
func (v Vote) Up(user string) (bool, error) {
	ctx := context.Background()

	tx := v.client.TxPipeline()
	tx.SAdd(ctx, v.voteUpKey, user)
	tx.SRem(ctx, v.voteDownKey, user) // 移除可能存在的反对票

	res, err := tx.Exec(ctx)
	if err != nil {
		return false, err
	}

	return res[0].(*redis.IntCmd).Val() == 1, nil
}

// Down 用户给某个对象投反对票
func (v Vote) Down(user string) (bool, error) {
	ctx := context.Background()

	tx := v.client.TxPipeline()
	tx.SAdd(ctx, v.voteDownKey, user)
	tx.SRem(ctx, v.voteUpKey, user) // 移除可能存在的赞成票

	res, err := tx.Exec(ctx)
	if err != nil {
		return false, err
	}
	return res[0].(*redis.IntCmd).Val() == 1, nil
}

// IsVoted 用户是否已经投过票
func (v Vote) IsVoted(user string) (bool, error) {
	ctx := context.Background()

	tx := v.client.TxPipeline()
	tx.SIsMember(ctx, v.voteUpKey, user)
	tx.SIsMember(ctx, v.voteDownKey, user)

	cmds, err := tx.Exec(ctx)
	if err != nil {
		return false, err
	}

	return cmds[0].(*redis.BoolCmd).Val() || cmds[1].(*redis.BoolCmd).Val(), nil
}

// UnVote 取消用户的投票
func (v Vote) UnVote(user string) (bool, error) {
	ctx := context.Background()

	tx := v.client.TxPipeline()
	tx.SRem(ctx, v.voteUpKey, user)
	tx.SRem(ctx, v.voteDownKey, user)

	res, err := tx.Exec(ctx)
	if err != nil {
		return false, err
	}

	return res[0].(*redis.IntCmd).Val() == 1 || res[1].(*redis.IntCmd).Val() == 1, nil
}

// UpCount 获取目前投赞成票的用户数
func (v Vote) UpCount() (int64, error) {
	return v.client.SCard(context.Background(), v.voteUpKey).Result()
}

// DownCount 获取目前投反对票的用户数
func (v Vote) DownCount() (int64, error) {
	return v.client.SCard(context.Background(), v.voteDownKey).Result()
}

// Total 获取目前参与投票的总用户数
func (v Vote) Total() (int64, error) {
	ctx := context.Background()

	tx := v.client.TxPipeline()
	tx.SCard(ctx, v.voteUpKey)
	tx.SCard(ctx, v.voteDownKey)

	cmds, err := tx.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return cmds[0].(*redis.IntCmd).Val() + cmds[1].(*redis.IntCmd).Val(), nil
}
