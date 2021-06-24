package domain_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComment_Validate(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name    string
		Comment *domain.Comment
	}{
		{
			"Invalid ID",
			&domain.Comment{
				ID:          -5,
				Body:        "gigantic body",
				AuthorEmail: "author@gigantic.com",
			},
		},
		{
			"Missing Body",
			&domain.Comment{
				ID:          5,
				AuthorEmail: "author@innate.com",
			},
		},
		{
			"Missing Author",
			&domain.Comment{
				ID:   5,
				Body: "futuristic body",
			},
		},
		{
			"Invalid Author",
			&domain.Comment{
				ID:          5,
				Body:        "mountainous body",
				AuthorEmail: "not a mountainous email",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			u, err := tc.Comment.Validate()
			assert.Error(t, err)
			assert.Nil(t, u)
		})
	}

	t.Run("Ok", func(t *testing.T) {
		t.Parallel()

		c := &domain.Comment{
			ID:          0,
			Body:        "coherent body",
			AuthorEmail: "user@coherent.com",
		}

		vc, err := c.Validate()
		require.NoError(t, err)
		assert.Same(t, c, vc)
	})
}
