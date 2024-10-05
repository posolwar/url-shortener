package router

import (
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/remove"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type UrlController interface {
	redirect.URLGetter
	save.URLSaver
	remove.URLRemover
}

type RouterApp struct {
	UrlController UrlController
	Cfg           *config.Config
}

func NewRouterApp(urlController UrlController, cfg *config.Config) RouterApp {
	return RouterApp{UrlController: urlController, Cfg: cfg}
}

func GetRouter(app RouterApp) *chi.Mux {
	router := chi.NewRouter()
	{
		router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
		router.Use(logger.New)           // Логирование всех запросов
		router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	}

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			app.Cfg.HTTPServer.User: app.Cfg.HTTPServer.Password,
		}))

		r.Delete("/{alias}", remove.RemoveURLHandler(app.UrlController))
		r.Post("/", save.AliasSaveHandler(app.UrlController))
	})

	router.Get("/{alias}", redirect.RedirectHandler(app.UrlController))

	return router
}
