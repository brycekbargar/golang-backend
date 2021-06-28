package inmemory

import (
	"sort"
	"strings"
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
)

// Create creates a new article.
func (r *implementation) CreateArticle(a *domain.Article) (*domain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for s := range r.articles {
		if s == strings.ToLower(a.Slug) {
			return nil, domain.ErrDuplicateArticle
		}
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail)]; !ok {
		return nil, domain.ErrNoAuthor
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
	}
	return r.GetArticleBySlug(a.Slug)
}

// LatestArticlesByCriteria lists articles paged/filtered by the given criteria.
func (r *implementation) LatestArticlesByCriteria(query domain.ListCriteria) (
	[]*domain.AuthoredArticle,
	error,
) {
	// i wish this was sql qq
	results := make([]*domain.AuthoredArticle, 0, query.Limit)

	off := 0
	lim := 0

	lf := strings.ToLower(query.FavoritedByUserEmail)
	faveUser, ok := r.users[lf]
	if lf != "" && !ok {
		return nil, domain.ErrUserNotFound
	}

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

		if lf != "" && !strings.Contains(faveUser.favorites, strings.ToLower(ar.slug)) {
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
func (r *implementation) GetArticleBySlug(s string) (*domain.AuthoredArticle, error) {
	if a, ok := r.articles[strings.ToLower(s)]; ok {

		aa, ok := r.users[strings.ToLower(a.author)]
		if !ok {
			return nil, domain.ErrNoAuthor
		}

		fc := 0
		for _, f := range r.users {
			if strings.Contains(f.favorites, strings.ToLower(a.slug)) {
				fc++
			}
		}

		return &domain.AuthoredArticle{
			Article: domain.Article{
				Slug:         a.slug,
				Title:        a.title,
				Description:  a.description,
				Body:         a.body,
				TagList:      strings.Split(a.tagList, ","),
				CreatedAtUTC: a.createdAtUTC,
				UpdatedAtUTC: a.updatedAtUTC,
				AuthorEmail:  a.author,
			},
			Author:        aa,
			FavoriteCount: fc,
		}, nil
	}

	return nil, domain.ErrArticleNotFound
}

// GetCommentsBySlug gets a single article and its comments with the given slug.
func (r *implementation) GetCommentsBySlug(s string) (*domain.CommentedArticle, error) {
	a, err := r.GetArticleBySlug(s)
	if err != nil {
		return nil, err
	}

	if ar, ok := r.articles[strings.ToLower(s)]; ok {
		cs := make([]*domain.Comment, 0, len(ar.comments))
		for _, c := range ar.comments {
			cs = append(cs, &domain.Comment{
				ID:           c.id,
				Body:         c.body,
				CreatedAtUTC: c.createdAtUTC,
				AuthorEmail:  c.author,
			})
		}

		return &domain.CommentedArticle{
			Article:  a.Article,
			Comments: cs,
		}, nil
	}
	return nil, domain.ErrArticleNotFound
}

// UpdateArticleBySlug finds a single article based on its slug
// then applies the provide mutations.
func (r *implementation) UpdateArticleBySlug(s string, update func(*domain.Article) (*domain.Article, error)) (*domain.AuthoredArticle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := r.GetArticleBySlug(s)
	if err != nil {
		return nil, err
	}
	prevSlug := strings.ToLower(f.Slug)

	a, err := update(&f.Article)
	if err != nil {
		return nil, err
	}

	if _, ok := r.users[strings.ToLower(a.AuthorEmail)]; !ok {
		return nil, domain.ErrNoAuthor
	}

	removed := r.articles[strings.ToLower(s)]
	delete(r.articles, strings.ToLower(s))

	for s := range r.articles {
		if s == strings.ToLower(a.Slug) {

			// Add the deleted article back if they've become a duplicate
			r.articles[strings.ToLower(removed.slug)] = removed
			return nil, domain.ErrDuplicateArticle
		}
	}

	if strings.ToLower(a.Slug) != prevSlug {
		for _, v := range r.users {
			// Make sure users favoriting this one get an updated key
			v.favorites = strings.ReplaceAll(v.favorites, prevSlug, strings.ToLower(a.Slug))
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
	}

	return r.GetArticleBySlug(a.Slug)
}

// UpdateCommentsBySlug finds a single article based on its slug
// then applies the provide mutations to its comments.
func (r *implementation) UpdateCommentsBySlug(s string, update func(*domain.CommentedArticle) (*domain.CommentedArticle, error)) (*domain.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := r.GetCommentsBySlug(s)
	if err != nil {
		return nil, err
	}

	a, err := update(f)
	if err != nil {
		return nil, err
	}

	id := 0
	for _, c := range a.Comments {
		if c.ID > id {
			id = c.ID
		}
	}

	ncs := make([]*domain.Comment, 0, len(a.Comments))
	cs := make([]*commentRecord, 0, len(a.Comments))
	for _, c := range a.Comments {
		if c.ID == 0 {
			id++
			c.ID = id
			c.CreatedAtUTC = time.Now().UTC()
			ncs = append(ncs, c)
		}
		cs = append(cs, &commentRecord{
			id:           c.ID,
			body:         c.Body,
			createdAtUTC: c.CreatedAtUTC,
			author:       c.AuthorEmail,
		})
	}
	r.articles[strings.ToLower(a.Slug)].comments = cs

	if len(ncs) > 0 {
		return ncs[0], nil
	}

	return nil, nil
}

// DeleteArticleBySlug deletes the article with the provide slug if it exists.
func (r *implementation) DeleteArticle(a *domain.Article) error {
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

// GetAuthorByEmail finds a single author based on their email address or nil if they don't exist.
func (r *implementation) GetAuthorByEmail(e string) domain.Author {
	if a, err := r.GetUserByEmail(e); err == nil {
		return a
	}

	return nil
}
