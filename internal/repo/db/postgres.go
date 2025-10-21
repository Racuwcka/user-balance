package db

import (
	"context"
	"errors"
	"math"

	"github.com/jackc/pgx/v5"

	httpserver "github.com/Racuwcka/user-balance.git/internal/http-server"
	"github.com/Racuwcka/user-balance.git/pkg/client/postgresclient"
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

	var newBalance int

	q1 := `
	INSERT INTO public.balances (user_id, balance)
	VALUES ($1, $2)
	ON CONFLICT (user_id)
	DO UPDATE SET balance = public.balances.balance + EXCLUDED.balance
	RETURNING balance
	`

	amountInt := FloatToInt(amount)
	if err = tx.QueryRow(ctx, q1, userId, amountInt).Scan(&newBalance); err != nil {
		return 0, err
	}

	q2 := `
	INSERT INTO public.transactions (user_id, amount, type)
	VALUES ($1, $2, $3)
	`

	if _, err = tx.Exec(ctx, q2, userId, amountInt, TransactionDeposit); err != nil {
		return 0, err
	}

	newBalanceFloat := IntToFloat(newBalance)
	return newBalanceFloat, nil
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

	var balance int

	q1 := `SELECT balance FROM public.balances WHERE user_id=$1 FOR UPDATE`

	if err = tx.QueryRow(ctx, q1, userId).Scan(&balance); err != nil {
		return 0, err
	}

	amountInt := FloatToInt(amount)
	if balance < amountInt {
		return 0, httpserver.ErrInsufficientBalance
	}

	newBalance := balance - amountInt

	q2 := `UPDATE public.balances SET balance=$1 WHERE user_id=$2`

	if _, err = tx.Exec(ctx, q2, newBalance, userId); err != nil {
		return 0, err
	}

	q3 := `
	INSERT INTO public.reserves (user_id, service_id, order_id, amount, status)
	VALUES ($1, $2, $3, $4, $5)
    `

	if _, err = tx.Exec(ctx, q3, userId, serviceId, orderId, amountInt, Wait); err != nil {
		return 0, err
	}

	q4 := `
	INSERT INTO public.transactions (user_id, service_id, order_id, amount, type)
	VALUES ($1, $2, $3, $4, $5)
	`

	if _, err = tx.Exec(ctx, q4, userId, serviceId, orderId, amountInt, TransactionReserve); err != nil {
		return 0, err
	}

	newBalanceFloat := IntToFloat(newBalance)
	return newBalanceFloat, nil
}

func (r *Repo) Balance(ctx context.Context, userId uint32) (float64, error) {
	var balance int

	q := `SELECT balance FROM public.balances WHERE user_id=$1`

	if err := r.client.QueryRow(ctx, q, userId).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	balanceFloat := IntToFloat(balance)
	return balanceFloat, nil
}

func (r *Repo) Revenue(ctx context.Context, userId uint32, serviceId uint32, orderId uint64, amount float64) error {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	var reserveID int
	var reserveAmount int

	q1 := `
	SELECT id, amount
	FROM public.reserves
	WHERE user_id=$1 AND service_id=$2 AND order_id=$3 AND status=$4
	ORDER BY created_at
	LIMIT 1
	FOR UPDATE
	`

	if err = tx.QueryRow(ctx, q1, userId, serviceId, orderId, Wait).Scan(&reserveID, &reserveAmount); err != nil {
		return httpserver.ErrReserveNotFound
	}

	amountInt := FloatToInt(amount)
	if amountInt != reserveAmount {
		return httpserver.ErrReserveMismatch
	}

	q2 := `
	UPDATE public.reserves
	SET status=$1
	WHERE id=$2
	`

	if _, err = tx.Exec(ctx, q2, Success, reserveID); err != nil {
		return err
	}

	q3 := `
	INSERT INTO public.transactions (user_id, service_id, order_id, amount, type)
	VALUES ($1, $2, $3, $4, $5)
	`

	if _, err = tx.Exec(ctx, q3, userId, serviceId, orderId, amountInt, TransactionRevenue); err != nil {
		return err
	}

	return nil
}

func FloatToInt(amount float64) int {
	return int(math.Round(amount * 100))
}

func IntToFloat(amount int) float64 {
	return float64(amount) / 100
}
