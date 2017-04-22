package storage

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type (
	DBHandler struct {
		config DBConfig
	}
	DBConfig struct {
		User string
		Password string
		Host string
		Port string
		Name string
	}
)

func (db *DBHandler) Initialize (config DBConfig) {
	db.config = config
}

func (db *DBHandler) GetInstance () (*sql.DB, error) {
	config := db.config
	conString := config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Name
	return sql.Open("mysql",conString)
}
