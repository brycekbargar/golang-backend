package inmemory

import (
	"sync"

	"github.com/brycekbargar/realworld-backend/domains/articledomain"
)

type articleRecord struct {
}

// articles is a (super inefficient) in-memory repository implementation for the articledomain.Repository.
type articles struct {
	mu   *sync.Mutex
	repo map[string]articleRecord
}

// NewArticles creates a new articledomain.Repository implementation.
func NewArticles() articledomain.Repository {
	return &articles{
		&sync.Mutex{},
		make(map[string]articleRecord),
	}
}
