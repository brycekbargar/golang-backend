package echohttp

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/brycekbargar/realworld-backend/ports"
	"github.com/brycekbargar/realworld-backend/ports/serialization"
)

type usersHandler struct {
	repo        domain.Repository
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	jc          ports.JWTConfig
}

func newUsersHandler(
	repo domain.Repository,
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	jc ports.JWTConfig,
) *usersHandler {
	return &usersHandler{
		repo,
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
	user, err := serialization.RegisterToUser(c.Bind)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	created, err := r.repo.CreateUser(user)
	if err != nil {
		if err == domain.ErrDuplicateUser {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				err)
		}
		return err
	}

	token, err := makeJwt(r, created.Email)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		serialization.UserToUser(created, token))
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

	authed, err := r.repo.GetUserByEmail(l.User.Email)
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

	return c.JSON(
		http.StatusOK,
		serialization.UserToUser(&authed.User, token))
}

func (r *usersHandler) user(c echo.Context) error {
	em, token, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	found, err := r.repo.GetUserByEmail(em)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(
		http.StatusOK,
		serialization.UserToUser(&found.User, token.Raw))
}

func (r *usersHandler) update(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	delta, err := serialization.UpdateToDelta(c.Bind)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	updated, err := r.repo.UpdateUserByEmail(
		em,
		func(u *domain.User) (*domain.User, error) {
			delta(u)
			return u.Validate()
		})
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	// Users can change their email so we need to make sure we're giving them a new token.
	token, err := makeJwt(r, updated.Email)
	if err != nil {
		return err
	}

	return c.JSON(
		http.StatusOK,
		serialization.UserToUser(updated, token))
}

func (r *usersHandler) profile(c echo.Context) (err error) {
	em, _, _ := c.(*userContext).identity()

	if len(c.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	found, err := r.repo.GetUserByUsername(c.Param("username"))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	following := serialization.NotFollowing
	if len(em) > 0 {
		cu, err := r.repo.GetUserByEmail(em)
		if err != nil {
			if err == domain.ErrUserNotFound {
				return echo.ErrNotFound
			}
			return err
		}

		following = serialization.MaybeFollowing(cu)
	}

	return c.JSON(
		http.StatusOK,
		serialization.UserToProfile(found, following))
}

func (r *usersHandler) follow(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	if len(c.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	fu, err := r.repo.GetUserByUsername(c.Param("username"))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	err = r.repo.UpdateFanboyByEmail(
		em,
		func(u *domain.Fanboy) (*domain.Fanboy, error) {
			u.StartFollowing(fu.Email)
			return u, nil
		})
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(
		http.StatusOK,
		serialization.UserToProfile(fu, serialization.Following))
}

func (r *usersHandler) unfollow(c echo.Context) error {
	em, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	if len(c.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	fu, err := r.repo.GetUserByUsername(c.Param("username"))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	err = r.repo.UpdateFanboyByEmail(
		em,
		func(u *domain.Fanboy) (*domain.Fanboy, error) {
			u.StopFollowing(fu.Email)
			return u, nil
		})
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return c.JSON(
		http.StatusOK,
		serialization.UserToProfile(fu, serialization.NotFollowing))
}
