package balance

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/dto"
	"github.com/Racuwcka/user-balance.git/internal/lib/api/response"
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
		responseErr(h.log, w, r, "param user_id is empty", http.StatusBadRequest, "invalid request")
		return
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil || userId == 0 {
		responseErr(h.log, w, r, "invalid param user_id", http.StatusUnprocessableEntity, "invalid user_id")
		return
	}

	balance, err := h.provider.Balance(context.Background(), uint32(userId))
	if err != nil {
		responseErr(h.log, w, r, "get user balance failed", http.StatusInternalServerError, "get balance failed")
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

func responseErr(log *slog.Logger, w http.ResponseWriter, r *http.Request, msgLog string, status int, msgError string) {
	log.Error(msgLog)
	render.Status(r, status)
	render.JSON(w, r, response.Error(msgError))
}
