package articledomain

import "errors"

// ErrNotFound indicates the requested article was not found.
var ErrNotFound = errors.New("user not found")

// ErrDuplicateValue indicates the requested article could not be created because another article has the same slug.
var ErrDuplicateValue = errors.New("article has a duplicate slug")

// ListCriteria is the set of optional parameters to page/filter the Articles.
type ListCriteria struct {
	Tag                  string
	AuthorEmails         []string
	FavoritedByUserEmail string
	Limit                int
	Offset               int
}

// Repository allows performing abstracted I/O operations on articles.
type Repository interface {
	// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
	LatestArticlesByCriteria(ListCriteria) ([]*AuthoredArticle, error)
	// GetArticleBySlug gets a single article with the given slug.
	GetArticleBySlug(string) (*AuthoredArticle, error)
}
