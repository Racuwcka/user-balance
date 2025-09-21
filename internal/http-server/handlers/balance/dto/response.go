package dto

type DepositResponse struct {
	UserId  uint32  `json:"user_id"`
	Balance float64 `json:"balance"`
}
