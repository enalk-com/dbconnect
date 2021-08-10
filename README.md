### dbconnect -> connect redis, postgreql, mongo
***

### Quickstart
***


#### config
```json
// ~/.config/yourapp/config.json
{
	"pq": [{
		"id": "dbID" /*could be any unique string.
		Will be used for fetching this db later.*/
		"user": "db user",
		"pwd": "db pwd",
		"host": "db host",
		"port": 5432,
		"db": "db name"
		/* check struct PQConfig for more options */
	}, ...],
	"redis": [{
		"id": "redis instance ID" /*could be any unique string.
		Will be used for fetching this instance later.*/
		"host": "redis instance host"
		"port": 6379 /*default: 6379*/
		"pwd": "redis instance pwd",
		"db": "0",
		/* check struct RedisConfig for more options */
		
	},...],

	"mongo": [{
		"id": "mongo instance id" /*could be any unique string.
		Will be used for fetching this instance later.*/
		"user": "db user", // optional
		"pwd": "db pwd", // optional
		"host": "db host", // optional
		"port": 27017, // optional
		"db": "db name", // optional
		"authSource": "auth db name", // optional
		"connectString": "mongodb://localhost:27017" /*  optional, however
		takes precedence over all other options when present */
	}]
}
```

#### Info

1. every config element needs to have an `id` which is later used
   to reference a connection via getter functions
2. connections are not made on `dbconnect.Init([]bytes("{}"))`; only
   at first call via getter functions
3. each connection is made only once in the lifetime of a server. At this
   moment there is no way to circumvent this. Since most connections made
   are pools - this is not really an issue unless the databases are unreachable
4. TLS connections are not supported at the moment for redis and mongo
