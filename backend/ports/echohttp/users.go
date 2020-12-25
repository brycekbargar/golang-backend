package echohttp

import "github.com/labstack/echo/v4"

type userHandler struct {
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	key         []byte
}

func newUserHandler(
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	key []byte) *userHandler {
	return &userHandler{
		authed,
		maybeAuthed,
		key,
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

func (r *userHandler) create(c echo.Context) (err error) {
	return nil
}

func (r *userHandler) login(c echo.Context) (err error) {
	return nil

}

func (r *userHandler) user(c echo.Context) (err error) {
	return nil
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
