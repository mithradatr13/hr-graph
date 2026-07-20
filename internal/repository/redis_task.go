package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"task-manager/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTaskCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisTaskCache(client *redis.Client, ttl time.Duration) *RedisTaskCache {
	return &RedisTaskCache{client: client, ttl: ttl}
}

func (c *RedisTaskCache) Get(ctx context.Context, id int64) (*domain.Task, error) {
	key := fmt.Sprintf("task:%d", id)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var task domain.Task
	if err := json.Unmarshal([]byte(val), &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *RedisTaskCache) Set(ctx context.Context, task *domain.Task) error {
	key := fmt.Sprintf("task:%d", task.ID)
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *RedisTaskCache) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("task:%d", id)
	return c.client.Del(ctx, key).Err()
}
