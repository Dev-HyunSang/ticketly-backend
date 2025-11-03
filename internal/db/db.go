package db

import (
	"database/sql"
	"fmt"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/dev-hyunsang/ticketly-backend/config"
	"github.com/dev-hyunsang/ticketly-backend/lib/ent"
)

func ConnectMySQL() (*ent.Client, error) {
	host := config.Getenv("DB_HOST")
	port := config.Getenv("DB_PORT")
	user := config.Getenv("DB_USER")
	password := config.Getenv("DB_PASSWORD")
	dbname := config.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	drv := entsql.OpenDB("mysql", db)

	return ent.NewClient(ent.Driver(drv)), nil
}
