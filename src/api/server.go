package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/F0rzend/simbirsoft_contest/src/common"
	"github.com/F0rzend/simbirsoft_contest/src/domain/account"
)

const JSONContentType = "application/json"

type Server struct {
	account *account.Handlers
}

func NewServer() *Server {
	return &Server{
		account: account.NewHandlers(),
	}
}

func (s *Server) GetHTTPHandler(logger *zerolog.Logger) http.Handler {
	translatedValidator := common.NewTranslatedValidator()

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
		middleware.AllowContentType(JSONContentType),
		render.SetContentType(render.ContentTypeJSON),

		hlog.NewHandler(*logger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Send()
		}),
		hlog.RemoteAddrHandler("ip"),
		hlog.RequestIDHandler("req_id", "Request-Id"),

		common.TranslatedValidatorCtxMiddleware(translatedValidator),
	)

	r.Get("/healthcheck", healthcheck)

	r.Post("/registration", s.account.Registration)

	return r
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
