package articledomain

import (
	"errors"
	"time"

	id "github.com/gosimple/slug"
)

// ErrRequiredArticleFields indicates when an Article is instantiated without all the required fields.
var ErrRequiredArticleFields = errors.New("slug, title, description, body, and author are required for articles")

// ErrInvalidSlug indicates when an Article is instantiated with an invalid slug.
var ErrInvalidSlug = errors.New("slug must be a valid slug")

// Article is an individual post in the application.
type Article struct {
	slug         string
	title        string
	description  string
	body         string
	tagList      []string
	createdAtUTC time.Time
	updatedAtUTC time.Time
	author       string
	comments     []*Comment
	favoritedBy  map[string]interface{}
}

// CommentedArticle is an individual post in the application with its comment information included.
type CommentedArticle = Article

// AuthoredArticle is an individual post in the application with its author information included.
type AuthoredArticle struct {
	Article
	Author
}

// Author is the author of an article.
type Author interface {
	Email() string
	Bio() string
	Image() string
}

// NewArticle creates a new Article with the provided information and defaults for the rest.
func NewArticle(title string, description string, body string, authorEmail string, tags ...string) (*Article, error) {
	if len(title) == 0 || len(description) == 0 || len(body) == 0 || len(authorEmail) == 0 {
		return nil, ErrRequiredArticleFields
	}

	slug := id.Make(title)
	if len(slug) == 0 || !id.IsSlug(slug) {
		return nil, ErrInvalidSlug
	}

	now := time.Now().UTC()
	return &Article{
		slug,
		title,
		description,
		body,
		tags,
		now,
		now,
		authorEmail,
		make([]*Comment, 0),
		make(map[string]interface{}),
	}, nil
}

// ExistingArticle creates an Article with the provided information.
func ExistingArticle(slug string, title string, description string, body string, tags []string, createdAt time.Time, updatedAt time.Time, author Author, comments []*Comment, favoritedBy []string) (*AuthoredArticle, error) {
	if len(slug) == 0 ||
		len(title) == 0 ||
		len(description) == 0 ||
		len(body) == 0 ||
		author == nil {
		return nil, ErrRequiredArticleFields
	}

	if !id.IsSlug(slug) {
		return nil, ErrInvalidSlug
	}

	favs := make(map[string]interface{})
	for _, f := range favoritedBy {
		favs[f] = nil
	}

	return &AuthoredArticle{
		Article{
			slug,
			title,
			description,
			body,
			tags,
			createdAt,
			updatedAt,
			author.Email(),
			comments,
			favs,
		},
		author,
	}, nil
}

// UpdatedArticle merges the provided Article with the (optional) new values provided.
func UpdatedArticle(article Article, title string, description string, body string) (*Article, error) {
	if len(title) > 0 {
		slug := id.Make(title)
		if len(slug) == 0 || !id.IsSlug(slug) {
			return nil, ErrInvalidSlug
		}
		article.slug = slug
		article.title = title
	}

	if len(description) > 0 {
		article.description = description
	}

	if len(body) > 0 {
		article.body = body
	}

	article.updatedAtUTC = time.Now()
	return &article, nil
}

// Slug is the article's identifier (derived from the title).
func (a Article) Slug() string {
	return a.slug
}

// Title is the article's user entered title.
func (a Article) Title() string {
	return a.title
}

// Description is something I don't understand about the domain...
func (a Article) Description() string {
	return a.description
}

// Body is the article's content.
func (a Article) Body() string {
	return a.body
}

// CreatedAtUTC is the system generated time (in utc) when the article was created.
func (a Article) CreatedAtUTC() time.Time {
	return a.createdAtUTC
}

// UpdatedAtUTC is the system generated time (in utc) when the article was last updated.
func (a Article) UpdatedAtUTC() time.Time {
	return a.updatedAtUTC
}

// AuthorEmail is the email (the identifier) of the user that created the Article.
func (a Article) AuthorEmail() string {
	return a.author
}

// FavoriteCount is the number of users that have Favorited this Article.
func (a Article) FavoriteCount() int {
	return len(a.favoritedBy)
}

// Tags is the slice of tags associated with the Article on creation.
func (a Article) Tags() []string {
	ts := make([]string, 0, len(a.tagList))
	copy(ts, a.tagList)
	return ts
}

// Comments is the slice of comments associated with the Article sorted in the order they were created.
func (a CommentedArticle) Comments() []*Comment {
	cs := make([]*Comment, len(a.comments))
	copy(cs, a.comments)
	return cs
}

// AddComment creates a new comment and adds it to this Article.
func (a *CommentedArticle) AddComment(body string, authorEmail string) error {
	id := 1
	for _, c := range a.comments {
		if c.id >= id {
			id = c.id + 1
		}
	}

	c, err := newComment(id, body, authorEmail)
	if err != nil {
		return err
	}

	a.comments = append(a.comments, c)
	return nil
}

// RemoveComment removes the comment (if it exists by id) from this Article.
func (a *CommentedArticle) RemoveComment(id int) {
	if id == 0 {
		return
	}

	for i, c := range a.comments {
		if c.id == id {
			a.comments = append(a.comments[:i], a.comments[i+1:]...)
			return
		}
	}
}

// IsAFavoriteOf checks to see if the given userEmail has favorited this article.
func (a *Article) IsAFavoriteOf(userEmail string) (ok bool) {
	_, ok = a.favoritedBy[userEmail]
	return
}

// Favorite marks this Article as a favorite of the given userEmail.
func (a *Article) Favorite(userEmail string) {
	a.favoritedBy[userEmail] = nil
}

// Unfavorite marks this Article as a no longer a favorite of the given userEmail.
func (a *Article) Unfavorite(userEmail string) {
	delete(a.favoritedBy, userEmail)
}
