package articledomain

import (
	"time"

	"github.com/gosimple/slug"
)

// Article is an individual post in the application.
type Article struct {
	Slug        string
	Title       string
	Description string
	Body        string
	TagList     []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Author      string
}

// NewArticle creates a new Article with the provide information and defaults for the rest.
func NewArticle(title string, description string, body string, authorEmail string, tags ...string) (*Article, error) {
	return &Article{
		slug.Make(title),
		title,
		description,
		body,
		tags,
		time.Now().UTC(),
		time.Now().UTC(),
		authorEmail,
	}, nil
}
