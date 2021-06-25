package serialization

import (
	"github.com/brycekbargar/realworld-backend/domain"
)

type register struct {
	User registerUser `json:"user"`
}
type registerUser struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Password *string `json:"password,omitempty"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// RegisterToUser converts a input serializable user to a domain user.
func RegisterToUser(
	bind func(interface{}) error,
) (*domain.User, error) {
	r := new(register)
	if err := bind(r); err != nil {
		return nil, err
	}

	return domain.NewUserWithPassword(
		r.User.Email,
		r.User.Username,
		*r.User.Password,
	)
}

// UpdateUserToDelta converts a input serializable user to a delta for a domain user.
func UpdateUserToDelta(
	bind func(interface{}) error,
) (func(*domain.User), error) {
	r := new(register)
	if err := bind(r); err != nil {
		return nil, err
	}

	return func(u *domain.User) {
		if r.User.Email != "" {
			u.Email = r.User.Email
		}
		if r.User.Username != "" {
			u.Username = r.User.Username
		}
		if r.User.Password != nil && *r.User.Password != "" {
			u.SetPassword(*r.User.Password)
		}
		if r.User.Bio != nil {
			u.Bio = *r.User.Bio
		}
		if r.User.Image != nil {
			u.Image = *r.User.Image
		}
	}, nil
}

type login struct {
	User loginUser `json:"user"`
}
type loginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginToCredentials converts a input serializable credential set to credentials.
func LoginToCredentials(
	bind func(interface{}) error,
) (email string, password string, err error) {
	l := new(login)
	if err := bind(l); err != nil {
		return "", "", err
	}

	return l.User.Email, l.User.Password, nil
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

// CreateToArticle converts a input serializable article to a domain article for the given author.
func CreateToArticle(
	bind func(interface{}) error,
	a domain.Author,
) (*domain.Article, error) {
	ar := new(create)
	if err := bind(ar); err != nil {
		return nil, err
	}

	return domain.NewArticle(
		ar.Article.Title,
		ar.Article.Description,
		ar.Article.Body,
		a.GetEmail(),
		ar.Article.TagList...,
	)
}

// UpdateArticleToDelta converts a input article user to a delta for a domain article.
func UpdateArticleToDelta(
	bind func(interface{}) error,
) (func(*domain.Article), error) {
	ar := new(create)
	if err := bind(ar); err != nil {
		return nil, err
	}

	return func(a *domain.Article) {
		if ar.Article.Title != "" {
			a.SetTitle(ar.Article.Title)
		}
		if ar.Article.Description != "" {
			a.Description = ar.Article.Description
		}
		if ar.Article.Body != "" {
			a.Body = ar.Article.Body
		}
	}, nil
}
