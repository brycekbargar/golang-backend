package echohttp

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	key         []byte
	method      jwt.SigningMethod
}

func newUserHandler(
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	key []byte,
	method jwt.SigningMethod) *userHandler {
	return &userHandler{
		authed,
		maybeAuthed,
		key,
		method,
	}
}

func (r *userHandler) routes(s *echo.Echo) {
	s.POST("/users", r.create)
	s.POST("/users/login", r.login)
	s.GET("/user", r.user, r.authed)
	s.PUT("/user", r.update, r.authed)

	s.GET("/profile/:name", r.profile, r.maybeAuthed)
	s.GET("/profile/:name/follow", r.follow, r.authed)
	s.DELETE("/profile/:name/follow", r.unfollow, r.authed)
}

type user struct {
	User userUser `json:"user"`
}
type userUser struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func (r *userHandler) create(c echo.Context) (err error) {
	return nil
}

type login struct {
	User loginUser `json:"user"`
}
type loginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *userHandler) login(c echo.Context) (err error) {
	l := new(login)
	if err = c.Bind(l); err != nil {
		return echo.ErrBadRequest
	}

	if l.User.Email != "jon@jonsnow.com" || l.User.Password != "shhh!" {
		return echo.ErrUnauthorized
	}

	token := jwt.New(r.method)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "JonSnow"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(r.key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    l.User.Email,
			Token:    t,
			Username: claims["name"].(string),
		},
	})
}

func (r *userHandler) user(c echo.Context) (err error) {
	ju := c.Get("user").(*jwt.Token)
	claims := ju.Claims.(jwt.MapClaims)
	return c.JSON(http.StatusOK, user{
		userUser{
			Token:    ju.Raw,
			Username: claims["name"].(string),
		},
	})
}
func (r *userHandler) update(c echo.Context) (err error) {
	return nil
}
func (r *userHandler) profile(c echo.Context) (err error) {
	return nil
}

func (r *userHandler) follow(c echo.Context) (err error) {
	return nil
}

func (r *userHandler) unfollow(c echo.Context) (err error) {
	return nil
}
