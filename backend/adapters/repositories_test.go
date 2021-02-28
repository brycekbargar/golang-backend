package adapters_test

import (
	"testing"

	"github.com/brycekbargar/realworld-backend/adapters"
	"github.com/brycekbargar/realworld-backend/adapters/inmemory"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

func subjects(t *testing.T) map[string]*adapters.RepositoryImplementation {
	t.Helper()
	subjects := map[string]*adapters.RepositoryImplementation{
		"inmemory": inmemory.NewInstance(),
	}

	for _, s := range subjects {
		for _, u := range userdomain.Fixture {
			s.Users.CreateUser(u)
		}
	}

	return subjects
}
