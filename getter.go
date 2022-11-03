package dbconnect

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetPQ returns a pointer to an pgxpool.Pool instance identified by input
func (c Conns) GetPQ(id string) (*pgxpool.Pool, error) {
	if c.pqMap == nil {
		return nil, fmt.Errorf("possibly no postgresql configurations provided")
	}

	if _, ok := c.pqMap[id]; !ok {
		return nil, fmt.Errorf("no postgreql configuration for ID: %s found", id)
	}

	return c.c.PQ[c.pqMap[id]].db()
}

// GetRoach returns a pointer to an pgxpool.Pool instance identified by input
func (c Conns) GetRoach(id string) (*pgxpool.Pool, error) {
	if c.roachMap == nil {
		return nil, fmt.Errorf("possibly no cockroachdb configurations provided")
	}

	if _, ok := c.roachMap[id]; !ok {
		return nil, fmt.Errorf("no cockroach configuration for ID: %s found", id)
	}

	return c.c.CockroachDB[c.roachMap[id]].db()
}

// GetRedisPool returns a pointer to a redis.Pool instance identified by input
func (c Conns) GetRedisPool(id string) (*redis.Pool, error) {
	if c.redisMap == nil {
		return nil, fmt.Errorf("possibly no redis configurations provided")
	}

	if _, ok := c.redisMap[id]; !ok {
		return nil, fmt.Errorf("no redis configuration for ID: %s found", id)
	}

	return c.c.Redis[c.redisMap[id]].pool()
}

// GetRedisConn is a conveninece function; returns a redis.Conn instance identified by input
// It is imperative that you close this connection after your work is done
// e.g.
//
//	conn, _ := dbconnect.GetRedisConn("redis_main")
//	defer conn.Close()
func (c Conns) GetRedisConn(id string) (redis.Conn, error) {
	p, err := c.GetRedisPool(id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nYou've requested a redis.Conn instance. Don't forget to close it after your work is done\n\n")
	return p.Get(), nil
}

// GetRedisPubSubConn is a convenience function; returns a redis.PubSubConn instance identified by input
func (c Conns) GetRedisPubSubConn(id string) (redis.PubSubConn, error) {
	conn, err := c.GetRedisConn(id)
	if err != nil {
		return redis.PubSubConn{}, err
	}

	return redis.PubSubConn{
		Conn: conn,
	}, nil
}

// GetMongoClient returns a pointer to a mongo.Client instance identified by input
func (c Conns) GetMongoClient(ctx context.Context, id string) (*mongo.Client, error) {
	if c.mongoMap == nil {
		return nil, fmt.Errorf("possibly no mongo configurations provided")
	}

	if _, ok := c.mongoMap[id]; !ok {
		return nil, fmt.Errorf("no mongo configuration for ID: %s found", id)
	}

	return c.c.Mongo[c.mongoMap[id]].client(ctx)
}

// GetMongoDB returns a pointer to a mongo.Database instance identified by input
func (c Conns) GetMongoDB(ctx context.Context, id string, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	if c.mongoMap == nil {
		return nil, fmt.Errorf("possibly no mongo configurations provided")
	}

	if _, ok := c.mongoMap[id]; !ok {
		return nil, fmt.Errorf("no mongo configuration for ID: %s found", id)
	}

	return c.c.Mongo[c.mongoMap[id]].db(ctx, opts...)
}
