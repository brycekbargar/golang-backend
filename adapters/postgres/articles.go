package postgres

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/jackc/pgconn"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateArticle creates a new article.
func (r *implementation) CreateArticle(a *domain.Article) (*domain.AuthoredArticle, error) {
	auth, err := r.getUserByEmail(a.AuthorEmail)
	if err == domain.ErrUserNotFound {
		return nil, domain.ErrNoAuthor
	}
	if err != nil {
		return nil, err
	}

	tags, err := json.Marshal(a.TagList)
	if err != nil {
		return nil, err
	}

	res := r.db.Omit("id").Create(&Article{
		Slug:        a.Slug,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		TagList:     datatypes.JSON(tags),
		Author:      *auth,
	})

	var pgErr *pgconn.PgError
	if errors.As(res.Error, &pgErr) && pgErr.Code == "23505" {
		return nil, domain.ErrDuplicateArticle
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetArticleBySlug(a.Slug)
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(domain.ListCriteria) ([]*domain.AuthoredArticle, error) {
	return nil, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(s string) (*domain.AuthoredArticle, error) {
	found, err := r.getArticleBySlug(s)
	if err != nil {
		return nil, err
	}

	var tags []string
	err = json.Unmarshal(found.TagList, &tags)
	if err != nil {
		return nil, err
	}

	return &domain.AuthoredArticle{
		Article: domain.Article{
			Slug:         found.Slug,
			Title:        found.Title,
			Description:  found.Description,
			Body:         found.Body,
			TagList:      tags,
			CreatedAtUTC: found.CreatedAt,
			UpdatedAtUTC: found.UpdatedAt,
			AuthorEmail:  found.Author.GetEmail(),
		},
		Author: found.Author,
	}, nil
}

func (r *implementation) getArticleBySlug(s string) (*Article, error) {
	var found Article
	res := r.db.
		Preload(clause.Associations).
		First(&found, "slug = ?", s)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ErrArticleNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return &found, nil
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(string) (*domain.CommentedArticle, error) {
	return nil, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(s string, update func(*domain.Article) (*domain.Article, error)) (*domain.AuthoredArticle, error) {
	a, err := r.GetArticleBySlug(s)
	if err != nil {
		return nil, err
	}

	article, err := update(&a.Article)
	if err != nil {
		return nil, err
	}

	found, err := r.getArticleBySlug(s)
	if err != nil {
		return nil, err
	}

	tags, err := json.Marshal(article.TagList)
	if err != nil {
		return nil, err
	}

	found.Slug = article.Slug
	found.Title = article.Title
	found.Description = article.Description
	found.Body = article.Body
	found.TagList = datatypes.JSON(tags)
	res := r.db.Save(found)

	var pgErr *pgconn.PgError
	if errors.As(res.Error, &pgErr) && pgErr.Code == "23505" {
		return nil, domain.ErrDuplicateArticle
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetArticleBySlug(article.Slug)
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(string, func(*domain.CommentedArticle) (*domain.CommentedArticle, error)) (*domain.Comment, error) {
	return nil, nil
}

// DeleteArticle deletes the article if it exists.
func (r *implementation) DeleteArticle(a *domain.Article) error {
	if a == nil {
		return nil
	}

	found, err := r.getArticleBySlug(a.Slug)
	if err != nil {
		return err
	}

	res := r.db.Delete(&found)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// DistinctTags returns a distinct list of tags on all articles
func (r *implementation) DistinctTags() ([]string, error) {
	var articles []*Article
	res := r.db.Select("tag_list").Find(&articles)
	if res.Error != nil {
		return nil, res.Error
	}

	tm := make(map[string]interface{})
	for _, a := range articles {
		var tags []string
		err := json.Unmarshal(a.TagList, &tags)
		if err != nil {
			return nil, err
		}

		for _, t := range tags {
			tm[strings.ToLower(t)] = nil
		}
	}

	tags := make([]string, 0, len(tm))
	for t := range tm {
		tags = append(tags, t)
	}

	return tags, nil
}
