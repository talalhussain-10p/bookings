package router

import (
	"io"
	"net/http"

	"github.com/bookings/pkg/model"
	"github.com/bookings/pkg/param"
	"github.com/bookings/pkg/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func HandlerHTTP(r *chi.Mux) {
	r.Use(noSurf)
	r.Use(sessionLoad)
	r.Use(middleware.Logger)
	r.Get("/", home)
	r.Get("/about", about)
}

func home(w http.ResponseWriter, r *http.Request) {
	p := param.Eject(r)
	remoteIP := r.RemoteAddr
	s := p.Session
	s.Put(r.Context(), "remote_ip", remoteIP)
	out, err := render.Client("home.page.tmpl", &model.TemplateData{}, true).RenderTemplate()
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(w, out)
}

func about(w http.ResponseWriter, r *http.Request) {
	p := param.Eject(r)
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"
	remoteIP := p.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	out, err := render.Client("about.page.tmpl", &model.TemplateData{
		StringMap: stringMap,
	}, true).RenderTemplate()
	if err != nil {
		return
	}
	defer out.Close()
	io.Copy(w, out)
}
