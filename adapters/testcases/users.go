package testcases

import (
	"fmt"
	"testing"

	"github.com/brycekbargar/realworld-backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Users_CreateUser(
	t *testing.T,
	r domain.Repository,
) {
	u1 := testUser("faithful")
	cu1, err := r.CreateUser(u1)
	require.NoError(t, err)
	assert.Equal(t, u1, cu1)

	u2 := testUser("kindhearted")
	cu2, err := r.CreateUser(u2)
	require.NoError(t, err)
	assert.Equal(t, u2, cu2)

	u1.Username = "icy username"
	_, err = r.CreateUser(u1)
	assert.ErrorIs(t, err, domain.ErrDuplicateUser)

	u2.Email = "user@icy.com"
	_, err = r.CreateUser(u2)
	assert.ErrorIs(t, err, domain.ErrDuplicateUser)

	r.UpdateUserByEmail(
		"user@faithful.com",
		func(u *domain.User) (*domain.User, error) {
			u.Bio = "noisy bio"
			u.Image = "http://noisy.com/profile.png"
			u.SetPassword("!4321tseT")
			return u, nil
		})
	fu, err := r.GetUserByEmail("user@faithful.com")
	require.NoError(t, err)
	assert.Equal(t, "noisy bio", fu.Bio)
	assert.Equal(t, "http://noisy.com/profile.png", fu.Image)
	hp, err := fu.HasPassword("!4321tseT")
	require.NoError(t, err)
	assert.True(t, hp, "because the password changed")
}

func Users_GetUserByEmail(
	t *testing.T,
	r domain.Repository,
) {
	u := testUser("finicky")
	_, err := r.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.GetUserByEmail(u.Email)
	require.NoError(t, err)
	assert.Equal(t, *u, fu.User)

	_, err = r.GetUserByEmail("user@light.com")
	assert.ErrorIs(t, err, domain.ErrUserNotFound)

	r.UpdateUserByEmail(
		"user@finicky.com",
		func(u *domain.User) (*domain.User, error) {
			u.Email = "user@nutty.com"
			return u, nil
		})
	fu, err = r.GetUserByEmail("user@nutty.com")
	require.NoError(t, err)
	assert.Equal(t, "finicky username", fu.Username)

	r.CreateUser(testUser("snobbish"))
	_, err = r.UpdateUserByEmail(
		"user@nutty.com",
		func(u *domain.User) (*domain.User, error) {
			u.Email = "user@snobbish.com"
			return u, nil
		})
	assert.ErrorIs(t, err, domain.ErrDuplicateUser)
	fu, err = r.GetUserByEmail("user@nutty.com")
	require.NoError(t, err)
}

func Users_GetUserByUsername(
	t *testing.T,
	r domain.Repository,
) {
	u := testUser("stormy")
	_, err := r.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.GetUserByUsername(u.Username)
	require.NoError(t, err)
	assert.Equal(t, u, fu)

	_, err = r.GetUserByUsername("dashing username")
	assert.ErrorIs(t, err, domain.ErrUserNotFound)

	r.UpdateUserByEmail(
		"user@stormy.com",
		func(u *domain.User) (*domain.User, error) {
			u.Username = "thirsty username"
			return u, nil
		})
	fu, err = r.GetUserByUsername("thirsty username")
	require.NoError(t, err)
	assert.Equal(t, "user@stormy.com", fu.Email)

	r.CreateUser(testUser("dusty"))
	_, err = r.UpdateUserByEmail(
		"user@stormy.com",
		func(u *domain.User) (*domain.User, error) {
			u.Username = "dusty username"
			return u, nil
		})
	assert.ErrorIs(t, err, domain.ErrDuplicateUser)
	fu, err = r.GetUserByUsername("thirsty username")
	assert.NoError(t, err)
}

func Users_UpdateFanboyByEmail_Following(
	t *testing.T,
	r domain.Repository,
) {
	u := testUser("gifted")
	_, err := r.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.GetUserByEmail(u.Email)
	require.NoError(t, err)
	assert.Empty(t, fu.FollowingEmails())

	for _, a := range []string{"important", "lumpy", "remarkable", "valuable"} {
		u := testUser(a)
		_, err := r.CreateUser(u)
		require.NoError(t, err)
	}

	err = r.UpdateFanboyByEmail(
		"user@gifted.com",
		func(f *domain.Fanboy) (*domain.Fanboy, error) {
			f.StartFollowing("user@lumpy.com")
			f.StartFollowing("user@remarkable.com")
			return f, nil
		})
	require.NoError(t, err)

	fu, err = r.GetUserByEmail("user@gifted.com")
	require.NoError(t, err)
	assert.Len(t, fu.FollowingEmails(), 2)
	assert.True(t, fu.IsFollowing("user@lumpy.com"))
	assert.False(t, fu.IsFollowing("user@valuable.com"))

	_, err = r.UpdateUserByEmail(
		"user@remarkable.com",
		func(u *domain.User) (*domain.User, error) {
			u.Email = "user@best.com"
			return u, nil
		})
	require.NoError(t, err)

	fu, err = r.GetUserByEmail("user@gifted.com")
	require.NoError(t, err)
	assert.True(t, fu.IsFollowing("user@best.com"))
	assert.False(t, fu.IsFollowing("user@remarkable.com"))

	err = r.UpdateFanboyByEmail(
		"user@gifted.com",
		func(f *domain.Fanboy) (*domain.Fanboy, error) {
			f.StopFollowing("user@best.com")
			f.StartFollowing("user@important.com")
			f.StartFollowing("user@valuable.com")
			f.StartFollowing("not an email")
			return f, nil
		})
	require.NoError(t, err)

	fu, err = r.GetUserByEmail("user@gifted.com")
	require.NoError(t, err)
	assert.Len(t, fu.FollowingEmails(), 3)
	assert.False(t, fu.IsFollowing("user@best.com"))
	assert.False(t, fu.IsFollowing("not an email"))
}

func Users_UpdateFanboyByEmail_Favorites(
	t *testing.T,
	r domain.Repository,
) {
	u := testUser("luxuriant")
	_, err := r.CreateUser(u)
	require.NoError(t, err)

	fu, err := r.GetUserByEmail(u.Email)
	require.NoError(t, err)
	assert.Empty(t, fu.FavoritedSlugs())

	r.CreateUser(testAuthor("brainy"))
	for _, adj := range []string{
		"callous",
		"aware",
		"magnificient",
	} {
		a := testArticle(adj)
		a.AuthorEmail = "author@brainy.com"

		_, err := r.CreateArticle(a)
		require.NoError(t, err)
	}

	err = r.UpdateFanboyByEmail(
		"user@luxuriant.com",
		func(f *domain.Fanboy) (*domain.Fanboy, error) {
			f.Favorite("callous-title")
			f.Favorite("aware-title")
			return f, nil
		})
	require.NoError(t, err)

	fu, err = r.GetUserByEmail("user@luxuriant.com")
	require.NoError(t, err)
	assert.Len(t, fu.FavoritedSlugs(), 2)
	assert.True(t, fu.Favors("callous-title"))
	assert.False(t, fu.Favors("magnificient-title"))

	fa, err := r.GetArticleBySlug("callous-title")
	require.NoError(t, err)
	assert.Equal(t, 1, fa.FavoriteCount)

	_, err = r.CreateUser(testUser("careful"))
	require.NoError(t, err)
	r.UpdateFanboyByEmail(
		"user@careful.com",
		func(f *domain.Fanboy) (*domain.Fanboy, error) {
			f.Favorite("callous-title")
			return f, nil
		})

	fa, err = r.GetArticleBySlug("callous-title")
	require.NoError(t, err)
	assert.Equal(t, 2, fa.FavoriteCount)

	_, err = r.UpdateArticleBySlug(
		"callous-title",
		func(a *domain.Article) (*domain.Article, error) {
			a.SetTitle("careful-title")
			return a, nil
		})
	require.NoError(t, err)

	fu, err = r.GetUserByEmail("user@luxuriant.com")
	require.NoError(t, err)
	assert.True(t, fu.Favors("careful-title"))
	assert.False(t, fu.Favors("callous-title"))

	err = r.UpdateFanboyByEmail(
		"user@luxuriant.com",
		func(f *domain.Fanboy) (*domain.Fanboy, error) {
			f.Unfavorite("aware-title")
			f.Favorite("magnificient-title")
			f.Favorite("magnificient-title")
			return f, nil
		})
	require.NoError(t, err)

	fu, err = r.GetUserByEmail("user@luxuriant.com")
	require.NoError(t, err)
	assert.Len(t, fu.FavoritedSlugs(), 2)
	assert.False(t, fu.Favors("aware-title"))
}

func testUser(adj string) *domain.User {
	u, _ := domain.NewUserWithPassword(
		fmt.Sprintf("user@%v.com", adj),
		fmt.Sprintf("%v username", adj),
		"Test1234!",
	)
	u.Bio = fmt.Sprintf("%v bio", adj)
	u.Image = fmt.Sprintf("http://%v.com/profile.png", adj)

	return u
}
