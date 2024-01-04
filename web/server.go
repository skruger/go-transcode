package web

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func NewServer(config *ServerConfig) (*echo.Echo, error) {
	loginMiddleware := getMiddleware()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplateRenderer("web/templates")
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(302, "/admin")

	})
	e.GET("/login", getLogin, loginMiddleware...)
	e.POST("/login", postLogin, loginMiddleware...)
	admin := e.Group("/admin", loginMiddleware...)
	admin.Use(loginRequiredMiddleware)
	admin.GET("", config.getAdmin)
	admin.GET("/hx/assets", config.getAssetList)

	admin.GET("/profiles", config.getProfiles)
	return e, nil
}

func getMiddleware() []echo.MiddlewareFunc {
	sessionPath := os.TempDir()
	sessionStore := sessions.NewFilesystemStore(sessionPath)
	csrfConfig := middleware.DefaultCSRFConfig
	csrfConfig.TokenLookup = "cookie:_csrf"
	return []echo.MiddlewareFunc{
		session.MiddlewareWithConfig(session.Config{Store: sessionStore}),
		middleware.CSRFWithConfig(csrfConfig),
	}
}
