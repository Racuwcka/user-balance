package postgres

import (
	"github.com/Racuwcka/shorter-url/pkg/client/postgresql"
	"github.com/Racuwcka/user-balance.git/storage"
)

type Repo struct {
	client postgresql.Client
}

var _ storage.Storage = (*Repo)(nil)

func New(c postgresql.Client) *Repo {
	return &Repo{
		client: c,
	}
}

func (r *Repo) Deposit(userId uint32, serviceId uint16, amount uint) error {
	return nil
}
