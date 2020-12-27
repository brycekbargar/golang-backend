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
