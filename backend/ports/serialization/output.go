package serialization

import (
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
)

type user struct {
	User userUser `json:"user"`
}
type userUser struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// UserToUser converts a domain user to an output serializable user.
func UserToUser(
	u *domain.User,
	t string,
) interface{} {
	return &user{
		userUser{
			Email:    u.Email,
			Token:    t,
			Username: u.Username,
			Bio:      optional(u.Bio),
			Image:    optional(u.Image),
		},
	}
}

func optional(s string) *string {
	return &s
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

// Following is a convenience func for when the user is following a profile.
func Following(*domain.User) bool { return true }

// NotFollowing is a convenience func for when the user is not following a profile.
func NotFollowing(*domain.User) bool { return false }

// MaybeFollowing is a convenience func for checking if the  user is following a profile.
func MaybeFollowing(f *domain.Fanboy) func(u *domain.User) bool {
	return func(u *domain.User) bool { return f.IsFollowing(u.Email) }
}

// UserToProfile converts a domain user to an output serializable profile.
func UserToProfile(
	u *domain.User,
	f func(*domain.User) bool,
) interface{} {
	return &profile{
		profileUser{
			Username:  u.Username,
			Bio:       u.Bio,
			Image:     u.Image,
			Following: f(u),
		},
	}
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

type article struct {
	Article interface{} `json:"article"`
}

type list struct {
	Articles      []interface{} `json:"articles"`
	ArticlesCount int           `json:"articlesCount"`
}

func internalArticle(
	a *domain.AuthoredArticle,
	cu *domain.Fanboy,
) interface{} {
	return &articleArticle{
		Slug:           a.Slug,
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		TagList:        a.TagList,
		CreatedAt:      a.CreatedAtUTC,
		UpdatedAt:      a.UpdatedAtUTC,
		FavoritesCount: a.FavoriteCount,
		Favorited:      cu != nil && cu.Favors(a.Slug),
		Author: author{
			Username:  a.GetUsername(),
			Bio:       a.GetBio(),
			Image:     a.GetImage(),
			Following: cu != nil && cu.IsFollowing(a.GetEmail()),
		},
	}
}

// AuthoredArticleToArticle converts a domain article into an output serialiable article for the current user.
func AuthoredArticleToArticle(
	a *domain.AuthoredArticle,
	cu *domain.Fanboy,
) interface{} {
	return &article{internalArticle(a, cu)}
}

// ManyAuthoredArticlesToArticles converts multiple domain articles into an output serialiable list of articles for the current user.
func ManyAuthoredArticlesToArticles(
	as []*domain.AuthoredArticle,
	cu *domain.Fanboy,
) interface{} {
	res := list{
		make([]interface{}, 0, len(as)),
		len(as),
	}
	for _, a := range as {
		res.Articles = append(res.Articles, internalArticle(a, cu))
	}

	return res
}

type commentComment struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    author    `json:"author"`
}

type comment struct {
	Comment interface{} `json:"comment"`
}

type commentList struct {
	Comments []interface{} `json:"comments"`
}

func internalComment(
	c *domain.Comment,
	a domain.Author,
	cu *domain.Fanboy,
) interface{} {
	return &commentComment{
		c.ID,
		c.CreatedAtUTC,
		c.CreatedAtUTC,
		c.Body,
		author{
			a.GetUsername(),
			a.GetBio(),
			a.GetImage(),
			cu.IsFollowing(a.GetEmail()),
		},
	}

}

func CommentToComment(
	c *domain.Comment,
	a domain.Author,
	cu *domain.Fanboy,
) interface{} {
	return &comment{internalComment(c, a, cu)}
}

func ArticleToCommentList(
	ar *domain.CommentedArticle,
	author func(string) domain.Author,
	cu *domain.Fanboy,
) interface{} {
	res := commentList{
		make([]interface{}, 0, len(ar.Comments)),
	}
	for _, c := range ar.Comments {
		if a := author(c.AuthorEmail); a != nil {
			res.Comments = append(res.Comments, internalComment(c, a, cu))
		}
	}

	return res
}

type tagList struct {
	Tags []string `json:"tags"`
}

// TagsToTagList converts a list of article tags to a output serializable list.
func TagsToTaglist(
	ts []string,
) interface{} {
	return &tagList{ts}
}
