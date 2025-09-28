package revenue

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Racuwcka/user-balance.git/internal/http-server/handlers/user-balance/dto"
	"github.com/Racuwcka/user-balance.git/internal/lib/api/response"
	"github.com/Racuwcka/user-balance.git/internal/lib/logger/sl"
	"github.com/Racuwcka/user-balance.git/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type provider interface {
	Revenue(ctx context.Context, userId uint32, serviceId uint32, orderId uint64, amount float64) error
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

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.log = h.log.With(
		slog.String("op", "handler.balance.revenue"),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	req := &dto.TransactionRequest{}
	if err := render.DecodeJSON(r.Body, req); err != nil {
		render.Status(r, http.StatusBadRequest)

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
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		h.log.Error("invalid request", sl.Err(err))
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, response.ValidationError(validateErr))
		return
	}

	if err := h.provider.Revenue(context.Background(), req.UserId, req.ServiceId, req.OrderId, req.Amount); err != nil {
		render.Status(r, http.StatusInternalServerError)

		if errors.Is(err, storage.ErrReserveNotFound) {
			h.log.Error("user reserve is not found", sl.Err(err))
			render.JSON(w, r, response.Error("user reserve is not found"))
			return
		}

		if errors.Is(err, storage.ErrReserveMismatch) {
			h.log.Error("reserve amount mismatch",
				slog.Any("err", err),
				slog.Float64("expected", req.Amount),
			)
			render.JSON(w, r, response.Error("user reserve amount mismatch"))
			return
		}

		h.log.Error("funds reservation failed", sl.Err(err))
		render.JSON(w, r, response.Error("funds reservation failed"))
		return
	}

	render.JSON(w, r, response.OK())
}
