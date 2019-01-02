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
	json      jsoniter.API
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

// Implementation example
func (rc *Client) DummyFunc(ctx context.Context, key string) apierr.Error {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	_, err := conn.Do("HGET", key)
	if err != nil {
		return apierr.RedisUnavailable(err)
	}

	return nil
}
