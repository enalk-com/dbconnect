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

type Conns struct {
	c        *Config
	pqMap    map[string]int
	mongoMap map[string]int
	redisMap map[string]int
	roachMap map[string]int
}

func New(p string) (*Conns, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var c Config
	switch path.Ext(p) {
	case ".json":
		if err := json.Unmarshal(bs, &c); err != nil {
			return nil, err
		}
	case ".toml":
		if _, err := toml.Decode(string(bs), &c); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid extension: \"%s\"", path.Ext(p))
	}

	conns := Conns{
		c: &c,
	}
	// create indices
	// Redis
	if c.Redis != nil && len(c.Redis) > 0 {
		if conns.redisMap == nil {
			conns.redisMap = map[string]int{}
		}
		for i, r := range c.Redis {
			conns.redisMap[r.ID] = i
		}
	}

	// mongo
	if c.Mongo != nil && len(c.Mongo) > 0 {
		if conns.mongoMap == nil {
			conns.mongoMap = map[string]int{}
		}
		for i, m := range c.Mongo {
			conns.mongoMap[m.ID] = i
		}
	}

	// pg
	if c.PQ != nil && len(c.PQ) > 0 {
		if conns.pqMap == nil {
			conns.pqMap = map[string]int{}
		}
		for i, pq := range c.PQ {
			conns.pqMap[pq.ID] = i
		}
	}

	// cockroachdb
	if c.CockroachDB != nil && len(c.CockroachDB) > 0 {
		if conns.roachMap == nil {
			conns.roachMap = map[string]int{}
		}

		for i, rc := range c.CockroachDB {
			conns.roachMap[rc.ID] = i
		}
	}
	return &conns, nil
}

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
