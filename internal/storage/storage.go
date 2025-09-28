package storage

import "errors"

var (
	ErrInsufficientBalance = errors.New("insufficient balance for reservation")
	ErrReserveNotFound     = errors.New("reserve not found")
	ErrReserveMismatch     = errors.New("reserve amount mismatch")
)
