package storage

import (
	"context"

	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	log "gitlab.com/proemergotech/log-go/v2"
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

func (rc *Redis) Close() error {
	return rc.redisPool.Close()
}

func (rc *Redis) closeConn(ctx context.Context, conn redis.Conn) {
	if err := conn.Close(); err != nil {
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
		return "", redisError{Err: err}.E()
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
		return redisError{Err: err}.E()
	}

	return nil
}

// Implementation example for get complex value
func (rc *Redis) GetComplexFunc(ctx context.Context, key string) (*DummyType, error) {
	conn := rc.redisPool.Get()
	defer rc.closeConn(ctx, conn)

	repl, err := conn.Do("GET", key)
	if err != nil {
		return nil, redisError{Err: err}.E()
	}
	if repl == nil {
		return nil, nil
	}

	result := &DummyType{}
	if err := rc.json.Unmarshal(repl.([]byte), result); err != nil {
		return nil, service.SemanticError{Err: err}.E()
	}

	return result, nil
}

type DummyType struct {
	Test string `json:"test_dummy" redis:"dummy_test"`
}
