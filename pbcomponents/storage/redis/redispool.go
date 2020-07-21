package redis

import (
	"errors"
	"keywea.com/cloud/pblib/pbconfig"
	"sync"
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

type RPool struct {
	name string
	pool *redigo.Pool

	closed bool
	rpmu sync.Mutex
}

func (rpool *RPool) ParseConfig(configor *pbconfig.Configor) PoolConfig {
	return rediS.parseConfig(*configor)
}

func (rpool *RPool) UpdatePool(config PoolConfig) (*RPool, error) {
	rmu.Lock()
	defer rmu.Unlock()

	oldConfig := rediS.poolConfigs[config.Name]

	if oldConfig.Network != config.Network || oldConfig.Server != config.Server ||
		oldConfig.ConnectionTimeout != config.ConnectionTimeout || oldConfig.ReadTimeout != config.ReadTimeout ||
		oldConfig.WriteTimeout != config.WriteTimeout || oldConfig.Password != config.Password ||
		oldConfig.DB != config.DB { // reset pool
		err := rpool.Destroy()
		if err != nil {
		}
		delete(rediS.redisPool, config.Name)
		delete(rediS.poolConfigs, config.Name)
		return rediS.create(config)
	}

	if config.MaxIdle > 0 && oldConfig.MaxIdle != config.MaxIdle {
		rpool.pool.MaxIdle = config.MaxIdle
	}
	if config.MaxActive > 0 && oldConfig.MaxActive != config.MaxActive {
		rpool.pool.MaxActive = config.MaxActive
	}
	if config.TestOnBorrow {
		rpool.pool.TestOnBorrow = testOnBorrowFunc
	} else {
		rpool.pool.TestOnBorrow = nil
	}

	if config.IdleTimeout > 0 && oldConfig.IdleTimeout != config.IdleTimeout {
		rpool.pool.IdleTimeout = time.Second * config.IdleTimeout
	}
	rpool.pool.Wait = config.Wait
	delete(rediS.poolConfigs, config.Name)
	rediS.poolConfigs[config.Name] = config

	return rpool, nil
}

func (rpool *RPool) GetConn() redigo.Conn {
	return rpool.pool.Get()
}

func (rpool *RPool) Destroy() error {
	rpool.rpmu.Lock()
	defer rpool.rpmu.Unlock()
	if rpool.closed {
		return nil
	}
	rpool.closed = true
	return rpool.pool.Close()
}

// command
func (rpool *RPool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(commandName, args...)
}

func (rpool *RPool) Expire(key string, ttl int) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_EXPIRE, key, ttl)
}

func (rpool *RPool) GetTTL(key string) (time.Duration, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	ttl, err := redigo.Int64(conn.Do(REDIS_CMD_TTL, key))
	return time.Duration(ttl) * time.Second, err
}

func (rpool *RPool) Delete(key string) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_DELETE, key)
}

func (rpool *RPool) Set(key string, data interface{}) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_SET, key, data)
}

func (rpool *RPool) Get(key string) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_GET, key)
}

func (rpool *RPool) GetString(key string) (string, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.String(conn.Do(REDIS_CMD_GET, key))
}

func (rpool *RPool) GetBytes(key string) ([]byte, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Bytes(conn.Do(REDIS_CMD_GET, key))
}

func (rpool *RPool) GetInt(key string) (int, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int(conn.Do(REDIS_CMD_GET, key))
}

func (rpool *RPool) GetInt64(key string) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_GET, key))
}

func (rpool *RPool) GetUint64(key string) (uint64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Uint64(conn.Do(REDIS_CMD_GET, key))
}

func (rpool *RPool) Keys(pattern string) ([]string, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Strings(conn.Do(REDIS_CMD_KEYS, pattern))
}

func (rpool *RPool) KeysByteSlices(pattern string) ([][]byte, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.ByteSlices(conn.Do(REDIS_CMD_KEYS, pattern))
}

func (rpool *RPool) HKeys(key string) ([]string, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Strings(conn.Do(REDIS_CMD_HKEYS, key))
}

func (rpool *RPool) Exists(key string) (bool, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	count, err := redigo.Int(conn.Do(REDIS_CMD_EXISTS, key))
	if count == 0 {
		return false, err
	} else {
		return true, err
	}
}

func (rpool *RPool) Incr(key string) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_INCR, key))
}

func (rpool *RPool) Decr(key string) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_DECR, key))
}

func (rpool *RPool) IncrBy(key string, incBy int64) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_INCRBY, key, incBy))
}

func (rpool *RPool) DecrBy(key string, decrBy int64) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_DECRBY, key))
}

func (rpool *RPool) IncrByFloat(key string, incBy float64) (float64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Float64(conn.Do(REDIS_CMD_INCRBYFLOAT, key, incBy))
}

func (rpool *RPool) DecrByFloat(key string, decrBy float64) (float64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Float64(conn.Do(REDIS_CMD_DECRBYFLOAT, key, decrBy))
}

func (rpool *RPool) Publish(key string, message interface{}) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_PUBLISH, key, message))
}

// hash map
func (rpool *RPool) HSet(key string, HKey string, data interface{}) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_HSET, key, HKey, data)
}

func (rpool *RPool) HGet(key string, HKey string) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_HGET, key, HKey)
}

func (rpool *RPool) HMGet(key string, hashKeys ...string) ([]interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	ret, err := conn.Do(REDIS_CMD_HMGET, key, hashKeys)
	if err != nil {
		return nil, err
	}
	reta, ok := ret.([]interface{})
	if !ok {
		return nil, errors.New("result not an array")
	}
	return reta, nil
}

func (rpool *RPool) HMSet(key string, hashKeys []string, vals []interface{}) (interface{}, error) {
	if len(hashKeys) == 0 || len(hashKeys) != len(vals) {
		var ret interface{}
		return ret, errors.New("bad length")
	}
	input := []interface{}{key}
	for i, v := range hashKeys {
		input = append(input, v, vals[i])
	}
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_HMSET, input...)
}

func (rpool *RPool) HGetString(key string, HKey string) (string, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.String(conn.Do(REDIS_CMD_HGET, key, HKey))
}
func (rpool *RPool) HGetFloat(key string, HKey string) (float64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	f, err := redigo.Float64(conn.Do(REDIS_CMD_HGET, key, HKey))
	return float64(f), err
}
func (rpool *RPool) HGetInt(key string, HKey string) (int, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int(conn.Do(REDIS_CMD_HGET, key, HKey))
}
func (rpool *RPool) HGetInt64(key string, HKey string) (int64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Int64(conn.Do(REDIS_CMD_HGET, key, HKey))
}
func (rpool *RPool) HGetUint64(key string, HKey string) (uint64, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Uint64(conn.Do(REDIS_CMD_HGET, key, HKey))
}
func (rpool *RPool) HGetBool(key string, HKey string) (bool, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Bool(conn.Do(REDIS_CMD_HGET, key, HKey))
}
func (rpool *RPool) HDel(key string, HKey string) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_HDEL, key, HKey)
}
func (rpool *RPool) HGetAll(key string) (interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return conn.Do(REDIS_CMD_HGETALL, key)
}

func (rpool *RPool) HGetAllValues(key string) ([]interface{}, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Values(conn.Do(REDIS_CMD_HGETALL, key))
}
func (rpool *RPool) HGetAllString(key string) ([]string, error) {
	conn := rpool.GetConn()
	defer conn.Close()
	return redigo.Strings(conn.Do(REDIS_CMD_HGETALL, key))
}
