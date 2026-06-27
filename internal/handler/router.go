// Package handler wires up HTTP routes and middleware.
package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/juanozorio/task-api/docs" // generated swagger docs
	"github.com/juanozorio/task-api/internal/config"
	"github.com/juanozorio/task-api/internal/service"
)

// NewRouter builds and returns the application HTTP router.
func NewRouter(svc service.TaskService, cfg config.ServerConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(jsonLogger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	// Health & readiness probes (used by Kubernetes)
	r.Get("/healthz", healthz)
	r.Get("/readyz", readyz)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		taskHandler := newTaskHandler(svc)

		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", taskHandler.create)
			r.Get("/", taskHandler.getAll)
			r.Get("/status", taskHandler.getByStatus)
			r.Get("/{id}", taskHandler.getByID)
			r.Put("/{id}", taskHandler.update)
		})
	})

	return r
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func readyz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
