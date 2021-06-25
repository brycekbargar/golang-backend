package echohttp

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/brycekbargar/realworld-backend/ports/serialization"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

func init() {
	slug.CustomSub = map[string]string{
		"feed": "f",
	}
}

type articlesHandler struct {
	repo        domain.Repository
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
}

func newArticlesHandler(
	repo domain.Repository,
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
) *articlesHandler {
	return &articlesHandler{
		repo,
		authed,
		maybeAuthed,
	}
}

func (h *articlesHandler) mapRoutes(g *echo.Group) {
	g.GET("/articles", h.list, h.maybeAuthed)
	g.GET("/articles/feed", h.feed, h.authed)
	g.GET("/articles/:slug", h.article, h.maybeAuthed)
	g.POST("/articles", h.create, h.authed)
	g.PUT("/articles/:slug", h.update, h.authed)
	g.DELETE("/articles/:slug", h.delete, h.authed)

	g.GET("/articles/:slug/comments", h.commentList, h.maybeAuthed)
	g.POST("/articles/:slug/comments", h.addComment, h.authed)
	g.DELETE("/articles/:slug/comments/:id", h.removeComment, h.authed)

	g.POST("/articles/:slug/favorite", h.favorite, h.authed)
	g.DELETE("/articles/:slug/favorite", h.unfavorite, h.authed)

	g.GET("/tags", h.tags)
}

type author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func (h *articlesHandler) list(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	var u *domain.Fanboy
	if ok {
		u, _ = h.repo.GetUserByEmail(em)
	}

	lc := domain.ListCriteria{
		Tag:   ctx.QueryParam("tag"),
		Limit: 20,
	}

	a := ctx.QueryParam("author")
	if len(a) > 0 {
		if ae, err := h.repo.GetUserByUsername(a); err == nil {
			lc.AuthorEmails = []string{ae.Email}
		}
	}
	f := ctx.QueryParam("favorited")
	if len(f) > 0 {
		if fe, err := h.repo.GetUserByUsername(f); err == nil {
			lc.FavoritedByUserEmail = fe.Email
		}
	}
	l := ctx.QueryParam("limit")
	if li, err := strconv.Atoi(l); err == nil {
		lc.Limit = li
	}
	o := ctx.QueryParam("offset")
	if oi, err := strconv.Atoi(o); err == nil {
		lc.Offset = oi
	}

	// get all articles
	al, err := h.repo.LatestArticlesByCriteria(lc)
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.ManyAuthoredArticlesToArticles(al, u))
}

func (h *articlesHandler) feed(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.repo.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	lc := domain.ListCriteria{
		Limit:        20,
		AuthorEmails: u.FollowingEmails(),
	}
	l := ctx.QueryParam("limit")
	if li, err := strconv.Atoi(l); err == nil {
		lc.Limit = li
	}
	o := ctx.QueryParam("offset")
	if oi, err := strconv.Atoi(o); err == nil {
		lc.Offset = oi
	}

	// Get the feed articles
	al, err := h.repo.LatestArticlesByCriteria(lc)
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.ManyAuthoredArticlesToArticles(al, u))
}

func (h *articlesHandler) article(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	var u *domain.Fanboy
	if ok {
		u, _ = h.repo.GetUserByEmail(em)
	}

	// get the article
	ar, err := h.repo.GetArticleBySlug(ctx.Param("slug"))
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.AuthoredArticleToArticle(ar, u))
}

func (h *articlesHandler) create(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.repo.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	article, err := serialization.CreateToArticle(ctx.Bind, u)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	created, err := h.repo.CreateArticle(article)
	if err != nil {
		if err == domain.ErrDuplicateArticle {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				err)
		}
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.AuthoredArticleToArticle(created, u))
}

func (h *articlesHandler) update(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.repo.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	delta, err := serialization.UpdateArticleToDelta(ctx.Bind)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	updated, err := h.repo.UpdateArticleBySlug(
		ctx.Param("slug"),
		func(a *domain.Article) (*domain.Article, error) {
			if a.AuthorEmail != em {
				return nil, errors.New("articles can only be updated by their author")
			}

			delta(a)
			return a.Validate()
		})
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.AuthoredArticleToArticle(updated, u))
}

func (h *articlesHandler) delete(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	ar, err := h.repo.GetArticleBySlug(ctx.Param("slug"))
	if err != nil {
		return err
	}
	if ar.AuthorEmail != em {
		return errors.New("articles can only be deleted by their author")
	}

	if err = h.repo.DeleteArticle(&ar.Article); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *articlesHandler) commentList(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	var u *domain.Fanboy
	if ok {
		u, _ = h.repo.GetUserByEmail(em)
	}

	ar, err := h.repo.GetCommentsBySlug(ctx.Param("slug"))
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.ArticleToCommentList(ar, h.repo.GetAuthorByEmail, u))
}

func (h *articlesHandler) addComment(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.repo.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	body, err := serialization.CommentToBody(ctx.Bind)
	if err != nil {
		return echo.ErrBadRequest
	}

	// Make the thing
	newc, err := h.repo.UpdateCommentsBySlug(
		ctx.Param("slug"),
		func(a *domain.CommentedArticle) (*domain.CommentedArticle, error) {
			return a, a.AddComment(body, em)
		})
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.CommentToComment(newc, u, u))
}

func (h *articlesHandler) removeComment(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	c := ctx.Param("id")
	cid, err := strconv.Atoi(c)
	if err != nil {
		return echo.ErrBadRequest
	}

	// Delete the thing
	_, err = h.repo.UpdateCommentsBySlug(
		ctx.Param("slug"),
		func(a *domain.CommentedArticle) (*domain.CommentedArticle, error) {
			for _, c := range a.Comments {
				if c.ID == cid && c.AuthorEmail != em {
					return nil, errors.New("comments can only be deleted by their author")
				}
			}

			a.RemoveComment(cid)
			return a, nil
		})
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *articlesHandler) favorite(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.repo.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	s := ctx.Param("slug")
	found, err := h.repo.GetArticleBySlug(s)
	if err != nil {
		return err
	}

	err = h.repo.UpdateFanboyByEmail(
		em,
		func(u *domain.Fanboy) (*domain.Fanboy, error) {
			u.Favorite(s)
			return u, nil
		})
	if err != nil {
		return err
	}

	u.Favorite(s)
	return ctx.JSON(
		http.StatusOK,
		serialization.AuthoredArticleToArticle(found, u))
}

func (h *articlesHandler) unfavorite(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.repo.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	s := ctx.Param("slug")
	found, err := h.repo.GetArticleBySlug(s)
	if err != nil {
		return err
	}

	err = h.repo.UpdateFanboyByEmail(
		em,
		func(u *domain.Fanboy) (*domain.Fanboy, error) {
			u.Favorite(s)
			return u, nil
		})
	if err != nil {
		return err
	}

	u.Unfavorite(s)
	return ctx.JSON(
		http.StatusOK,
		serialization.AuthoredArticleToArticle(found, u))
}

func (h *articlesHandler) tags(ctx echo.Context) error {
	tags, err := h.repo.DistinctTags()
	if err != nil {
		return err
	}

	return ctx.JSON(
		http.StatusOK,
		serialization.TagsToTaglist(tags))
}
