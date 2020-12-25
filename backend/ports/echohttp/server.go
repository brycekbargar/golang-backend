package echohttp

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start starts the given server after performing Echo specific setup.
func Start(port int, secret string) error {
	s := echo.New()
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())

	key := []byte(secret)
	fullAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: key,
		AuthScheme: "Token",
	})
	maybeAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: key,
		AuthScheme: "Token",
		Skipper: func(c echo.Context) bool {
			// Partially auth'd endpoints have different behavior when the user is logged in
			// We want to make sure that only truly anon requests skip auth in these scenarios
			auth := c.Request().Header.Get("Authorization")
			return len(auth) == 0 || len(strings.TrimPrefix(auth, "Token ")) == 0
		},
	})

	newUserHandler(fullAuth, maybeAuth, key).routes(s)

	return s.Start(":" + strconv.Itoa(port))
}
