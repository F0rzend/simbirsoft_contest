package common

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type ctxKey struct {
	id string
}

var (
	translatedValidatorCtxKey = &ctxKey{"translatedValidator"}
	authorizedUserIDCtxKey    = &ctxKey{"authorizedUserID"}
)

type Middleware func(next http.Handler) http.Handler

func TranslatedValidatorCtxMiddleware(translatedValidator *TranslatedValidator) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), translatedValidatorCtxKey, translatedValidator))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func TranslatedValidatorFromRequest(r *http.Request) (*TranslatedValidator, error) {
	translatedValidator, ok := r.Context().Value(translatedValidatorCtxKey).(*TranslatedValidator)
	if !ok {
		return nil, errors.New("Cannot get validator from request context")
	}

	return translatedValidator, nil
}

type Checker func(ctx context.Context, login string, password string) (int, error)

func BasicAuthMiddleware(checkFn Checker) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			login, password, ok := r.BasicAuth()
			if !ok {
				RenderError(w, r, NewAuthorizationError("Invalid header format"))
				return
			}

			id, err := checkFn(r.Context(), login, password)
			if err != nil {
				RenderError(w, r, NewAuthorizationError("Wrong login or password."))
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), authorizedUserIDCtxKey, id))

			next.ServeHTTP(w, r)
		})
	}
}

// func UserIDFromRequest(r *http.Request) (int, error) {
//	id, ok := r.Context().Value(translatedValidatorCtxKey).(int)
//	if !ok {
//		return 0, errors.New("Cannot get user id from request context")
//	}
//
//	return id, nil
//}
