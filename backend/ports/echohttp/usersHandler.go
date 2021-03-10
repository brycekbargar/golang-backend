package echohttp

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/brycekbargar/realworld-backend/ports"
)

type usersHandler struct {
	users       userdomain.Repository
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	jc          ports.JWTConfig
}

func newUsersHandler(
	users userdomain.Repository,
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	jc ports.JWTConfig) *usersHandler {
	return &usersHandler{
		users,
		authed,
		maybeAuthed,
		jc,
	}
}

func (r *usersHandler) mapRoutes(g *echo.Group) {
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

func makeJwt(r *usersHandler, e string) (string, error) {
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

func (r *usersHandler) create(c echo.Context) error {
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
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	if _, err := r.users.CreateUser(created); err != nil {
		if err == userdomain.ErrDuplicateValue {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				err)
		}
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

func (r *usersHandler) login(c echo.Context) error {
	l := new(login)
	if err := c.Bind(l); err != nil {
		return echo.ErrBadRequest
	}

	authed, err := r.users.GetUserByEmail(l.User.Email)
	if err != nil {
		return echo.ErrUnauthorized
	}

	if pw, err := authed.HasPassword(l.User.Password); !pw || err != nil {
		return echo.ErrUnauthorized
	}

	token, err := makeJwt(r, authed.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    authed.Email,
			Username: authed.Username,
			Token:    token,
			Bio:      optional(authed.Bio),
			Image:    optional(authed.Image),
		},
	})
}

func (r *usersHandler) user(c echo.Context) error {
	em, token, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	found, err := r.users.GetUserByEmail(em)
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    found.Email,
			Username: found.Username,
			Token:    token.Raw,
			Bio:      optional(found.Bio),
			Image:    optional(found.Image),
		},
	})
}

func (r *usersHandler) update(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	b := new(user)
	if err := c.Bind(b); err != nil {
		return echo.ErrBadRequest
	}

	found, err := r.users.UpdateUserByEmail(
		em,
		func(u *userdomain.User) (*userdomain.User, error) {
			if b.User.Email != "" {
				u.Email = b.User.Email
			}
			if b.User.Username != "" {
				u.Username = b.User.Username
			}
			if b.User.Password != nil && *b.User.Password != "" {
				u.SetPassword(*b.User.Password)
			}
			if b.User.Bio != nil {
				u.Bio = *b.User.Bio
			}
			if b.User.Image != nil {
				u.Image = *b.User.Image
			}
			return u.Validate()
		})
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	// Users can change their email so we need to make sure we're giving them a new token.
	token, err := makeJwt(r, found.Email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user{
		userUser{
			Email:    found.Email,
			Username: found.Username,
			Token:    token,
			Bio:      optional(found.Bio),
			Image:    optional(found.Image),
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

func (r *usersHandler) profile(c echo.Context) (err error) {
	em, _, _ := c.(*userContext).identity()

	if len(c.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	found, err := r.users.GetUserByUsername(c.Param("username"))
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	following := false
	if len(em) > 0 {
		cu, err := r.users.GetUserByEmail(em)
		if err != nil {
			if err == userdomain.ErrNotFound {
				return echo.ErrNotFound
			}
			return err
		}

		following = cu.IsFollowing(found.Email)
	}

	return c.JSON(http.StatusOK, profile{
		profileUser{
			Username:  found.Username,
			Bio:       found.Bio,
			Image:     found.Image,
			Following: following,
		},
	})
}

func (r *usersHandler) follow(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	if len(c.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	fu, err := r.users.GetUserByUsername(c.Param("username"))
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	err = r.users.UpdateFanboyByEmail(
		em,
		func(u *userdomain.Fanboy) (*userdomain.Fanboy, error) {
			u.StartFollowing(fu.Email)
			return u, nil
		})
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, profile{
		profileUser{
			Username:  fu.Username,
			Bio:       fu.Bio,
			Image:     fu.Image,
			Following: true,
		},
	})
}

func (r *usersHandler) unfollow(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	if len(c.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	fu, err := r.users.GetUserByUsername(c.Param("username"))
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	err = r.users.UpdateFanboyByEmail(
		em,
		func(u *userdomain.Fanboy) (*userdomain.Fanboy, error) {
			u.StopFollowing(fu.Email)
			return u, nil
		})
	if err != nil {
		if err == userdomain.ErrNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(http.StatusOK, profile{
		profileUser{
			Username:  fu.Username,
			Bio:       fu.Bio,
			Image:     fu.Image,
			Following: false,
		},
	})
}

func optional(s string) *string {
	return &s
}
