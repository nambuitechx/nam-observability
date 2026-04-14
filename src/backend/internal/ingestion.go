package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nambuitechx/nam-observaility/internal/configs"
)

type Router struct {
	Port string
	Router *chi.Mux
}

func NewRouter() *Router{

	go func() {
		for {
			checkHealth()
			time.Sleep(5 * time.Second)
		}
	}()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configs.Logger.Info("incoming request",
			"path", r.URL.Path,
			"method", r.Method,
		)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "healthy",
		})
	}))

	// Env
	envConfig := configs.NewEnvConfig()

	return &Router{
		Port: envConfig.Port,
		Router: r,
	}
}

func (r *Router) Shutdown() {
	
}

func checkHealth() {
	configs.Logger.Info("server health check")
}
