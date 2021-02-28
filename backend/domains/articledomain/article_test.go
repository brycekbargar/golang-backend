package articledomain_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArticle(t *testing.T) {
	t.Parallel()
	f := articledomain.Fixture

	cases := []struct {
		Name          string
		Title         string
		Description   string
		Body          string
		Author        string
		Tags          []string
		ExpectedError error
	}{
		{
			"Created (with slug)",
			f[0].Title(),
			f[0].Description(),
			f[0].Body(),
			f[0].AuthorEmail(),
			f[0].Tags(),
			nil,
		},
		{
			"Invalid Slug",
			"*$*#()%)",
			f[0].Description(),
			f[0].Body(),
			f[0].AuthorEmail(),
			f[0].Tags(),
			articledomain.ErrInvalidSlug,
		},
		{
			"Missing Title",
			"",
			f[0].Description(),
			f[0].Body(),
			f[0].AuthorEmail(),
			f[0].Tags(),
			articledomain.ErrRequiredArticleFields,
		},
		{
			"Missing Description",
			f[0].Title(),
			"",
			f[0].Body(),
			f[0].AuthorEmail(),
			f[0].Tags(),
			articledomain.ErrRequiredArticleFields,
		},
		{
			"Missing Body",
			f[0].Title(),
			f[0].Description(),
			"",
			f[0].AuthorEmail(),
			f[0].Tags(),
			articledomain.ErrRequiredArticleFields,
		},
		{
			"Created (with no tags)",
			f[0].Title(),
			f[0].Description(),
			f[0].Body(),
			f[0].AuthorEmail(),
			make([]string, 0),
			nil,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			a, err := articledomain.NewArticle(tc.Title, tc.Description, tc.Body, tc.Author, tc.Tags...)
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.Title, a.Title())
			assert.Equal(t, tc.Description, a.Description())
			assert.Equal(t, tc.Body, a.Body())
			assert.Equal(t, tc.Tags, a.Tags())
			require.NotEmpty(t, a.Slug())
		})
	}
}

func TestUpdateArticle(t *testing.T) {
	t.Parallel()
	f := articledomain.Fixture

	cases := []struct {
		Name                string
		UpdatedTitle        string
		ExpectedTitle       string
		ExpectedSlug        string
		UpdatedDescription  string
		ExpectedDescription string
		UpdatedBody         string
		ExpectedBody        string
		ExpectedError       error
	}{
		{
			"All New Values",
			"whimsical title",
			"whimsical title",
			"whimsical-title",
			"whimsical description",
			"whimsical description",
			"whimsical body",
			"whimsical body",
			nil,
		},
		{
			"No New Values",
			"",
			f[2].Title(),
			f[2].Slug(),
			"",
			f[2].Description(),
			"",
			f[2].Body(),
			nil,
		},
		{
			"Invalid Slug",
			"%*($)",
			"n/a",
			"n/a",
			"",
			"n/a",
			"",
			"n/a",
			articledomain.ErrInvalidSlug,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			a, err := articledomain.UpdatedArticle(
				*f[2],
				tc.UpdatedTitle,
				tc.UpdatedDescription,
				tc.UpdatedBody,
			)

			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.ExpectedTitle, a.Title())
			assert.Equal(t, tc.ExpectedSlug, a.Slug())
			assert.Equal(t, tc.ExpectedDescription, a.Description())
			assert.Equal(t, tc.ExpectedBody, a.Body())
			assert.True(t, a.UpdatedAtUTC().After(a.CreatedAtUTC()))
			assert.Equal(t, f[2].UpdatedAtUTC(), f[2].CreatedAtUTC(),
				"because the original should be unmodified during update")
		})
	}
}

func TestAddComment(t *testing.T) {
	t.Parallel()
	f := articledomain.Fixture
	err := f[0].AddComment("mysterious comment", "mysterious author")

	require.NoError(t, err)

	assert.Equal(t, 2, len(f[0].Comments()))
	for _, c := range f[0].Comments() {
		// the last comment in this fixture'd article has id 6
		if c.ID() == 7 {
			assert.Equal(t, "mysterious comment", c.Body())
			assert.Equal(t, "mysterious author", c.AuthorEmail())
			assert.NotEmpty(t, c.CreatedAtUTC())
			assert.Equal(t, c.UpdatedAtUTC(), c.CreatedAtUTC())
			return
		}
	}

	assert.Fail(t, "because created comment wasn't found")
}

func TestRemoveComment(t *testing.T) {
	t.Parallel()
	f := articledomain.Fixture
	f[1].RemoveComment(4)
	assert.Equal(t, 2, len(f[1].Comments()))
}

func TestIsFavoriteOf(t *testing.T) {
	t.Parallel()

	f := articledomain.Fixture
	assert.True(t, f[1].IsAFavoriteOf("complex@parched.com"))
	f[1].Unfavorite("complex@parched.com")
	assert.False(t, f[1].IsAFavoriteOf("complex@parched.com"))
	assert.Equal(t, 2, f[1].FavoriteCount())
	f[1].Favorite("lucky@unsuitable.com")
	assert.True(t, f[1].IsAFavoriteOf("lucky@unsuitable.com"))
	assert.Equal(t, 3, f[1].FavoriteCount())
}
