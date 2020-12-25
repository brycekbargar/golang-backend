package echohttp

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Start starts the given server after performing Echo specific setup.
func Start(port int, secret string) error {
	s := echo.New()

	key := []byte(secret)

	fullAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: key,
	})

	maybeAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: key,
		Skipper: func(c echo.Context) bool {
			// Partially auth'd endpoints have different behavior when the user is logged in
			// We want to make sure that only truly anon requests skip auth in these scenarios
			return len(c.Request().Header.Get("Authorization")) == 0
		},
	})

	newUserHandler(fullAuth, maybeAuth, key).routes(s)

	return s.Start(":" + strconv.Itoa(port))
}
