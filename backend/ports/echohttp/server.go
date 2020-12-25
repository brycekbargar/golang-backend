package echohttp

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Start starts the given server after performing Echo specific setup.
func Start(port int) error {
	s := echo.New()

	newUserHandler(nil, nil).routes(s)

	return s.Start(":" + strconv.Itoa(port))
}
