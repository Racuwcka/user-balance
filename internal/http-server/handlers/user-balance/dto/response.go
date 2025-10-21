package dto

import "github.com/Racuwcka/user-balance.git/internal/lib/api/response"

type UserResponse struct {
	response.Response
	UserId  uint32  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type ReserveResponse struct {
	response.Response
	UserId   uint32  `json:"user_id"`
	Reserved float64 `json:"reserved"`
	Balance  float64 `json:"balance"`
}
