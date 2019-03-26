package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/log-go"
)

type Redis struct {
	redisPool *redis.Pool
	json      jsoniter.API //Use this to be able to save objects as value and use redis tags instead of json ones
}

func NewRedis(redisPool *redis.Pool, json jsoniter.API) *Redis {
	return &Redis{
		redisPool: redisPool,
		json:      json,
	}
}

func NewRedisPool(cfg *config.Config) (*redis.Pool, error) {
	redisPoolIdleTimeout, err := time.ParseDuration(cfg.RedisStorePoolIdleTimeout)
	if err != nil {
		return nil, service.SemanticError{Msg: "invalid value for redis_pool_idle_timeout, must be duration"}.E()
	}

	return &redis.Pool{
		MaxIdle:     cfg.RedisStorePoolMaxIdle,
		IdleTimeout: redisPoolIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", cfg.RedisStoreHost, cfg.RedisStorePort), redis.DialDatabase(cfg.RedisStoreDatabase))
		},
	}, nil
}

func (rc *Redis) Close() error {
	return rc.redisPool.Close()
}

func (rc *Redis) closeConn(ctx context.Context, conn redis.Conn) {
	err := conn.Close()
	if err != nil {
		err = service.SemanticError{Msg: "failed closing redis connection, this might result in memory leak", Err: err}.E()
		log.Warn(ctx, err.Error(), "error", err)
	}
}

// Implementation example for get simple value
func (rc *Redis) GetSimpleFunc(ctx context.Context, key string) (string, error) {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", redisUnavailableError{Err: err}.E()
	}

	return value, nil
}

// Implementation example for save complex value
func (rc *Redis) SaveComplexFunc(ctx context.Context, key string, value DummyType) error {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	body, err := rc.json.Marshal(value)
	if err != nil {
		return service.SemanticError{Err: err}.E()
	}

	_, err = conn.Do("SET", key, body)
	if err != nil {
		return redisUnavailableError{Err: err}.E()
	}

	return nil
}

// Implementation example for get complex value
func (rc *Redis) GetComplexFunc(ctx context.Context, key string) (*DummyType, error) {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	repl, err := conn.Do("GET", key)
	if err != nil {
		return nil, redisUnavailableError{Err: err}.E()
	}
	if repl == nil {
		return nil, nil
	}

	result := &DummyType{}
	err = rc.json.Unmarshal(repl.([]byte), result)
	if err != nil {
		return nil, service.SemanticError{Err: err}.E()
	}

	return result, nil
}

type DummyType struct {
	Test string `json:"test_dummy" redis:"dummy_test"`
}
