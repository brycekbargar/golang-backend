package serialization

import "github.com/brycekbargar/realworld-backend/domain"

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
