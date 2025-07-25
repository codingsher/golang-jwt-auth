package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codingsher/user-jwt-auth/internal/config"
	"github.com/codingsher/user-jwt-auth/internal/http/handlers/user"
	"github.com/codingsher/user-jwt-auth/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initializated", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	router := http.NewServeMux()
	router.HandleFunc("POST /api/register_user", user.RegisterUser(storage))
	router.HandleFunc("POST /api/login_user", user.LoginUser(storage))


	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server failed to start!")
		}
	}()

	<-done

	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("serve failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown success...")
}
