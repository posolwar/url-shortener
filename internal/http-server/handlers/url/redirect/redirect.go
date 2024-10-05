package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"url-shortener/pkg/logger/sl"
	"url-shortener/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

//go:generate go run github.com/vektra/mockery/v3 --name=URLGetter
type URLGetter interface {
	GetURL(ctx context.Context, alias string) (*url.URL, error)
}

func RedirectHandler(urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.RedirectHandler"

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log := slog.Default().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			// Не передан alias, сообщаем об этом клиенту
			log.Info("alias not provided")
			http.Error(w, "alias not provided", http.StatusBadRequest)

			return
		}

		u, err := urlGetter.GetURL(r.Context(), alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			// Не нашли URL, сообщаем об этом клиенту
			log.Info("url not found", "alias", alias, "err", err.Error())
			http.Error(w, "not found", http.StatusNotFound)

			return
		}
		if err != nil {
			// Не удалось осуществить поиск
			log.Error("failed to get url", sl.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)

			return
		}

		log.Debug("got url",
			slog.String("alias", alias),
			slog.String("url", u.String()),
		)

		http.Redirect(w, r, u.String(), http.StatusFound)
	}
}
