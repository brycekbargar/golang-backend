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

	return s.server.Start(":" + strconv.Itoa(port))
}
