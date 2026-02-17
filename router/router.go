package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimdlwr "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/render"
	"github.com/spf13/viper"

	pkglog "sb-scanner/pkg/logger"
	"sb-scanner/pkg/repository"
	"sb-scanner/router/controller"
)

// @title                       SB Scanner API
// @version                     v1
// @BasePath                    /
func New(v *viper.Viper, debug bool) http.Handler {
	logger := pkglog.GetLogger().With("pkg", "router")
	r := chi.NewRouter()

	repo, err := repository.NewRepository(v.GetString("db.url"), v.GetString("db.name"))
	if err != nil {
		logger.Error("failed to initialize repository", "err", err)
		os.Exit(1)
	}

	r.Use(chimdlwr.RealIP)
	r.Use(chimdlwr.RequestID)
	r.Use(httplog.RequestLogger(pkglog.GetRequestLogger(), &httplog.Options{
		Schema:          httplog.SchemaECS.Concise(true),
		RecoverPanics:   true,
		LogRequestBody:  func(req *http.Request) bool { return debug },
		LogResponseBody: func(r *http.Request) bool { return debug },
	}))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	ctrl := controller.NewController(repo)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/commit", ctrl.GetCommits)
		})
		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{\"message\":\"not found\"}\n"))
		})
	})

	fileServer := http.FileServer(http.Dir("./public/dist"))
	r.Handle("/*", http.StripPrefix("/", fileServer))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/dist/index.html")
	})

	return r
}
