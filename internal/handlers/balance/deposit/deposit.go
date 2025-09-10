package deposit

import (
	"encoding/json"
	"net/http"

	"github.com/Racuwcka/user-balance.git/internal/handlers"
	"github.com/Racuwcka/user-balance.git/internal/handlers/balance/dto"
)

type depositService interface {
	Deposit(userId uint32, serviceId uint16, amount uint) error
}
type Handler struct {
	depositService depositService
}

func New(d depositService) *Handler {
	return &Handler{
		depositService: d,
	}
}
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &dto.DepositRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.InvalidRequest(w, err, http.StatusBadRequest)
		return
	}

	err := dto.Validate.Struct(req)
	if err != nil {
		handlers.InvalidValidation(w, err, http.StatusUnprocessableEntity)
		return
	}

	handlers.GetSuccessResponse(w, []byte{})
}
