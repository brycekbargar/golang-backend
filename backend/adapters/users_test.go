package adapters_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/adapters/inmemory"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var subjects = map[string]userdomain.Repository{
	"inmemory": inmemory.NewUsers(),
}

func TestCreate_RoundTrips(t *testing.T) {
	t.Parallel()

	cu, err := userdomain.NewUserWithPassword("user@panicky.com", "panicky user", "panicky password")
	require.NoError(t, err)

	for k, r := range subjects {
		r := r
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			err := r.Create(cu)
			require.NoError(t, err)

			t.Run("By Email", func(t *testing.T) {
				t.Parallel()

				fu, err := r.GetUserByEmail("user@panicky.com")
				require.NoError(t, err)

				assert.Equal(t, "user@panicky.com", fu.Email())
				assert.Equal(t, "panicky user", fu.Username())
				h, err := fu.HasPassword("panicky password")
				require.NoError(t, err)
				assert.True(t, h)
			})

			t.Run("By Username", func(t *testing.T) {
				t.Parallel()

				fu, err := r.GetUserByUsername("panicky user")
				require.NoError(t, err)

				assert.Equal(t, "user@panicky.com", fu.Email())
				assert.Equal(t, "panicky user", fu.Username())
				h, err := fu.HasPassword("panicky password")
				require.NoError(t, err)
				assert.True(t, h)
			})
		})
	}
}

func TestCreate_DuplicateUser(t *testing.T) {
	t.Parallel()

	cu, err := userdomain.NewUserWithPassword("user@smart.com", "smart user", "smart password")
	require.NoError(t, err)

	for k, r := range subjects {
		r := r
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			err := r.Create(cu)
			require.NoError(t, err)

			t.Run("By Username", func(t *testing.T) {
				t.Parallel()

				du, err := userdomain.UpdatedUser(*cu, "user@ablaze.com", "", nil, nil, "")
				require.NoError(t, err)
				err = r.Create(du)
				assert.EqualError(t, err, userdomain.ErrDuplicateValue.Error(),
					"because usernames are unique")
			})

			t.Run("By Email", func(t *testing.T) {
				t.Parallel()
				du, err := userdomain.UpdatedUser(*cu, "", "ablaze user", nil, nil, "")
				require.NoError(t, err)
				err = r.Create(du)
				assert.EqualError(t, err, userdomain.ErrDuplicateValue.Error(),
					"because emails are unique")
			})
		})
	}
}

func TestUpdateUserByEmail(t *testing.T) {
	t.Parallel()
	f := userdomain.Fixture

	for k, r := range subjects {
		r := r
		t.Run(k, func(t *testing.T) {
			for _, u := range f {
				err := r.Create(u)
				require.NoError(t, err)
			}
			t.Parallel()

			err := r.UpdateUserByEmail("user@grumpy.com",
				func(u *userdomain.User) (*userdomain.User, error) {
					return u, nil
				})
			require.EqualError(t, err, userdomain.ErrNotFound.Error())

			err = r.UpdateUserByEmail(f[0].Email(),
				func(u *userdomain.User) (*userdomain.User, error) {
					return userdomain.ExistingUser(
						"user@aback.com",
						"aback username",
						"",
						"",
						[]*userdomain.User{
							f[1],
							f[2],
						},
						f[0].Password(),
					)
				})
			require.NoError(t, err)

			u, err := r.GetUserByEmail("user@aback.com")
			assert.NoError(t, err, "because a user was updated to this email")
			_, err = r.GetUserByUsername("aback username")
			assert.NoError(t, err, "because a user was updated to this username")
			_, err = r.GetUserByEmail(f[0].Email())
			assert.EqualError(t, err, userdomain.ErrNotFound.Error(),
				"because this user's email was updated")
			_, err = r.GetUserByUsername(f[0].Username())
			assert.EqualError(t, err, userdomain.ErrNotFound.Error(),
				"because this user's username was updated")

			assert.Len(t, u.FollowingEmails(), 2)
			assert.Contains(t, u.FollowingEmails(), f[1].Email())
			assert.Contains(t, u.FollowingEmails(), f[2].Email())

			err = r.UpdateUserByEmail(f[1].Email(),
				func(u *userdomain.User) (*userdomain.User, error) {
					nu, err := userdomain.NewUserWithPassword(
						"user@squealing.com",
						"squealing username",
						"squealing password")
					require.NoError(t, err)

					u.StartFollowing(nu)
					return u, nil
				})
			require.EqualError(t, err, userdomain.ErrNotFound.Error(),
				"because followed users must exist")
		})
	}
}
