package adapters

import (
	"github.com/brycekbargar/realworld-backend/domains/articledomain"
	"github.com/brycekbargar/realworld-backend/domains/userdomain"
)

// RepositoryImplementation is a shared container for the various implementation of the domain repositories.
type RepositoryImplementation struct {
	Name     string
	Users    userdomain.Repository
	Articles articledomain.Repository
}
