package testcases

import (
	"fmt"
	"testing"
	"time"

	"github.com/brycekbargar/realworld-backend/adapters"
	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Articles_CreateArticle(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	a := testArticle("hospitable")
	_, err := r.Articles.CreateArticle(a)
	assert.ErrorIs(t, err, articledomain.ErrNoAuthor)

	u := testAuthor("hospitable")
	r.Users.CreateUser(u)

	now := time.Now().UTC()
	ca, err := r.Articles.CreateArticle(a)
	require.NoError(t, err)

	assert.Equal(t, a.Slug, ca.Slug)
	assert.Equal(t, a.Title, ca.Title)
	assert.Equal(t, a.Body, ca.Body)
	assert.Empty(t, a.FavoritedBy)

	assert.Equal(t, a.AuthorEmail, ca.AuthorEmail)
	assert.Equal(t, a.AuthorEmail, ca.Author.Email())
	assert.Equal(t, u.Bio, ca.Author.Bio())
	assert.Equal(t, u.Image, ca.Author.Image())

	assert.True(t, ca.CreatedAtUTC.After(now))
	assert.Equal(t, ca.CreatedAtUTC, ca.UpdatedAtUTC)

	r.Users.UpdateUserByEmail(
		"author@hospitable.com",
		func(u *userdomain.User) (*userdomain.User, error) {
			u.Email = "author@whole.com"
			return u, nil
		})
	fa, err := r.Articles.GetArticleBySlug("hospitable-title")
	require.NoError(t, err)
	assert.Equal(t, "author@whole.com", fa.AuthorEmail)
	assert.Equal(t, "author@whole.com", fa.Author.Email())

	r.Articles.UpdateArticleBySlug(
		"hospitable-title",
		func(a *articledomain.Article) (*articledomain.Article, error) {
			a.SetTitle("venomous title")
			return a, nil
		})
	_, err = r.Articles.GetArticleBySlug("hospitable-title")
	require.ErrorIs(t, err, articledomain.ErrNotFound)
	fa2, err := r.Articles.GetArticleBySlug("venomous-title")
	require.NoError(t, err)
	assert.Equal(t, fa.CreatedAtUTC, fa2.CreatedAtUTC)
	assert.True(t, fa.UpdatedAtUTC.Before(fa2.UpdatedAtUTC))
}

func testAuthor(adj string) *userdomain.User {
	a := testUser(adj)
	a.Email = fmt.Sprintf("author@%v.com", adj)
	return a
}

func testArticle(adj string) (a *articledomain.Article) {
	a, _ = articledomain.NewArticle(
		fmt.Sprintf("%v title", adj),
		fmt.Sprintf("%v description", adj),
		fmt.Sprintf("%v body", adj),
		fmt.Sprintf("author@%v.com", adj),
		fmt.Sprintf("%v one", adj),
		fmt.Sprintf("%v two", adj),
		fmt.Sprintf("%v three", adj),
	)
	return
}
