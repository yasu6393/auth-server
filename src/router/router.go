package router

import (
	"io"
	"github.com/labstack/echo"
	"service"
	"util"
	"html/template"
)

type (
	Template struct {
	    templates *template.Template
	}
	ServerConfig struct {
		Host string `json:host`
		Listen string `json:listen`
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func Initialize(dbConfig util.DBConfig, redisConfig util.RedisConfig) *echo.Echo {
	tpl := &Template {
		templates: template.Must(template.ParseGlob("public/html/*.html")),
	}
    e := echo.New()
    e.Renderer = tpl
    
    serv := new(service.AuthHandler)
    serv.Initialize(dbConfig, redisConfig)
    // Routes
    v1 := e.Group("/v1")
    {
        v1.GET("/login", serv.LoginView)
        v1.POST("/dologin", serv.Login)
        v1.GET("/loginresult", serv.LoginResultView)
        v1.GET("/createtoken/:key", serv.CreateToken)
    }
    return e
}
