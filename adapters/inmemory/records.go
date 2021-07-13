package inmemory

import (
	"sync"
	"time"

	"github.com/brycekbargar/realworld-backend/domain"
)

// NewInstance creates a new instance of the In-Memory store with the repository interface implementations
func NewInstance() domain.Repository {
	i := &implementation{
		&sync.Mutex{},
		make(map[string]*userRecord),
		make(map[string]*articleRecord),
	}
	return i
}

type implementation struct {
	mu       *sync.Mutex
	users    map[string]*userRecord
	articles map[string]*articleRecord
}

type userRecord struct {
	email     string
	username  string
	bio       string
	image     string
	following string
	favorites string
	password  []byte
}

func (u userRecord) GetUsername() string {
	return u.username
}
func (u userRecord) GetEmail() string {
	return u.email
}
func (u userRecord) GetBio() string {
	return u.bio
}
func (u userRecord) GetImage() string {
	return u.image
}

type articleRecord struct {
	slug         string
	title        string
	description  string
	body         string
	tagList      string
	createdAtUTC time.Time
	updatedAtUTC time.Time
	author       string
	comments     []commentRecord
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
