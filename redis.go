package dbconnect

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisConfig defines the parameters of a redis connection
type RedisConfig struct {
	ID string `json:"id" toml:"id"`
	// Allowed "tcp", "unix". Default: "tcp"
	Network string `json:"network" toml:"network"`
	// The host to connect to
	// Format: IP:PORT
	Host string `json:"host" toml:"host"`
	Port int    `json:"port" toml:"port"`
	// Password to use when connecting with redis
	Pwd string `json:"pwd" toml:"pwd"`
	// RawURL defines a URL using the Redis URI scheme.
	// URLs should follow the draft IANA specification for the scheme
	// (https://www.iana.org/assignments/uri-schemes/prov/redis).
	// Addr when specified is preferred over this
	RawURL string `json:"raw_url" toml:"raw_url"`
	// DialTimeoutSeconds specifies the timeout for connecting to the Redis server
	DialTimeoutSeconds int `json:"dial_timeout_seconds" toml:"dial_timeout_seconds"`
	// DB defines which redis database to connect to
	DB int `json:"db" toml:"db"`
	// DialKeepAlive specifies the keep-alive period for TCP connections
	// to the Redis server when no DialNetDial option is specified.
	// If zero, keep-alives are not enabled. If no DialKeepAlive option
	// is specified then the default of 5 minutes is used to ensure
	// that half-closed TCP sessions are detected.
	KeepAliveMins int `json:"keep_alive_mins" toml:"keep_alive_mins"`
	// ReadTimeoutSeconds specifies the timeout for reading a single command.
	ReadTimeoutSeconds int `json:"read_timeout_seconds" toml:"read_timeout_seconds"`
	// WriteTimeoutSeconds specifies the timeout for writing a single command.
	WriteTimeoutSeconds int `json:"write_timeout_seconds" toml:"write_timeout_seconds"`
	// Maximum number of idle connections in the pool.
	MaxIdle int `json:"max_idle" toml:"max_idle"`
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int `json:"max_active" toml:"max_active"`
	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
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
