package dto

import (
	"github.com/Racuwcka/user-balance.git/internal/lib/api/response"
)

type DepositRequest struct {
	response.Response
	UserId uint32  `validate:"required,min=1" json:"user_id"`
	Amount float64 `validate:"required,min=1,max=1000000" json:"amount"`
}

type TransactionRequest struct {
	response.Response
	UserId    uint32  `validate:"required,min=1" json:"user_id"`
	ServiceId uint32  `validate:"required,min=1" json:"service_id"`
	OrderId   uint64  `validate:"required,min=1" json:"order_id"`
	Amount    float64 `validate:"required,min=1,max=1000000" json:"amount"`
}
