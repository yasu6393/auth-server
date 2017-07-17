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

func (db *DBHandler) Initialize (user string, password string, host string, port string, name string) {
	db.config.User = user
	db.config.Password = password
	db.config.Host = host
	db.config.Port = port
	db.config.Name = name
}

func (db *DBHandler) GetInstance () (*sql.DB, error) {
	config := db.config
	conString := config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Name
	return sql.Open("mysql",conString)
}
