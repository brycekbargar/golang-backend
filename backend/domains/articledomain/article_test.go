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
			assert.NotEmpty(t, a.CreatedAtUTC())
			assert.Equal(t, a.UpdatedAtUTC(), a.CreatedAtUTC())
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
