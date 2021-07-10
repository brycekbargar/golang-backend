package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gosimple/slug"
)

func init() {
	govalidator.TagMap["slug"] = govalidator.Validator(func(str string) bool {
		return slug.IsSlug(str)
	})
}

// Article is an individual post in the application.
type Article struct {
	Slug         string `valid:"required,slug"`
	Title        string `valid:"required"`
	Description  string `valid:"required"`
	Body         string `valid:"required"`
	TagList      []string
	CreatedAtUTC time.Time
	UpdatedAtUTC time.Time
	AuthorEmail  string `valid:"required,email"`
}

// CommentedArticle is an individual post in the application with its comment information included.
type CommentedArticle struct {
	Article
	Comments []Comment
}

// AuthoredArticle is an individual post in the application with its author information included.
type AuthoredArticle struct {
	Article
	Author
	FavoriteCount int
}

// Author is the author of an article.
type Author interface {
	GetUsername() string
	GetEmail() string
	GetBio() string
	GetImage() string
}

func (u User) GetUsername() string {
	return u.Username
}
func (u User) GetEmail() string {
	return u.Email
}
func (u User) GetBio() string {
	return u.Bio
}
func (u User) GetImage() string {
	return u.Image
}

// NewArticle creates a new Article with the provided information and defaults for the rest.
func NewArticle(title string, description string, body string, authorEmail string, tags ...string) (*Article, error) {
	return (&Article{
		Slug:        slug.Make(title),
		Title:       title,
		Description: description,
		Body:        body,
		TagList:     tags,
		AuthorEmail: authorEmail,
	}).Validate()
}

// Validate returns the provided Article if it is valid, otherwise error will contain validation errors.
func (a *Article) Validate() (*Article, error) {
	if v, err := govalidator.ValidateStruct(a); !v {
		return nil, err
	}

	return a, nil
}

// SetTitle sets the title and slugifies it too.
func (a *Article) SetTitle(title string) {
	a.Slug = slug.Make(title)
	a.Title = title
}

// AddComment creates a new comment and adds it to this Article.
func (a *CommentedArticle) AddComment(body string, authorEmail string) error {
	c, err := NewComment(body, authorEmail)
	if err != nil {
		return err
	}

	a.Comments = append(a.Comments, *c)
	return nil
}

// RemoveComment removes the comment (if it exists by id) from this Article.
func (a *CommentedArticle) RemoveComment(id int) {
	for i, c := range a.Comments {
		if c.ID == id {
			a.Comments = append(a.Comments[:i], a.Comments[i+1:]...)
			return
		}
	}
}
