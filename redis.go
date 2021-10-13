package dbconnect

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisConfig defines the parameters of a redis connection
type RedisConfig struct {
	ID string `json:"id" toml:"id"`
	// Allowed "tcp", "unix". Default: "tcp"
	Network string `json:"network" toml:"network"`
	Host    string `json:"host" toml:"host"`
	Port    int    `json:"port" toml:"port"`
	Pwd     string `json:"pwd" toml:"pwd"`
	// RawURL defines a URL using the Redis URI scheme.
	// URLs should follow the draft IANA specification for the scheme
	// (https://www.iana.org/assignments/uri-schemes/prov/redis).
	// Addr when specified is preferred over this
	RawURL              string `json:"raw_url" toml:"raw_url"`
	DialTimeoutSeconds  int    `json:"dial_timeout_seconds" toml:"dial_timeout_seconds"`
	DB                  int    `json:"db" toml:"db"`
	KeepAliveMins       int    `json:"keep_alive_mins" toml:"keep_alive_mins"`
	ReadTimeoutSeconds  int    `json:"read_timeout_seconds" toml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `json:"write_timeout_seconds" toml:"write_timeout_seconds"`
	MaxIdle             int    `json:"max_idle" toml:"max_idle"`
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive       int `json:"max_active" toml:"max_active"`
	IdleTimeoutMins int `json:"idle_timeout_mins" toml:"idle_timeout_mins"`
	// If Wait is true and the pool is at the MaxActive limit, then Get() waits
	// for a connection to be returned to the pool before returning.
	Wait bool `json:"wait" toml:"wait"`
	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetimeSeconds int `json:"max_conn_lifetime_seconds" toml:"max_conn_lifetime_seconds"`
	_pool                  *redis.Pool
	once                   sync.Once
}

func (rc *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", rc.Host, rc.Port)
}

func (rc *RedisConfig) assert() error {
	if rc.Addr() == "" && rc.RawURL == "" {
		return fmt.Errorf("both addr and raw_url cannot be empty")
	}
	return nil
}

func (rc *RedisConfig) expandEnv() {
	rc.ID = os.ExpandEnv(rc.ID)
	rc.Network = os.ExpandEnv(rc.Network)
	rc.Host = os.ExpandEnv(rc.Host)
	rc.Pwd = os.ExpandEnv(rc.Pwd)
	rc.RawURL = os.ExpandEnv(rc.RawURL)
}

func (rc *RedisConfig) defaults() {
	if rc.Network == "" {
		rc.Network = "tcp"
	}

	if rc.DialTimeoutSeconds == 0 {
		rc.DialTimeoutSeconds = 5
	}

	if rc.KeepAliveMins == 0 {
		rc.KeepAliveMins = 5
	}

	if rc.ReadTimeoutSeconds == 0 {
		rc.ReadTimeoutSeconds = 3
	}

	if rc.WriteTimeoutSeconds == 0 {
		if rc.ReadTimeoutSeconds != 0 {
			rc.WriteTimeoutSeconds = rc.ReadTimeoutSeconds
		} else {
			rc.WriteTimeoutSeconds = 3
		}
	}
}

func (rc *RedisConfig) connect() {
	rc.once.Do(func() {
		rc.expandEnv()
		if err := rc.assert(); err != nil {
			log.Fatal(err.Error())
		}
		rc.defaults()
		rc._pool = &redis.Pool{
			MaxIdle:         rc.MaxIdle,
			MaxActive:       rc.MaxActive,
			IdleTimeout:     time.Duration(rc.IdleTimeoutMins) * time.Minute,
			Wait:            rc.Wait,
			MaxConnLifetime: time.Duration(rc.MaxConnLifetimeSeconds) * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				return nil
			},

			Dial: func() (redis.Conn, error) {
				dops := []redis.DialOption{
					redis.DialConnectTimeout(time.Duration(rc.DialTimeoutSeconds) * time.Second),
					redis.DialKeepAlive(time.Duration(rc.KeepAliveMins) * time.Minute),
					redis.DialReadTimeout(time.Duration(rc.ReadTimeoutSeconds) * time.Second),
					redis.DialWriteTimeout(time.Duration(rc.WriteTimeoutSeconds) * time.Second),
					redis.DialDatabase(rc.DB),
				}
				if rc.Pwd != "" {
					dops = append(dops, redis.DialPassword(rc.Pwd))
				}

				if rc.Addr() == "" && rc.RawURL != "" {
					return redis.DialURL(rc.RawURL, dops...)
				}

				return redis.Dial(rc.Network, rc.Addr(), dops...)
			},
		}
	})
}

func (rc *RedisConfig) pool() (*redis.Pool, error) {
	rc.connect()
	if rc._pool == nil {
		return nil, fmt.Errorf("invalid redis configuration. pool could not be created")
	}
	return rc._pool, nil
}
