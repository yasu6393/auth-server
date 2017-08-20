package router

import (
	"fmt"
	"io"
	"github.com/labstack/echo"
	"github.com/ipfans/echo-session"
	"html/template"
	"controler/auth"
)

type (
	Template struct {
	    templates *template.Template
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func Initialize() *echo.Echo {
	fmt.Println("Initialize Router")
	tpl := &Template {
		templates: template.Must(template.ParseGlob("public/html/*.html")),
	}
    e := echo.New()
    e.Renderer = tpl

	store, err := session.NewRedisStore(32, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		panic(err)
	}
    e.Use(session.Sessions("GLOBAL_SESSION", store))

    // Routes
    idp := new (auth.IdP)
    relp := new (auth.RelP)
    v1 := e.Group("/v1")
    {
        v1.GET("/startauth", relp.StartAuth)
        v1.GET("/authorize", idp.Authorize)
        v1.POST("/login", idp.Login)
        v1.POST("/callback", relp.Callback)
        v1.GET("/auth/hello", idp.ChapStart)
        v1.POST("/auth/authorize", idp.ChapAuth)
    }
    return e
}
