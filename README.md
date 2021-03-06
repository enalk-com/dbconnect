### dbconnect -> connect redis, postgreql, mongo
***

### Quickstart
***


#### config

```toml
# config-filename.toml
[[pq]]
id="database ID" # can be any unique string identifier.
# you will use this later to get a sql.DB instance
user="db user"
pwd="db pwd"
host="db host"
port="db port"
db="db name to connect to"
# check struct PQConfig for more options
# more [[pq]] blocks can be added

[[mongo]]
id="mongo instance id" #could be any unique string.
# Will be used for fetching this instance later.\
user="db user", # optional
pwd="db pwd", # optional
host="db host", # optional
port=27017, # optional
db="db name", # optional
authSource="auth db name", # optional
connectString="mongodb://localhost:27017" # optional, however
# takes precedence over all other options when present
# more [[mongo]] blocks can be added for multiple dbs

[[redis]]
id="redis instance ID" # could be any unique string.
# Will be used for fetching this instance later.
host="redis instance host"
port=6379 # default: 6379
pwd="redis instance pwd",
db="0",
# check struct RedisConfig for more options
# more [[redis]] blocks can be added for multiple dbs
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