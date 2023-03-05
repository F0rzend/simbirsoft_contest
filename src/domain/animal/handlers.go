package animal

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

func (h *Handlers) GetAnimal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	animalID, err := common.GetInt64FromRequest(r, "animal_id")
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	entity, err := h.service.GetAnimal(ctx, animalID)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	response := &Response{
		ID:                 entity.ID,
		Types:              entity.Types,
		Weight:             entity.Weight,
		Length:             entity.Length,
		Height:             entity.Height,
		Gender:             entity.Gender,
		LifeStatus:         entity.LifeStatus,
		ChippingDateTime:   entity.ChippingDateTime,
		ChipperID:          entity.ChipperID,
		ChippingLocationID: entity.ChippingLocationID,
		VisitedLocations:   entity.VisitedLocations,
		DeathDateTime:      entity.DeathDateTime,
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
		return
	}
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := NewSearchParams(r)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	entities, err := h.service.Search(
		ctx,
		params.StartDateTime,
		params.EndDateTime,
		params.ChipperID,
		params.ChippingLocationID,
		params.LifeStatus,
		params.Gender,
		params.From,
		params.Size,
	)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	response := make(ResponseList, len(entities))

	for i, entity := range entities {
		response[i] = &Response{
			ID:                 entity.ID,
			Types:              entity.Types,
			Weight:             entity.Weight,
			Length:             entity.Length,
			Height:             entity.Height,
			Gender:             entity.Gender,
			LifeStatus:         entity.LifeStatus,
			ChippingDateTime:   entity.ChippingDateTime,
			ChipperID:          entity.ChipperID,
			ChippingLocationID: entity.ChippingLocationID,
			VisitedLocations:   entity.VisitedLocations,
			DeathDateTime:      entity.DeathDateTime,
		}
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
	}
}
