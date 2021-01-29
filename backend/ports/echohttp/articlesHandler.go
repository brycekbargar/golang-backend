package echohttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

func init() {
	slug.CustomSub["feed"] = "f"
}

type articlesHandler struct {
	users       userdomain.Repository
	articles    articledomain.Repository
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
}

func newArticlesHandler(
	users userdomain.Repository,
	articles articledomain.Repository,
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc,
) *articlesHandler {
	return &articlesHandler{
		users,
		articles,
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

type articleArticle struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         author    `json:"author"`
}

type list struct {
	Articles      []articleArticle `json:"articles"`
	ArticlesCount int              `json:"articlesCount"`
}

func (h *articlesHandler) list(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	var u *userdomain.User
	if ok {
		u, _ = h.users.GetUserByEmail(em)
	}

	lc := articledomain.ListCriteria{
		Tag:   ctx.QueryParam("tag"),
		Limit: 20,
	}

	a := ctx.QueryParam("author")
	if len(a) > 0 {
		if ae, err := h.users.GetUserByUsername(a); err == nil {
			lc.AuthorEmails = []string{ae.Email()}
		}
	}
	f := ctx.QueryParam("favorited")
	if len(f) > 0 {
		if fe, err := h.users.GetUserByUsername(f); err == nil {
			lc.FavoritedByUserEmail = fe.Email()
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
	al, err := h.articles.LatestArticlesByCriteria(lc)
	if err != nil {
		return err
	}

	res := list{
		make([]articleArticle, len(al)),
		len(al),
	}

	for _, a := range al {
		aa := articleArticle{
			Slug:           a.Slug(),
			Title:          a.Title(),
			Description:    a.Description(),
			Body:           a.Body(),
			TagList:        a.Tags(),
			CreatedAt:      a.CreatedAtUTC(),
			UpdatedAt:      a.UpdatedAtUTC(),
			FavoritesCount: a.FavoriteCount(),
			Author: author{
				Username: a.Email(),
				Bio:      a.Bio(),
				Image:    a.Image(),
			},
		}
		if u != nil {
			aa.Favorited = a.IsAFavoriteOf(u.Email())
			aa.Author.Following = u.IsFollowing(a.Email())
		}

		res.Articles = append(res.Articles, aa)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h *articlesHandler) feed(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.users.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	lc := articledomain.ListCriteria{
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
	al, err := h.articles.LatestArticlesByCriteria(lc)
	if err != nil {
		return err
	}

	res := list{
		make([]articleArticle, len(al)),
		len(al),
	}

	for _, a := range al {
		aa := articleArticle{
			Slug:           a.Slug(),
			Title:          a.Title(),
			Description:    a.Description(),
			Body:           a.Body(),
			TagList:        a.Tags(),
			CreatedAt:      a.CreatedAtUTC(),
			UpdatedAt:      a.UpdatedAtUTC(),
			FavoritesCount: a.FavoriteCount(),
			Author: author{
				Username: a.Email(),
				Bio:      a.Bio(),
				Image:    a.Image(),
			},
		}
		if u != nil {
			aa.Favorited = a.IsAFavoriteOf(u.Email())
			aa.Author.Following = u.IsFollowing(a.Email())
		}

		res.Articles = append(res.Articles, aa)
	}

	return ctx.JSON(http.StatusOK, res)
}

type article struct {
	Article articleArticle `json:"article"`
}

func (h *articlesHandler) article(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.users.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	// get the article
	ar, err := h.articles.GetArticleBySlug(ctx.Param("slug"))
	if err != nil {
		return err
	}

	res := article{
		articleArticle{
			Slug:           ar.Slug(),
			Title:          ar.Title(),
			Description:    ar.Description(),
			Body:           ar.Body(),
			TagList:        ar.Tags(),
			CreatedAt:      ar.CreatedAtUTC(),
			UpdatedAt:      ar.UpdatedAtUTC(),
			FavoritesCount: ar.FavoriteCount(),
			Author: author{
				Username: ar.Email(),
				Bio:      ar.Bio(),
				Image:    ar.Image(),
			},
		},
	}
	if u != nil {
		res.Article.Favorited = ar.IsAFavoriteOf(u.Email())
		res.Article.Author.Following = u.IsFollowing(ar.Email())
	}

	return ctx.JSON(http.StatusOK, res)
}

type createArticle struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList,omitempty"`
}
type create struct {
	Article createArticle `json:"article"`
}

func (h *articlesHandler) create(ctx echo.Context) error {
	em, _, ok := ctx.(*userContext).identity()
	u, err := h.users.GetUserByEmail(em)
	if !ok || err != nil {
		return identityNotOk
	}

	a := new(create)
	if err := ctx.Bind(a); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	created, err := articledomain.NewArticle(
		a.Article.Title,
		a.Article.Description,
		a.Article.Body,
		em,
		a.Article.TagList...,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			err)
	}

	if err := h.articles.Create(created); err != nil {
		if err == userdomain.ErrDuplicateValue {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				err)
		}
		return err
	}

	return ctx.JSON(http.StatusCreated, article{
		articleArticle{
			Slug:           created.Slug(),
			Title:          created.Title(),
			Description:    created.Description(),
			Body:           created.Body(),
			TagList:        created.Tags(),
			CreatedAt:      created.CreatedAtUTC(),
			UpdatedAt:      created.UpdatedAtUTC(),
			Favorited:      false,
			FavoritesCount: created.FavoriteCount(),
			Author: author{
				Username:  u.Email(),
				Bio:       u.Bio(),
				Image:     u.Image(),
				Following: u.IsFollowing(u.Email()),
			},
		},
	})
}

func (h *articlesHandler) update(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	a := new(create)
	if err := c.Bind(a); err != nil {
		return echo.ErrBadRequest
	}

	// Update the thing

	return c.JSON(http.StatusOK, article{})
}

func (h *articlesHandler) delete(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// Delete the thing

	return c.NoContent(http.StatusOK)
}

type commentComment struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    author    `json:"author"`
}

type commentList struct {
	Comments []commentComment `json:"comments"`
}

func (h *articlesHandler) commentList(c echo.Context) error {
	em, _, _ := c.(*userContext).identity()

	// get all comments
	if len(em) > 0 {
		// set the following logic?
	}

	return c.JSON(http.StatusOK, commentList{
		make([]commentComment, 0),
	})
}

type comment struct {
	Comment commentComment `json:"comment"`
}

type addCommentComment struct {
	Body string `json:"body"`
}

type addComment struct {
	Comment addCommentComment `json:"comment"`
}

func (h *articlesHandler) addComment(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	a := new(create)
	if err := c.Bind(a); err != nil {
		return echo.ErrBadRequest
	}

	// Make the thing

	return c.JSON(http.StatusCreated, comment{})
}

func (h *articlesHandler) removeComment(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// Delete the thing

	return c.NoContent(http.StatusOK)
}

func (h *articlesHandler) favorite(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// favorite

	return c.JSON(http.StatusOK, article{})
}

func (h *articlesHandler) unfavorite(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// unfavorite

	return c.JSON(http.StatusOK, article{})
}

type tagList struct {
	Tags []string `json:"tags"`
}

func (h *articlesHandler) tags(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// read all the tags

	return c.JSON(http.StatusOK, tagList{})
}
