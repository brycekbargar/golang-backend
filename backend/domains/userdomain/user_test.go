package userdomain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

func TestNewUserWithPassword(t *testing.T) {
	t.Parallel()
	f := userdomain.Fixture

	cases := []struct {
		Name          string
		Email         string
		Username      string
		Password      string
		ExpectedError error
	}{
		{
			"Created (with hashed password)",
			f[0].Email(),
			f[0].Username(),
			f[0].Password(),
			nil,
		},
		{
			"Invalid Email",
			"definitely not an email",
			f[0].Username(),
			f[0].Password(),
			userdomain.ErrInvalidEmail,
		},
		{
			"Missing Email",
			"",
			f[0].Username(),
			f[0].Password(),
			userdomain.ErrRequiredUserFields,
		},
		{
			"Missing Username",
			f[0].Email(),
			"",
			f[0].Password(),
			userdomain.ErrRequiredUserFields,
		},
		{
			"Missing Password",
			f[0].Email(),
			f[0].Username(),
			"",
			userdomain.ErrRequiredUserFields,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			u, err := userdomain.NewUserWithPassword(tc.Email, tc.Username, tc.Password)
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.Email, u.Email())
				assert.Equal(t, tc.Username, u.Username())
				require.NotEmpty(t, u.Password())
				assert.NotEqual(t, tc.Password, u.Password())
			}
		})
	}
}

func TestExistingUser(t *testing.T) {
	t.Parallel()
	f := userdomain.Fixture

	cases := []struct {
		Name          string
		Email         string
		Username      string
		Bio           string
		Image         string
		Password      string
		ExpectedError error
	}{
		{
			"Created",
			f[0].Email(),
			f[0].Username(),
			f[0].Bio(),
			f[0].Image(),
			f[0].Password(),
			nil,
		},
		{
			"Missing Optional Fields",
			f[0].Email(),
			f[0].Username(),
			"",
			"",
			f[0].Password(),
			nil,
		},
		{
			"Invalid Email",
			"definitely not an email",
			f[0].Username(),
			f[0].Bio(),
			f[0].Image(),
			f[0].Password(),
			userdomain.ErrInvalidEmail,
		},
		{
			"Missing Email",
			"",
			f[0].Username(),
			f[0].Bio(),
			f[0].Image(),
			f[0].Password(),
			userdomain.ErrRequiredUserFields,
		},
		{
			"Missing Username",
			f[0].Email(),
			"",
			f[0].Bio(),
			f[0].Image(),
			f[0].Password(),
			userdomain.ErrRequiredUserFields,
		},
		{
			"Missing Password",
			f[0].Email(),
			f[0].Username(),
			f[0].Bio(),
			f[0].Image(),
			"",
			userdomain.ErrRequiredUserFields,
		},
		{
			"Invalid Image",
			f[0].Email(),
			f[0].Username(),
			f[0].Bio(),
			"definitely not an image",
			f[0].Password(),
			userdomain.ErrInvalidImage,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			u, err := userdomain.ExistingUser(tc.Email, tc.Username, tc.Bio, tc.Image, nil, tc.Password)
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.Email, u.Email())
			assert.Equal(t, tc.Username, u.Username())
			assert.Equal(t, tc.Bio, u.Bio())
			assert.Equal(t, tc.Image, u.Image())
			assert.NotNil(t, u.FollowingEmails()) // More thorough tests are elsewhere
			assert.Equal(t, tc.Password, u.Password())
		})
	}
}

func TestUpdated_User(t *testing.T) {
	t.Parallel()
	f := userdomain.Fixture

	cases := []struct {
		Name             string
		UpdatedEmail     string
		ExpectedEmail    string
		UpdatedUsername  string
		ExpectedUsername string
		UpdatedBio       *string
		ExpectedBio      string
		UpdatedImage     *string
		ExpectedImage    string
		UpdatedPassword  string
		ExpectedPassword string
		ExpectedError    error
	}{
		{
			"All New Values",
			"user@spotless.com",
			"user@spotless.com",
			"spotless username",
			"spotless username",
			optional("spotless bio"),
			"spotless bio",
			optional("http://spotless.com/image"),
			"http://spotless.com/image",
			"spotless password",
			"spotless password",
			nil,
		},
		{
			"No New Values",
			"",
			f[3].Email(),
			"",
			f[3].Username(),
			nil,
			f[3].Bio(),
			nil,
			f[3].Image(),
			"",
			"Test1234!",
			nil,
		},
		{
			"Invalid email",
			"definitely not an email",
			"na@anything.com",
			"n/a",
			"n/a",
			nil,
			"n/a",
			nil,
			"n/a",
			"n/a",
			"n/a",
			userdomain.ErrInvalidEmail,
		},
		{
			"Invalid image",
			"na@anything.com",
			"na@anything.com",
			"n/a",
			"n/a",
			nil,
			"n/a",
			optional("definitely not an image"),
			"n/a",
			"n/a",
			"n/a",
			userdomain.ErrInvalidImage,
		},
		{
			"Empty image and bio",
			"na@anything.com",
			"na@anything.com",
			"n/a",
			"n/a",
			optional(""),
			"",
			optional(""),
			"",
			"n/a",
			"n/a",
			nil,
		},
	}

	f[3].StartFollowing(f[0])
	f[3].StartFollowing(f[2])
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			u, err := userdomain.UpdatedUser(*f[3],
				tc.UpdatedEmail,
				tc.UpdatedUsername,
				tc.UpdatedBio,
				tc.UpdatedImage,
				tc.UpdatedPassword)
			if tc.ExpectedError != nil {
				require.EqualError(t, err, tc.ExpectedError.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.ExpectedEmail, u.Email())
			assert.Equal(t, tc.ExpectedUsername, u.Username())
			assert.Equal(t, tc.ExpectedBio, u.Bio())
			assert.Equal(t, tc.ExpectedImage, u.Image())
			assert.Len(t, u.FollowingEmails(), 2)
			h, err := u.HasPassword(tc.ExpectedPassword)
			require.NoError(t, err)
			assert.True(t, h)
		})
	}
}

func TestHasPassword(t *testing.T) {
	t.Parallel()

	pw := "impolite password"
	u, err := userdomain.NewUserWithPassword("user@impolite.com", "impolite user", pw)
	require.NoError(t, err)

	t.Run("Same Password", func(t *testing.T) {
		t.Parallel()

		h, err := u.HasPassword(pw)
		require.NoError(t, err)
		assert.True(t, h)
	})

	t.Run("Different Password", func(t *testing.T) {
		t.Parallel()

		h, err := u.HasPassword("literally anything else")
		require.NoError(t, err)
		assert.False(t, h)
	})

	t.Run("Bad Hash", func(t *testing.T) {
		t.Parallel()

		u, err := userdomain.ExistingUser("user@frail.com", "frail user", "", "", nil, "frail password")
		require.NoError(t, err)

		_, err = u.HasPassword("doesn't matter")
		require.Error(t, err, "because ExistingUser expects a password hash which wasn't provided")
	})
}

func TestUser_Following(t *testing.T) {
	t.Parallel()
	f := userdomain.Fixture

	u, err := userdomain.ExistingUser(
		f[0].Email(),
		f[0].Username(),
		f[0].Bio(),
		f[0].Image(),
		[]*userdomain.User{
			f[2],
			f[3],
		},
		f[0].Password(),
	)
	require.NoError(t, err)

	fem := u.FollowingEmails()
	assert.Len(t, fem, 2, "because we initialized the existing user as following 2 other users")
	for _, em := range fem {
		assert.Contains(t, []string{f[2].Email(), f[3].Email()}, em)
	}
	assert.True(t, u.IsFollowing(f[2]))
	assert.True(t, u.IsFollowing(f[3]))

	u.StartFollowing(f[1])
	fem = u.FollowingEmails()
	assert.Len(t, fem, 3, "because they're no longer following a new user")
	assert.Contains(t, fem, f[1].Email())
	assert.True(t, u.IsFollowing(f[1]))

	u.StartFollowing(f[1])
	fem = u.FollowingEmails()
	assert.Len(t, fem, 3, "because following is idempotent")

	u.StopFollowing(f[2])
	fem = u.FollowingEmails()
	assert.Len(t, fem, 2, "because they're no longer following a user")
	assert.NotContains(t, fem, f[2].Email())
	assert.False(t, u.IsFollowing(f[2]))

	u.StopFollowing(f[2])
	fem = u.FollowingEmails()
	assert.Len(t, fem, 2, "because unfollowing is idempotent")

	assert.NotPanics(t, func() {
		u := u
		u.IsFollowing(nil)
	}, "because checking for nil followers should be ok")

	assert.NotPanics(t, func() {
		u := u
		u.StartFollowing(nil)
	}, "because following nil should be ok")

	assert.NotPanics(t, func() {
		u := u
		u.StopFollowing(nil)
	}, "because unfollowing nil should be ok")

	u.StopFollowing(f[1])
	u.StopFollowing(f[3])
	u.StopFollowing(f[3])
	fem = u.FollowingEmails()
	assert.Empty(t, fem, "because we've unfollowed everyone")
	assert.False(t, u.IsFollowing(f[1]))
	assert.False(t, u.IsFollowing(f[3]))
}

func optional(s string) *string {
	return &s
}
