package auth

import (
	"fmt"
	"regexp"
	"net/http"
	"net/url"
	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"common/storage"
	"common/util"
//	"time"
)

type (

	IdP struct {
	}

	Error struct {
		Code string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	}
)


func (idp *IdP)Authorize (c echo.Context) error {
	fmt.Println("Authorize called")

	db := new (storage.DBHandler)
	db.Initialize("root", "", "localhost", "3306", "userdb")
	con, err := db.GetInstance()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "DB Error", Details: err.Error()})
	}
	defer con.Close() // 関数がリターンする直前に呼び出される

	sess := session.Default(c)

	// response_typeはコード固定とする
	response_type := c.Request().FormValue("response_type")
	if response_type != "code" {

	}

	// client_idのチェック
	client_id := c.Request().FormValue("client_id")
	sql := "SELECT client_secret, redirect_uri FROM client WHERE client_id=?"
	rows, err := con.Query(sql, client_id)
	if err != nil || rows.Err() != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "SQL Error", Details: err.Error()})
	}
	defer rows.Close()
	if rows.Next() {
		var client_secret, redirect_uri string
		err = rows.Scan(&client_secret, &redirect_uri)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "SQL Error", Details: err.Error()})
		}

		if c.Request().FormValue("redirect_uri") != redirect_uri {
			return c.JSON(http.StatusBadRequest, Error{Code: "E002", Message: "Bad Request", Details: "redirect_uri is invalid"})
		}
		sess.Set("redirect_uri", redirect_uri)
		fmt.Println(redirect_uri)
	} else {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "No Client", Details: "client_id is invalid"})
	}

	// scopeには、"open_id"が必須
	scope := c.Request().FormValue("scope")
	reg_scope, _ := regexp.Compile(".*openid.*")
	fmt.Println(reg_scope.String())
	if reg_scope == nil {
		return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: "regex error"})
	}

	if "" == reg_scope.FindString(scope) {
		return c.JSON(http.StatusBadRequest, Error{Code: "E002", Message: "Parameter Error", Details: "scope is invalid"})
	}

	// stateはCSRFの防止に使用する
	state := c.Request().FormValue("state")
	if "" == state {
		return c.JSON(http.StatusBadRequest, Error{Code: "E002", Message: "Parameter Error", Details: "state is invalid"})
	}
	sess.Set("state", state)

	// nonceはリプレイアタックの防止に使用する
	nonce := c.Request().FormValue("nonce")
	if "" == nonce {
		return c.JSON(http.StatusBadRequest, Error{Code: "E002", Message: "Parameter Error", Details: "nonce is invalid"})
	}
	sess.Set("nonce", nonce)

	sess.Save()
	// それぞれのチェックがOkだったら、ログインビューを返却する
	return c.Render(http.StatusOK, "login", nil)
}

func (idp *IdP)Login(c echo.Context) error {
	sess := session.Default(c)
	uh := new (util.UtilHandler)

	userid := c.Request().FormValue("userid")
	password := c.Request().FormValue("password")
	fmt.Println(userid, password)

	db := new (storage.DBHandler)
	db.Initialize("root", "", "localhost", "3306", "userdb")
	con, err := db.GetInstance()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "DB Error", Details: err.Error()})
	}
	defer con.Close() // 関数がリターンする直前に呼び出される

	sql := "SELECT * FROM user WHERE user=? and password=?"
	rows, err := con.Query(sql, userid, password)
	if err != nil || rows.Err() != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "SQL Error", Details: err.Error()})
	}
	defer rows.Close()

	dst := "http://localhost:8080/v1/loginresult"
	if rows.Next() {
		fmt.Println("login success")
		redirect_uri := sess.Get("redirect_uri").(string)
		state := sess.Get("state").(string)

		location_url, errUrl := url.Parse(redirect_uri)
		if errUrl != nil {
			return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: "url can not parse"})
		}
		query := location_url.Query()
		query.Set("state", state)
		query.Set("code", uh.RandString128())
		location_url.RawQuery = query.Encode()

		dst = location_url.String()
	} else {
		fmt.Println("login failed")
	}
	return c.Redirect(http.StatusTemporaryRedirect, dst)
}