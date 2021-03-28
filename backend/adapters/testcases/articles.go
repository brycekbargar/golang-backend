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
	a1, err := articledomain.NewArticle(
		"hospitable title",
		"hospitable description",
		"hospitable body",
		"author@hospitable.com",
		"hospitable one", "hospitable two", "hospitable three",
	)
	_, err = r.Articles.CreateArticle(a1)
	assert.ErrorIs(t, err, articledomain.ErrNoAuthor)

	u := testAuthor("hospitable")
	r.Users.CreateUser(u)

	now := time.Now().UTC()
	ca1, err := r.Articles.CreateArticle(a1)
	require.NoError(t, err)

	assert.Equal(t, a1.Slug, ca1.Slug)
	assert.Equal(t, a1.Title, ca1.Title)
	assert.Equal(t, a1.Body, ca1.Body)
	assert.Empty(t, a1.FavoritedBy)

	assert.Equal(t, u.Email, ca1.AuthorEmail)
	assert.Equal(t, u.Email, ca1.Author.Email())
	assert.Equal(t, u.Bio, ca1.Author.Bio())
	assert.Equal(t, u.Image, ca1.Author.Image())

	assert.True(t, ca1.CreatedAtUTC.After(now))
	assert.Equal(t, ca1.CreatedAtUTC, ca1.UpdatedAtUTC)
}

func testAuthor(adj string) *userdomain.User {
	a := testUser(adj)
	a.Email = fmt.Sprintf("author@%v.com", adj)
	return a
}
