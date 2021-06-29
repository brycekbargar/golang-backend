package postgres

import "github.com/brycekbargar/realworld-backend/domain"

// CreateUser creates a new user.
func (r *implementation) CreateUser(*domain.User) (*domain.User, error) {
	return nil, nil
}

// GetUserByEmail finds a single user based on their email address.
func (r *implementation) GetUserByEmail(string) (*domain.Fanboy, error) {
	return nil, nil
}

// GetAuthorByEmail finds a single author based on their email address or nil if they don't exist.
func (r *implementation) GetAuthorByEmail(string) domain.Author {
	return nil
}

// GetUserByUsername finds a single user based on their username.
func (r *implementation) GetUserByUsername(string) (*domain.User, error) {
	return nil, nil
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *implementation) UpdateUserByEmail(string, func(*domain.User) (*domain.User, error)) (*domain.User, error) {
	return nil, nil
}

// UpdateFanboyByEmail finds a single user based on their email address,
// then applies the provide mutations (probably to the follower list).
func (r *implementation) UpdateFanboyByEmail(string, func(*domain.Fanboy) (*domain.Fanboy, error)) error {
	return nil
}
