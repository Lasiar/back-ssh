package model

import (
	"database/sql"
	"fmt"
	"github.com/kshvakov/clickhouse"
	"time"
)

func PingClick(db *sql.DB) error {
	for i := 0; i < 3; i++ {
		err := db.Ping()
		if err != nil {
			time.Sleep(1)
			if i == 2 {
				continue
			}
			if exception, ok := err.(*clickhouse.Exception); ok {
				return fmt.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			} else {
				return err
			}
		}
	}
	return nil
}
