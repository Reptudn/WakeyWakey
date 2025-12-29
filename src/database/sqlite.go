package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

// Init opens the SQLite database once with sensible defaults.
func Init(dbPath string) (*sql.DB, error) {
	var err error

	once.Do(func() {
		dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)", dbPath)

		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			return
		}

		err = db.Ping()
	})

	return db, err
}

func Get() *sql.DB {
	return db
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}