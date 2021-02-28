// Package inmemory is an in-memory implementation of the adapters layer.
package inmemory

import (
	"strings"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

// Create creates a new user.
func (r *implementation) CreateUser(u *userdomain.User) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for e, v := range r.users {
		if e == strings.ToLower(u.Email()) ||
			strings.ToLower(v.username) == strings.ToLower(u.Username()) {
			return nil, userdomain.ErrDuplicateValue
		}
	}

	r.users[strings.ToLower(u.Email())] = userRecord{
		u.Email(),
		u.Username(),
		u.Bio(),
		u.Image(),
		strings.Join(u.FollowingEmails(), ","),
		u.Password(),
	}
	return r.GetUserByEmail(u.Email())
}

// GetUserByEmail finds a single user based on their username.
func (r *implementation) GetUserByEmail(e string) (*userdomain.User, error) {
	return r.getUserByEmail(e, true)
}

// getUserByEmail finds a single user based on their username (without infinitely recursing)
func (r *implementation) getUserByEmail(e string, recurse bool) (*userdomain.User, error) {
	if u, ok := r.users[strings.ToLower(e)]; ok {

		var uf []*userdomain.User
		if recurse && len(u.following) > 0 {
			f := strings.Split(u.following, ",")
			uf = make([]*userdomain.User, 0, len(f))

			for _, fe := range f {
				fu, err := r.getUserByEmail(fe, false)

				if err != nil {
					return nil, userdomain.ErrNotFound
				}
				uf = append(uf, fu)
			}
		}

		return userdomain.ExistingUser(
			u.email,
			u.username,
			u.bio,
			u.image,
			uf,
			u.password)
	}

	return nil, userdomain.ErrNotFound
}

// GetUserByUsername finds a single user based on their email address.
func (r *implementation) GetUserByUsername(un string) (*userdomain.User, error) {
	for k, v := range r.users {
		if strings.ToLower(v.username) == strings.ToLower(un) {
			return r.GetUserByEmail(k)
		}
	}

	return nil, userdomain.ErrNotFound
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *implementation) UpdateUserByEmail(e string, uf func(*userdomain.User) (*userdomain.User, error)) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, err := r.GetUserByEmail(e)
	if err != nil {
		return nil, err
	}

	u, err = uf(u)
	if err != nil {
		return nil, err
	}

	for _, fe := range u.FollowingEmails() {
		if _, err := r.getUserByEmail(fe, false); err != nil {
			return nil, err
		}
	}

	ru := r.users[strings.ToLower(e)]
	delete(r.users, strings.ToLower(e))

	for e, v := range r.users {
		if e == strings.ToLower(u.Email()) ||
			strings.ToLower(v.username) == strings.ToLower(u.Username()) {

			// Add the deleted user back if they've become a duplicate
			r.users[strings.ToLower(e)] = ru
			return nil, userdomain.ErrDuplicateValue
		}
	}

	r.users[strings.ToLower(u.Email())] = userRecord{
		u.Email(),
		u.Username(),
		u.Bio(),
		u.Image(),
		strings.Join(u.FollowingEmails(), ","),
		u.Password(),
	}

	return r.GetUserByEmail(u.Email())
}
