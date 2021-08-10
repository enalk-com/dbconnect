package dbconnect

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

var pqMap map[string]*PQConfig
var mongoMap map[string]*MongoConfig
var redisMap map[string]*RedisConfig
var once sync.Once = fetchOnce()

func fetchOnce() sync.Once {
	return sync.Once{}
}

// Config defines the overarching container for all supported databases
type Config struct {
	Redis []*RedisConfig `json:"redis,omitempty"`
	PQ    []*PQConfig    `json:"pq,omitempty"`
	Mongo []*MongoConfig `json:"mongo,omitempty"`
	mu    sync.Mutex
}

func initialize(v []byte) error {
	var c Config
	if err := json.Unmarshal(v, &c); err != nil {
		return err
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
	return nil
}

// Init initializes the provided configuration into a map for easy fetching later
func Init(v []byte) error {
	if err := initialize(v); err != nil {
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
	if err := Init(bs); err != nil {
		return err
	}
	return nil
}

// UpdateConfig updates a pre-initialized configuration.
// You can overwrite existing config IDs by providing a new
// configuration for the same ID.<db-type>
func UpdateConfig(v []byte) error {
	once = fetchOnce()
	if err := initialize(v); err != nil {
		return err
	}

	return nil
}
