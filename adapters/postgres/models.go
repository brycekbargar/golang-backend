package postgres

import (
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewInstance creates a new instance of the postgres store with the repository interface implementations. Panics on error.
func NewInstance(dsn string) domain.Repository {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&User{},
		&Password{},
		&Article{},
		&Comment{},
	)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return implementation{
		db,
	}
}

type implementation struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Email      string `gorm:"uniqueIndex"`
	Username   string `gorm:"uniqueIndex"`
	Bio        string
	Image      string
	PasswordID uint
	Password   Password
	Following  []User
	Favorites  []Article
}

type Password struct {
	ID    uint
	Value []byte
}

type Article struct {
	gorm.Model
	Slug        string `gorm:"uniqueIndex"`
	Title       string
	Description string
	Body        string
	TagList     datatypes.JSON
	Author      User
	Comments    []Comment
}

type Comment struct {
	gorm.Model
	Body   string
	Author User
}
