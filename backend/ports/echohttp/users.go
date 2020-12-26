package echohttp

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/brycekbargar/realworld-backend/ports"
)

type userHandler struct {
	users       userdomain.Repository
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	jc          ports.JWTConfig
}

func newUserHandler(
	users userdomain.Repository,
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	jc ports.JWTConfig) *userHandler {
	return &userHandler{
		users,
		authed,
		maybeAuthed,
		jc,
	}
}

func (r *userHandler) routes(g *echo.Group) {
	g.POST("/users", r.create)
	g.POST("/users/login", r.login)
	g.GET("/user", r.user, r.authed)
	g.PUT("/user", r.update, r.authed)

	g.GET("/profile/:name", r.profile, r.maybeAuthed)
	g.GET("/profile/:name/follow", r.follow, r.authed)
	g.DELETE("/profile/:name/follow", r.unfollow, r.authed)
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
	token := jwt.New(r.jc.Method)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = e
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(r.jc.Key)
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

	authed, err := r.users.GetLoginUserByEmail(l.User.Email)
	if err != nil {
		return err
	}

	if pw, err := authed.HasPassword(l.User.Password); !pw || err != nil {
		return echo.ErrUnauthorized
	}

	token, err := makeJwt(r, authed.Email())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    authed.Email(),
			Username: authed.Username(),
			Token:    token,
			Bio:      authed.Bio(),
			Image:    authed.Image(),
		},
	})
}

func (r *userHandler) user(c echo.Context) (err error) {
	ju := c.Get("user").(*jwt.Token)
	claims := ju.Claims.(jwt.MapClaims)

	found, err := r.users.GetUserByEmail(claims["email"].(string))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    found.Email(),
			Username: found.Username(),
			Token:    ju.Raw,
			Bio:      found.Bio(),
			Image:    found.Image(),
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
