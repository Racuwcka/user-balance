package storage

import "context"

type Storage interface {
	Deposit(ctx context.Context, userId uint32, amount float64) (float64, error)
}
