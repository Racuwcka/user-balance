package postgres

import (
	"context"

	"github.com/Racuwcka/shorter-url/pkg/client/postgresql"
	"github.com/Racuwcka/user-balance.git/internal/storage"
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

func (r *Repo) Deposit(ctx context.Context, userId uint32, amount float64) (float64, error) {
	var newBalance float64

	q := `
	INSERT INTO public.balances (user_id, amount)
	VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET amount = public.balances.amount + EXCLUDED.amount
	RETURNING amount
	`

	err := r.client.QueryRow(ctx, q, userId, amount).Scan(&newBalance)
	if err != nil {
		return 0, err
	}

	return newBalance, nil
}
