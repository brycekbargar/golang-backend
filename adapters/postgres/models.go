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
		sqlx.MustConnect("postgres", dsn),
	}
}

type Migrateable interface {
	MustMigrate() domain.Repository
}

func (r *implementation) MustMigrate() domain.Repository {
	ctx := context.TODO()
	tx := r.db.MustBeginTx(ctx, nil)

	// TODO: Use an actual migration framework
	err := tx.QueryRowContext(
		context.TODO(),
		"SELECT * FROM schema_version").Scan()
	if err == sql.ErrNoRows {
		_, err = tx.ExecContext(ctx, `
CREATE TABLE schema_version (
	version varchar(40) NOT NULL,
	applied timestamp without time zone default (now() at time zone 'utc')
)
INSERT INTO schema_version VALUES 
	('0.0.1.0')

-- TODO: Migrate the thing
`)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
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
