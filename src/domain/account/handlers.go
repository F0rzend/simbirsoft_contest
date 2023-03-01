package account

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Handlers struct {
	service *Service
}

func NewHandlers() *Handlers {
	return &Handlers{
		service: NewService(),
	}
}

func (h *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	request := new(RegistrationRequest)

	if err := render.Bind(r, request); err != nil {
		common.RenderError(w, r, err)
	}
	log.Debug().
		Int("problem_content_type", int(render.GetContentType("application/problem+json"))).
		Int("request_content_type", int(render.GetRequestContentType(r))).
		Send()
}
