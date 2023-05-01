# Mini Wallet

## Dependencies
- [typicalw](https://github.com/typical-go/typical-go)
- [uber-dig](https://github.com/uber-go/dig) (for Dependency Injection)
- PostgreSQL
- [golang-migrate](https://github.com/golang-migrate/migrate)
- Redis
- Kafka
- JWT

## Installation

1. If you don't have golang migrate, you can install first using this command
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
2. You can start with this command
```bash
go mod tidy
```
3. To up migration you can use this command
```bash
make migration-up CONNECTION="postgresql://username:password@host:port/db-name"
```
4. Don't forget to change `.env.example` file to be `.env` and **CHANGE ALL VALUE**
5. After change `.env` value you can generate new DI and Env Config using this command
```bash
./typicalw g
```
6. Don't forget to add this topic to your kafka
```
wallet.created-transaction
```

## Running
1. To start the program, you can easily run with this command
```bash
./typicalw r
```
2. If you want to access swagger docs, you can open this URL
```
http://localhost:8089/swagger/index.html
```

## Overview Process
### Get Existing Wallet Amount
This process using cache (Redis) to make less connection to database. Using simple update value in cache. You can check on this flowchart

![redis.png](./docs/redis.png)

### Update Cache after Deposit & Withdrawal Wallet
This process implement queue using kafka. After user doing deposit or withdrawal, system will send kafka message to notifiy there is new transaction from wallet. This process also tackling requirement to not directly update value after user using their balance.

![kafka.png](./docs/kafka.png)

After current system receiver message from kafka, system will refresh cache to used by get wallet API

![kafka-redis.png](./docs/kafka-redis.png)

## Other
- This code using `JSend` as Rest response format
- This code implementing `domain design` to separate scope
- This code contain `in-code comment`, other engineer will be more easy to get understanding of code
- This code use `postgres DB locking` to prevent race condition