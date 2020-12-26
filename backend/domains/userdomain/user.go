package userdomain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHash is an indicator that a string is a bcrypt hashed value.
type PasswordHash = string

// ErrorRequiredNewUserFields indicates when a NewUser is attempted to be created without all the required fields.
var ErrorRequiredNewUserFields = errors.New("password is required to create a user")

// ErrorRequiredUserFields indicates when a NewUser is instantiated without all the required fields.
var ErrorRequiredUserFields = errors.New("email and username are required for users")

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
	if len(email) == 0 || len(username) == 0 {
		return nil, ErrorRequiredUserFields
	}
	if len(password) == 0 {
		return nil, ErrorRequiredNewUserFields
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
	if len(email) == 0 || len(username) == 0 {
		return nil, ErrorRequiredUserFields
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
func UpdatedUser(user *User, email string, username string, bio *string, image *string, password string) (*User, error) {
	if len(email) > 0 {
		user.email = email
	}
	if len(username) > 0 {
		user.username = username
	}
	if bio != nil {
		user.bio = *bio
	}
	if image != nil {
		user.image = *image
	}
	if len(password) > 0 {
		pw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			return nil, err
		}
		user.password = pw
	}

	return user, nil
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

// HasPassword checks if the provide password string matches the stored hash for the user.
func (u *User) HasPassword(password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(u.password, []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}

// IsFollowing checks if the provided user is currently being followed by this user.
func (u *User) IsFollowing(fu *User) bool {
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
	if u.IsFollowing(fu) {
		return
	}
	u.following = append(u.following, fu)
}

// StopFollowing tracks that the provided user should be unfollowed.
// This method is idempotent (but possibly not thread-safe).
func (u *User) StopFollowing(su *User) {
	for i, f := range u.following {
		if f.email == su.email {
			u.following = append(u.following[:i], u.following[i+1:]...)
		}
	}
}
