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
		if e == strings.ToLower(u.Email) ||
			strings.ToLower(v.username) == strings.ToLower(u.Username) {
			return nil, userdomain.ErrDuplicateValue
		}
	}

	r.users[strings.ToLower(u.Email)] = &userRecord{
		u.Email,
		u.Username,
		u.Bio,
		u.Image,
		"",
		u.Password,
	}

	f, err := r.GetUserByEmail(u.Email)
	return &f.User, err
}

// GetUserByEmail finds a single user based on their username.
func (r *implementation) GetUserByEmail(e string) (*userdomain.Fanboy, error) {
	if u, ok := r.users[strings.ToLower(e)]; ok {

		emails := strings.Split(u.following, ",")
		follows := make(map[string]interface{}, len(emails))
		for _, em := range emails {
			follows[em] = nil
		}

		return &userdomain.Fanboy{
			User: userdomain.User{
				Email:    u.email,
				Username: u.username,
				Bio:      u.bio,
				Image:    u.image,
				Password: u.password,
			},
			Following: follows,
		}, nil
	}

	return nil, userdomain.ErrNotFound
}

// GetUserByUsername finds a single user based on their username.
func (r *implementation) GetUserByUsername(un string) (*userdomain.User, error) {
	for k, v := range r.users {
		if strings.ToLower(v.username) == strings.ToLower(un) {
			f, err := r.GetUserByEmail(k)
			return &f.User, err
		}
	}

	return nil, userdomain.ErrNotFound
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *implementation) UpdateUserByEmail(e string, update func(*userdomain.User) (*userdomain.User, error)) (*userdomain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := r.GetUserByEmail(e)
	if err != nil {
		return nil, err
	}
	prevEm := strings.ToLower(f.Email)

	u, err := update(&f.User)
	if err != nil {
		return nil, err
	}

	removed := r.users[strings.ToLower(e)]
	delete(r.users, strings.ToLower(e))

	for e, v := range r.users {
		if e == strings.ToLower(u.Email) ||
			strings.ToLower(v.username) == strings.ToLower(u.Username) {

			// Add the deleted user back if they've become a duplicate
			r.users[strings.ToLower(e)] = removed
			return nil, userdomain.ErrDuplicateValue
		}
	}

	if strings.ToLower(u.Email) != prevEm {
		for _, v := range r.users {
			// Make sure users following this one get an updated key
			v.following = strings.ReplaceAll(v.following, prevEm, u.Email)
		}
	}

	follows := make([]string, len(f.Following))
	for k := range f.Following {
		follows = append(follows, k)
	}

	r.users[strings.ToLower(u.Email)] = &userRecord{
		u.Email,
		u.Username,
		u.Bio,
		u.Image,
		strings.ToLower(strings.Join(follows, ",")),
		u.Password,
	}

	f, err = r.GetUserByEmail(u.Email)
	return &f.User, err
}

func (r *implementation) UpdateFanboyByEmail(e string, update func(*userdomain.Fanboy) (*userdomain.Fanboy, error)) error {
	var uf *userdomain.Fanboy
	err := func() error {
		r.mu.Lock()
		defer r.mu.Unlock()

		f, err := r.GetUserByEmail(e)
		if err != nil {
			return err
		}

		uf, err = update(f)
		if err != nil {
			return err
		}

		follows := make([]string, len(uf.Following))
		for k := range uf.Following {
			if k != "" {
				follows = append(follows, k)
			}
		}

		fr, ok := r.users[strings.ToLower(e)]
		if !ok {
			return userdomain.ErrNotFound
		}
		fr.following = strings.ToLower(strings.Join(follows, ","))

		return nil
	}()

	if err != nil {
		return err
	}

	_, err = r.UpdateUserByEmail(e, func(*userdomain.User) (*userdomain.User, error) {
		return &uf.User, nil
	})
	return err
}
