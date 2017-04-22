package router

import (
	"fmt"
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
		Views string
	}
    DBConfig service.DBConfig
    RedisConfig service.RedisConfig
)


func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func Initialize(serverConfig ServerConfig, configDB DBConfig, configRedis RedisConfig) *echo.Echo {
	fmt.Println("Initialize Router")
	tpl := &Template {
		templates: template.Must(template.ParseGlob(serverConfig.Views)),
	}
    e := echo.New()
    e.Renderer = tpl
    
    serv := new(service.AuthHandler)
    
	var redisConfig service.RedisConfig
	var dbConfig service.DBConfig
	util.StructCast(&configRedis, &redisConfig)
	util.StructCast(&configDB, &dbConfig)

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
