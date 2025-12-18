package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	host    string
	db      int
	expires time.Duration
}

func NewCache(host string, db int, expires time.Duration) *Cache {
	return &Cache{
		host:    host,
		db:      db,
		expires: expires}
}

func GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	client := GetClient()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.Set(ctx, key, v, c.expires*time.Second).Err()
	return err
}

func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	client := GetClient()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(value), dest)
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	client := GetClient()
	_, err := client.Del(ctx, key).Result()
	return err
}
