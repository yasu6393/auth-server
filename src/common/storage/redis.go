package storage

import (
	"github.com/go-redis/redis"
)

type (
	RedisHandler struct {
		config RedisConfig
	}
	RedisConfig struct {
		Addr string
		Port string
		Password string
		DB int
	}
)

func (r *RedisHandler) Initialize(addr string, port string, password string, db int) {
	r.config.Addr = addr
	r.config.Port = port
	r.config.Password = password
	r.config.DB = db
}

func (r *RedisHandler) GetInstance() *redis.Client {
	config := r.config
	Client := redis.NewClient(&redis.Options{
			Addr: config.Addr + ":" + config.Port,
			Password: config.Password,
			DB: config.DB,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		return nil
	}
	return Client
}