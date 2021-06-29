package postgres_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/adapters/postgres"
	"github.com/brycekbargar/realworld-backend/adapters/testcases"
)

var uut = postgres.MustNewInstance("")

func Test_Users(t *testing.T) {
	t.Parallel()

	t.Run("Create and Update User", func(t *testing.T) {
		t.Parallel()
		testcases.Users_CreateUser(t, uut)
	})
	t.Run("Get and Update User By Email", func(t *testing.T) {
		t.Parallel()
		testcases.Users_GetUserByEmail(t, uut)
	})
	t.Run("Get and Update User By Username", func(t *testing.T) {
		t.Parallel()
		testcases.Users_GetUserByUsername(t, uut)
	})
	t.Run("Fanboy Following Users", func(t *testing.T) {
		t.Parallel()
		testcases.Users_UpdateFanboyByEmail_Following(t, uut)
	})
	t.Run("Fanboy Favoriting Articles", func(t *testing.T) {
		t.Parallel()
		testcases.Users_UpdateFanboyByEmail_Favorites(t, uut)
	})
}

func Test_Articles(t *testing.T) {
	t.Parallel()

	t.Run("Create and Update Article", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_CreateArticle(t, uut)
	})
	t.Run("Get and Update Article", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_GetArticleBySlug(t, uut)
	})
	t.Run("Delete Article", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_DeleteArticle(t, uut)
	})
	t.Run("Query Articles", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_LatestArticlesByCriteria(t, uut)
	})
	t.Run("Create and Delete Comments", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_UpdateCommentsBySlug(t, uut)
	})
	t.Run("Query Tags", func(t *testing.T) {
		t.Parallel()
		testcases.Articles_DistinctTags(t, uut)
	})
}
