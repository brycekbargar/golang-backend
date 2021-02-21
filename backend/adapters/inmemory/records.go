package inmemory

import (
	"sync"
)

var (
	mu = &sync.Mutex{}
	// Users is the instance of an in-memory implementation of the userdomain.Repository
	Users = &users{
		make(map[string]userRecord),
	}
	// Articles is the instance of an in-memory implementation of the articledomain.Repository
	Articles = &articles{
		make(map[string]articleRecord),
	}
)

type userRecord struct {
	email     string
	username  string
	bio       string
	image     string
	following string
	password  string
}

// users is a (super inefficient) in-memory repository implementation for the usersdomain.Repository.
type users struct {
	repo map[string]userRecord
}

type articleRecord struct {
}

// articles is a (super inefficient) in-memory repository implementation for the articledomain.Repository.
type articles struct {
	repo map[string]articleRecord
}
