// Package echohttp is a port into the application.
// (currently it's the only port and has a bunch of business logic :scream:)
// It contains a single http server using the Echo library.
package echohttp

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/brycekbargar/realworld-backend/ports"
)

// Start starts the given server after performing Echo specific setup.
func Start(
	jc ports.JWTConfig,
	port int,
	users userdomain.Repository,
	articles articledomain.Repository,
) error {
	s := echo.New()
	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			uc := &userContext{c}
			return next(uc)
		}
	})
	s.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			// TODO: Figure out how not to leak the route details?
			// TODO: Also figure out how this is _really_ done.
			// These routes take a password in the parameters so we want to leave them out of logs.
			if strings.HasPrefix(strings.ToLower(c.Path()), "/api/users/login") {
				return true
			}
			if strings.HasPrefix(strings.ToLower(c.Path()), "/api/user") &&
				(c.Request().Method == http.MethodPost || c.Request().Method == http.MethodPut) {
				return true
			}

			return false
		},
	}))

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
			// We want to make sure that anon requests skip auth in these scenarios
			auth := c.Request().Header.Get("Authorization")
			return len(strings.TrimPrefix(strings.ToLower(auth), "token ")) == 0
		},
	})

	api := s.Group("/api")
	newUsersHandler(users, fullAuth, maybeAuth, jc).mapRoutes(api)
	newArticlesHandler(users, articles, fullAuth, maybeAuth).mapRoutes(api)

	return s.Start(":" + strconv.Itoa(port))
}
