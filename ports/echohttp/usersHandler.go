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

func (h *usersHandler) create(ctx echo.Context) error {
	user, err := serialization.RegisterToUser(ctx.Bind)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	created, err := h.repo.CreateUser(ctx.Request().Context(), user)
	if err != nil {
		if err == domain.ErrDuplicateUser {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				err)
		}
		return err
	}

	token, err := makeJwt(h, created.Email)
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToUser(created, token))
}

func (h *usersHandler) login(ctx echo.Context) error {
	em, pw, err := serialization.LoginToCredentials(ctx.Bind)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	authed, err := h.repo.GetUserByEmail(ctx.Request().Context(), em)
	if err != nil {
		return echo.ErrUnauthorized
	}

	if ok, err := authed.HasPassword(pw); !ok || err != nil {
		return echo.ErrUnauthorized
	}

	token, err := makeJwt(h, authed.Email)
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToUser(&authed.User, token))
}

func (h *usersHandler) user(ctx echo.Context) error {
	em, token, ok := ctx.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	found, err := h.repo.GetUserByEmail(ctx.Request().Context(), em)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToUser(&found.User, token.Raw))
}

func (h *usersHandler) update(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	delta, err := serialization.UpdateUserToDelta(ctx.Bind)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	updated, err := h.repo.UpdateUserByEmail(ctx.Request().Context(),
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
	token, err := makeJwt(h, updated.Email)
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToUser(updated, token))
}

func (h *usersHandler) profile(ctx echo.Context) (err error) {
	em, _, _ := ctx.(*userContext).identity()

	if len(ctx.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	found, err := h.repo.GetUserByUsername(ctx.Request().Context(), ctx.Param("username"))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	following := serialization.NotFollowing
	if len(em) > 0 {
		cu, err := h.repo.GetUserByEmail(ctx.Request().Context(), em)
		if err != nil {
			if err == domain.ErrUserNotFound {
				return echo.ErrNotFound
			}
			return err
		}

		following = serialization.MaybeFollowing(cu)
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToProfile(found, following))
}

func (h *usersHandler) follow(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	if len(ctx.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	found, err := h.repo.GetUserByUsername(ctx.Request().Context(), ctx.Param("username"))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	err = h.repo.UpdateFanboyByEmail(ctx.Request().Context(),
		em,
		func(u *domain.Fanboy) (*domain.Fanboy, error) {
			u.StartFollowing(found.Email)
			return u, nil
		})
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToProfile(found, serialization.Following))
}

func (h *usersHandler) unfollow(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	if len(ctx.Param("username")) == 0 {
		return echo.ErrBadRequest
	}

	found, err := h.repo.GetUserByUsername(ctx.Request().Context(), ctx.Param("username"))
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	err = h.repo.UpdateFanboyByEmail(ctx.Request().Context(),
		em,
		func(u *domain.Fanboy) (*domain.Fanboy, error) {
			u.StopFollowing(found.Email)
			return u, nil
		})
	if err != nil {
		if err == domain.ErrUserNotFound {
			return echo.ErrNotFound
		}
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.UserToProfile(found, serialization.NotFollowing))
}
