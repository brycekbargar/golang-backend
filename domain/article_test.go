package domain_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewArticle(t *testing.T) {
	t.Parallel()

	t.Run("Title is slugified", func(t *testing.T) {
		t.Parallel()

		a, err := domain.NewArticle(
			"tough title",
			"tough description",
			"tough body",
			"author@tough.com",
			"tough tag 1",
			"tough tag 2")
		require.NoError(t, err)
		assert.NotEqual(t, "tough title", a.Slug)

	})

	t.Run("Validation happens", func(t *testing.T) {
		t.Parallel()

		a, err := domain.NewArticle(
			"purple title",
			"purple description",
			"purple body",
			"purple is not an email",
			"purple tag 1",
			"purple tag 2")
		assert.Error(t, err)
		assert.Nil(t, a)
	})
}

func TestArticle_SetTitle(t *testing.T) {
	t.Parallel()

	a, err := domain.NewArticle(
		"fine title",
		"fine description",
		"fine body",
		"author@fine.com",
		"fine tag 1",
		"fine tag 2")
	require.NoError(t, err)
	assert.Equal(t, "fine-title", a.Slug)

	a.SetTitle("puzzling title")
	assert.Equal(t, "puzzling-title", a.Slug)
}
func TestArticle_Validate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name    string
		Article *domain.Article
	}{
		{
			"Invalid Slug",
			&domain.Article{
				Slug:        "%#&@*( def not a wanting slug",
				Title:       "wanting title",
				Description: "wanting description",
				Body:        "wanting body",
				AuthorEmail: "author@wanting.com",
			},
		},
		{
			"Missing Slug",
			&domain.Article{
				Title:       "garrulous title",
				Description: "garrulous description",
				Body:        "garrulous body",
				AuthorEmail: "author@garrulous.com",
			},
		},
		{
			"Missing Title",
			&domain.Article{
				Slug:        "feeble-slug",
				Description: "feeble description",
				Body:        "feeble body",
				AuthorEmail: "author@feeble.com",
			},
		},
		{
			"Missing Description",
			&domain.Article{
				Slug:        "majestic-slug",
				Title:       "majestic title",
				Body:        "majestic body",
				AuthorEmail: "author@majestic.com",
			},
		},
		{
			"Missing Body",
			&domain.Article{
				Slug:        "thundering-slug",
				Title:       "thundering title",
				Description: "thundering description",
				AuthorEmail: "author@thundering.com",
			},
		},
		{
			"Missing Author",
			&domain.Article{
				Slug:        "noxious-slug",
				Title:       "noxious title",
				Description: "noxious description",
				Body:        "noxious body",
			},
		},
		{
			"Invalid Author",
			&domain.Article{
				Slug:        "daffy-slug",
				Title:       "daffy title",
				Description: "daffy description",
				Body:        "daffy body",
				AuthorEmail: "not a daffy email",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			u, err := tc.Article.Validate()
			assert.Error(t, err)
			assert.Nil(t, u)
		})
	}

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		a := &domain.Article{
			Slug:        "perpetual-slug",
			Title:       "perpetual title",
			Description: "perpetual description",
			Body:        "perpetual body",
			AuthorEmail: "author@perpetual.com",
		}

		va, err := a.Validate()
		require.NoError(t, err)
		assert.Same(t, a, va)
	})
}

func TestArticle_Comments(t *testing.T) {
	t.Parallel()

	ca := domain.CommentedArticle{
		Article: domain.Article{},
		Comments: []domain.Comment{
			{ID: 5},
			{ID: 8},
			{ID: 13},
			{ID: 21},
		},
	}

	err := ca.AddComment("mysterious title", "author@mysterious.com")
	require.NoError(t, err)
	assert.Len(t, ca.Comments, 5)

	err = ca.AddComment("", "")
	assert.Error(t, err)
	assert.Len(t, ca.Comments, 5)

	ca.RemoveComment(8)
	assert.Len(t, ca.Comments, 4)
	ca.RemoveComment(8)
	assert.Len(t, ca.Comments, 4)
}
