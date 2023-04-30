# mini-wallet

<!-- DO NOT EDIT. This file generated due to '@envconfig' annotation -->

## Configuration List
| Field Name | Default | Required | 
|---|---|:---:|
| APP_NAME | Mini Wallet | Yes |
| APP_KEY | mini-wallet | Yes |
| APP_ENV | local | Yes |
| APP_ADDRESS | :8089 | Yes |
| CACHE_DEFAULT_MAX_AGE | 30s |  |
| CACHE_PREFIX_KEY | cache_ |  |
| CACHE_HOST | localhost | Yes |
| CACHE_PORT | 6379 | Yes |
| CACHE_PASS |  |  |
| PG_DBNAME | dbname | Yes |
| PG_DBUSER | dbuser | Yes |
| PG_DBPASS | dbpass | Yes |
| PG_HOST | localhost | Yes |
| PG_PORT | 9999 | Yes |
| PG_MAX_OPEN_CONNS | 30 | Yes |
| PG_MAX_IDLE_CONNS | 6 | Yes |
| PG_CONN_MAX_LIFETIME | 30m | Yes |
| APP_ADDRESS | :8089 | Yes |
| APP_READ_TIMEOUT | 5s |  |
| APP_WRITE_TIMEOUT | 10s |  |
| APP_DEBUG | true |  |
| BROKER_KAFKA_ADDRES | :9092 |  |

## DotEnv example
```
APP_NAME=Mini Wallet
APP_KEY=mini-wallet
APP_ENV=local
APP_ADDRESS=:8089
CACHE_DEFAULT_MAX_AGE=30s
CACHE_PREFIX_KEY=cache_
CACHE_HOST=localhost
CACHE_PORT=6379
CACHE_PASS=
PG_DBNAME=dbname
PG_DBUSER=dbuser
PG_DBPASS=dbpass
PG_HOST=localhost
PG_PORT=9999
PG_MAX_OPEN_CONNS=30
PG_MAX_IDLE_CONNS=6
PG_CONN_MAX_LIFETIME=30m
APP_ADDRESS=:8089
APP_READ_TIMEOUT=5s
APP_WRITE_TIMEOUT=10s
APP_DEBUG=true
BROKER_KAFKA_ADDRES=:9092
```

