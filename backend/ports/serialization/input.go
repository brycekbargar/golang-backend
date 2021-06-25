package serialization

import "github.com/brycekbargar/realworld-backend/domain"

type register struct {
	User registerUser `json:"user"`
}
type registerUser struct {
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Password *string `json:"password,omitempty"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// RegisterToUser converts a input serializable user to a domain user.
func RegisterToUser(
	bind func(interface{}) error,
) (*domain.User, error) {
	r := new(register)
	if err := bind(r); err != nil {
		return nil, err
	}

	return domain.NewUserWithPassword(
		r.User.Email,
		r.User.Username,
		*r.User.Password,
	)
}

// RegisterToUser converts a input serializable user to a delta for a domain user.
func UpdateToDelta(
	bind func(interface{}) error,
) (func(*domain.User), error) {
	r := new(register)
	if err := bind(r); err != nil {
		return nil, err
	}

	return func(u *domain.User) {
		if r.User.Email != "" {
			u.Email = r.User.Email
		}
		if r.User.Username != "" {
			u.Username = r.User.Username
		}
		if r.User.Password != nil && *r.User.Password != "" {
			u.SetPassword(*r.User.Password)
		}
		if r.User.Bio != nil {
			u.Bio = *r.User.Bio
		}
		if r.User.Image != nil {
			u.Image = *r.User.Image
		}
	}, nil
}
