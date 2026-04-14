package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/nambuitechx/nam-observaility/internal"
	"github.com/nambuitechx/nam-observaility/internal/configs"
)

func main() {
	r := internal.NewRouter()
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", r.Port),
		Handler: r.Router,
	}

	go func() {
		configs.Logger.Info("server is running", "port", r.Port)
		if err := srv.ListenAndServe(); err != nil {
			configs.Logger.Error("server error", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<- ctx.Done()
	r.Shutdown()
	srv.Shutdown(context.Background())
	
	configs.Logger.Info("server grafully shuts down")
}
