package inmemory

import (
	"github.com/brycekbargar/realworld-backend/domains/articledomain"
)

// Create creates a new article.
func (r *articles) Create(*articledomain.Article) (*articledomain.AuthoredArticle, error)

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *articles) LatestArticlesByCriteria(articledomain.ListCriteria) ([]*articledomain.AuthoredArticle, error)

// GetArticleBySlug gets a single article with the given slug.
func (r *articles) GetArticleBySlug(string) (*articledomain.AuthoredArticle, error)

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *articles) GetCommentsBySlug(string) (*articledomain.CommentedArticle, error)

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *articles) UpdateArticleBySlug(string, func(*articledomain.Article) (*articledomain.Article, error)) (*articledomain.AuthoredArticle, error)

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *articles) UpdateCommentsBySlug(string, func(*articledomain.CommentedArticle) (*articledomain.CommentedArticle, error)) (*articledomain.Comment, error)

// DeleteArticleBySlug deletes the article with the provide slug if it exists.
func (r *articles) Delete(*articledomain.Article) error

// DistinctTags returns a distinct list of tags on articles
func (r *articles) DistinctTags() ([]string, error)
