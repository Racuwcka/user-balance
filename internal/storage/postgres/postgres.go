package postgres

import (
	"context"
	"errors"

	"github.com/Racuwcka/user-balance.git/internal/storage"
	"github.com/Racuwcka/user-balance.git/pkg/client/postgresclient"
	"github.com/jackc/pgx/v5"
)

type TransactionType string

const (
	TransactionDeposit TransactionType = "deposit"
	TransactionReserve TransactionType = "reserve"
	TransactionRevenue TransactionType = "revenue"
)

type ReserveStatus string

const (
	Success  ReserveStatus = "success"
	Wait     ReserveStatus = "wait"
	Rejected ReserveStatus = "rejected"
)

type Repo struct {
	client postgresclient.Client
}

var _ storage.Storage = (*Repo)(nil)

func New(c postgresclient.Client) *Repo {
	return &Repo{
		client: c,
	}
}

func (r *Repo) Deposit(ctx context.Context, userId uint32, amount float64) (float64, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var newBalance float64

	q1 := `
	INSERT INTO public.balances (user_id, balance)
	VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET balance = public.balances.balance + EXCLUDED.balance
	RETURNING balance
	`

	if err = tx.QueryRow(ctx, q1, userId, amount).Scan(&newBalance); err != nil {
		return 0, err
	}

	q2 := `
	INSERT INTO public.transactions (user_id, amount, type)
	VALUES ($1, $2, $3)
	`

	if _, err = tx.Exec(ctx, q2, userId, amount, TransactionDeposit); err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (r *Repo) Reserve(ctx context.Context, userId uint32, serviceId uint32, orderId uint64, amount float64) (float64, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var balance float64

	q1 := `SELECT balance FROM public.balances WHERE user_id=$1`

	if err = tx.QueryRow(ctx, q1, userId).Scan(&balance); err != nil {
		return 0, err
	}

	if balance < amount {
		return 0, storage.ErrInsufficientBalance
	}

	newBalance := balance - amount

	q2 := `UPDATE public.balances SET balance=$1 WHERE user_id=$2`

	if _, err = tx.Exec(ctx, q2, newBalance, userId); err != nil {
		return 0, err
	}

	q3 := `
	INSERT INTO public.reserves (user_id, service_id, order_id, amount, status)
	VALUES ($1, $2, $3, $4, $5)
    `

	if _, err = tx.Exec(ctx, q3, userId, serviceId, orderId, amount, Wait); err != nil {
		return 0, err
	}

	q4 := `
	INSERT INTO public.transactions (user_id, service_id, order_id, amount, type)
	VALUES ($1, $2, $3, $4, $5)
	`

	if _, err = tx.Exec(ctx, q4, userId, serviceId, orderId, amount, TransactionReserve); err != nil {
		return 0, err
	}

	return newBalance, nil
}

func (r *Repo) Balance(ctx context.Context, userId uint32) (float64, error) {
	var balance float64

	q := `SELECT balance FROM public.balances WHERE user_id=$1`

	if err := r.client.QueryRow(ctx, q, userId).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return balance, nil
}

func (r *Repo) Revenue(ctx context.Context, userId uint32, serviceId uint32, orderId uint64, amount float64) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return nil
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var reserveAmount float64

	q1 := `
	SELECT amount FROM public.reserves
	WHERE user_id=$1 AND service_id=$2 AND order_id=$3 AND status=$4
	FOR UPDATE
	`

	if err = tx.QueryRow(ctx, q1, userId, serviceId, orderId, Wait).Scan(&reserveAmount); err != nil {
		return storage.ErrReserveNotFound
	}

	if amount != reserveAmount {
		return storage.ErrReserveMismatch
	}

	q2 := `
	UPDATE public.reserves SET status=$1
	WHERE user_id=$2 AND service_id=$3 AND order_id=$4 AND status=$5
	`

	if _, err = tx.Exec(ctx, q2, Success, userId, serviceId, orderId, Wait); err != nil {
		return err
	}

	q3 := `
	INSERT INTO public.transactions (user_id, service_id, order_id, amount, type)
	VALUES ($1, $2, $3, $4, $5)
	`

	if _, err = tx.Exec(ctx, q3, userId, serviceId, orderId, amount, TransactionRevenue); err != nil {
		return nil
	}

	return nil
}
