package model

import (
	"database/sql"
	"github.com/go-redis/redis"
	"log"
)

func NewClick(config string) *sql.DB {
	db, err := sql.Open("clickhouse", config)
	if err != nil {
		log.Fatal(err)
	}
	PingClick(db)
	return db
}

func NewRedis() *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := db.Ping().Result()
	if err != nil {
		log.Println(err)
	}
	return db
}
