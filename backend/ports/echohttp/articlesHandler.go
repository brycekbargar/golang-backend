package echohttp

import (
	"github.com/labstack/echo/v4"

	"github.com/brycekbargar/realworld-backend/ports"
)

type articlesHandler struct {
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
	jc          ports.JWTConfig
}

func newArticlesHandler(
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
	jc ports.JWTConfig) *articlesHandler {
	return &articlesHandler{
		authed,
		maybeAuthed,
		jc,
	}
}

func (r *articlesHandler) mapRoutes(g *echo.Group) {
	g.GET("/articles", r.list, r.maybeAuthed)
}

func (r *articlesHandler) list(c echo.Context) error {
	em, _, _ := c.(*userContext).identity()

	// get all articles
	if len(em) > 0 {
		// set the following logic?
	}

	return nil
}
