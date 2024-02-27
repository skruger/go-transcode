package web

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func contextMap(c echo.Context) echo.Map {
	return echo.Map{
		"csrf": c.Get("csrf"),
		"context": c,
	}
}

func getLogin(c echo.Context) error  {
	//tmpl, err := template.New("login").ParseFiles("web/login.tmpl")
	//if err != nil {
	//	log.Printf("couldn't load template! %s", err)
	//}
	//c.Response().Header().Set("Content-Type", "text/html")
	//err = tmpl.ExecuteTemplate(c.Response().Writer, "login.tmpl", contextMap(c))
	//if err != nil {
	//	log.Printf("couldn't render template! %s", err)
	//}
	//return c.NoContent(200)
	//

	return c.Render(200, "base.html,login.html", contextMap(c))
	//return c.String(200, "OK")
}

func postLogin(c echo.Context) error {

	return c.Redirect(302, "/admin/")
}

func loginRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path: "/",
			MaxAge: 86400,

		}
		//user := sess.Values["user_id"]
		//log.Printf("User: %s", user)
		sess.Save(c.Request(), c.Response())
		// Uncomment this to enforce logins when that is ready
		//if user == nil {
		//	redirect := fmt.Sprintf("/login?next=%s", url.QueryEscape(c.Path()))
		//	return c.Redirect(302, redirect)
		//}
		return next(c)
	}
}
