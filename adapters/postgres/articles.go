package postgres

import "github.com/brycekbargar/realworld-backend/domain"

// CreateArticle creates a new article.
func (r *implementation) CreateArticle(*domain.Article) (*domain.AuthoredArticle, error) {
	return nil, nil
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(domain.ListCriteria) ([]*domain.AuthoredArticle, error) {
	return nil, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(string) (*domain.AuthoredArticle, error) {
	return nil, nil
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(string) (*domain.CommentedArticle, error) {
	return nil, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(string, func(*domain.Article) (*domain.Article, error)) (*domain.AuthoredArticle, error) {
	return nil, nil
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(string, func(*domain.CommentedArticle) (*domain.CommentedArticle, error)) (*domain.Comment, error) {
	return nil, nil
}

// DeleteArticle deletes the article if it exists.
func (r *implementation) DeleteArticle(*domain.Article) error {
	return nil
}

// DistinctTags returns a distinct list of tags on all articles
func (r *implementation) DistinctTags() ([]string, error) {
	return nil, nil
}
