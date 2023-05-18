package param

import (
	"context"
	"errors"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

type keyInContext string

var (
	paramKey keyInContext = "param"
)

type Param struct {
	AppENV      bool
	Session     *scs.SessionManager
	CSRFHandler *nosurf.CSRFHandler
	Router      chi.Router
}

func Inject(router chi.Router, p *Param) {
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), paramKey, p)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func Eject(r *http.Request) *Param {
	if v := r.Context().Value(paramKey); v != nil {
		if p, ok := v.(*Param); ok {
			return p
		}
	}
	return nil
}

func EjectParamFromContext(ctx context.Context) (*Param, error) {
	v := ctx.Value(paramKey)
	if p, ok := v.(*Param); ok {
		return p, nil
	}
	return nil, errors.New("Param not found in context")
}
