package main

import (
	"router"
	"util"
)

type (
	ConfigData struct {
		DB util.DBConfig `json:db`
		Redis util.RedisConfig `json:redis`
	}
)

var config ConfigData

func main() {
	var fh util.FileHandler
	err := fh.LoadJson("config.json", &config)
	if err != nil {
		panic(err)
	}
	
	e := router.Initialize(config.DB, config.Redis)
    e.Start(":8080")
}