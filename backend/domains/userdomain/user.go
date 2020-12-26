package userdomain

import "errors"

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
	following []User
}

// UserWithPassword is a user that is being created with a password.
// This will probably die at some point but I'm unsure atm how go idiomatically handles this.
type UserWithPassword struct {
	User
	password string
}

// NewUserWithPassword creates a new User with the provide information.
// password being a parameter (and also later a prop) is awful and will be removed in the future
func NewUserWithPassword(email string, username string, password string) (*UserWithPassword, error) {
	if len(password) == 0 {
		return nil, ErrorRequiredNewUserFields
	}

	user, err := NewUser(email, username, "", "")
	if err != nil {
		return nil, err
	}

	return &UserWithPassword{
		*user,
		password,
	}, nil
}

// NewUser creates a new User with the provide information.
// password being a parameter (and also later a prop) is awful and will be removed in the future
func NewUser(email string, username string, bio string, image string) (*User, error) {
	if len(email) == 0 || len(username) == 0 {
		return nil, ErrorRequiredUserFields
	}

	return &User{
		email,
		username,
		bio,
		image,
		make([]User, 0),
	}, nil
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

// Password gets the user's password during creation.
func (u UserWithPassword) Password() string {
	return u.password
}
