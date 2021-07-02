package postgres

import (
	"errors"
	"strings"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateUser creates a new user.
func (r *implementation) CreateUser(u *domain.User) (*domain.User, error) {
	res := r.db.Create(&User{
		Email:    u.Email,
		Username: u.Username,
		Bio:      u.Bio,
		Image:    u.Image,
		Password: Password{Value: u.Password},
	})

	var pgErr *pgconn.PgError
	if errors.As(res.Error, &pgErr) && pgErr.Code == "23505" {
		return nil, domain.ErrDuplicateUser
	}
	if res.Error != nil {
		return nil, res.Error
	}

	f, err := r.GetUserByEmail(u.Email)
	return &f.User, err
}

// GetUserByEmail finds a single user based on their email address.
func (r *implementation) GetUserByEmail(em string) (*domain.Fanboy, error) {
	found, err := r.getUserByEmail(em)
	if err != nil {
		return nil, err
	}

	follows := make(map[string]interface{}, len(found.Following))
	for _, u := range found.Following {
		follows[strings.ToLower(u.Email)] = nil
	}

	favorites := make(map[string]interface{}, len(found.Favorites))
	for _, a := range found.Favorites {
		favorites[strings.ToLower(a.Slug)] = nil
	}

	return &domain.Fanboy{
		User: domain.User{
			Email:    found.Email,
			Username: found.Username,
			Bio:      found.Bio,
			Image:    found.Image,
			Password: found.Password.Value,
		},
		Following: follows,
		Favorites: favorites,
	}, nil
}

func (r *implementation) getUserByEmail(em string) (*User, error) {
	var found User
	res := r.db.
		Preload(clause.Associations).
		First(&found, "email = ?", em)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return &found, nil
}

// GetAuthorByEmail finds a single author based on their email address or nil if they don't exist.
func (r *implementation) GetAuthorByEmail(em string) domain.Author {
	auth, err := r.getUserByEmail(em)
	if err != nil {
		return nil
	}
	return auth
}

// GetUserByUsername finds a single user based on their username.
func (r *implementation) GetUserByUsername(un string) (*domain.User, error) {
	var found User
	res := r.db.
		Preload(clause.Associations).
		First(&found, "username = ?", un)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}

	return &domain.User{
		Email:    found.Email,
		Username: found.Username,
		Bio:      found.Bio,
		Image:    found.Image,
		Password: found.Password.Value,
	}, nil
}

// UpdateUserByEmail finds a single user based on their email address,
// then applies the provide mutations.
func (r *implementation) UpdateUserByEmail(em string, update func(*domain.User) (*domain.User, error)) (*domain.User, error) {
	f, err := r.GetUserByEmail(em)
	if err != nil {
		return nil, err
	}

	user, err := update(&f.User)
	if err != nil {
		return nil, err
	}

	found, err := r.getUserByEmail(em)
	if err != nil {
		return nil, err
	}

	found.Email = user.Email
	found.Username = user.Username
	found.Bio = user.Bio
	found.Image = user.Image
	found.Password = Password{Value: user.Password}
	res := r.db.Save(found)

	var pgErr *pgconn.PgError
	if errors.As(res.Error, &pgErr) && pgErr.Code == "23505" {
		return nil, domain.ErrDuplicateUser
	}
	if res.Error != nil {
		return nil, res.Error
	}

	f, err = r.GetUserByEmail(user.Email)
	return &f.User, err
}

// UpdateFanboyByEmail finds a single user based on their email address,
// then applies the provide mutations (probably to the follower list).
func (r *implementation) UpdateFanboyByEmail(em string, update func(*domain.Fanboy) (*domain.Fanboy, error)) error {
	return nil
}
