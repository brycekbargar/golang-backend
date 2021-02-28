package inmemory

import (
	"strings"
	"time"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

// Create creates a new article.
func (r *implementation) CreateArticle(a *articledomain.Article) (*articledomain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for s := range r.articles {
		if s == strings.ToLower(a.Slug()) {
			return nil, articledomain.ErrDuplicateValue
		}
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail())]; !ok {
		return nil, articledomain.ErrNoAuthor
	}

	now := time.Now().UTC()
	r.articles[strings.ToLower(a.Slug())] = articleRecord{
		a.Slug(),
		a.Title(),
		a.Description(),
		a.Body(),
		a.Tags(),
		now,
		now,
		a.AuthorEmail(),
		make([]*commentRecord, 0),
		map[string]interface{}{},
	}
	return r.GetArticleBySlug(a.Slug())
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(articledomain.ListCriteria) ([]*articledomain.AuthoredArticle, error) {
	return nil, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(s string) (*articledomain.AuthoredArticle, error) {
	if a, ok := r.articles[strings.ToLower(s)]; ok {

		aa, err := r.getUserByEmail(a.author, false)
		if err == userdomain.ErrNotFound {
			return nil, articledomain.ErrNoAuthor
		}
		if err != nil {
			return nil, err
		}

		return articledomain.ExistingArticle(
			a.slug,
			a.title,
			a.description,
			a.body,
			a.tagList,
			a.createdAtUTC,
			a.updatedAtUTC,
			aa,
			make([]*articledomain.Comment, 0),
			make([]string, 0),
		)
	}

	return nil, articledomain.ErrNotFound
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(string) (*articledomain.CommentedArticle, error) {
	return nil, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(string, func(*articledomain.Article) (*articledomain.Article, error)) (*articledomain.AuthoredArticle, error) {
	return nil, nil
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(string, func(*articledomain.CommentedArticle) (*articledomain.CommentedArticle, error)) (*articledomain.Comment, error) {
	return nil, nil
}

// DeleteArticleBySlug deletes the article with the provide slug if it exists.
func (r *implementation) DeleteArticle(*articledomain.Article) error {
	return nil
}

// DistinctTags returns a distinct list of tags on articles
func (r *implementation) DistinctTags() ([]string, error) {
	return nil, nil
}
