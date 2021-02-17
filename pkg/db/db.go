package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	retryCount = 1
)

type DB = sqlx.DB

// NewDB returns connected Client
func NewDB(driverName, dsn string) (*DB, error) {
	var conn *sqlx.DB
	var err error
	for i := 0; i < retryCount; i++ {
		conn, err = sqlx.Connect(driverName, dsn)
		if err != nil {
			log.Print(err)
			time.Sleep(time.Second)
			continue
		}
		return conn, nil
	}
	log.Print(dsn)
	return nil, errors.Wrap(err, "DB connect")
}

// DSN returns data source name for connection
func DSN(user, password, db, host string, port int) string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s", user, password, host, port, db)
}
