package util

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type (
	DBHandler struct {
		config DBConfig
	}
	DBConfig struct {
		User string `json:user`
		Password string `json:password`
		Host string `json:"host"`
		Port string `json:"port"`
		Name string `json:"name"`
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
