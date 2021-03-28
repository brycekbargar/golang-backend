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

	_, err = r.Articles.CreateArticle(testArticle("hospitable"))
	assert.ErrorIs(t, err, articledomain.ErrDuplicateValue)

	now = time.Now().UTC()
	ua := testArticle("enchanting")
	ua.SetTitle("hospitable title")
	ua.AuthorEmail = u.Email
	_, err = r.Articles.UpdateArticleBySlug(
		"hospitable-title",
		func(a *articledomain.Article) (*articledomain.Article, error) {
			return ua, nil
		})
	require.NoError(t, err)
	fa, err := r.Articles.GetArticleBySlug("hospitable-title")
	assert.NoError(t, err)
	assert.Equal(t, ua.Slug, fa.Slug)
	assert.Equal(t, ua.Title, fa.Title)
	assert.Equal(t, ua.Body, fa.Body)

	assert.True(t, fa.CreatedAtUTC.Before(now))
	assert.True(t, fa.UpdatedAtUTC.After(now))

	r.Users.UpdateUserByEmail(
		"author@hospitable.com",
		func(u *userdomain.User) (*userdomain.User, error) {
			u.Email = "author@whole.com"
			return u, nil
		})
	fa, err = r.Articles.GetArticleBySlug("hospitable-title")
	require.NoError(t, err)
	assert.Equal(t, "author@whole.com", fa.AuthorEmail)
	assert.Equal(t, "author@whole.com", fa.Author.Email())
}

func Articles_GetArticleBySlug(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	u := testAuthor("observant")
	r.Users.CreateUser(u)

	now := time.Now().UTC()
	a := testArticle("observant")
	_, err := r.Articles.CreateArticle(a)
	require.NoError(t, err)

	_, err = r.Articles.GetArticleBySlug("befitting-title")
	assert.ErrorIs(t, err, articledomain.ErrNotFound)

	fa, err := r.Articles.GetArticleBySlug(a.Slug)
	assert.NoError(t, err)
	assert.Equal(t, a.Slug, fa.Slug)
	assert.Equal(t, a.Title, fa.Title)
	assert.Equal(t, a.Body, fa.Body)

	assert.Equal(t, a.AuthorEmail, fa.AuthorEmail)
	assert.Equal(t, a.AuthorEmail, fa.Author.Email())
	assert.Equal(t, u.Bio, fa.Author.Bio())
	assert.Equal(t, u.Image, fa.Author.Image())

	assert.True(t, fa.CreatedAtUTC.After(now))
	assert.Equal(t, fa.CreatedAtUTC, fa.UpdatedAtUTC)

	r.Articles.UpdateArticleBySlug(
		"observant-title",
		func(a *articledomain.Article) (*articledomain.Article, error) {
			a.SetTitle("silent title")
			return a, nil
		})
	_, err = r.Articles.GetArticleBySlug("observant-title")
	assert.ErrorIs(t, err, articledomain.ErrNotFound)
	_, err = r.Articles.GetArticleBySlug("silent-title")
	assert.NoError(t, err)

	r.Users.CreateUser(testAuthor("modern"))
	r.Articles.CreateArticle(testArticle("modern"))
	_, err = r.Articles.UpdateArticleBySlug(
		"silent-title",
		func(a *articledomain.Article) (*articledomain.Article, error) {
			a.SetTitle("modern title")
			return a, nil
		})
	assert.ErrorIs(t, err, articledomain.ErrDuplicateValue)
	fa, err = r.Articles.GetArticleBySlug("silent-title")
	assert.NoError(t, err)
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
