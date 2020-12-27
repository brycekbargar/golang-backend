package echohttp

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var identityNotOk = echo.NewHTTPError(
	http.StatusUnauthorized,
	"email claim was not found in provide jwt token")

// userContext is the echo.Context + the currently logged in user based on the jwt token.
// If the request is made anonymously email will be nil
type userContext struct {
	echo.Context
}

func (uc *userContext) identity() (string, *jwt.Token, bool) {
	ju := uc.Get("user")
	if ju == nil {
		return "", nil, false
	}

	jt := ju.(*jwt.Token)
	email := jt.Claims.(jwt.MapClaims)["email"].(string)
	if len(email) == 0 {
		return "", jt, false
	}

	return email, jt, true
}

// userContextCreate creates a new custom context with the current user's email address.
func userContextCreate(c echo.Context) (*userContext, error) {
	return &userContext{c}, nil
}
