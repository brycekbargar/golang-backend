package postgres

import (
	"errors"
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

// CreateArticle creates a new article.
func (r *implementation) CreateArticle(a *domain.Article) (*domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(ctx, `
INSERT INTO articles (slug, title, description, body, tags, author_id)
	(SELECT $2, $3, $4, $5, $6, u.id
	FROM users u WHERE u.email = $1)`,
		a.AuthorEmail, a.Slug, a.Title, a.Description, a.Body, a.TagList)
	if err != nil {
		tx.Rollback(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrDuplicateArticle
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, domain.ErrNoAuthor
		}

		return nil, err
	}
	if res.RowsAffected() != 1 {
		tx.Rollback(ctx)
		return nil, domain.ErrNoAuthor
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetArticleBySlug(a.Slug)
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(lc domain.ListCriteria) ([]domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	var slugs []string
	err = pgxscan.Select(ctx, tx, &slugs, `
WITH faves AS (
	SELECT a.id, fu.email
	FROM articles a
	LEFT JOIN favorited_articles fa ON
		fa.article_id = a.id
	LEFT JOIN users fu ON
		fa.user_id = fu.id
), slugs AS (
	SELECT DISTINCT
		a.slug, a.updated
	FROM articles a
	INNER JOIN users u ON
		a.author_id = u.id
	LEFT JOIN faves f ON
		a.id = f.id
	WHERE (length($3) = 0 OR $3 = ANY(a.tags))
	AND ($4::text[] IS NULL OR array_length($4::text[], 1) = 0 OR u.email = ANY($4))
	AND (length($5) = 0 OR f.email = $5)
	ORDER BY a.updated DESC
)
SELECT slug FROM slugs
LIMIT $1 OFFSET $2
`,
		lc.Limit, lc.Offset, lc.Tag, lc.AuthorEmails, lc.FavoritedByUserEmail)
	if err != nil {
		return nil, err
	}

	latest, err := getArticleBySlug(tx, slugs...)
	if err == domain.ErrArticleNotFound {
		return make([]domain.AuthoredArticle, 0), nil
	}
	if err != nil {
		return nil, err
	}

	return latest, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(s string) (*domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	found, err := getArticleBySlug(tx, s)
	if err != nil {
		return nil, err
	}

	return &found[0], nil
}

func getArticleBySlug(q pgxscan.Querier, s ...string) ([]domain.AuthoredArticle, error) {
	var found []domain.AuthoredArticle
	err := pgxscan.Select(ctx, q, &found, `
WITH faves AS (
	SELECT 
		a.id
		,COUNT(fa.article_id) AS count
	FROM articles a
	LEFT JOIN favorited_articles fa ON
		a.id = fa.article_id
	GROUP BY a.id
)
SELECT
	a.slug
	,a.title
	,a.description
	,a.body
	,a.tags as tag_list
	,a.created AS created_at_utc
	,a.updated AS updated_at_utc
	,u.email AS author_email
	,f.count AS favorite_count
FROM 
	articles a
	,users u
	,faves f
WHERE a.slug = ANY($1)
AND a.author_id = u.id
AND a.id = f.id
ORDER BY a.updated DESC
`,
		s)
	if errors.Is(err, pgx.ErrNoRows) ||
		(err == nil && len(found) == 0) {
		return nil, domain.ErrArticleNotFound
	}
	if err != nil {
		return nil, err
	}

	// TODO: Cache the authors? Somehow read map them using scany?
	for i := range found {
		a := &found[i]
		a.Author, _ = getUserByEmail(q, a.AuthorEmail)
	}

	return found, nil
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(s string) (*domain.CommentedArticle, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	found, err := getArticleBySlug(tx, s)
	if err != nil {
		return nil, err
	}

	var comments []domain.Comment
	err = pgxscan.Select(ctx, tx, &comments, `
SELECT c.id, c.body, c.created as created_at_utc, u.email as author_email
	FROM articles a, article_comments c, users u
	WHERE a.slug = $1
	AND a.id = c.article_id
	AND u.id = c.author_id
`, s)
	if err != nil {
		return nil, err
	}

	return &domain.CommentedArticle{
		Article:  found[0].Article,
		Comments: comments,
	}, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(s string, update func(*domain.Article) (*domain.Article, error)) (*domain.AuthoredArticle, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	as, err := getArticleBySlug(tx, s)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	a, err := update(&as[0].Article)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	res, err := tx.Exec(ctx, `
UPDATE articles
	SET slug = $3, title = $4, description = $5, body = $6, updated = now() at time zone 'utc', author_id = u.id
 	FROM users u
	WHERE slug = $1
	AND u.email = $2
	`, s, a.AuthorEmail, a.Slug, a.Title, a.Description, a.Body)

	if err != nil {
		tx.Rollback(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrDuplicateArticle
		}
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return nil, domain.ErrNoAuthor
		}

		return nil, err
	}

	if res.RowsAffected() != 1 {
		tx.Rollback(ctx)
		return nil, domain.ErrNoAuthor
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetArticleBySlug(a.Slug)
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(s string, update func(*domain.CommentedArticle) (*domain.CommentedArticle, error)) (*domain.Comment, error) {
	a, err := r.GetCommentsBySlug(s)
	if err != nil {
		return nil, err
	}

	a, err = update(a)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	var new *domain.Comment
	ids := make([]int, 0, len(a.Comments))
	for _, c := range a.Comments {
		if c.ID <= 0 {
			new = &c
		} else {
			ids = append(ids, c.ID)
		}
	}

	_, err = tx.Exec(ctx, `
DELETE FROM article_comments
	USING articles a 
	WHERE a.slug = $1
	AND a.id = article_id
	AND article_comments.id <> ANY($2)
`,
		s, ids)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	if new != nil {
		var id int
		var created time.Time
		err = tx.QueryRow(ctx, `
INSERT INTO article_comments (article_id, author_id, body)
	(SELECT a.id, u.id, $3
		FROM articles a, users u
		WHERE a.slug = $1
		AND u.email = $2)
	RETURNING id, created`,
			a.Slug, new.AuthorEmail, new.Body).Scan(&id, &created)

		if err != nil {
			tx.Rollback(ctx)

			return nil, err
		}

		new.ID = id
		new.CreatedAtUTC = created

		if err = tx.Commit(ctx); err != nil {
			return nil, err
		}
		return new, nil
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return nil, nil
}

// DeleteArticle deletes the article if it exists.
func (r *implementation) DeleteArticle(a *domain.Article) error {
	if a == nil {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	res, err := tx.Exec(ctx, `
DELETE FROM articles a
	WHERE a.slug = $1`, a.Slug)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	if res.RowsAffected() != 1 {
		tx.Rollback(ctx)
		return domain.ErrArticleNotFound
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// DistinctTags returns a distinct list of tags on all articles
func (r *implementation) DistinctTags() ([]string, error) {
	var tags []string
	err := pgxscan.Select(ctx, r.db, &tags, `
SELECT DISTINCT UNNEST(tags) FROM articles
`)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
