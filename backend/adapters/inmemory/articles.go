package inmemory

import (
	"sort"
	"strings"
	"time"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
)

// Create creates a new article.
func (r *implementation) CreateArticle(a *articledomain.Article) (*articledomain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for s := range r.articles {
		if s == strings.ToLower(a.Slug) {
			return nil, articledomain.ErrDuplicateValue
		}
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail)]; !ok {
		return nil, articledomain.ErrNoAuthor
	}

	now := time.Now().UTC()
	r.articles[strings.ToLower(a.Slug)] = &articleRecord{
		a.Slug,
		a.Title,
		a.Description,
		a.Body,
		strings.Join(a.TagList, ","),
		now,
		now,
		a.AuthorEmail,
		make([]*commentRecord, 0),
		map[string]interface{}{},
	}
	return r.GetArticleBySlug(a.Slug)
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(query articledomain.ListCriteria) (
	[]*articledomain.AuthoredArticle,
	error,
) {
	// i wish this was sql qq
	results := make([]*articledomain.AuthoredArticle, 0, query.Limit)

	off := 0
	lim := 0

	lf := strings.ToLower(query.FavoritedByUserEmail)
	lt := strings.ToLower(query.Tag)
	am := make(map[string]interface{}, len(query.AuthorEmails))
	for _, ae := range query.AuthorEmails {
		am[strings.ToLower(ae)] = nil
	}

	ordered := make([]*articleRecord, 0, len(r.articles))
	for _, ar := range r.articles {
		ordered = append(ordered, ar)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].createdAtUTC.After(ordered[j].createdAtUTC)
	})

	for _, ar := range ordered {
		if off < query.Offset {
			off++
			continue
		}

		_, a := am[ar.author]
		if len(query.AuthorEmails) > 0 && !a {
			continue
		}

		_, f := ar.favoritedBy[lf]
		if lf != "" && !f {
			continue
		}

		if lt != "" && !strings.Contains(strings.ToLower(ar.tagList), lt) {
			continue
		}

		da, err := r.GetArticleBySlug(ar.slug)
		if err != nil {
			continue
		}
		results = append(results, da)

		lim++
		if lim >= query.Limit {
			break
		}
	}

	return results, nil
}

// GetArticleBySlug gets a single article with the given slug.
func (r *implementation) GetArticleBySlug(s string) (*articledomain.AuthoredArticle, error) {
	if a, ok := r.articles[strings.ToLower(s)]; ok {

		aa, ok := r.users[strings.ToLower(a.author)]
		if !ok {
			return nil, articledomain.ErrNoAuthor
		}

		return &articledomain.AuthoredArticle{
			Article: articledomain.Article{
				Slug:         a.slug,
				Title:        a.title,
				Description:  a.description,
				Body:         a.body,
				TagList:      strings.Split(a.tagList, ","),
				CreatedAtUTC: a.createdAtUTC,
				UpdatedAtUTC: a.updatedAtUTC,
				AuthorEmail:  a.author,
				FavoritedBy:  a.favoritedBy,
			},
			Author: aa,
		}, nil
	}

	return nil, articledomain.ErrNotFound
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(string) (*articledomain.CommentedArticle, error) {
	return nil, nil
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(s string, update func(*articledomain.Article) (*articledomain.Article, error)) (*articledomain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := r.GetArticleBySlug(s)
	if err != nil {
		return nil, err
	}

	a, err := update(&f.Article)
	if err != nil {
		return nil, err
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail)]; !ok {
		return nil, articledomain.ErrNoAuthor
	}

	removed := r.articles[strings.ToLower(s)]
	delete(r.articles, strings.ToLower(s))

	for s := range r.articles {
		if s == strings.ToLower(a.Slug) {

			// Add the deleted article back if they've become a duplicate
			r.articles[strings.ToLower(removed.slug)] = removed
			return nil, articledomain.ErrDuplicateValue
		}
	}

	now := time.Now().UTC()
	r.articles[strings.ToLower(a.Slug)] = &articleRecord{
		a.Slug,
		a.Title,
		a.Description,
		a.Body,
		strings.Join(a.TagList, ","),
		a.CreatedAtUTC,
		now,
		a.AuthorEmail,
		make([]*commentRecord, 0),
		map[string]interface{}{},
	}

	return r.GetArticleBySlug(a.Slug)
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(string, func(*articledomain.CommentedArticle) (*articledomain.CommentedArticle, error)) (*articledomain.Comment, error) {
	return nil, nil
}

// DeleteArticleBySlug deletes the article with the provide slug if it exists.
func (r *implementation) DeleteArticle(a *articledomain.Article) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if a == nil {
		return nil
	}

	delete(r.articles, strings.ToLower(a.Slug))
	return nil
}

// DistinctTags returns a distinct list of tags on articles
func (r *implementation) DistinctTags() ([]string, error) {
	tm := make(map[string]interface{})
	for _, ar := range r.articles {
		for _, t := range strings.Split(ar.tagList, ",") {
			tm[strings.ToLower(t)] = nil
		}
	}

	tags := make([]string, 0, len(tm))
	for t := range tm {
		tags = append(tags, t)
	}

	return tags, nil
}
