package inmemory

import (
	"strings"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

/* type Repository interface {
	// Create creates a new user.
	Create(*User) error
	// GetUserByEmail finds a single user based on their email address.
	GetUserByEmail(string) (*User, error)
	// GetUserByUsername finds a single user based on their email address.
	GetUserByUsername(string) (*User, error)
	// UpdateUserByEmail finds a single user based on their email address,
	// then applies the provide mutations.
	UpdateUserByEmail(string, func(*User) (*User, error)) error
}
*/

// Users is an in-memory repository implementation for the usersdomain.Repository.
type Users struct {
	repository []*userdomain.User
}

// NewUsers creates a new userdomain.Repository implementation for the users domain
func NewUsers() *Users {
	return &Users{
		make([]*userdomain.User, 20),
	}
}

// Create creates a new user.
func (r *Users) Create(u *userdomain.User) error {
	_, err := r.GetUserByEmail(u.Email())
	if err == nil || err == userdomain.ErrorNotFound {
		r.repository = append(r.repository, u)
		return nil
	}

	if err != userdomain.ErrorNotFound {
		return userdomain.ErrorDuplicateValue
	}

	return err
}

// GetUserByEmail finds a single user based on their username.
func (r *Users) GetUserByEmail(e string) (*userdomain.User, error) {
	for _, ru := range r.repository {
		if strings.ToLower(ru.Email()) == strings.ToLower(e) {
			return ru, nil
		}
	}
	return nil, userdomain.ErrorNotFound
}

// GetUserByUsername finds a single user based on their email address.
func (r *Users) GetUserByUsername(un string) (*userdomain.User, error) {
	for _, ru := range r.repository {
		if strings.ToLower(ru.Username()) == strings.ToLower(un) {
			return ru, nil
		}
	}
	return nil, userdomain.ErrorNotFound
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *Users) UpdateUserByEmail(e string, uf func(*userdomain.User) (*userdomain.User, error)) error {
	u, err := r.GetUserByEmail(e)
	if err != nil {
		return err
	}

	u, err = uf(u)
	if err != nil {
		return err
	}

	// This is kind of a hack with pointers and everything being in memory
	return nil
}
