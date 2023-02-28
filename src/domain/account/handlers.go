package account

import "net/http"

type Handlers struct {
	service *Service
}

func NewHandlers() *Handlers {
	return &Handlers{
		service: NewService(),
	}
}

func (h *Handlers) Registration(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
