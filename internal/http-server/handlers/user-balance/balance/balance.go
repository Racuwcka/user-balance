package balance

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/dto"
	"github.com/Racuwcka/user-balance.git/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type provider interface {
	Balance(ctx context.Context, userId uint32) (float64, error)
}

type Handler struct {
	log      *slog.Logger
	provider provider
}

func New(log *slog.Logger, p provider) *Handler {
	return &Handler{
		log:      log,
		provider: p,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", "handler.balance.get"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	userIdStr := chi.URLParam(r, "user_id")
	if userIdStr == "" {
		h.log.Error("param user_id is empty")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request"))
		return
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil || userId == 0 {
		h.log.Error("invalid param user_id")
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, response.Error("invalid user_id"))
		return
	}

	balance, err := h.provider.Balance(context.Background(), uint32(userId))
	if err != nil {
		h.log.Error("get user balance failed")
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("get balance failed"))
		return
	}

	responseOK(w, r, uint32(userId), balance)
}

func responseOK(w http.ResponseWriter, r *http.Request, userId uint32, balance float64) {
	render.JSON(w, r, dto.UserResponse{
		Response: response.OK(),
		UserId:   userId,
		Balance:  balance,
	})
}
