package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/nambuitechx/nam-observaility/internal"
	"github.com/nambuitechx/nam-observaility/internal/configs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	r := internal.NewRouter()
	h := otelhttp.NewHandler(r.Router, "http-server")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", r.Port),
		Handler: h,
	}

	go func() {
		configs.Logger.Info("server is running", "port", r.Port)
		if err := srv.ListenAndServe(); err != nil {
			configs.Logger.Error("server error", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	// Drain HTTP first so in-flight spans export, then flush/shutdown the tracer provider.
	_ = srv.Shutdown(context.Background())
	r.Shutdown()

	configs.Logger.Info("server grafully shuts down")
}
