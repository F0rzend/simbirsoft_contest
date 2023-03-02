package account

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Handlers struct {
	service *Service
}

func NewHandlers(config *common.Config) (*Handlers, error) {
	accountService, err := NewService(config)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		service: accountService,
	}, nil
}

func (h *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request := new(RegistrationRequest)

	if err := render.Bind(r, request); err != nil {
		common.RenderError(w, r, err)
		return
	}

	entity, err := h.service.Register(
		ctx,
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
	)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	response := &RegistrationResponse{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email.Address,
	}

	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
		return
	}
}
