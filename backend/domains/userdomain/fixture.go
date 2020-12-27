package userdomain

import "golang.org/x/crypto/bcrypt"

var enc, _ = bcrypt.GenerateFromPassword([]byte("Test1234!"), 14)

// Fixture is a slice of valid user objects (mostly for test data purposes).
var Fixture = []*User{
	{
		"user@comprehensive.com",
		"comprehensive username",
		"comprehensive bio",
		"http://comprehensive.com/image.png",
		make([]*User, 0),
		enc,
	},
	{
		"user@limping.com",
		"limping username",
		"limping bio",
		"http://limping.com/image.png",
		make([]*User, 0),
		enc,
	},
	{
		"user@public.com",
		"public username",
		"public bio",
		"http://public.com/image.png",
		make([]*User, 0),
		enc,
	},
	{
		"user@jaded.com",
		"jaded username",
		"jaded bio",
		"http://jaded.com/image.png",
		make([]*User, 0),
		enc,
	},
}
