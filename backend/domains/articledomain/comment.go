package articledomain

import (
	"errors"
	"time"
)

// ErrRequiredCommentFields indicates when a Comment is instantiated without all the required fields.
var ErrRequiredCommentFields = errors.New("id, body, and author are required for comments")

// Comment is an individual comment associated with a single Article.
type Comment struct {
	id           int
	body         string
	createdAtUTC time.Time
	author       string
}

func newComment(id int, body string, author string) (*Comment, error) {
	if id == 0 || len(body) == 0 || len(author) == 0 {
		return nil, ErrRequiredCommentFields
	}

	return &Comment{
		id,
		body,
		time.Now().UTC(),
		author,
	}, nil
}

// ExistingComment creates a comment with the provided information.
func ExistingComment(id int, body string, author string, createdAt time.Time) (*Comment, error) {
	if id == 0 || len(body) == 0 || len(author) == 0 {
		return nil, ErrRequiredCommentFields
	}

	return &Comment{
		id,
		body,
		createdAt,
		author,
	}, nil
}

// ID is the identifier (not globally unique) of the comment on the parent article.
func (c Comment) ID() int {
	return c.id
}

// Body is the comment's content.
func (c Comment) Body() string {
	return c.body
}

// CreatedAtUTC is the system generated time (in utc) when the comment was created.
func (c Comment) CreatedAtUTC() string {
	return c.CreatedAtUTC()
}

// UpdatedAtUTC is the system generated time (in utc) when the article was last updated.
// Currently updating comments isn't supported so this will always be the CreateAtUTC time.
func (c Comment) UpdatedAtUTC() string {
	return c.CreatedAtUTC()
}

// AuthorEmail is the email (the identifier) of the user that created the Comment.
func (c Comment) AuthorEmail() string {
	return c.author
}
