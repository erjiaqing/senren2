package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var (
	DB    *sql.DB
	Cache *cache.Cache
)

func init() {
	var err error
	sqlurn := os.Getenv("SQL_CONN")
	if sqlurn == "" {
		sqlurn = "root:123456@tcp(127.0.0.1:3306)/senren2?charset=utf8mb4,utf8&parseTime=true"
	}
	DB, err = sql.Open("mysql", sqlurn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Infof("Connected to database")
	Cache = cache.New(30*time.Second, 10*time.Minute)
}
