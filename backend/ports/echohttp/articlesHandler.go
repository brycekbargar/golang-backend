package echohttp

import (
	"net/http"
	"time"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
)

func init() {
	slug.CustomSub["feed"] = "f"
}

type articlesHandler struct {
	authed      echo.MiddlewareFunc
	maybeAuthed echo.MiddlewareFunc
}

func newArticlesHandler(
	authed echo.MiddlewareFunc,
	maybeAuthed echo.MiddlewareFunc) *articlesHandler {
	return &articlesHandler{
		authed,
		maybeAuthed,
	}
}

func (r *articlesHandler) mapRoutes(g *echo.Group) {
	g.GET("/articles", r.list, r.maybeAuthed)
	g.GET("/articles/feed", r.feed, r.authed)
	g.GET("/articles/:slug", r.article, r.maybeAuthed)
	g.POST("/articles", r.create, r.authed)
	g.PUT("/articles/:slug", r.update, r.authed)
	g.DELETE("/articles/:slug", r.delete, r.authed)

	g.GET("/articles/:slug/comments", r.commentList, r.maybeAuthed)
	g.POST("/articles/:slug/comments", r.addComment, r.authed)
	g.DELETE("/articles/:slug/comments/:id", r.removeComment, r.authed)

	g.POST("/articles/:slug/favorite", r.favorite, r.authed)
	g.DELETE("/articles/:slug/favorite", r.unfavorite, r.authed)

	g.GET("/tags", r.tags)
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

func (r *articlesHandler) list(c echo.Context) error {
	em, _, _ := c.(*userContext).identity()

	// get all articles
	if len(em) > 0 {
		// set the following/favorited logic?
	}

	return c.JSON(http.StatusOK, list{
		make([]articleArticle, 0),
		0,
	})
}

func (r *articlesHandler) feed(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// Get the feed articles

	return c.JSON(http.StatusOK, list{
		make([]articleArticle, 0),
		0,
	})
}

type article struct {
	Article articleArticle `json:"article"`
}

func (r *articlesHandler) article(c echo.Context) error {
	em, _, _ := c.(*userContext).identity()

	// get the article
	if len(em) > 0 {
		// set the following/favorited logic?
	}

	return c.JSON(http.StatusOK, article{})
}

type createArticle struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList,omitempty"`
}
type create struct {
	Article articleArticle `json:"article"`
}

func (r *articlesHandler) create(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	a := new(create)
	if err := c.Bind(a); err != nil {
		return echo.ErrBadRequest
	}

	// Make the thing

	return c.JSON(http.StatusCreated, article{})
}

func (r *articlesHandler) update(c echo.Context) error {
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

func (r *articlesHandler) delete(c echo.Context) error {
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

func (r *articlesHandler) commentList(c echo.Context) error {
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

func (r *articlesHandler) addComment(c echo.Context) error {
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

func (r *articlesHandler) removeComment(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// Delete the thing

	return c.NoContent(http.StatusOK)
}

func (r *articlesHandler) favorite(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// favorite

	return c.JSON(http.StatusOK, article{})
}

func (r *articlesHandler) unfavorite(c echo.Context) error {
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

func (r *articlesHandler) tags(c echo.Context) error {
	_, _, ok := c.(*userContext).identity()
	if !ok {
		return identityNotOk
	}

	// read all the tags

	return c.JSON(http.StatusOK, tagList{})
}
