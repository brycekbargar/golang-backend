// Package userdomain represents the Model for Users. This includes IAM and Profiles.
// Notably this doesn't include Authorship of articles.
// In this simple implementation there are going to be a couple of hack to share the userdomain models with the articlesdomain.
// In a real application there should be some sort of published event to make sure the authors/commentors on articles are in sync.
// Also, in a real application writing IAM from scratch seems real bad.
// The Profiles "domain" is lumped in here because there are only a couple of profile related actions.
// I could definitely see this domain becoming exclusively "profile" as the capabilities are expanded and IAM is outsourced to something like Auth0.
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
