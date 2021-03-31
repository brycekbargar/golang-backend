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
func Articles_DeleteArticle(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	r.Users.CreateUser(testAuthor("deranged"))
	r.Articles.CreateArticle(testArticle("deranged"))

	a, err := r.Articles.GetArticleBySlug("deranged-title")
	require.NoError(t, err)
	err = r.Articles.DeleteArticle(&a.Article)
	assert.NoError(t, err)
	a, err = r.Articles.GetArticleBySlug("deranged-title")
	assert.ErrorIs(t, err, articledomain.ErrNotFound)

	err = r.Articles.DeleteArticle(nil)
	assert.NoError(t, err)
}

func Articles_LatestArticlesByCriteria(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	tt := "Articles_LatestArticlesByCriteria"

	authors := []*userdomain.User{
		testAuthor("frail"),
		testAuthor("simple"),
		testAuthor("reminiscent"),
	}
	for _, a := range authors {
		_, err := r.Users.CreateUser(a)
		require.NoError(t, err)
	}

	source := make([]*articledomain.AuthoredArticle, 0, 13)
	for i, adj := range []string{
		"bright",
		"colorful",
		"sour",
		"fantastic",
		"mellow",
		"splendid",
		"gruesome",
		"madly",
		"kind",
		"organic",
		"public",
		"flagrant",
		"waggish",
	} {
		a := testArticle(adj)
		a.AuthorEmail = authors[i%3].Email
		a.TagList = append(a.TagList, tt)

		_, err := r.Articles.CreateArticle(a)
		require.NoError(t, err)

		fa, err := r.Articles.GetArticleBySlug(a.Slug)
		require.NoError(t, err)
		source = append(source, fa)
	}

	cases := []struct {
		Name  string
		Query articledomain.ListCriteria
		Try   func([]*articledomain.AuthoredArticle, error)
	}{
		{
			"All Tagged",
			articledomain.ListCriteria{
				Tag:    tt,
				Offset: 0,
				Limit:  13,
			},
			func(all []*articledomain.AuthoredArticle, err error) {
				require.NoError(t, err)

				assert.ElementsMatch(t, source, all)
			},
		},
		{
			"Offset Limit Sort",
			articledomain.ListCriteria{
				Tag:    tt,
				Offset: 2,
				Limit:  5,
			},
			func(some []*articledomain.AuthoredArticle, err error) {
				require.NoError(t, err)

				assert.Len(t, some, 5)
				assert.Equal(t, "public-title", some[0].Slug)
				assert.Equal(t, "gruesome-title", some[4].Slug)
			},
		},
		{
			"Authored By",
			articledomain.ListCriteria{
				Tag:          tt,
				Offset:       0,
				Limit:        13,
				AuthorEmails: []string{"author@frail.com", "author@simple.com"},
			},
			func(authored []*articledomain.AuthoredArticle, err error) {
				require.NoError(t, err)

				assert.NotEmpty(t, authored)
				for _, a := range authored {
					assert.NotEqual(t, "author@reminiscent.com", a.AuthorEmail)
				}
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.Try(r.Articles.LatestArticlesByCriteria(tc.Query))
		})
	}
}
func Articles_DistinctTags(
	t *testing.T,
	r *adapters.RepositoryImplementation,
) {
	tt := "Articles_DistinctTags"

	r.Users.CreateUser(testAuthor("threatening"))
	for _, adj := range []string{
		"stimulating",
		"exultant",
		"helpless",
	} {
		a := testArticle(adj)
		a.AuthorEmail = "author@threatening.com"
		a.TagList = append(a.TagList, tt)

		_, err := r.Articles.CreateArticle(a)
		require.NoError(t, err)
	}

	tags, err := r.Articles.DistinctTags()
	require.NoError(t, err)

	assert.Contains(t, tags, "stimulating one", "stimulating two")
	assert.Contains(t, tags, "helpless three")
	htt := false
	for i := range tags {
		if tags[i] == tt {
			if htt {
				assert.Fail(t, "the testing tag was not distinct")
			}
			htt = true
		}
	}
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
