package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) DecreaseStock(ctx context.Context, key string, quantity uint32) error {
	luaScript := `
    local currentStock = redis.call('GET', KEYS[1])
    if currentStock == false then
        return redis.error_reply("库存键不存在")
    end
    currentStock = tonumber(currentStock)
    local purchaseAmount = tonumber(ARGV[1])

    if currentStock >= purchaseAmount then
        redis.call('DECRBY', KEYS[1], purchaseAmount)
		return redis.call('GET', KEYS[1])
    else
        return redis.error_reply("库存不足")
    end
    `
	_, err := r.client.Eval(ctx, luaScript, []string{key}, quantity).Result()
	return err
}
