package types

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Handlers struct {
	service *Service
}

func NewHandlers(di *common.DependencyInjectionContainer) *Handlers {
	return &Handlers{
		service: NewService(di),
	}
}

func (h *Handlers) GetAnimalType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	animalTypeID, err := common.GetInt64FromRequest(r, "type_id")
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	entity, err := h.service.GetAnimalType(ctx, animalTypeID)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	response := &Response{
		ID:   entity.ID,
		Type: entity.Type,
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
		return
	}
}
