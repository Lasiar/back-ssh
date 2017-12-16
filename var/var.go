package lib

import (
	"database/sql"
	"github.com/go-redis/redis"
)

var (
	ClickDB *sql.DB
	RedisDB *redis.Client
)

type PointCount struct {
	Count int `json:"count"`
}

type InfoPointJs struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Success   bool   `json:"success"`
}
