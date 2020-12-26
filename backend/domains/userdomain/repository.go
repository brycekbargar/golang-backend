package userdomain

import "errors"

// ErrorNotFound indicates the requested user was not found.
var ErrorNotFound = errors.New("user not found")

// ErrorDuplicateValue indicates the requested user could not be created because they already exist.
var ErrorDuplicateValue = errors.New("created user has a duplicate username or email address")

// Repository allows performing abstracted I/O operations on users.
type Repository interface {
	// Create creates a new user.
	Create(*User) error
	// GetUserByEmail finds a single user based on their email address.
	GetUserByEmail(string) (*User, error)
	// GetUserByUsername finds a single user based on their username.
	GetUserByUsername(string) (*User, error)
	// UpdateUserByEmail finds a single user based on their email address,
	// then applies the provide mutations.
	UpdateUserByEmail(string, func(*User) (*User, error)) error
}
