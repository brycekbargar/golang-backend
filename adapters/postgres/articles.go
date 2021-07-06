package postgres

import (
	"database/sql"
	"errors"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// CreateArticle creates a new article.
func (r *implementation) CreateArticle(a *domain.Article) (*domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
INSERT INTO articles (slug, title, description, body, author_id)
	(SELECT $2, $3, $4, $5, u.id
	FROM users u WHERE u.email = $1)`)
	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrDuplicateArticle
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, domain.ErrNoAuthor
		}

		return nil, err
	}
	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		tx.Rollback()
		return nil, domain.ErrNoAuthor
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetArticleBySlug(a.Slug)
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(domain.ListCriteria) ([]*domain.AuthoredArticle, error) {
	return nil, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(s string) (*domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	found, err := getArticleBySlug(s, tx)
	if err != nil {
		return nil, err
	}

	auth, err := getUserByEmail(found.AuthorEmail, tx)
	if err != nil {
		return nil, err
	}

	return &domain.AuthoredArticle{
		Article: *found,
		Author:  auth,
	}, nil
}

func getArticleBySlug(s string, q queryer) (*domain.Article, error) {
	var found *domain.Article
	err := q.GetContext(ctx, &found, `
SELECT a.slug, a.title, a.description, a.body, a.created AS createdatutc, a.updated AS updatedatutc, u.email AS authoremail
	FROM articles a, users u 
	WHERE a.slug = $1
	AND a.author_id = u.idk`, s)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrArticleNotFound
	}
	if err != nil {
		return nil, err
	}

	return found, nil
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(string) (*domain.CommentedArticle, error) {
	return nil, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(s string, update func(*domain.Article) (*domain.Article, error)) (*domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	a, err := getArticleBySlug(s, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	a, err = update(a)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res, err := tx.ExecContext(ctx, `
UPDATE articles a FROM users u
	SET a.slug = $3, a.title = $4, a.description = $5, a.body = $6, a.updated = now() at time zone 'utc', a.author_id = u.id
	WHERE a.slug = $1
	AND u.email = a.$2
	`, s, a.AuthorEmail, a.Slug, a.Title, a.Description, a.Body)

	if err != nil {
		tx.Rollback()

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrDuplicateArticle
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, domain.ErrNoAuthor
		}

		return nil, err
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		tx.Rollback()
		return nil, domain.ErrNoAuthor
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetArticleBySlug(a.Slug)
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `
DELETE FROM articles a
	WHERE a.slug = $1`, a.Slug)
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		tx.Rollback()
		return domain.ErrArticleNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DistinctTags returns a distinct list of tags on all articles
func (r *implementation) DistinctTags() ([]string, error) {
	return nil, nil
}
