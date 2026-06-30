package db

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _client *gorm.DB

func init() {
	var err error
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnvInt("DB_PORT", 43306)
	user := getEnv("DB_USER", "qiu")
	pass := getEnv("DB_PASS", "123456")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, pass, host, port, Ride.String())
	_client, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetClient() *gorm.DB {
	return _client
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
