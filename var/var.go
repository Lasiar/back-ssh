package lib

import (
	"database/sql"
	"github.com/go-redis/redis"
)

var (
	ClickDB *sql.DB
	RedisDB *redis.Client
	Keys    []int
)

type PointCount struct {
	Count int `json:"count"`
}

type InfoPointJs struct {
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Success   bool   `json:"success"`
}

type PointAllList struct {
	Point []int `json:"point"`
}

type ListenJson struct {
	Point      int             `json:"point"`
	Statistics [][]interface{} `json:"statistics"`
}

type Json struct {
	Point      int             `json:"point"`
	Statistics [][]interface{} `json:"statistics"`
}
