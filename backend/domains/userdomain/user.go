package userdomain

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHash is an indicator that a string is a bcrypt hashed value.
type PasswordHash = []byte

// User is an individual user in the application.
// A user can be both the current client logged in (usually id'd by email)
// and also an proile of someone that is followed (usually id'd by username).
type User struct {
	Email    string `valid:"required,email"`
	Username string `valid:"required"`
	Bio      string
	Image    string       `valid:"url,optional"`
	Password PasswordHash `valid:"required"`
}

// Fanboy is User with the Users they follow by email
type Fanboy struct {
	User
	Following map[string]interface{}
}

// NewUserWithPassword creates a new partially-hydrated User with the provide information.
func NewUserWithPassword(email string, username string, password string) (*User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	return (&User{
		Email:    email,
		Username: username,
		Password: pw,
	}).Validate()
}

// Validate returns the provided User if they are valid, otherwise error will contain validation errors.
func (u *User) Validate() (*User, error) {
	if v, err := govalidator.ValidateStruct(u); !v {
		return nil, err
	}

	return u, nil
}

// SetPassword sets the password hash from the plain-text value
func (u *User) SetPassword(password string) error {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	u.Password = pw
	return nil
}

// HasPassword checks if the provided password string matches the hash for the user.
func (u *User) HasPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// IsFollowing checks if the provided user is currently being followed by this user.
func (u *Fanboy) FollowingEmails() []string {
	emails := make([]string, 0)
	for em := range u.Following {
		if em != "" {
			emails = append(emails, em)
		}
	}
	return emails
}

// IsFollowing checks if the provided user is currently being followed by this user.
func (u *Fanboy) IsFollowing(email string) bool {
	if !govalidator.IsEmail(email) {
		return false
	}

	_, ok := u.Following[strings.ToLower(email)]
	return ok
}

// StartFollowing tracks that the provided user should be followed.
func (u *Fanboy) StartFollowing(email string) {
	if !govalidator.IsEmail(email) {
		return
	}
	u.Following[strings.ToLower(email)] = nil
}

// StopFollowing tracks that the provided user should be unfollowed.
func (u *Fanboy) StopFollowing(email string) {
	delete(u.Following, strings.ToLower(email))
}
