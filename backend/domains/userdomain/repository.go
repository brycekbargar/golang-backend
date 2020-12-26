package userdomain

import "errors"

// ErrorNotFound indicates the requested user was not found.
var ErrorNotFound = errors.New("user not found")

// ErrorDuplicateValue indicates the requested user could not be created because they already exist.
var ErrorDuplicateValue = errors.New("created user has a duplicate username or email address")

// Repository allows performing abstracted I/O operations on users.
type Repository interface {
	// Create creates a new user.
	Create(*UserWithPassword) error
	// GetUserByEmail queries for a user based on their email address.
	GetUserByEmail(string) (*User, error)
	// GetLoginUserByEmail queries for a user based on their email address.
	// The returned user allows for password checking and updating.
	GetLoginUserByEmail(string) (*UserWithPassword, error)
}
