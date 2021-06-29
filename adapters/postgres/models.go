package postgres

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

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
