package userdomain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

func TestNewUserWithPassword(t *testing.T) {
	t.Parallel()

	t.Run("Password is hashed", func(t *testing.T) {
		t.Parallel()

		u, err := userdomain.NewUserWithPassword("user@uneven.com", "uneven user", "Test1234!")
		require.NoError(t, err)
		assert.NotEqual(t, []byte("Test1234!"), u.Password)

	})

	t.Run("Validation happens", func(t *testing.T) {
		t.Parallel()

		u, err := userdomain.NewUserWithPassword("not a cuddly email", "cuddly user", "Test1234!")
		assert.NotNil(t, err)
		assert.Nil(t, u)
	})
}

func TestHasPassword(t *testing.T) {
	t.Parallel()

	t.Run("Password can be changed", func(t *testing.T) {
		t.Parallel()
		u, err := userdomain.NewUserWithPassword("user@impolite.com", "impolite user", "Test1234!")
		require.NoError(t, err)

		hp, err := u.HasPassword("Test1233!")
		require.NoError(t, err)
		assert.True(t, hp)

		hp, err = u.HasPassword("Completely different password")
		require.NoError(t, err)
		assert.False(t, hp)

		require.NoError(t, u.SetPassword("Completely different password"))
		hp, err = u.HasPassword("Completely different password")
		require.NoError(t, err)
		assert.True(t, hp)
	})

	t.Run("Bad password hash errors", func(t *testing.T) {
		t.Parallel()
		u := userdomain.User{
			Password: []byte("definitely not a valid password hash"),
		}

		hp, err := u.HasPassword("Test1234!")
		require.NoError(t, err)
		assert.True(t, hp)
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name string
		User *userdomain.User
	}{
		{
			"Invalid Email",
			&userdomain.User{
				Email:    "excited user email",
				Username: "excited user",
				Password: []byte("required but not validated"),
			},
		},
		{
			"Missing Email",
			&userdomain.User{
				Username: "gray user",
				Password: []byte("required but not validated"),
			},
		},
		{
			"Missing Username",
			&userdomain.User{
				Email:    "user@melodic.com",
				Password: []byte("required but not validated"),
			},
		},
		{
			"Missing Password",
			&userdomain.User{
				Email:    "user@adjoining.com",
				Username: "adjoining user",
			},
		},
		{
			"Invalid image",
			&userdomain.User{
				Email:    "user@good.com",
				Username: "good user",
				Image:    "http://this is not an image url",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			u, err := tc.User.Validate()
			assert.Error(t, err)
			assert.Nil(t, u)
		})
	}

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		u := &userdomain.User{
			Email:    "user@needy.com",
			Username: "needy user",
			Password: []byte("required but not validated"),
			Bio:      "needy bio",
			Image:    "https://profileimages.com/needy.gif",
		}

		vu, err := u.Validate()
		require.NoError(t, err)
		assert.Same(t, u, vu)
	})
}

func TestUser_Following(t *testing.T) {
	t.Parallel()

	f := userdomain.Fanboy{
		Following: map[string]interface{}{
			"user@craven.com":     nil,
			"user@cooing.com":     nil,
			"user@unsuitable.com": nil,
		},
	}
	assert.True(t, f.IsFollowing("user@craven.com"))
	assert.True(t, f.IsFollowing("user@cooing.com"))

	assert.False(t, f.IsFollowing("user@thoughtful.com"))
	f.StartFollowing("user@thoughtful.com")
	assert.True(t, f.IsFollowing("user@thoughtful.com"))

	f.StartFollowing("user@craven.com")
	assert.True(t, f.IsFollowing("user@craven.com"),
		"because following is idempotent")

	f.StopFollowing("user@craven.com")
	f.StopFollowing("user@craven.com")
	assert.False(t, f.IsFollowing("user@craven.com"),
		"because following is idempotent")

	f.StartFollowing("definitely not an email")
	assert.False(t, f.IsFollowing("definitely not an email"),
		"because only valid emails can be followed")
}
