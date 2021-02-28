package adapters_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/adapters/inmemory"
	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var aSubjects = map[string]articledomain.Repository{
	"inmemory": inmemory.Articles,
}

func init() {
	for _, r := range uSubjects {
		for _, u := range userdomain.Fixture {
			r.Create(u)
		}
	}
}

func TestArticlesCreate_RoundTrips(t *testing.T) {
	t.Parallel()

	na, err := articledomain.NewArticle(
		"misty title",
		"misty description",
		"misty body",
		userdomain.Fixture[0].Email(),
		"misty tag 1", "misty tag 2", "misty tag 3")
	require.NoError(t, err)

	for k, r := range aSubjects {
		r := r
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			ca, err := r.Create(na)
			require.NoError(t, err)

			assert.Equal(t, na.Slug(), ca.Slug())
			assert.Equal(t, na.Title(), ca.Title())
			assert.Equal(t, na.Description(), ca.Description())
			assert.Equal(t, na.Body(), ca.Body())
			assert.Equal(t, na.AuthorEmail(), ca.Email())
			assert.ElementsMatch(t, na.Tags(), ca.Tags())
			assert.Equal(t, ca.CreatedAtUTC(), ca.UpdatedAtUTC())
			assert.NotEmpty(t, ca.CreatedAtUTC())
			assert.Empty(t, ca.Comments())
			assert.Zero(t, ca.FavoriteCount())

			fa, err := r.GetArticleBySlug(ca.Slug())
			require.NoError(t, err)

			assert.Equal(t, na.Slug(), fa.Slug())
			assert.Equal(t, na.Title(), fa.Title())
			assert.Equal(t, na.Description(), fa.Description())
			assert.Equal(t, na.Body(), fa.Body())
			assert.Equal(t, na.AuthorEmail(), fa.Email())
			assert.ElementsMatch(t, na.Tags(), fa.Tags())
			assert.Equal(t, ca.CreatedAtUTC(), fa.UpdatedAtUTC())
			assert.Equal(t, ca.UpdatedAtUTC(), fa.UpdatedAtUTC())
			assert.Empty(t, fa.Comments())
			assert.Zero(t, fa.FavoriteCount())
		})

	}
}
