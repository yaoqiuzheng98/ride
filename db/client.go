package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var _client *sql.DB

func init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		"qiu", "123456", "127.0.0.1", 43306, Ride.String())
	_client, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = _client.Ping(); err != nil {
		panic(err)
	}
}

func GetClient() *sql.DB {
	return _client
}
