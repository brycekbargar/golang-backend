package serialization

import "github.com/brycekbargar/realworld-backend/domain"

type user struct {
	User userUser `json:"user"`
}
type userUser struct {
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
) interface{} {
	return &user{
		userUser{
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

type profile struct {
	Profile profileUser `json:"profile"`
}
type profileUser struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

// Following is a convenience func for when the user is following a profile.
func Following(*domain.User) bool { return true }

// NotFollowing is a convenience func for when the user is not following a profile.
func NotFollowing(*domain.User) bool { return false }

// MaybeFollowing is a convenience func for checking if the  user is following a profile.
func MaybeFollowing(f *domain.Fanboy) func(u *domain.User) bool {
	return func(u *domain.User) bool { return f.IsFollowing(u.Email) }
}

// UserToProfile converts a domain user to an output serializable profile.
func UserToProfile(
	u *domain.User,
	f func(*domain.User) bool,
) interface{} {
	return &profile{
		profileUser{
			Username:  u.Username,
			Bio:       u.Bio,
			Image:     u.Image,
			Following: f(u),
		},
	}
}
