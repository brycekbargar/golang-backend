package serialization

import "github.com/brycekbargar/realworld-backend/domain"

// User is the json wrapper for a single user.
type User struct {
	User UserUser `json:"user"`
}

// UserUser is the contract for user operations.
type UserUser struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    *string `json:"image"`
}

// UserToUser converts a domain user to an output serializable user.
func UserToUser(
	u *domain.User,
	t string,
) *User {
	return &User{
		UserUser{
			Email:    u.Email,
			Token:    t,
			Username: u.Username,
			Bio:      optional(u.Bio),
			Image:    optional(u.Image),
		},
	}
}

func optional(s string) *string {
	return &s
}

// Profile is the json wrapper for a single profile.
type Profile struct {
	Profile ProfileUser `json:"profile"`
}

// ProfileUser is the contract for profile operations.
type ProfileUser struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func Following(*domain.User) bool    { return true }
func NotFollowing(*domain.User) bool { return false }
func MaybeFollowing(f *domain.Fanboy) func(u *domain.User) bool {
	return func(u *domain.User) bool { return f.IsFollowing(u.Email) }
}

// UserToProfile converts a domain user to an output serializable profile.
func UserToProfile(
	u *domain.User,
	f func(*domain.User) bool,
) *Profile {
	return &Profile{
		ProfileUser{
			Username:  u.Username,
			Bio:       u.Bio,
			Image:     u.Image,
			Following: f(u),
		},
	}
}
