package inmemory

import (
	"sync"
	"time"

	"github.com/brycekbargar/realworld-backend/adapters"
)

// NewInstance creates a new instance of the In-Memory store with the repository interface implementations
func NewInstance() *adapters.RepositoryImplementation {
	i := &implementation{
		&sync.Mutex{},
		make(map[string]userRecord),
		make(map[string]articleRecord),
	}
	return &adapters.RepositoryImplementation{Users: i, Articles: i}
}

type implementation struct {
	mu       *sync.Mutex
	users    map[string]userRecord
	articles map[string]articleRecord
}

type userRecord struct {
	email     string
	username  string
	bio       string
	image     string
	following string
	password  []byte
}

type articleRecord struct {
	slug         string
	title        string
	description  string
	body         string
	tagList      []string
	createdAtUTC time.Time
	updatedAtUTC time.Time
	author       string
	comments     []*commentRecord
	favoritedBy  map[string]interface{}
}

type commentRecord struct {
	id           int
	body         string
	createdAtUTC time.Time
	author       string
}

// articles is a (super inefficient) in-memory repository implementation for the articledomain.Repository.
type articles struct {
}
