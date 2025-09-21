package deposit

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Racuwcka/user-balance.git/internal/http-server/handlers"
	"github.com/Racuwcka/user-balance.git/internal/http-server/handlers/balance/dto"
)

type provider interface {
	Deposit(ctx context.Context, userId uint32, amount float64) (float64, error)
}

type Handler struct {
	provider provider
}

func New(p provider) *Handler {
	return &Handler{
		provider: p,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &dto.DepositRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.InvalidRequest(w, err)
	}

	err := dto.Validate.Struct(req)
	if err != nil {
		handlers.InvalidValidation(w, err)
	}

	balance, err := h.provider.Deposit(context.Background(), req.UserId, req.Amount)
	if err != nil {
		handlers.GetErrorResponse(w, "funds transfer error", err, http.StatusInternalServerError)
		return
	}

	res := dto.DepositResponse{
		UserId:  req.UserId,
		Balance: balance,
	}

	raw, err := json.Marshal(res)
	if err != nil {
		handlers.InvalidResponse(w, err)
	}

	handlers.GetSuccessResponse(w, raw)
}
