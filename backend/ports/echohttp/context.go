package echohttp

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// userContext is the echo.Context + the currently logged in user based on the jwt token.
// If the request is made anonymously email will be nil
type userContext struct {
	echo.Context
	email *string
	token *jwt.Token
}

// userContextCreate creates a new custom context with the current user's email address.
func userContextCreate(c echo.Context) (*userContext, error) {
	jt := c.Get("user").(*jwt.Token)
	if jt == nil {
		return &userContext{c, nil, nil}, nil
	}

	email := jt.Claims.(jwt.MapClaims)["email"].(string)
	if len(email) == 0 {
		return nil, errors.New("jwt claims didn't contain the current user's email address")
	}

	return &userContext{
		c,
		&email,
		jt,
	}, nil
}
