package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nambuitechx/nam-observaility/internal/configs"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Router struct {
	Router         *chi.Mux
	Port           string
	TracerShutdown func(context.Context) error
}

func NewRouter() *Router {
	envConfig := configs.NewEnvConfig()
	tracerShutdown := InitTracer(envConfig)

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
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := otel.Tracer("backend")

		ctx, span := tr.Start(r.Context(), "incoming request")
		defer span.End()

		// Use the span on ctx (same trace as otelhttp parent). TraceID.String() is empty if invalid.
		sc := trace.SpanContextFromContext(ctx)
		traceID := sc.TraceID().String()
		spanID := sc.SpanID().String()

		configs.Logger.InfoContext(ctx, "incoming request",
			"path", r.URL.Path,
			"method", r.Method,
			"trace_id", traceID,
			"span_id", spanID,
		)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "healthy",
		})
	}))

	go func() {
		for {
			checkHealth()
			time.Sleep(5 * time.Second)
		}
	}()

	return &Router{
		Router:         r,
		Port:           envConfig.Port,
		TracerShutdown: tracerShutdown,
	}
}

func (r *Router) Shutdown() {
	if err := r.TracerShutdown(context.Background()); err != nil {
		log.Fatal("Failed to shutdown otel tracer")
	}
}

func checkHealth() {
	configs.Logger.Info("server health check")
}

func InitTracer(envConfig *configs.EnvConfig) func(context.Context) error {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(fmt.Sprintf("%s:%s", envConfig.AlloyHost, envConfig.AlloyPort)),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("backend"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown
}
