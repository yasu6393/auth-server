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

func (ah *AuthHandler)Login(c echo.Context) error {
	var response errorResponse
	userid := c.Request().FormValue("userid")
	password := c.Request().FormValue("password")
	con, err := ah.db.GetInstance()
	if err != nil {
		response.Code = "E001"
		response.Message = "DB Connection Error"
		response.Debug = err.Error()
		return c.JSON(http.StatusInternalServerError, &response)
	}
	defer con.Close() // 関数がリターンする直前に呼び出される

	sql := "SELECT * FROM user WHERE user=? and password=?"
	rows, err := con.Query(sql, userid, password)
	if err != nil || rows.Err() != nil {
		response.Code = "E001"
		response.Message = "SQL Error"
		response.Debug = err.Error()
		return c.JSON(http.StatusInternalServerError, &response)
	}
	defer rows.Close()

	url := ""
	if rows.Next() {
		fmt.Println("login success")
		key := ah.util.RandString128()
		client := ah.redis.GetInstance()
		defer client.Close()
		var expire time.Duration = 2 * time.Minute;
		err := client.Set(key, userid, expire).Err()
		if err != nil {
		fmt.Println("Redis Set Error")
			panic(err)
		}
		url = "http://localhost:8080/v1/loginresult?key=" + key
	} else {
		fmt.Println("login failed")
		url = "http://localhost:8080/v1/loginresult"
	}
	return c.Redirect(http.StatusSeeOther, url)
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

func (ah *AuthHandler)LoginView (c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
	
}

func (ah *AuthHandler)LoginResultView (c echo.Context) error {
	return c.Render(http.StatusOK, "loginresult", nil)
}