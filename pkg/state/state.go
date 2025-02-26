package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Jiali-Xing/hotelApp/internal/config"

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
		config.DebugLog("Redis client initialized at %s", redisAddr)
	})
}

func GetState[T interface{}](ctx context.Context, key string) (T, error) {
	var value T

	initRedisClient()

	// Retrieve state directly from Redis
	result, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist
		config.DebugLog("Key not found: %s", key)
		return value, errors.New("key not found")
	} else if err != nil {
		// Other Redis error
		return value, err
	}

	// Unmarshal the JSON-encoded value into the value of type T
	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		return value, err
	}
	config.DebugLog("Retrieved value for key %s", key)
	return value, nil
}

func GetBulkState[T interface{}](ctx context.Context, keys []string) ([]T, error) {
	// Initialize Redis client
	initRedisClient()

	// Retrieve multiple keys from Redis
	items, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	values := make(map[string]*T)
	for i, item := range items {
		var value T
		if item == nil {
			return nil, errors.New(fmt.Sprintf("key %s not found", keys[i]))
		}
		err = json.Unmarshal([]byte(item.(string)), &value)
		if err != nil {
			return nil, err
		}
		values[keys[i]] = &value
	}

	var returnValues []T
	for _, key := range keys {
		returnValues = append(returnValues, *values[key])
	}
	return returnValues, nil
}

func GetBulkStateDefault[T interface{}](ctx context.Context, keys []string, defVal T) ([]T, error) {
	// Initialize Redis client
	initRedisClient()

	// Retrieve multiple keys from Redis
	items, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	values := make(map[string]*T)
	for i, item := range items {
		var value T
		if item == nil {
			value = defVal
		} else {
			err = json.Unmarshal([]byte(item.(string)), &value)
			if err != nil {
				return nil, err
			}
		}
		values[keys[i]] = &value
	}

	var returnValues []T
	for _, key := range keys {
		returnValues = append(returnValues, *values[key])
	}
	return returnValues, nil
}

func SetState(ctx context.Context, key string, value interface{}) error {
	initRedisClient()

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Save state directly to Redis
	err = rdb.Set(ctx, key, valueBytes, 0).Err()
	if err != nil {
		return err
	}
	config.DebugLog("Set value for key %s", key)
	return nil
}

func SetBulkState(ctx context.Context, kvs map[string]interface{}) error {
	// Initialize Redis client
	initRedisClient()

	pipe := rdb.Pipeline()

	for k, v := range kvs {
		valueBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		pipe.Set(ctx, k, valueBytes, 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
