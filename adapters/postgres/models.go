package postgres

import (
	"context"

	"github.com/brycekbargar/realworld-backend/domain"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// MustNewInstance creates a new instance of the postgres store with the repository interface implementations. Panics on error.
func MustNewInstance(dsn string) Migrateable {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		panic(err)
	}
	if err = pool.Ping(ctx); err != nil {
		panic(err)
	}

	return &implementation{pool}
}

type Migrateable interface {
	MustMigrate() domain.Repository
}

func (r *implementation) MustMigrate() domain.Repository {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(err)
	}

	// TODO: Use an actual migration framework
	_, err = tx.Exec(ctx, `
CREATE TABLE IF NOT EXISTS schema_version (
	version varchar(40) NOT NULL,
	applied timestamp without time zone default (now() at time zone 'utc')
)`)
	if err != nil {
		panic(err)
	}

	rows, err := tx.Query(
		ctx,
		"SELECT * FROM schema_version")
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	if !rows.Next() {
		_, err = tx.Exec(ctx, `
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
	tags 		text[],
	created	 	timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
	updated	 	timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
	author_id 	integer NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE favorited_articles (
	user_id 	integer NOT NULL REFERENCES users ON DELETE CASCADE,
	article_id 	integer NOT NULL REFERENCES articles ON DELETE CASCADE,
	UNIQUE (user_id, article_id)
);

CREATE TABLE article_comments (
	id 			serial PRIMARY KEY,
	article_id 	integer NOT NULL REFERENCES articles ON DELETE CASCADE,
	author_id 	integer NOT NULL REFERENCES users ON DELETE CASCADE,
	body		text,
	created	 	timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc')
);
`)
	}
	if err != nil {
		panic(err)
	}

	tx.Commit(ctx)
	return r
}

type implementation struct {
	db *pgxpool.Pool
}

type queryer interface {
	GetContext(context.Context, interface{}, string, ...interface{}) error
}
