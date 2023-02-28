package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/F0rzend/simbirsoft_contest/src/domain/account"
)

type Server struct {
	account *account.Handlers
}

func NewServer() *Server {
	return &Server{
		account: account.NewHandlers(),
	}
}

func (s *Server) GetHTTPHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/healthcheck", healthcheck)

	return r
}

func healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
