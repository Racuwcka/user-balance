package dto

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

type DepositRequest struct {
	UserId uint32  `validate:"required,min=1" json:"user_id"`
	Amount float64 `validate:"required,min=1,max=1000000" json:"amount"`
}
