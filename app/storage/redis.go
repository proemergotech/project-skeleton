package storage

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"gitlab.com/proemergotech/log-go/v3"
	"gitlab.com/proemergotech/uuid-go"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

// todo: remove
//  Example for lua script
// KEYS[1]: group:%s:dummy:%s:data
// KEYS[2]: group:%s:uuid:%s:dummy
// ARGV[1]: <dummy_data>
// ARGV[2]: group:%s:dummy:%s
var saveScript = redis.NewScript(2, `
redis.call("SET", KEYS[1], ARGV[1])
redis.call("SET", KEYS[2], ARGV[2])
`)

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

func (r *Redis) Close() error {
	return r.redisPool.Close()
}

func (r *Redis) closeConn(ctx context.Context, conn redis.Conn) {
	if err := conn.Close(); err != nil {
		err = service.SemanticError{Msg: "failed closing redis connection, this might result in memory leak", Err: err}.E()
		log.Warn(ctx, err.Error(), "error", err)
	}
}

// todo: remove
//  Examples for key generation
func (r *Redis) dummyID(group string, uuid uuid.UUID) string {
	return fmt.Sprintf("group:%s:dummy:%s", group, uuid)
}

func (r *Redis) dummyKey(dummyID string) string {
	return fmt.Sprintf("%s:data", dummyID)
}

func (r *Redis) dummyTestKey(group string, test uuid.UUID) string {
	return fmt.Sprintf("group:%s:test:%s:dummy", group, test)
}

// todo: remove
//  Implementation example for get simple value
func (r *Redis) GetSimpleFunc(ctx context.Context, key string) (string, error) {
	conn := r.redisPool.Get()
	defer r.closeConn(ctx, conn)

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", redisError{Err: err}.E()
	}

	return value, nil
}

// todo: remove
//  Implementation example for save complex value
func (r *Redis) SaveDummy(ctx context.Context, dummy *service.DummyType) error {
	conn := r.redisPool.Get()
	defer r.closeConn(ctx, conn)

	body, err := r.json.Marshal(dummy)
	if err != nil {
		return service.SemanticError{Err: err}.E()
	}

	dummyID := r.dummyID(dummy.Group, dummy.UUID)
	_, err = saveScript.Do(
		conn,
		r.dummyKey(dummyID),
		r.dummyTestKey(dummy.Group, dummy.TestUUID),
		body,
		dummyID,
	)
	if err != nil {
		return redisError{Err: err}.E()
	}

	return nil
}
