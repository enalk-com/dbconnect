package dbconnect

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig defines all the parameters to be used for establishing
// a mongodb connection
type MongoConfig struct {
	ID               string `json:"id,omitempty" toml:"id,omitempty"`
	DB               string `json:"db,omitempty" toml:"db,omitempty"`
	User             string `json:"user,omitempty" toml:"user,omitempty"`
	Pwd              string `json:"pwd,omitempty" toml:"pwd,omitempty"`
	AuthSource       string `json:"authSource,omitempty" toml:"authSource,omitempty"`
	Host             string `json:"host,omitempty" toml:"host,omitempty"`
	Port             int    `json:"port,omitempty" toml:"port,omitempty"`
	ConnectionString string `json:"connectionString,omitempty" toml:"connectionString,omitempty"`
	_client          *mongo.Client
	once             sync.Once
}

func (mc *MongoConfig) defaults() {
	if mc.Host == "" {
		mc.Host = "localhost"
	}

	if mc.Port == 0 {
		mc.Port = 27017
	}
}

func (mc *MongoConfig) prepareURI() {
	if mc.ConnectionString != "" {
		return
	}
	mc.defaults()
	uri := fmt.Sprintf("%s:%d/", mc.Host, mc.Port)
	if mc.User != "" || mc.Pwd != "" {
		uri = fmt.Sprintf("%s:%s@%s", mc.User, mc.Pwd, uri)
	}
	if mc.DB != "" {
		uri = fmt.Sprintf("%s%s", uri, mc.DB)
	}

	if mc.AuthSource != "" {
		uri = fmt.Sprintf("%s?authSource=%s", uri, mc.AuthSource)
	}

	mc.ConnectionString = fmt.Sprintf("mongodb://%s", uri)
	log.Printf("ConnectionString: %s\n", mc.ConnectionString)
}

func (mc *MongoConfig) connect(ctx context.Context) {
	mc.once.Do(func() {
		if mc._client != nil {
			return
		}
		mc.prepareURI()
		opts := options.Client().ApplyURI(mc.ConnectionString)
		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			log.Printf("error creating client: %s\n", err.Error())
			log.Fatal(err.Error())
		}

		if err := client.Ping(ctx, nil); err != nil {
			log.Printf("error pinging with client: %s", err.Error())
			log.Fatal(err.Error())
		}

		mc._client = client
	})
}

// Client returns a *mongo.Client connection
func (mc *MongoConfig) client(ctx context.Context) (*mongo.Client, error) {
	mc.connect(ctx)
	if mc._client == nil {
		return nil, fmt.Errorf("invalid mongo configuration. client was not created")
	}
	return mc._client, nil
}

// Database returns a *mongo.Database object if db key was specified during configuration
func (mc *MongoConfig) db(ctx context.Context, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	if mc.DB == "" {
		return nil, fmt.Errorf("empty database name")
	}
	c, err := mc.client(ctx)
	if err != nil {
		return nil, err
	}
	return c.Database(mc.DB, opts...), nil
}
