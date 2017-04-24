/**
 * OAuth2.0のリソースオーナー
 *
 * OAuth2.0のリソースオーナー用のメソッドを定義する
 * 定義されるメソッド
 * - LoginView:認証画面表示
 * - Login:認証
 * Golang versions 1.8.1 After
 *
 * @category   CategoryName
 * @package    oauth2/resource/owner
 * @author     Yasushi Furusawa
 * @copyright  Copyright (c) 2017 OchaFramework Project.
 * @license    http://www.opensource.org/licenses/mit-license.html  MIT License
 * @version    0.1
**/
package owner

import (
//	"bytes"
	"fmt"
	"time"
	"net/http"
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql"
	"util"
	"storage"
)

type (
	ResourceOwner struct {
		util *util.UtilHandler
		redis *storage.RedisHandler
		db *storage.DBHandler
	}
	errorResponse struct {
		Code string `json: code`
		Message string `json: message`
		Debug string `json: debug`
	}
    DBConfig storage.DBConfig
    RedisConfig storage.RedisConfig
)


func (ro *ResourceOwner)Initialize(configDB DBConfig, configRedis RedisConfig) {
	fmt.Println("Initialize Service")
	ro.util = new(util.UtilHandler)
	ro.redis = new(storage.RedisHandler)
	var redisConfig storage.RedisConfig
	util.StructCast(&configRedis, &redisConfig)
	ro.redis.Initialize(redisConfig)

	ro.db = new(storage.DBHandler)
	var dbConfig storage.DBConfig
	util.StructCast(&configDB, &dbConfig)
	ro.db.Initialize(dbConfig)
}

func (ro *ResourceOwner)Token (c echo.Context) error {
	var errResp errorResponse
	grant_type := c.Request().FormValue("grant_type")
	client_id := c.Request().FormValue("client_id")
	client_secret := c.Request().FormValue("client_secret")
	redirect_uri := c.Request().FormValue("redirect_uri")

	if grant_type != "authorization_code" {
		errResp.Code = "E0201"
		errResp.Message = "不正なリクエストです。"
		errResp.Debug = "Unknown Client"
		return c.JSON(http.StatusBadRequest, errResp)
	}

	if redirect_uri == "" {
		errResp.Code = "E0202"
		errResp.Message = "不正なリクエストです。"
		errResp.Debug = "Unknown Redirect URI"
		return c.JSON(http.StatusBadRequest, errResp)
	}


	con, err := ro.db.GetInstance()
	if err != nil {
		errResp.Code = "E001"
		errResp.Message = "DB Connection Error"
		errResp.Debug = err.Error()
		return c.JSON(http.StatusInternalServerError, &errResp)
	}
	defer con.Close()

	var tmp_secret string
	sql := "SELECT client_secret FROM client WHERE client_id=?"
	if err = con.QueryRow(sql, client_id).Scan(&tmp_secret); err != nil {
		errResp.Code = "E001"
		errResp.Message = "SQL Error"
		errResp.Debug = err.Error()
		return c.JSON(http.StatusInternalServerError, &errResp)
	}
	if tmp_secret != client_secret {
		errResp.Code = "E1202"
		errResp.Message = "不正なリクエストです。"
		errResp.Debug = "Unknown Client Secret"
		return c.JSON(http.StatusBadRequest, errResp)
	}

	return c.Redirect(http.StatusSeeOther, redirect_uri)
}

func (ro *ResourceOwner)LoginView (c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)

}

func (ro *ResourceOwner)Login(c echo.Context) error {
	var response errorResponse
	userid := c.Request().FormValue("userid")
	password := c.Request().FormValue("password")
	con, err := ro.db.GetInstance()
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
		key := ro.util.RandString128()
		client := ro.redis.GetInstance()
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


func (ro *ResourceOwner)LoginResultView (c echo.Context) error {
	return c.Render(http.StatusOK, "loginresult", nil)
}
