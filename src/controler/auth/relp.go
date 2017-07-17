package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"github.com/labstack/echo"
	"common/util"
	"common/storage"
	"encoding/json"
	"time"
)


type (
	RelP struct {
	}

	AuthCert struct {
		State string `json:state`
		Nonce string	`json:nonce`
	}
)


func (relp *RelP)StartAuth (c echo.Context) error {
	uh := new (util.UtilHandler)
	client_id := "authapp"				// TODO:設定に持てるようにする
	scope := "openid"						// TODO:設定に持てるようにする
	auth_cert := new(AuthCert)
	auth_cert.State = uh.RandString128()
	auth_cert.Nonce = uh.RandString128()

	// コードが返却された場合の検証用に、stateとnonceを保持しておく
	rh := new (storage.RedisHandler)
	rh.Initialize ("localhost", "6379", "", 0)
	redis := rh.GetInstance()
	var expire time.Duration = 5 * time.Minute
	data, _ := json.Marshal(auth_cert)
	fmt.Println(data)
	err := redis.Set(client_id, data, expire).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: err.Error()})
	}

	location_url, errUrl := url.Parse("http://localhost:8080/v1/authorize")
	if errUrl != nil {
		return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: "url can not parse"})
	}
	query := location_url.Query()
	query.Set("response_type", "code")
	query.Set("scope", scope)
	query.Set("client_id", client_id)
	query.Set("state", auth_cert.State)
	query.Set("nonce", auth_cert.Nonce)
	query.Set("redirect_uri", "http://localhost:8080/v1/callback")
	location_url.RawQuery = query.Encode()
	fmt.Println(location_url)

	//return c.String(http.StatusOK, location_url.String())
	return c.Redirect(http.StatusTemporaryRedirect, location_url.String())
}


func (relp *RelP)Callback (c echo.Context) error {
//	uh := new (util.UtilHandler)
	client_id := "authapp"				// TODO:設定に持てるようにする
	auth_cert := new(AuthCert)

	rh := new (storage.RedisHandler)
	rh.Initialize ("localhost", "6379", "", 0)
	redis := rh.GetInstance()
	data, err := redis.Get(client_id).Bytes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: err.Error()})
	}
	json.Unmarshal(data, auth_cert)
	fmt.Println(auth_cert)

	if c.Request().FormValue("state") != auth_cert.State {
		return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: "state invalid"})
	}

	return c.String(http.StatusOK, "OK")
}