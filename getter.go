package dbconnect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPQ returns a pointer to an pgxpool.Pool instance identified by input
func GetPQ(id string) (*pgxpool.Pool, error) {
	if pqMap == nil {
		return nil, fmt.Errorf("possibly no postgresql configurations provided")
	}

	if _, ok := pqMap[id]; !ok {
		return nil, fmt.Errorf("no postgreql configuration for ID: %s found", id)
	}

	return pqMap[id].db()
}

// GetPQConfig returns the configuration struct used to create this connection
func GetPQConfig(id string) ([]byte, error) {
	if pqMap == nil {
		return nil, fmt.Errorf("possibly no postgresql configurations provided")
	}

	if _, ok := pqMap[id]; !ok {
		return nil, fmt.Errorf("no postgreql configuration for ID: %s found", id)
	}

	b, err := json.Marshal(pqMap[id])
	if err != nil {
		return nil, fmt.Errorf("error marshalling configuration: %s", err.Error())
	}

	return b, nil
}

// GetRoach returns a pointer to an pgxpool.Pool instance identified by input
func GetRoach(id string) (*pgxpool.Pool, error) {
	if roachMap == nil {
		return nil, fmt.Errorf("possibly no cockroachdb configurations provided")
	}

	if _, ok := roachMap[id]; !ok {
		return nil, fmt.Errorf("no cockroach configuration for ID: %s found", id)
	}

	return roachMap[id].db()
}

// GetRoachConfig returns the configuration struct used to create this connection
func GetRoachConfig(id string) ([]byte, error) {
	if roachMap == nil {
		return nil, fmt.Errorf("possibly no cockroach db configurations provided")
	}

	if _, ok := roachMap[id]; !ok {
		return nil, fmt.Errorf("no cockroach db configuration for ID: %s found", id)
	}

	b, err := json.Marshal(roachMap[id])
	if err != nil {
		return nil, fmt.Errorf("error marshalling configuration: %s", err.Error())
	}

	return b, nil
}

// GetRedisPool returns a pointer to a redis.Pool instance identified by input
func GetRedisPool(id string) (*redis.Pool, error) {
	if redisMap == nil {
		return nil, fmt.Errorf("possibly no redis configurations provided")
	}

	if _, ok := redisMap[id]; !ok {
		return nil, fmt.Errorf("no redis configuration for ID: %s found", id)
	}

	return redisMap[id].pool()
}

// GetRedisConn is a conveninece function; returns a redis.Conn instance identified by input
// It is imperative that you close this connection after your work is done
// e.g.
//   conn, _ := dbconnect.GetRedisConn("redis_main")
//   defer conn.Close()
func GetRedisConn(id string) (redis.Conn, error) {
	p, err := GetRedisPool(id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nYou've requested a redis.Conn instance. Don't forget to close it after your work is done\n\n")
	return p.Get(), nil
}

// GetRedisPubSubConn is a convenience function; returns a redis.PubSubConn instance identified by input
func GetRedisPubSubConn(id string) (redis.PubSubConn, error) {
	c, err := GetRedisConn(id)
	if err != nil {
		return redis.PubSubConn{}, err
	}

	return redis.PubSubConn{
		Conn: c,
	}, nil
}

// GetRedisConfig returns the configuration used to create this connection
func GetRedisConfig(id string) ([]byte, error) {
	if redisMap == nil {
		return nil, fmt.Errorf("possibly no redis configurations provided")
	}

	if _, ok := redisMap[id]; !ok {
		return nil, fmt.Errorf("no redis configuration for ID: %s found", id)
	}

	b, err := json.Marshal(redisMap[id])
	if err != nil {
		return nil, fmt.Errorf("error marshalling configuration: %s", err.Error())
	}
	return b, nil
}

// GetMongoClient returns a pointer to a mongo.Client instance identified by input
func GetMongoClient(ctx context.Context, id string) (*mongo.Client, error) {
	if mongoMap == nil {
		return nil, fmt.Errorf("possibly no mongo configurations provided")
	}

	if _, ok := mongoMap[id]; !ok {
		return nil, fmt.Errorf("no mongo configuration for ID: %s found", id)
	}

	return mongoMap[id].client(ctx)
}

// GetMongoDB returns a pointer to a mongo.Database instance identified by input
func GetMongoDB(ctx context.Context, id string, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	if mongoMap == nil {
		return nil, fmt.Errorf("possibly no mongo configurations provided")
	}

	if _, ok := mongoMap[id]; !ok {
		return nil, fmt.Errorf("no mongo configuration for ID: %s found", id)
	}

	return mongoMap[id].db(ctx, opts...)
}

// GetMongoConfig returns the configuration used to create this connection
func GetMongoConfig(id string) ([]byte, error) {
	if mongoMap == nil {
		return nil, fmt.Errorf("possibly no mongo configurations provided")
	}

	if _, ok := mongoMap[id]; !ok {
		return nil, fmt.Errorf("no mongo configuration for ID: %s found", id)
	}

	b, err := json.Marshal(mongoMap[id])
	if err != nil {
		return nil, fmt.Errorf("error marshalling configuration: %s", err.Error())
	}
	return b, nil
}
