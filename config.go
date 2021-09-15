package dbconnect

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/BurntSushi/toml"
)

var pqMap map[string]*PQConfig
var mongoMap map[string]*MongoConfig
var redisMap map[string]*RedisConfig
var roachMap map[string]*RoachConfig
var once sync.Once = fetchOnce()

func fetchOnce() sync.Once {
	return sync.Once{}
}

// Config defines the overarching container for all supported databases
type Config struct {
	Redis       []*RedisConfig `json:"redis,omitempty" toml:"redis,omitempty"`
	PQ          []*PQConfig    `json:"pq,omitempty" toml:"pq,omitempty"`
	Mongo       []*MongoConfig `json:"mongo,omitempty" toml:"mongo,omitempty"`
	CockroachDB []*RoachConfig `json:"cockroachdb,omitempty" toml:"cockroachdb,omitempty"`
	mu          sync.Mutex
}

func initialize(v []byte, ext string) error {
	var c Config

	switch ext {
	case ".json":
		if err := json.Unmarshal(v, &c); err != nil {
			return err
		}
	case ".toml":
		if _, err := toml.Decode(string(v), &c); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid extension: \"%s\"", ext)
	}

	// Redis
	if c.Redis != nil && len(c.Redis) > 0 {
		if redisMap == nil {
			redisMap = map[string]*RedisConfig{}
		}
		for _, r := range c.Redis {
			redisMap[r.ID] = r
		}
	}

	// mongo
	if c.Mongo != nil && len(c.Mongo) > 0 {
		if mongoMap == nil {
			mongoMap = map[string]*MongoConfig{}
		}
		for _, m := range c.Mongo {
			mongoMap[m.ID] = m
		}
	}

	// pg
	if c.PQ != nil && len(c.PQ) > 0 {
		if pqMap == nil {
			pqMap = map[string]*PQConfig{}
		}
		for _, pq := range c.PQ {
			pqMap[pq.ID] = pq
		}
	}

	// cockroachdb
	if c.CockroachDB != nil && len(c.CockroachDB) > 0 {
		if roachMap == nil {
			roachMap = map[string]*RoachConfig{}
		}

		for _, rc := range c.CockroachDB {
			roachMap[rc.ID] = rc
		}
	}
	return nil
}

// Init initializes the provided configuration into a map for easy fetching later
func InitJSON(v []byte) error {
	if err := initialize(v, ".json"); err != nil {
		return err
	}
	return nil
}

// Init initializes the provided configuration into a map for easy fetching later
func InitTOML(v []byte) error {
	if err := initialize(v, ".toml"); err != nil {
		return err
	}
	return nil
}

func InitFile(p string) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err := initialize(bs, path.Ext(p)); err != nil {
		return err
	}
	return nil
}
