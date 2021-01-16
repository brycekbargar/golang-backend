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

	return &Article{
		slug,
		title,
		description,
		body,
		tags,
		time.Now().UTC(),
		time.Now().UTC(),
		authorEmail,
	}, nil
}

// ExistingArticle creates an Article with the provided information.
func ExistingArticle(slug string, title string, description string, body string, createdAt time.Time, updatedAt time.Time, authorEmail string, tags ...string) (*Article, error) {
	if len(slug) == 0 ||
		len(title) == 0 ||
		len(description) == 0 ||
		len(body) == 0 ||
		len(authorEmail) == 0 {
		return nil, ErrRequiredArticleFields
	}

	if !id.IsSlug(slug) {
		return nil, ErrInvalidSlug
	}

	return &Article{
		slug,
		title,
		description,
		body,
		tags,
		createdAt,
		updatedAt,
		authorEmail,
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

// Tags is the slice of tags associated with the Article on creation.
func (a Article) Tags() []string {
	ts := make([]string, 0, len(a.tagList))
	copy(ts, a.tagList)
	return ts
}
