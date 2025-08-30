package cache

import (
	"github.com/gomodule/redigo/redis"

	"github.com/amonaco/goauth/libs/config"
)

var pool *redis.Pool

// Start initializes the connections to redis
func Start() {
	conf := config.Get()

	pool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.DialURL(conf.Redis.Address)
		if err != nil {
			return nil, err
		}

		return c, err
	}, conf.Redis.MaxConn)
}

// Close closes the connections to redis
func Close() {
	pool.Close()
}

// Get retreives the value of a key from redis
func Get(key string) (string, error) {
	conn := pool.Get()
	return redis.String(conn.Do("GET", key))
}

// Set sets the value of a key in redis with a time to live in seconds
func Set(key, value string, ttl int) error {
	conn := pool.Get()
	_, err := conn.Do("SET", key, value, "EX", ttl)
	return err
}

// Del deletes a key from redis store
func Del(key string) error {
	conn := pool.Get()
	_, err := conn.Do("DEL", key)
	return err
}

// Gets and deletes a key in single transaction
func GetDel(key string) (string, error) {
	var res string
	conn := pool.Get()
	conn.Send("MULTI")
	conn.Send("GET", key)
	conn.Send("DEL", key)

	reply, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return "", err
	}

	_, err = redis.Scan(reply, &res)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Pushes to a list and sets its expiry (can't be done in a single operation)
func PushExpire(key string, value string, ttl int) error {
	conn := pool.Get()
	conn.Send("MULTI")
	conn.Send("LPUSH", key, value)
	conn.Send("EXPIRE", key, ttl)

	_, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return err
	}

	return nil
}

// Deletes an item from list
func LRem(key string, token string) error {
	conn := pool.Get()
	_, err := conn.Do("LREM", key, "0", token)
	return err
}
