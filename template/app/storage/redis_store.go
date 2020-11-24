//%: {{ if .RedisStore }}
package storage

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/uuid"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

//%: {{ if .Examples }}
// todo: remove
//  Example for lua script
// KEYS[1]: group:%s:dummy:%s:data
// KEYS[2]: group:%s:uuid:%s:dummy
// ARGV[1]: <dummy_data>
// ARGV[2]: group:%s:dummy:%s
var saveStoreScript = redis.NewScript(2, `
redis.call("SET", KEYS[1], ARGV[1])
redis.call("SET", KEYS[2], ARGV[2])
`)

//%: {{ end }}

type RedisStore struct {
	redisPool *redis.Pool
	json      jsoniter.API //Use this to be able to save objects as value and use redis tags instead of json ones
}

func NewRedisStore(redisPool *redis.Pool, json jsoniter.API) *RedisStore {
	return &RedisStore{
		redisPool: redisPool,
		json:      json,
	}
}

func (r *RedisStore) Close() error {
	return r.redisPool.Close()
}

func (r *RedisStore) closeConn(ctx context.Context, conn redis.Conn) {
	if err := conn.Close(); err != nil {
		//%:{{ `
		err = skeleton.SemanticError{Msg: "failed closing redis connection, this might result in memory leak", Err: err}.E()
		//%: ` | replace "skeleton" .SchemaPackage }}
		log.Warn(ctx, err.Error(), "error", err)
	}
}

//%: {{ if .Examples }}
// todo: remove
//  Examples for key generation
func (r *RedisStore) dummyID(group string, dummyUUID uuid.UUID) string {
	return fmt.Sprintf("group:%s:dummy:%s", group, dummyUUID)
}

func (r *RedisStore) dummyKey(dummyID string) string {
	return fmt.Sprintf("%s:data", dummyID)
}

func (r *RedisStore) dummyTestKey(group string, testUUID uuid.UUID) string {
	return fmt.Sprintf("group:%s:test:%s:dummy", group, testUUID)
}

// todo: remove
//  Implementation example for get simple value
func (r *RedisStore) GetSimpleFunc(ctx context.Context, key string) (string, error) {
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
//%:{{ `
func (r *RedisStore) SaveDummy(ctx context.Context, dummy *skeleton.DummyType) error {
	//%: ` | replace "skeleton" .SchemaPackage | trim }}
	conn := r.redisPool.Get()
	defer r.closeConn(ctx, conn)

	body, err := r.json.Marshal(dummy)
	if err != nil {
		//%:{{ `
		return skeleton.SemanticError{Err: err}.E()
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
	}

	dummyID := r.dummyID(dummy.Group, dummy.DummyUUID)
	_, err = saveStoreScript.Do(
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

//%: {{ end }}

//%: {{ end }}
