package inmemory

import (
	"strings"
	"time"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
)

// Create creates a new article.
func (r *implementation) CreateArticle(a *articledomain.Article) (*articledomain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for s := range r.articles {
		if s == strings.ToLower(a.Slug) {
			return nil, articledomain.ErrDuplicateValue
		}
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail)]; !ok {
		return nil, articledomain.ErrNoAuthor
	}

	now := time.Now().UTC()
	r.articles[strings.ToLower(a.Slug)] = &articleRecord{
		a.Slug,
		a.Title,
		a.Description,
		a.Body,
		a.TagList,
		now,
		now,
		a.AuthorEmail,
		make([]*commentRecord, 0),
		map[string]interface{}{},
	}
	return r.GetArticleBySlug(a.Slug)
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(articledomain.ListCriteria) ([]*articledomain.AuthoredArticle, error) {
	return nil, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(s string) (*articledomain.AuthoredArticle, error) {
	if a, ok := r.articles[strings.ToLower(s)]; ok {

		aa, ok := r.users[strings.ToLower(a.author)]
		if !ok {
			return nil, articledomain.ErrNoAuthor
		}

		return &articledomain.AuthoredArticle{
			Article: articledomain.Article{
				Slug:         a.slug,
				Title:        a.title,
				Description:  a.description,
				Body:         a.body,
				TagList:      a.tagList,
				CreatedAtUTC: a.createdAtUTC,
				UpdatedAtUTC: a.updatedAtUTC,
				AuthorEmail:  a.author,
				FavoritedBy:  a.favoritedBy,
			},
			Author: aa,
		}, nil
	}

	return nil, articledomain.ErrNotFound
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(string) (*articledomain.CommentedArticle, error) {
	return nil, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(s string, update func(*articledomain.Article) (*articledomain.Article, error)) (*articledomain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := r.GetArticleBySlug(s)
	if err != nil {
		return nil, err
	}

	a, err := update(&f.Article)
	if err != nil {
		return nil, err
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail)]; !ok {
		return nil, articledomain.ErrNoAuthor
	}

	removed := r.articles[strings.ToLower(s)]
	delete(r.articles, strings.ToLower(s))

	for s := range r.articles {
		if s == strings.ToLower(a.Slug) {

			// Add the deleted article back if they've become a duplicate
			r.articles[strings.ToLower(removed.slug)] = removed
			return nil, articledomain.ErrDuplicateValue
		}
	}

	now := time.Now().UTC()
	r.articles[strings.ToLower(a.Slug)] = &articleRecord{
		a.Slug,
		a.Title,
		a.Description,
		a.Body,
		a.TagList,
		a.CreatedAtUTC,
		now,
		a.AuthorEmail,
		make([]*commentRecord, 0),
		map[string]interface{}{},
	}

	return r.GetArticleBySlug(a.Slug)
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(string, func(*articledomain.CommentedArticle) (*articledomain.CommentedArticle, error)) (*articledomain.Comment, error) {
	return nil, nil
}

// DeleteArticleBySlug deletes the article with the provide slug if it exists.
func (r *implementation) DeleteArticle(a *articledomain.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if a == nil {
		return nil
	}

	delete(r.articles, strings.ToLower(a.Slug))
	return nil
}

// DistinctTags returns a distinct list of tags on articles
func (r *implementation) DistinctTags() ([]string, error) {
	return nil, nil
}
