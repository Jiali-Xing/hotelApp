package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func initRedisClient() {
	once.Do(func() {
		// Get Redis address from environment variable
		redisAddr := os.Getenv("REDIS_ADDR")
		if redisAddr == "" {
			redisAddr = "localhost:6379" // Default to localhost if not set
		}

		// Initialize Redis client
		rdb = redis.NewClient(&redis.Options{
			Addr: redisAddr,
		})
	})
}

func GetState[T interface{}](ctx context.Context, key string) (T, error) {
	var value T

	initRedisClient()

	// Retrieve state directly from Redis
	result, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist
		return value, errors.New("key not found")
	} else if err != nil {
		// Other Redis error
		log.Panic(err)
	}

	// Unmarshal the JSON-encoded value into the value of type T
	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		log.Panic(err)
	}

	return value, nil
}

func GetBulkState[T interface{}](ctx context.Context, keys []string) ([]T, error) {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Redis server address
	})

	// Retrieve multiple keys from Redis
	items, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		log.Panic(err)
	}

	values := make(map[string]*T)
	for i, item := range items {
		var value T
		if item == nil {
			return nil, errors.New(fmt.Sprintf("key %s not found", keys[i]))
		}
		err = json.Unmarshal([]byte(item.(string)), &value)
		if err != nil {
			log.Panic(err)
		}
		values[keys[i]] = &value
	}

	var returnValues []T
	for _, key := range keys {
		returnValues = append(returnValues, *values[key])
	}
	return returnValues, nil
}

func GetBulkStateDefault[T interface{}](ctx context.Context, keys []string, defVal T) []T {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Redis server address
	})

	// Retrieve multiple keys from Redis
	items, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		log.Panic(err)
	}

	values := make(map[string]*T)
	for i, item := range items {
		var value T
		if item == nil {
			value = defVal
		} else {
			err = json.Unmarshal([]byte(item.(string)), &value)
			if err != nil {
				log.Panic(err)
			}
		}
		values[keys[i]] = &value
	}

	var returnValues []T
	for _, key := range keys {
		returnValues = append(returnValues, *values[key])
	}
	return returnValues
}

func SetState(ctx context.Context, key string, value interface{}) {
	initRedisClient()

	valueBytes, err := json.Marshal(value)
	if err != nil {
		log.Panic(err)
	}

	// Save state directly to Redis
	err = rdb.Set(ctx, key, valueBytes, 0).Err()
	if err != nil {
		log.Panic(err)
	}
}

func SetBulkState(ctx context.Context, kvs map[string]interface{}) {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Redis server address
	})

	pipe := rdb.Pipeline()

	for k, v := range kvs {
		valueBytes, err := json.Marshal(v)
		if err != nil {
			log.Panic(err)
		}
		pipe.Set(ctx, k, valueBytes, 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Panic(err)
	}

}
