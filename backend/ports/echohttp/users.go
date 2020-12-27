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

	g.GET("/profiles/:username", r.profile, r.maybeAuthed)
	g.POST("/profiles/:username/follow", r.follow, r.authed)
	g.DELETE("/profiles/:username/follow", r.unfollow, r.authed)
}

type user struct {
	User userUser `json:"user"`
}
type userUser struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
	Password *string `json:"password,omitempty"`
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
			Username: u.User.Username,
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

func (r *userHandler) login(c echo.Context) error {
	l := new(login)
	if err := c.Bind(l); err != nil {
		return echo.ErrBadRequest
	}

	authed, err := r.users.GetUserByEmail(l.User.Email)
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
			Bio:      optional(authed.Bio()),
			Image:    optional(authed.Image()),
		},
	})
}

func (r *userHandler) user(c echo.Context) error {
	em, token, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	found, err := r.users.GetUserByEmail(em)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    found.Email(),
			Username: found.Username(),
			Token:    token.Raw,
			Bio:      optional(found.Bio()),
			Image:    optional(found.Image()),
		},
	})
}

func (r *userHandler) update(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	b := new(user)
	if err := c.Bind(b); err != nil {
		return echo.ErrBadRequest
	}

	err := r.users.UpdateUserByEmail(
		em,
		func(u *userdomain.User) (*userdomain.User, error) {
			return userdomain.UpdatedUser(*u,
				b.User.Email,
				b.User.Username,
				b.User.Bio,
				b.User.Image,
				*b.User.Password)
		})
	if err != nil {
		return err
	}

	found, err := r.users.GetUserByEmail(em)
	if err != nil {
		return err
	}

	// Users can change their email so we need to make sure we're giving them a new token.
	token, err := makeJwt(r, found.Email())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    found.Email(),
			Username: found.Username(),
			Token:    token,
			Bio:      optional(found.Bio()),
			Image:    optional(found.Image()),
		},
	})
}

type profile struct {
	Profile profileUser `json:"profile"`
}
type profileUser struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (r *userHandler) profile(c echo.Context) (err error) {
	em, _, _ := c.(*userContext).identity()

	found, err := r.users.GetUserByUsername(c.Param("username"))
	if err != nil {
		return err
	}

	following := false
	if len(em) > 0 {
		cu, err := r.users.GetUserByEmail(em)
		if err != nil {
			return err
		}

		following = cu.IsFollowing(found)
	}

	return c.JSON(http.StatusOK, profile{
		profileUser{
			Username:  found.Username(),
			Bio:       found.Bio(),
			Image:     found.Image(),
			Following: following,
		},
	})
}

func (r *userHandler) follow(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	fu, err := r.users.GetUserByUsername(c.Param("username"))
	if err != nil {
		return err
	}

	err = r.users.UpdateUserByEmail(
		em,
		func(u *userdomain.User) (*userdomain.User, error) {
			u.StartFollowing(fu)
			return u, nil
		})
	if err != nil {
		return err
	}

	found, err := r.users.GetUserByEmail(fu.Email())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, profile{
		profileUser{
			Username:  found.Username(),
			Bio:       found.Bio(),
			Image:     found.Image(),
			Following: found.IsFollowing(fu),
		},
	})
}

func (r *userHandler) unfollow(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	fu, err := r.users.GetUserByUsername(c.Param("username"))
	if err != nil {
		return err
	}

	err = r.users.UpdateUserByEmail(
		em,
		func(u *userdomain.User) (*userdomain.User, error) {
			u.StopFollowing(fu)
			return u, nil
		})
	if err != nil {
		return err
	}

	found, err := r.users.GetUserByEmail(fu.Email())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, profile{
		profileUser{
			Username:  found.Username(),
			Bio:       found.Bio(),
			Image:     found.Image(),
			Following: found.IsFollowing(fu),
		},
	})
}

func optional(s string) *string {
	return &s
}
