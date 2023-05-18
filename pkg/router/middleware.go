package router

import (
	"net/http"

	"github.com/bookings/pkg/param"
)

// NoSurf add CSRF protection to all POST requests
func noSurf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := param.Eject(r)
		p.CSRFHandler.SetBaseCookie(
			http.Cookie{
				HttpOnly: true,
				Path:     "/",
				Secure:   p.AppENV,
				SameSite: http.SameSiteLaxMode,
			},
		)
		next.ServeHTTP(w, r)
	})
}

// seesionLoad loads and saves the session on every request
func sessionLoad(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := param.Eject(r)
		p.Session.LoadAndSave(next).ServeHTTP(w, r)
	})
}
