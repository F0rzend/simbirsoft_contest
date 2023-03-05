package visits

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

func (h *Handlers) GetAnimalLocations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := common.GetInt64FromRequest(r, "animal_id")
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	params, err := NewGetAnimalLocationsQuery(r)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	entities, err := h.service.GetAnimalLocations(
		ctx,
		id,
		params.StartDateTime,
		params.EndDateTime,
		params.From,
		params.Size,
	)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	response := make(ResponsesList, len(entities))
	for i, entity := range entities {
		response[i] = Response{
			ID:                           entity.ID,
			DateTimeOfVisitLocationPoint: entity.DateTimeOfVisitLocationPoint,
			LocationPointID:              entity.LocationPointID,
		}
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
		return
	}
}
