package service

import (
//	"bytes"
	"fmt"
	"time"
	"net/http"
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-redis/redis"
	"util"
	"storage"
)

type (
	AuthHandler struct {
		util *util.UtilHandler
		redis *storage.RedisHandler
		db *storage.DBHandler
	}
	errorResponse struct {
		Code string `json: code`
		Message string `json: message`
		Debug string `json: debug`
	}
    createtokenResponse struct {
    	AuthKey string `json:"auth_key"`
    }
    validatetokenRequest struct {
    	AuthKey string `json:"auth_key"`
    }
    validatetokenResponse struct {
    	enable string `json:"enable"`
    }
    DBConfig storage.DBConfig
    RedisConfig storage.RedisConfig
)

/*
func (ah *AuthHandler)getBody(c echo.Context) string {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(c.Request().Body)
	return buffer.String()
}
*/

func (ah *AuthHandler)Initialize(configDB DBConfig, configRedis RedisConfig) {
	fmt.Println("Initialize Service")
	ah.util = new(util.UtilHandler)
	ah.redis = new(storage.RedisHandler)
	var redisConfig storage.RedisConfig
	util.StructCast(&configRedis, &redisConfig)
	ah.redis.Initialize(redisConfig)

	ah.db = new(storage.DBHandler)
	var dbConfig storage.DBConfig
	util.StructCast(&configDB, &dbConfig)
	ah.db.Initialize(dbConfig)
}

func (ah *AuthHandler)CreateToken (c echo.Context) error {
	var errres errorResponse
	key := c.Param("key")

	client := ah.redis.GetInstance()
	defer client.Close()
	userid, err := client.Get(key).Result()
	if err == redis.Nil {
		errres.Code = "E101"
		errres.Message = "Redis Error"
		errres.Debug = err.Error()
		return c.JSON(http.StatusInternalServerError, &errres)
	} else if err != nil {
		errres.Code = "E102"
		errres.Message = "Redis Unknown Error"
		errres.Debug = err.Error()
		return c.JSON(http.StatusInternalServerError, &errres)
	}
	client.Del(key)

	// Auth Keyを生成してRedisに保存＆クライアントに返却する。
	auth_key := ah.util.RandString128()

	var expire time.Duration = 14 * 24 * time.Hour
	err = client.Set(auth_key, userid, expire).Err()
	if err != nil {
		panic(err)
	}
	response := new(createtokenResponse)
	response.AuthKey = auth_key

	return c.JSON(http.StatusOK, response)
}
