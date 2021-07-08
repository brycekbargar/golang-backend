package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

var ctx = context.TODO()

// CreateUser creates a new user.
func (r *implementation) CreateUser(u *domain.User) (*domain.User, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	var id int
	err = tx.QueryRow(ctx, `
INSERT INTO users (email, username, bio, image) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`,
		u.Email, u.Username, u.Bio, u.Image).Scan(&id)

	if err != nil {
		tx.Rollback(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrDuplicateUser
		}

		return nil, err
	}

	// TODO: Use salts and pg stuff instead of the bcrypt server side implementation
	_, err = tx.Exec(ctx, `
INSERT INTO user_passwords (id, hash) 
	VALUES ($1, $2)
`, id, u.Password)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return getUserByEmail(u.Email, r.db)
}

// GetUserByEmail finds a single user based on their email address.
func (r *implementation) GetUserByEmail(em string) (*domain.Fanboy, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Commit(ctx)

	found, err := getUserByEmail(em, tx)
	if err != nil {
		return nil, err
	}

	var follows []string
	err = pgxscan.Select(ctx, tx, &follows, `
SELECT f.email
	FROM users u, followed_users fu, users f
	WHERE u.email = $1
	AND u.id = fu.follower_id
	AND f.id = fu.followed_id
`, em)
	if err != nil {
		return nil, err
	}

	following := make(map[string]interface{}, len(follows))
	for _, u := range follows {
		following[strings.ToLower(u)] = nil
	}

	var favors []string
	err = pgxscan.Select(ctx, tx, &favors, `
SELECT a.slug
	FROM users u, favorited_articles fa, articles a
	WHERE u.email = $1
	AND u.id = fa.user_id
	AND a.id = fa.article_id
`, em)
	if err != nil {
		return nil, err
	}

	favorites := make(map[string]interface{}, len(favors))
	for _, a := range favors {
		favorites[strings.ToLower(a)] = nil
	}

	return &domain.Fanboy{
		User:      *found,
		Following: following,
		Favorites: favorites,
	}, nil
}

func getUserByEmail(em string, q pgxscan.Querier) (*domain.User, error) {
	found := new(domain.User)
	err := pgxscan.Get(ctx, q, found, `
SELECT u.email, u.username, u.bio, u.image, p.hash as password
	FROM users u, user_passwords p
	WHERE u.email = $1 
	AND u.id = p.id`, em)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return found, nil
}

// GetAuthorByEmail finds a single author based on their email address or nil if they don't exist.
func (r *implementation) GetAuthorByEmail(em string) domain.Author {
	auth, err := getUserByEmail(em, r.db)
	if err != nil {
		return nil
	}
	return auth
}

// GetUserByUsername finds a single user based on their username.
func (r *implementation) GetUserByUsername(un string) (*domain.User, error) {
	found := new(domain.User)
	err := pgxscan.Get(ctx, r.db, found, `
SELECT u.email, u.username, u.bio, u.image, p.hash as password
	FROM users u, user_passwords p
	WHERE u.username = $1 
	AND u.id = p.id`, un)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return found, nil
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *implementation) UpdateUserByEmail(em string, update func(*domain.User) (*domain.User, error)) (*domain.User, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	u, err := getUserByEmail(em, tx)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	u, err = update(u)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	var id int
	err = tx.QueryRow(ctx, `
UPDATE users 
	SET email = $2, username = $3, bio = $4, image = $5
	WHERE email = $1
	RETURNING id`,
		em, u.Email, u.Username, u.Bio, u.Image).Scan(&id)

	if err != nil {
		tx.Rollback(ctx)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, domain.ErrDuplicateUser
		}

		return nil, err
	}

	// TODO: Use salts and pg stuff instead of the bcrypt server side implementation
	_, err = tx.Exec(ctx, `
UPDATE user_passwords
	SET hash = $2
	WHERE id = $1
`, id, u.Password)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return getUserByEmail(u.Email, r.db)
}

// UpdateFanboyByEmail finds a single user based on their email address,
// then applies the provide mutations (probably to the follower list).
func (r *implementation) UpdateFanboyByEmail(em string, update func(*domain.Fanboy) (*domain.Fanboy, error)) error {
	return nil
}
