package inmemory

import (
	"strings"
	"sync"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

type userRecord struct {
	email     string
	username  string
	bio       string
	image     string
	following string
	password  string
}

// Users is a (super inefficient) in-memory repository implementation for the usersdomain.Repository.
type Users struct {
	mu   *sync.Mutex
	repo map[string]userRecord
}

// NewUsers creates a new userdomain.Repository implementation for the users domain
func NewUsers() *Users {
	return &Users{
		&sync.Mutex{},
		make(map[string]userRecord),
	}
}

// Create creates a new user.
func (r *Users) Create(u *userdomain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for e, v := range r.repo {
		if e == strings.ToLower(u.Email()) ||
			strings.ToLower(v.username) == strings.ToLower(u.Username()) {
			return userdomain.ErrDuplicateValue
		}
	}

	r.repo[strings.ToLower(u.Email())] = userRecord{
		u.Email(),
		u.Username(),
		u.Bio(),
		u.Image(),
		strings.Join(u.FollowingEmails(), ","),
		u.Password(),
	}
	return nil
}

// GetUserByEmail finds a single user based on their username.
func (r *Users) GetUserByEmail(e string) (*userdomain.User, error) {
	return r.getUserByEmail(e, true)
}

// getUserByEmail finds a single user based on their username (without infinitely recursing)
func (r *Users) getUserByEmail(e string, recurse bool) (*userdomain.User, error) {
	if u, ok := r.repo[strings.ToLower(e)]; ok {

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
func (r *Users) GetUserByUsername(un string) (*userdomain.User, error) {
	for k, v := range r.repo {
		if strings.ToLower(v.username) == strings.ToLower(un) {
			return r.GetUserByEmail(k)
		}
	}

	return nil, userdomain.ErrNotFound
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *Users) UpdateUserByEmail(e string, uf func(*userdomain.User) (*userdomain.User, error)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, err := r.GetUserByEmail(e)
	if err != nil {
		return err
	}

	u, err = uf(u)
	if err != nil {
		return err
	}

	for _, fe := range u.FollowingEmails() {
		if _, err := r.getUserByEmail(fe, false); err != nil {
			return err
		}
	}

	ru := r.repo[strings.ToLower(e)]
	delete(r.repo, strings.ToLower(e))

	for e, v := range r.repo {
		if e == strings.ToLower(u.Email()) ||
			strings.ToLower(v.username) == strings.ToLower(u.Username()) {

			// Add the deleted user back if they've become a duplicate
			r.repo[strings.ToLower(e)] = ru
			return userdomain.ErrDuplicateValue
		}
	}

	r.repo[strings.ToLower(u.Email())] = userRecord{
		u.Email(),
		u.Username(),
		u.Bio(),
		u.Image(),
		strings.Join(u.FollowingEmails(), ","),
		u.Password(),
	}

	return nil
}
