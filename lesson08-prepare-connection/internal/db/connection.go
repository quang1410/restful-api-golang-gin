package db

import "lesson08-prepare-connection/internal/config"

var DB string

func InitDB() error {
	connStr := config.NewConfig().DNS()

	DB = connStr

	return nil
}
