package remove

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

//go:generate go run github.com/vektra/mockery/v3 --name=URLRemover
type URLRemover interface {
	DeleteURL(ctx context.Context, alias string) error
}

func RemoveURLHandler(removeController URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.RemoveURLHandler"

		log := slog.Default().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias not provided")
			http.Error(w, "alias not provided", http.StatusBadRequest)
			return
		}

		err := removeController.DeleteURL(r.Context(), alias)
		if err != nil {
			log.Error("failed to remove url", slog.String("err", err.Error()))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Info("url removed", "alias", alias)

		w.WriteHeader(http.StatusNoContent)
	}
}
