package postgres

import (
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MustNewInstance creates a new instance of the postgres store with the repository interface implementations. Panics on error.
func MustNewInstance(dsn string) domain.Repository {
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

	return &implementation{
		db,
	}
}

type implementation struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Email     string `gorm:"index:,unique"`
	Username  string `gorm:"index:,unique"`
	Bio       string
	Image     string
	Password  Password
	Following []*User    `gorm:"many2many:user_following"`
	Favorites []*Article `gorm:"many2many:user_favorites"`
}

func (u User) GetEmail() string {
	return u.Email
}
func (u User) GetUsername() string {
	return u.Username
}
func (u User) GetBio() string {
	return u.Bio
}
func (u User) GetImage() string {
	return u.Image
}

type Password struct {
	ID     uint `gorm:"primarykey"`
	UserID uint
	Value  []byte
}

type Article struct {
	gorm.Model
	Slug        string `gorm:"index:,unique"`
	Title       string
	Description string
	Body        string
	TagList     datatypes.JSON
	Author      User `gorm:"foreignkey:ID"`
	Comments    []*Comment
}

type Comment struct {
	gorm.Model
	ArticleID uint
	Body      string
	Author    User `gorm:"foreignkey:ID"`
}
