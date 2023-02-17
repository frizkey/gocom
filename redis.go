package gocom

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisKV struct {
	ctx    context.Context
	client *redis.Client
}

func (o *RedisKV) Set(key string, val interface{}, ttl ...time.Duration) error {

	targetTTL := time.Second * 0

	if len(ttl) > 0 {
		targetTTL = ttl[0]
	}

	return o.client.Set(o.ctx, key, val, targetTTL).Err()
}

func (o *RedisKV) SetNX(key string, val interface{}, ttl ...time.Duration) bool {

	targetTTL := time.Second * 0

	if len(ttl) > 0 {
		targetTTL = ttl[0]
	}

	cmd := o.client.SetNX(o.ctx, key, val, targetTTL)

	if cmd.Err() == nil {

		return cmd.Val()
	}

	return false
}

func (o *RedisKV) Get(key string) interface{} {

	cmd := o.client.Get(o.ctx, key)

	if cmd.Err() == nil {
		return cmd.Val()
	}

	return nil
}

func (o *RedisKV) GetString(key string) string {

	cmd := o.client.Get(o.ctx, key)

	if cmd.Err() == nil {
		return cmd.Val()
	}

	return ""
}

func (o *RedisKV) GetInt(key string) int {

	cmd := o.client.Get(o.ctx, key)

	if cmd.Err() == nil {
		val, err := strconv.Atoi(cmd.Val())

		if err == nil {
			return val
		}
	}

	return 0
}

func (o *RedisKV) Del(key string) error {

	cmd := o.client.Del(o.ctx, key)

	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (o *RedisKV) LPush(key string, val interface{}) error {

	return nil
}

func (o *RedisKV) LPop(key string) interface{} {

	return nil
}

func (o *RedisKV) LPopString(key string) string {

	return ""
}

func (o *RedisKV) LPopInt(key string) int {

	return 0
}

func (o *RedisKV) RPush(key string, val interface{}) error {

	return nil
}

func (o *RedisKV) RPop(key string) interface{} {

	return nil
}

func (o *RedisKV) RPopString(key string) string {

	return ""
}

func (o *RedisKV) RPopInt(key string) int {

	return 0
}

func (o *RedisKV) Len(key string) int {

	return 0
}

func (o *RedisKV) AtIndex(key string, index int) interface{} {

	return nil
}

func (o *RedisKV) AtIndexString(key string, index int) string {

	return ""
}

func (o *RedisKV) AtIndexInt(key string, index int) int {

	return 0
}

//-----------------------------------------------------------------------

func init() {
	RegKVCreator("redis", func(url string) (KV, error) {
		ret := &RedisKV{
			ctx: context.Background(),
		}

		opt, err := redis.ParseURL(url)

		if err == nil {
			ret.client = redis.NewClient(opt)
		}

		return ret, nil
	})
}
