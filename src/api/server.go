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
	"github.com/F0rzend/simbirsoft_contest/src/domain/animal"
	"github.com/F0rzend/simbirsoft_contest/src/domain/locations"
	"github.com/F0rzend/simbirsoft_contest/src/domain/types"
	"github.com/F0rzend/simbirsoft_contest/src/domain/visits"
)

const JSONContentType = "application/json"

type Server struct {
	account   *account.Handlers
	animal    *animal.Handlers
	types     *types.Handlers
	locations *locations.Handlers
	visits    *visits.Handlers
}

func NewServer(di *common.DependencyInjectionContainer) *Server {
	return &Server{
		account:   account.NewHandlers(di),
		animal:    animal.NewHandlers(di),
		types:     types.NewHandlers(di),
		locations: locations.NewHandlers(di),
		visits:    visits.NewHandlers(di),
	}
}

func (s *Server) GetHTTPHandler(logger *zerolog.Logger) (http.Handler, error) {
	translatedValidator, err := common.NewTranslatedValidator()
	if err != nil {
		return nil, err
	}

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

	r.Post("/registration", s.account.Registration)

	r.Group(func(r chi.Router) {
		r.Use(common.BasicAuthMiddleware(s.account.Auth))

		r.Get("/accounts/search", s.account.Search)
		r.Get("/accounts/{account_id}", s.account.GetAccount)

		r.Get("/animals/types/{type_id}", s.types.GetAnimalType)

		r.Get("/animals/search", s.animal.Search)
		r.Get("/animals/{animal_id}", s.animal.GetAnimal)

		r.Get("/animals/{animal_id}/locations", s.visits.GetAnimalLocations)

		r.Get("/locations/{location_id}", s.locations.GetLocation)
	})

	return r, nil
}
