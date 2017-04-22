package main

import (
	"fmt"
	"router"
	"encoding/json"
	"io/ioutil"
	"os"
	"util"
)

type (
	ConfigData struct {
		Server ServerConfig `json:server`
		DB DBConfig `json:db`
		Redis RedisConfig `json:redis`
	}
	ServerConfig struct {
		Views string `json:views`
		Listen string `json:listen`
	}
	RedisConfig struct {
		Addr string `json:addr`
		Port string `json:port`
		Password string `json:password`
		DB int `json:db`
	}
	DBConfig struct {
		User string `json:user`
		Password string `json:password`
		Host string `json:host`
		Port string `json:port`
		Name string `json:name`
	}
	
	FileHandler struct {
	}
)

func (fh FileHandler) LoadJson(filePath string, v interface{}) error {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, &v)
}

func (fh FileHandler) OutputFile(b []byte, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(b)

	return err
}


func main() {
	fmt.Println("Start Server")
	var config ConfigData
	var fh FileHandler
	err := fh.LoadJson("config.json", &config)
	if err != nil {
		panic(err)
	}
	
	var serverConfig router.ServerConfig
	var redisConfig router.RedisConfig
	var dbConfig router.DBConfig
	util.StructCast(&config.Server, &serverConfig)
	util.StructCast(&config.Redis, &redisConfig)
	util.StructCast(&config.DB, &dbConfig)
	
	e := router.Initialize(serverConfig, dbConfig, redisConfig)
    e.Start(config.Server.Listen)
}