package adapters_test

import (
	"fmt"
	"testing"

	"github.com/brycekbargar/realworld-backend/adapters"
	"github.com/brycekbargar/realworld-backend/adapters/inmemory"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRepositories(t *testing.T) []*adapters.RepositoryImplementation {
	t.Helper()

	return []*adapters.RepositoryImplementation{
		inmemory.NewInstance(),
	}
}

func Test_Users(t *testing.T) {
	t.Helper()

	for _, r := range createRepositories(t) {
		t.Parallel()

		r := r
		t.Run(r.Name, func(t *testing.T) {
			t.Parallel()

			t.Run("Create and Update User", func(t *testing.T) {
				t.Parallel()
				test_Users_CreateUser(t, r)
			})
			t.Run("Get and Update User By Email", func(t *testing.T) {
				t.Parallel()
				test_Users_GetUserByEmail(t, r)
			})
			t.Run("Get and Update User By Username", func(t *testing.T) {
				t.Parallel()
				test_Users_GetUserByUsername(t, r)
			})
			t.Run("Update Fanboy", func(t *testing.T) {
				t.Parallel()
				test_Users_UpdateFanboyByEmail(t, r)
			})
		})
	}

}

func test_Users_CreateUser(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	u1, err := userdomain.NewUserWithPassword(
		"user@faithful.com",
		"faithful username",
		"Test1234!",
	)
	cu1, err := r.Users.CreateUser(u1)
	require.NoError(t, err)
	assert.Equal(t, u1, cu1)

	u2, err := userdomain.NewUserWithPassword(
		"user@kindhearted.com",
		"kindhearted username",
		"Test1234!",
	)
	cu2, err := r.Users.CreateUser(u2)
	require.NoError(t, err)
	assert.Equal(t, u2, cu2)

	u1.Username = "icy username"
	_, err = r.Users.CreateUser(u1)
	assert.ErrorIs(t, err, userdomain.ErrDuplicateValue)

	u2.Email = "user@icy.com"
	_, err = r.Users.CreateUser(u2)
	assert.ErrorIs(t, err, userdomain.ErrDuplicateValue)

	r.Users.UpdateUserByEmail(
		"user@faithful.com",
		func(u *userdomain.User) (*userdomain.User, error) {
			u.Bio = "noisy bio"
			u.Image = "http://noisy.com/profile.png"
			u.SetPassword("!4321tseT")
			return u, nil
		})
	fu, err := r.Users.GetUserByEmail("user@faithful.com")
	require.NoError(t, err)
	assert.Equal(t, "noisy bio", fu.Bio)
	assert.Equal(t, "http://noisy.com/profile.png", fu.Image)
	hp, err := fu.HasPassword("!4321tseT")
	require.NoError(t, err)
	assert.True(t, hp)
}

func test_Users_GetUserByEmail(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	u, err := userdomain.NewUserWithPassword(
		"user@finicky.com",
		"finicky username",
		"Test1234!",
	)
	_, err = r.Users.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.Users.GetUserByEmail(u.Email)
	require.NoError(t, err)
	assert.Equal(t, *u, fu.User)

	_, err = r.Users.GetUserByEmail("user@light.com")
	assert.ErrorIs(t, err, userdomain.ErrNotFound)

	r.Users.UpdateUserByEmail(
		"user@finicky.com",
		func(u *userdomain.User) (*userdomain.User, error) {
			u.Email = "user@nutty.com"
			return u, nil
		})
	fu, err = r.Users.GetUserByEmail("user@nutty.com")
	require.NoError(t, err)
	assert.Equal(t, "finicky username", fu.Username)
}

func test_Users_GetUserByUsername(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	u, err := userdomain.NewUserWithPassword(
		"user@stormy.com",
		"stormy username",
		"Test1234!",
	)
	_, err = r.Users.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.Users.GetUserByUsername(u.Username)
	require.NoError(t, err)
	assert.Equal(t, u, fu)

	_, err = r.Users.GetUserByUsername("dashing username")
	assert.ErrorIs(t, err, userdomain.ErrNotFound)

	r.Users.UpdateUserByEmail(
		"user@stormy.com",
		func(u *userdomain.User) (*userdomain.User, error) {
			u.Username = "thirsty username"
			return u, nil
		})
	fu, err = r.Users.GetUserByUsername("thirsty username")
	require.NoError(t, err)
	assert.Equal(t, "user@stormy.com", fu.Email)
}

func test_Users_UpdateFanboyByEmail(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	u, err := userdomain.NewUserWithPassword(
		"user@gifted.com",
		"gifted username",
		"Test1234!",
	)
	_, err = r.Users.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.Users.GetUserByEmail(u.Email)
	require.NoError(t, err)
	assert.Empty(t, fu.FollowingEmails())

	for _, a := range []string{"important", "lumpy", "remarkable", "valuable"} {
		u, err := userdomain.NewUserWithPassword(
			fmt.Sprintf("user@%v.com", a),
			fmt.Sprintf("%v username", a),
			"Test1234!",
		)
		_, err = r.Users.CreateUser(u)
		require.NoError(t, err)
	}

	r.Users.UpdateFanboyByEmail(
		"user@gifted.com",
		func(f *userdomain.Fanboy) (*userdomain.Fanboy, error) {
			f.StartFollowing("user@lumpy.com")
			f.StartFollowing("user@remarkable.com")
			return f, nil
		})
	fu, err = r.Users.GetUserByEmail("user@gifted.com")
	require.NoError(t, err)
	assert.Len(t, fu.FollowingEmails(), 2)
	assert.True(t, fu.IsFollowing("user@lumpy.com"))
	assert.False(t, fu.IsFollowing("user@valuable.com"))

	r.Users.UpdateUserByEmail(
		"user@remarkable.com",
		func(u *userdomain.User) (*userdomain.User, error) {
			u.Email = "user@best.com"
			return u, nil
		})
	fu, err = r.Users.GetUserByEmail("user@gifted.com")
	require.NoError(t, err)
	assert.True(t, fu.IsFollowing("user@best.com"))
	assert.False(t, fu.IsFollowing("user@remarkable.com"))

	r.Users.UpdateFanboyByEmail(
		"user@gifted.com",
		func(f *userdomain.Fanboy) (*userdomain.Fanboy, error) {
			f.StopFollowing("user@best.com")
			f.StartFollowing("user@important.com")
			f.StartFollowing("user@valuable.com")
			f.StartFollowing("not an email")
			return f, nil
		})
	fu, err = r.Users.GetUserByEmail("user@gifted.com")
	require.NoError(t, err)
	assert.Len(t, fu.FollowingEmails(), 3)
	assert.False(t, fu.IsFollowing("user@best.com"))
	assert.False(t, fu.IsFollowing("not an email"))
}
