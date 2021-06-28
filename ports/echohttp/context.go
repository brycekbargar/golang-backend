package echohttp

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// identityNotOk is the common message for when the identity method returns "not ok".
var identityNotOk = echo.NewHTTPError(
	http.StatusUnauthorized,
	"email claim was not found in provide jwt token")

// userContext is the echo.Context + the currently logged in user based on the jwt token.
// If the request is made anonymously email will be nil.
type userContext struct {
	echo.Context
}

// identity gets the current identity of the user (as their email).
// It also returns the jwt where the email came from.
// This method only returns ok when the user could be fully identified.
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
