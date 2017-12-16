package global

import (
	"database/sql"
)

var (
	ClickDB *sql.DB
)

type PointCount struct {
	Count int `json:"count"`
}

