package echohttp

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

type userHandler struct {
	users       userdomain.Repository
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	key         []byte
	method      jwt.SigningMethod
}

func newUserHandler(
	users userdomain.Repository,
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	key []byte,
	method jwt.SigningMethod) *userHandler {
	return &userHandler{
		users,
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

type register struct {
	User registerUser `json:"user"`
}
type registerUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func makeJwt(r *userHandler, e string) (string, error) {
	token := jwt.New(r.method)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = e
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(r.key)
	if err != nil {
		return "", err
	}

	return t, nil
}

func (r *userHandler) create(c echo.Context) error {
	u := new(register)
	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	created, err := userdomain.NewUserWithPassword(
		u.User.Email,
		u.User.Username,
		u.User.Password,
	)
	if err != nil {
		return err
	}

	if err := r.users.Create(created); err != nil {
		return err
	}

	token, err := makeJwt(r, u.User.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    u.User.Email,
			Token:    token,
			Username: u.User.Password,
		},
	})
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

	token, err := makeJwt(r, l.User.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email: l.User.Email,
			Token: token,
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
