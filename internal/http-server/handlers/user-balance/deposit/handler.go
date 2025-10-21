package deposit

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/dto"
	"github.com/Racuwcka/user-balance.git/internal/lib/api/response"
	"github.com/Racuwcka/user-balance.git/internal/lib/logger/sl"
)

type provider interface {
	Deposit(ctx context.Context, userId uint32, amount float64) (float64, error)
}

type Handler struct {
	log         *slog.Logger
	validateErr validator.ValidationErrors
	provider    provider
}

func New(log *slog.Logger, p provider) *Handler {
	return &Handler{
		log:      log,
		provider: p,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", "handler.balance.deposit"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	req := &dto.DepositRequest{}
	if err := render.DecodeJSON(r.Body, req); err != nil {
		if errors.Is(err, io.EOF) {
			h.log.Error("request body is empty")
			render.JSON(w, r, response.Error("empty request"))
			return
		}

		h.log.Error("failed to decode request body", sl.Err(err))
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	h.log.Info("request body decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		errors.As(err, &h.validateErr)

		h.log.Error("invalid request", sl.Err(err))
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, response.ValidationError(h.validateErr))
		return
	}

	balance, err := h.provider.Deposit(context.Background(), req.UserId, req.Amount)
	if err != nil {
		h.log.Error("funds deposit failed", sl.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("funds deposit failed"))
		return
	}

	responseOK(w, r, req.UserId, balance)
}

func responseOK(w http.ResponseWriter, r *http.Request, userId uint32, balance float64) {
	render.JSON(w, r, dto.UserResponse{
		Response: response.OK(),
		UserId:   userId,
		Balance:  balance,
	})
}
