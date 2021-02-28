package articledomain

import "errors"

// ErrNoAuthor indicates when the author of an Article can't be found.
var ErrNoAuthor = errors.New("author not found")

// ErrNotFound indicates the requested article was not found.
var ErrNotFound = errors.New("article not found")

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
	// Create creates a new article.
	Create(*Article) (*AuthoredArticle, error)
	// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
	LatestArticlesByCriteria(ListCriteria) ([]*AuthoredArticle, error)
	// GetArticleBySlug gets a single article with the given slug.
	GetArticleBySlug(string) (*AuthoredArticle, error)
	// GetCommentsBySlug gets a single article and its comments with the given slug.
	GetCommentsBySlug(string) (*CommentedArticle, error)
	// UpdateArticleBySlug finds a single article based on its slug
	// then applies the provide mutations.
	UpdateArticleBySlug(string, func(*Article) (*Article, error)) (*AuthoredArticle, error)
	// UpdateCommentsBySlug finds a single article based on its slug
	// then applies the provide mutations to its comments.
	UpdateCommentsBySlug(string, func(*CommentedArticle) (*CommentedArticle, error)) (*Comment, error)
	// DeleteArticleBySlug deletes the article with the provide slug if it exists.
	Delete(*Article) error
	// DistinctTags returns a distinct list of tags on articles
	DistinctTags() ([]string, error)
}
