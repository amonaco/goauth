package cache

import (
	"errors"

	"github.com/gomodule/redigo/redis"

	"github.com/amonaco/goauth/lib/config"
)

var pool *redis.Pool

// Start initializes the connections to Redis.
func Start() {
	conf := config.Get()

	pool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.DialURL(conf.Redis)
		if err != nil {
			return nil, err
		}
		return c, nil
	}, conf.RedisMaxConn)
}

// Close closes the connections to Redis.
func Close() {
	if pool == nil {
		return
	}
	pool.Close()
	pool = nil
}

// Get retrieves a value from Redis.
func Get(key string) (string, error) {
	if pool == nil {
		return "", errors.New("cache: redis pool not initialized")
	}

	conn := pool.Get()
	defer conn.Close()
	return redis.String(conn.Do("GET", key))
}

// Set stores a value in Redis with a TTL in seconds.
func Set(key, value string, ttl int) error {
	if pool == nil {
		return errors.New("cache: redis pool not initialized")
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, "EX", ttl)
	return err
}

// Del deletes a key from Redis.
func Del(key string) error {
	if pool == nil {
		return errors.New("cache: redis pool not initialized")
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

// GetDel gets and deletes a key in a single transaction.
func GetDel(key string) (string, error) {
	if pool == nil {
		return "", errors.New("cache: redis pool not initialized")
	}

	conn := pool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("GET", key)
	conn.Send("DEL", key)

	reply, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return "", err
	}

	var res string
	_, err = redis.Scan(reply, &res)
	if err != nil {
		return "", err
	}

	return res, nil
}

// PushExpire pushes a value to a list and sets the TTL.
func PushExpire(key string, value string, ttl int) error {
	if pool == nil {
		return errors.New("cache: redis pool not initialized")
	}

	conn := pool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("LPUSH", key, value)
	conn.Send("EXPIRE", key, ttl)

	_, err := redis.Values(conn.Do("EXEC"))
	return err
}

// LRem deletes an item from a list.
func LRem(key string, token string) error {
	if pool == nil {
		return errors.New("cache: redis pool not initialized")
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("LREM", key, "0", token)
	return err
}
