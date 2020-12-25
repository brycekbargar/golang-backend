package echohttp

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Server is an Echo server.
type Server struct {
	server *echo.Echo
}

// Start starts the given server after performing Echo specific setup.
func (s *Server) Start(port int) error {
	s.server = echo.New()
	s.server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	s.server.POST("/users", nil)       // create (anon)
	s.server.POST("/users/login", nil) // login (anon)
	s.server.GET("/user", nil)         // get (auth'd only and current user)
	s.server.PUT("/user", nil)         // update (auth'd only and current user)

	s.server.GET("/profile/:name", nil)           // get (auth'd or anonymous w/ same results)
	s.server.GET("/profile/:name/follow", nil)    // get (auth'd only and current user)
	s.server.DELETE("/profile/:name/follow", nil) // get (auth'd only and current user)

	return s.server.Start(":" + strconv.Itoa(port))
}
