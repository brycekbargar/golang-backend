package userdomain

import (
	"errors"
	"net/url"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHash is an indicator that a string is a bcrypt hashed value.
type PasswordHash = string

// ErrInvalidEmail indicates when a User is instantiated with an invalid email.
var ErrInvalidEmail = errors.New("email must be a valid email")

// ErrInvalidImage indicates when a User is instantiated with an invalid image url.
var ErrInvalidImage = errors.New("image must be a valid url")

// ErrRequiredUserFields indicates when a User is instantiated without all the required fields.
var ErrRequiredUserFields = errors.New("email, username, and password are required for users")

// https://golangcode.com/validate-an-email-address/
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// User is an individual user in the application.
// A user can be both the current client logged in (usually id'd by username)
// and also an author of a post or someone to follow.
type User struct {
	email     string
	username  string
	bio       string
	image     string
	following []*User
	password  []byte
}

// NewUserWithPassword creates a new User with the provide information.
func NewUserWithPassword(email string, username string, password string) (*User, error) {
	if len(email) == 0 || len(username) == 0 || len(password) == 0 {
		return nil, ErrRequiredUserFields
	}

	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	pw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	return &User{
		email:    email,
		username: username,
		password: pw,
	}, nil
}

// ExistingUser creates a User with the provide information.
func ExistingUser(email string, username string, bio string, image string, following []*User, password PasswordHash) (*User, error) {
	if len(email) == 0 || len(username) == 0 || len(password) == 0 {
		return nil, ErrRequiredUserFields
	}

	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}

	if len(image) > 0 {
		if _, err := url.ParseRequestURI(image); err != nil {
			return nil, ErrInvalidImage
		}
	}

	return &User{
		email,
		username,
		bio,
		image,
		following,
		[]byte(password),
	}, nil
}

// UpdatedUser merged the provided user with the (optional) new values provided.
func UpdatedUser(user User, email string, username string, bio *string, image *string, password string) (*User, error) {
	if len(email) > 0 {
		if !emailRegex.MatchString(email) {
			return nil, ErrInvalidEmail
		}
		user.email = email
	}
	if len(username) > 0 {
		user.username = username
	}
	if bio != nil {
		user.bio = *bio
	}
	if image != nil {
		if len(*image) > 0 {
			if _, err := url.ParseRequestURI(*image); err != nil {
				return nil, ErrInvalidImage
			}
		}
		user.image = *image
	}
	if len(password) > 0 {
		pw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			return nil, err
		}
		user.password = pw
	}

	f := make([]*User, len(user.following))
	copy(f, user.following)
	user.following = f

	return &user, nil
}

// Email is user's email address, which acts as their id.
func (u User) Email() string {
	return u.email
}

// Username is how they are displayed to other users and acts as a secondary id.
func (u User) Username() string {
	return u.username
}

// Bio is an optional blurb a user enters about themselves.
func (u User) Bio() string {
	return u.bio
}

// Image is the optional href to the user's profile picture.
func (u User) Image() string {
	return u.image
}

// Password gets the user's hashed password.
func (u User) Password() PasswordHash {
	return string(u.password)
}

// FollowingEmails gets the emails of the other users this user is following.
func (u User) FollowingEmails() []string {
	es := make([]string, 0, len(u.following))
	for _, fu := range u.following {
		es = append(es, fu.email)
	}

	return es
}

// HasPassword checks if the provide password string matches the stored hash for the user.
func (u *User) HasPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// IsFollowing checks if the provided user is currently being followed by this user.
func (u *User) IsFollowing(fu *User) bool {
	if fu == nil {
		return false
	}

	for _, f := range u.following {
		if f.email == fu.email {
			return true
		}
	}

	return false
}

// StartFollowing tracks that the provided user should be followed.
// This method is idempotent (but possibly not thread-safe).
func (u *User) StartFollowing(fu *User) {
	if fu == nil || u.IsFollowing(fu) {
		return
	}
	u.following = append(u.following, fu)
}

// StopFollowing tracks that the provided user should be unfollowed.
// This method is idempotent (but possibly not thread-safe).
func (u *User) StopFollowing(su *User) {
	if su == nil {
		return
	}

	for i, f := range u.following {
		if f.email == su.email {
			u.following = append(u.following[:i], u.following[i+1:]...)
		}
	}
}
