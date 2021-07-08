package postgres

import (
	"context"
	"database/sql"

	"github.com/brycekbargar/realworld-backend/domain"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// MustNewInstance creates a new instance of the postgres store with the repository interface implementations. Panics on error.
func MustNewInstance(dsn string) Migrateable {
	return &implementation{
		sqlx.MustConnect("pgx", dsn),
	}
}

type Migrateable interface {
	MustMigrate() domain.Repository
}

func (r *implementation) MustMigrate() domain.Repository {
	ctx := context.TODO()
	tx := r.db.MustBeginTx(ctx, nil)

	// TODO: Use an actual migration framework
	tx.MustExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_version (
	version varchar(40) NOT NULL,
	applied timestamp without time zone default (now() at time zone 'utc')
)`)

	err := tx.QueryRowContext(
		context.TODO(),
		"SELECT * FROM schema_version").Scan()

	if err == sql.ErrNoRows {
		tx.MustExecContext(ctx, `
INSERT INTO schema_version VALUES 
	('0.0.1.0');

CREATE TABLE users (
	id 			serial PRIMARY KEY,
	email		text NOT NULL UNIQUE,
	username	text NOT NULL UNIQUE,
	bio			text,
	image		text
);
CREATE TABLE user_passwords (
	id		integer PRIMARY KEY REFERENCES users ON DELETE CASCADE,
	hash	text NOT NULL
);

CREATE TABLE followed_users (
	follower_id integer NOT NULL REFERENCES users ON DELETE CASCADE,
	followed_id integer NOT NULL REFERENCES users ON DELETE CASCADE,
	UNIQUE (follower_id, followed_id)
);

CREATE TABLE articles (
	id 			serial PRIMARY KEY,
	slug		text NOT NULL UNIQUE,
	title		text NOT NULL,
	description	text,
	body 		text,
	created	 	timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
	updated	 	timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
	author_id 	integer NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE favorited_articles (
	user_id integer NOT NULL REFERENCES users ON DELETE CASCADE,
	article_id integer NOT NULL REFERENCES articles ON DELETE CASCADE,
	UNIQUE (user_id, article_id)
);`)
	}

	tx.Commit()
	return r
}

type implementation struct {
	db *sqlx.DB
}

type queryer interface {
	GetContext(context.Context, interface{}, string, ...interface{}) error
}
