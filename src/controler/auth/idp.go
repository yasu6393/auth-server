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
  "crypto/sha256"
  "encoding/json"
//	"time"
)

type (

  IdP struct {
  }

	ChapData struct {
		Userid string
		Salt string
	}

  Error struct {
    Code string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details"`
  }
)

// OpenIdConnect 認証画面表示
// 既知のクライアントからのアクセス以外は受け付けない
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
    return c.JSON(http.StatusBadRequest, Error{Code: "E900", Message: "Not Supported", Details: response_type + " Not Supported"})
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
  defer con.Close()

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

// CHAP形式の認証
// 認証開始
func(idp *IdP)ChapStart(c echo.Context) error {
	fmt.Println("chap start")
	sess := session.Default(c)

	userid := c.QueryParam("userid")
	uh := new (util.UtilHandler)
	salt := uh.RandString64()

	authdata, err := json.Marshal(&ChapData{userid, salt})
	if err != nil {
      return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: err.Error()})
	}
	sess.Set("authdata", 	authdata)

	if err = sess.Save(); err != nil {
      return c.JSON(http.StatusInternalServerError, Error{Code: "E001", Message: "Server Error", Details: err.Error()})
	}

//	chapdata := sess.Get("authdata").(ChapData)
//	fmt.Println(chapdata)

	seed := []byte("furu12345ec56d7e08c3bd7d36ce863d1cd73d3fdff770f086c0db8e1b6a8fc9f0b8dbb")
	tmp := sha256.Sum256(seed)
	fmt.Println(tmp)

	return c.JSON(http.StatusOK, &struct {
		Salt string `json:"salt"`
	}{salt})
}

// パスワード認証
// sha256sum暗号化パスワードを受け付ける
func(idp *IdP)ChapAuth(c echo.Context) error {
	fmt.Println("ChapAuth start")
	sess := session.Default(c)

	requestbody := &struct {
		Pwd string `json:"password"`
	}{}

	if err := c.Bind(requestbody); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E999", Message: "Param Error", Details: err.Error()})
	}

	if sess.Get("authdata") == nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E999", Message: "Unknown Error", Details: "Invalid Access"})
	}
	var authdata ChapData
	if err := json.Unmarshal(sess.Get("authdata").([]byte), &authdata); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E999", Message: "Unknown Error", Details: err.Error()})
	}
	sess.Delete("authdata")

	userid := authdata.Userid
	salt := authdata.Salt

	// 指定したユーザーIDが存在するかチェックする
	db := new (storage.DBHandler)
	db.Initialize("root", "", "localhost", "3306", "userdb")
	con, err := db.GetInstance()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "DB Error", Details: err.Error()})
	}
	defer con.Close()

	var pwd string
	sql := "SELECT password FROM user WHERE user=?"
	if err := con.QueryRow(sql, userid).Scan(&pwd); err != nil {
		return c.JSON(http.StatusBadRequest, Error{Code: "E001", Message: "SQL Error", Details: err.Error()})
	}

	seed := []byte(pwd + salt)
	tmp := sha256.Sum256(seed)
	crypto_pwd := fmt.Sprintf("%x", tmp)

	if requestbody.Pwd != crypto_pwd {
		return c.JSON(http.StatusBadRequest, Error{Code: "E009", Message: "Auth Error", Details: "Auth Error"})
	}

	// トークン生成
	jwt_1 := "token"	// TODO:JWTの生成

	return c.JSON(http.StatusOK, &struct {
		Jwt string `json:"jwt"`
	}{jwt_1})
}