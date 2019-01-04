package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/log-go"
)

type Client struct {
	redisPool *redis.Pool
	json      jsoniter.API //Use this to be able to save objects as value and use redis tags instead of json ones
}

func NewClient(redisPool *redis.Pool, json jsoniter.API) *Client {
	return &Client{
		redisPool: redisPool,
		json:      json,
	}
}

func NewRedisPool(cfg *config.Config) (*redis.Pool, error) {
	redisPoolIdleTimeout, err := time.ParseDuration(cfg.RedisStorePoolIdleTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "invalid value for redis_pool_idle_timeout, must be duration")
	}

	return &redis.Pool{
		MaxIdle:     cfg.RedisStorePoolMaxIdle,
		IdleTimeout: redisPoolIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", cfg.RedisStoreHost, cfg.RedisStorePort), redis.DialDatabase(cfg.RedisStoreDatabase))
		},
	}, nil
}

func (rc *Client) closeConn(ctx context.Context, conn redis.Conn) {
	err := conn.Close()
	if err != nil {
		err = errors.Wrap(err, "Failed closing redis connection, this might result in memory leek")
		log.Warn(ctx, err.Error(), "error", err)
	}
}

// Implementation example for get simple value
func (rc *Client) GetSimpleFunc(ctx context.Context, key string) (string, apierr.Error) {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", apierr.RedisUnavailable(err)
	}

	return value, nil
}

// Implementation example for save complex value
func (rc *Client) SaveComplexFunc(ctx context.Context, key string, value DummyType) apierr.Error {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	body, err := rc.json.Marshal(value)
	if err != nil {
		return apierr.Semantic(err)
	}

	_, err = conn.Do("SET", key, body)
	if err != nil {
		return apierr.RedisUnavailable(err)
	}

	return nil
}

// Implementation example for get complex value
func (rc *Client) GetComplexFunc(ctx context.Context, key string) (*DummyType, apierr.Error) {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	repl, err := conn.Do("GET", key)
	if err != nil {
		return nil, apierr.RedisUnavailable(err)
	}
	if repl == nil {
		return nil, nil
	}

	result := &DummyType{}
	err = rc.json.Unmarshal(repl.([]byte), result)
	if err != nil {
		return nil, apierr.Semantic(err)
	}

	return result, nil
}

type DummyType struct {
	Test string `json:"test_dummy" redis:"dummy_test"`
}
