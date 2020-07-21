package redis

import (
	"errors"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"keywea.com/cloud/pblib/pbconfig"
	"keywea.com/cloud/pblib/pb/events"
	"sync"
	"time"
)

const (
	REDIS_CMD_SADD            = "SADD"
	REDIS_CMD_SCARD           = "SCARD"
	REDIS_CMD_SISMEMBER       = "SISMEMBER"
	REDIS_CMD_SMEMBERS        = "SMEMBERS"
	REDIS_CMD_SREM            = "SREM"
	REDIS_CMD_HSET            = "HSET"
	REDIS_CMD_HGET            = "HGET"
	REDIS_CMD_HMSET           = "HMSET"
	REDIS_CMD_HMGET           = "HMGET"
	REDIS_CMD_HDEL            = "HDEL"
	REDIS_CMD_HGETALL         = "HGETALL"
	REDIS_CMD_SET             = "SET"
	REDIS_CMD_SETNX           = "SETNX"
	REDIS_CMD_SETEX           = "SETEX"
	REDIS_CMD_GET             = "GET"
	REDIS_CMD_TTL             = "TTL"
	REDIS_CMD_STRLEN          = "STRLEN"
	REDIS_CMD_EXPIRE          = "EXPIRE"
	REDIS_CMD_DELETE          = "DEL"
	REDIS_CMD_KEYS            = "KEYS"
	REDIS_CMD_HKEYS           = "HKEYS"
	REDIS_CMD_EXISTS          = "EXISTS"
	REDIS_CMD_PERSIST         = "PERSIST"
	REDIS_CMD_ZADD            = "ZADD"
	REDIS_CMD_ZREM            = "ZREM"
	REDIS_CMD_ZRANGE          = "ZRANGE"
	REDIS_CMD_ZRANGE_BY_SCORE = "ZRANGEBYSCORE"
	REDIS_CMD_WITHSCORES      = "WITHSCORES"
	REDIS_CMD_INCR            = "INCR"
	REDIS_CMD_DECR            = "DECR"
	REDIS_CMD_INCRBY          = "INCRBY"
	REDIS_CMD_DECRBY          = "DECRBY"
	REDIS_CMD_INCRBYFLOAT     = "INCRBYFLOAT"
	REDIS_CMD_DECRBYFLOAT     = "DECRBYFLOAT"
	REDIS_CMD_PUBLISH     	  = "PUBLISH"
)

var (
	errNotFoundRedisPool = errors.New("Redis Pool Not Found")

	dialFunc = func(network, address string, dialOptions []redigo.DialOption) func() (redigo.Conn, error) {
		return func() (redigo.Conn, error) {
			conn, err := redigo.Dial(network, address, dialOptions...)
			return conn, err
		}
	}

	testOnBorrowFunc = func(c redigo.Conn, t time.Time) error {
		_, err := c.Do("ping")
		if err != nil {
			return err
		}
		return nil
	}

	testFunc = func(p *redigo.Pool) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("%+v", r))
			}

		}()
		c := p.Get()
		defer c.Close()
		return c.Err()
	}

	rediS *pbredis
	rmu sync.Mutex
)

type PoolConfig struct {
	Name 	          string
	Network           string
	Server 			  string
	ConnectionTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	Password          string
	DB                int
	MaxIdle           int
	MaxActive         int
	TestOnBorrow      bool
	IdleTimeout       time.Duration
	Wait              bool
}

type pbredis struct {
	redisPool map[string]*RPool
	poolConfigs map[string]PoolConfig

	mu sync.Mutex
}

func NewPool(name string, configor pbconfig.Configor) (*RPool, error) {
	rmu.Lock()
	if rediS == nil {
		rediS = &pbredis{
			redisPool: make(map[string]*RPool),
			poolConfigs: make(map[string]PoolConfig),
		}
		events.AddShutdownHook(func() error {
			rediS.Destroy()
			return nil
		}, events.SHUTDOWN_INDEX_REDIS)
	}
	rmu.Unlock()
	return rediS.createPool(name, configor)
}

func (r *pbredis) createPool(name string, configor pbconfig.Configor) (*RPool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if configor == nil {
		return nil, fmt.Errorf("Redis New Pool=%s create Error on nil configor", name)
	}
	config := r.parseConfig(configor)
	config.Name = name
	if rpool, ok := r.redisPool[config.Name]; ok { // pool name exists
		return rpool.UpdatePool(config)
	}
	return r.create(config)
}

func (r *pbredis) create(config PoolConfig) (*RPool, error) {
	pool := &redigo.Pool{}
	dialOptions := []redigo.DialOption{}
	if config.ConnectionTimeout > 0 {
		dialOptions = append(dialOptions, redigo.DialConnectTimeout(time.Second*config.ConnectionTimeout))
	}
	if config.ReadTimeout > 0 {
		dialOptions = append(dialOptions, redigo.DialReadTimeout(time.Second*config.ReadTimeout))
	}
	if config.WriteTimeout > 0 {
		dialOptions = append(dialOptions, redigo.DialWriteTimeout(time.Second*config.WriteTimeout))
	}
	if config.Password != "" {
		dialOptions = append(dialOptions, redigo.DialPassword(config.Password))
	}
	dialOptions = append(dialOptions, redigo.DialDatabase(config.DB))

	pool.Dial = dialFunc(config.Network, config.Server, dialOptions)

	if config.MaxIdle > 0 {
		pool.MaxIdle = config.MaxIdle
	}
	if config.MaxActive > 0 {
		pool.MaxActive = config.MaxActive
	}
	if config.TestOnBorrow {
		pool.TestOnBorrow = testOnBorrowFunc
	}

	if config.IdleTimeout > 0 {
		pool.IdleTimeout = time.Second * config.IdleTimeout
	}
	pool.Wait = config.Wait

	if err := testFunc(pool); err != nil {
		return nil, err
	}

	r.redisPool[config.Name] = &RPool{
		name: config.Name,
		pool: pool,
	}
	r.poolConfigs[config.Name] = config

	return r.redisPool[config.Name], nil
}

func (r *pbredis) parseConfig(configor pbconfig.Configor) PoolConfig {
	connectionTimeout, _ := configor.GetInt("connectionTimeout", 10)
	readTimeout, _ := configor.GetInt("readTimeout", 30)
	writeTimeout, _ := configor.GetInt("writeTimeout", 20)
	maxIdle, _ := configor.GetInt("maxIdle", 10)
	maxActive, _ := configor.GetInt("maxActive", 500)
	idleTimeout, _ := configor.GetInt("idleTimeout", 300) // 5min
	testOnBorrow, _ := configor.GetBool("testOnBorrow")
	wait, _ := configor.GetBool("wait")
	db, _ := configor.GetInt("db", 0)

	return PoolConfig{
		Network: configor.GetString("network", "tcp"),
		Server: configor.GetString("server", "127.0.0.1:6379"),
		ConnectionTimeout: time.Duration(connectionTimeout),
		ReadTimeout: time.Duration(readTimeout),
		WriteTimeout: time.Duration(writeTimeout),
		Password: configor.GetString("password", ""),
		DB: db,
		MaxIdle: maxIdle,
		MaxActive: maxActive,
		IdleTimeout: time.Duration(idleTimeout),
		TestOnBorrow: testOnBorrow,
		Wait: wait,
	}
}

func (r *pbredis) Get(name string) redigo.Conn {
	p, ok := r.redisPool[name]
	if !ok {
		panic(errNotFoundRedisPool)
	}
	return p.GetConn()
}

func (r *pbredis) Destroy() {
	for _, v := range r.redisPool {
		err := v.Destroy()
		if err != nil {
		}
	}
}

// 部分command实现
func Expire(RConn *redigo.Conn, key string, ttl int) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_EXPIRE, key, ttl)
}
func Persist(RConn *redigo.Conn, key string) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_PERSIST, key)
}
func Delete(RConn *redigo.Conn, key string) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_DELETE, key)
}
func Set(RConn *redigo.Conn, key string, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_SET, key, data)
}
func SetNX(RConn *redigo.Conn, key string, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_SETNX, key, data)
}
func SetEx(RConn *redigo.Conn, key string, ttl int, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_SETEX, key, ttl, data)
}
func Get(RConn *redigo.Conn, key string) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_GET, key)
}
func GetTTL(RConn *redigo.Conn, key string) (time.Duration, error) {
	ttl, err := redigo.Int64((*RConn).Do(REDIS_CMD_TTL, key))
	return time.Duration(ttl) * time.Second, err
}
func GetString(RConn *redigo.Conn, key string) (string, error) {
	return redigo.String((*RConn).Do(REDIS_CMD_GET, key))
}
func GetInt(RConn *redigo.Conn, key string) (int, error) {
	return redigo.Int((*RConn).Do(REDIS_CMD_GET, key))
}
func GetInt64(RConn *redigo.Conn, key string) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_GET, key))
}
func GetUInt64(RConn *redigo.Conn, key string) (uint64, error) {
	return redigo.Uint64((*RConn).Do(REDIS_CMD_GET, key))
}
func GetStringLength(RConn *redigo.Conn, key string) (int, error) {
	return redigo.Int((*RConn).Do(REDIS_CMD_STRLEN, key))
}
func ZAdd(RConn *redigo.Conn, key string, score float64, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_ZADD, key, score, data)
}
func ZRem(RConn *redigo.Conn, key string, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_ZREM, key, data)
}
func ZRange(RConn *redigo.Conn, key string, start int, end int, withScores bool) ([]interface{}, error) {
	if withScores {
		return redigo.Values((*RConn).Do(REDIS_CMD_ZRANGE, key, start, end, REDIS_CMD_WITHSCORES))
	}
	return redigo.Values((*RConn).Do(REDIS_CMD_ZRANGE, key, start, end))
}
func SAdd(RConn *redigo.Conn, setName string, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_SADD, setName, data)
}
func SCard(RConn *redigo.Conn, setName string) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_SCARD, setName))
}
func SIsMember(RConn *redigo.Conn, setName string, data interface{}) (bool, error) {
	return redigo.Bool((*RConn).Do(REDIS_CMD_SISMEMBER, setName, data))
}
func SMembers(RConn *redigo.Conn, setName string) ([]string, error) {
	return redigo.Strings((*RConn).Do(REDIS_CMD_SMEMBERS, setName))
}
func SRem(RConn *redigo.Conn, setName string, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_SREM, setName, data)
}
func HSet(RConn *redigo.Conn, key string, HKey string, data interface{}) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_HSET, key, HKey, data)
}

func HGet(RConn *redigo.Conn, key string, HKey string) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_HGET, key, HKey)
}

func HMGet(RConn *redigo.Conn, key string, hashKeys ...string) ([]interface{}, error) {
	ret, err := (*RConn).Do(REDIS_CMD_HMGET, key, hashKeys)
	if err != nil {
		return nil, err
	}
	reta, ok := ret.([]interface{})
	if !ok {
		return nil, errors.New("result not an array")
	}
	return reta, nil
}

func HMSet(RConn *redigo.Conn, key string, hashKeys []string, vals []interface{}) (interface{}, error) {
	if len(hashKeys) == 0 || len(hashKeys) != len(vals) {
		var ret interface{}
		return ret, errors.New("bad length")
	}
	input := []interface{}{key}
	for i, v := range hashKeys {
		input = append(input, v, vals[i])
	}
	return (*RConn).Do(REDIS_CMD_HMSET, input...)
}

func HGetString(RConn *redigo.Conn, key string, HKey string) (string, error) {
	return redigo.String((*RConn).Do(REDIS_CMD_HGET, key, HKey))
}
func HGetFloat(RConn *redigo.Conn, key string, HKey string) (float64, error) {
	f, err := redigo.Float64((*RConn).Do(REDIS_CMD_HGET, key, HKey))
	return float64(f), err
}
func HGetInt(RConn *redigo.Conn, key string, HKey string) (int, error) {
	return redigo.Int((*RConn).Do(REDIS_CMD_HGET, key, HKey))
}
func HGetInt64(RConn *redigo.Conn, key string, HKey string) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_HGET, key, HKey))
}
func HGetBool(RConn *redigo.Conn, key string, HKey string) (bool, error) {
	return redigo.Bool((*RConn).Do(REDIS_CMD_HGET, key, HKey))
}
func HDel(RConn *redigo.Conn, key string, HKey string) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_HDEL, key, HKey)
}
func HGetAll(RConn *redigo.Conn, key string) (interface{}, error) {
	return (*RConn).Do(REDIS_CMD_HGETALL, key)
}

func HGetAllValues(RConn *redigo.Conn, key string) ([]interface{}, error) {
	return redigo.Values((*RConn).Do(REDIS_CMD_HGETALL, key))
}
func HGetAllString(RConn *redigo.Conn, key string) ([]string, error) {
	return redigo.Strings((*RConn).Do(REDIS_CMD_HGETALL, key))
}

func Keys(RConn *redigo.Conn, pattern string) ([]string, error) {
	return redigo.Strings((*RConn).Do(REDIS_CMD_KEYS, pattern))
}

func HKeys(RConn *redigo.Conn, key string) ([]string, error) {
	return redigo.Strings((*RConn).Do(REDIS_CMD_HKEYS, key))
}

func Exists(RConn *redigo.Conn, key string) (bool, error) {
	count, err := redigo.Int((*RConn).Do(REDIS_CMD_EXISTS, key))
	if count == 0 {
		return false, err
	} else {
		return true, err
	}
}

func Incr(RConn *redigo.Conn, key string) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_INCR, key))
}

func Decr(RConn *redigo.Conn, key string) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_DECR, key))
}

func IncrBy(RConn *redigo.Conn, key string, incBy int64) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_INCRBY, key, incBy))
}

func DecrBy(RConn *redigo.Conn, key string, decrBy int64) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_DECRBY, key))
}

func IncrByFloat(RConn *redigo.Conn, key string, incBy float64) (float64, error) {
	return redigo.Float64((*RConn).Do(REDIS_CMD_INCRBYFLOAT, key, incBy))
}

func DecrByFloat(RConn *redigo.Conn, key string, decrBy float64) (float64, error) {
	return redigo.Float64((*RConn).Do(REDIS_CMD_DECRBYFLOAT, key, decrBy))
}

func Publish(RConn *redigo.Conn, key string, message interface{}) (int64, error) {
	return redigo.Int64((*RConn).Do(REDIS_CMD_PUBLISH, key, message))
}