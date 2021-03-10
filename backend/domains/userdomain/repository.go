package userdomain

import "errors"

// ErrNotFound indicates the requested user was not found.
var ErrNotFound = errors.New("user not found")

// ErrDuplicateValue indicates the requested user could not be created because they already exist.
var ErrDuplicateValue = errors.New("user has a duplicate username or email address")

// Repository allows performing abstracted I/O operations on users.
type Repository interface {
	// Create creates a new user.
	CreateUser(*User) (*User, error)
	// GetUserByEmail finds a single user based on their email address.
	GetUserByEmail(string) (*Fanboy, error)
	// GetUserByUsername finds a single user based on their username.
	GetUserByUsername(string) (*User, error)
	// UpdateUserByEmail finds a single user based on their email address,
	// then applies the provide mutations.
	UpdateUserByEmail(string, func(*User) (*User, error)) (*User, error)
	// UpdateFanboyByEmail finds a single user based on their email address,
	// then applies the provide mutations (probably to the follower list).
	UpdateFanboyByEmail(string, func(*Fanboy) (*Fanboy, error)) error
}
