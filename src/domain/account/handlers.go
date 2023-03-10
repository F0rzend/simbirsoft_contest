package account

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *Handlers) Registration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if _, _, authorized := r.BasicAuth(); authorized {
		common.RenderError(w, r, common.NewForbiddenError("You must be logged out"))
	}

	request := new(RegistrationRequest)
	if err := common.Bind(r, request); err != nil {
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

	response := &Response{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
	}

	render.Status(r, http.StatusCreated)

	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
		return
	}
}

func (h *Handlers) Auth(ctx context.Context, email, password string) (int, error) {
	return h.service.Auth(ctx, email, password)
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := NewSearchParameters(r)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	entities, err := h.service.Search(
		ctx,
		params.FirstName,
		params.LastName,
		params.Email,
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
			ID:        entity.ID,
			FirstName: entity.FirstName,
			LastName:  entity.LastName,
			Email:     entity.Email,
		}
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
	}
}

func (h *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := getAccountIDFromRequest(r)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	entity, err := h.service.GetAccount(ctx, id)
	if err != nil {
		common.RenderError(w, r, err)
		return
	}

	response := &Response{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Email:     entity.Email,
	}

	render.Status(r, http.StatusOK)
	if err := render.Render(w, r, response); err != nil {
		common.RenderError(w, r, err)
		return
	}
}

func getAccountIDFromRequest(r *http.Request) (uint, error) {
	param := chi.URLParam(r, "account_id")

	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, common.NewValidationError(common.InvalidRequestParameter{
			Name:   "account_id",
			Reason: "account_id must be a number",
		})
	}

	tv, err := common.TranslatedValidatorFromRequest(r)
	if err != nil {
		return 0, err
	}

	if err := tv.ValidateVar("id", id, "required,gt=0"); err != nil {
		return 0, err
	}

	return uint(id), err
}
