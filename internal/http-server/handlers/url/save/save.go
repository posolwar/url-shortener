package save

import (
	// для краткости даем короткий алиас пакету

	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/pkg/logger/sl"
	"url-shortener/pkg/random"
	"url-shortener/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

// TODO: move to config when needed
const aliasLength = 6

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v3 --name=URLSaver
type URLSaver interface {
	SaveURL(ctx context.Context, url *url.URL, alias string) (id int64, err error)
}

func AliasSaveHandler(URLSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.AliasSaveHandler"

		// Добавляем к текущму объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log := slog.Default().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаем объект запроса и анмаршаллим в него запрос
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "empty request",
			})

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "failed to decode request",
			})

			return
		}

		// Лучше больше логов, чем меньше - лишнее мы легко сможем почистить,
		// при необходимости. А вот недостающую информацию мы уже не получим.
		log.Debug("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаем в него структуру, которую нужно провалидировать
		if err := validator.New().Struct(req); err != nil {
			// Приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.Error(validateErr.Error()))

			return
		}

		// Генерация алиаса
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		reqURL, err := url.Parse(req.URL)
		if err != nil {
			log.Error("failed to parse url", sl.Err(err))

			render.JSON(w, r, resp.Error("invalid url"))

			return
		}

		id, err := URLSaver.SaveURL(r.Context(), reqURL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			// Отдельно обрабатываем ситуацию,
			// когда запись с таким Alias уже существует
			log.Info("url already exists", slog.String("url", reqURL.String()))
			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
