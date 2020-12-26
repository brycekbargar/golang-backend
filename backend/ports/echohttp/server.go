package echohttp

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/brycekbargar/realworld-backend/ports"
)

// Start starts the given server after performing Echo specific setup.
func Start(
	jc ports.JWTConfig,
	port int,
	users userdomain.Repository) error {
	s := echo.New()
	s.Use(middleware.Logger())

	fullAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    jc.Key,
		SigningMethod: jc.Method.Name,
		AuthScheme:    "Token",
	})
	maybeAuth := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    jc.Key,
		SigningMethod: jc.Method.Name,
		AuthScheme:    "Token",
		Skipper: func(c echo.Context) bool {
			// Partially auth'd endpoints have different behavior when the user is logged in
			// We want to make sure that only truly anon requests skip auth in these scenarios
			auth := c.Request().Header.Get("Authorization")
			return len(auth) == 0 || len(strings.TrimPrefix(auth, "Token ")) == 0
		},
	})

	newUserHandler(users, fullAuth, maybeAuth, jc).routes(s)

	return s.Start(":" + strconv.Itoa(port))
}
