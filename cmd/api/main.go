package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"Igorjr19/rinha-de-backend-2026/internal/api"
	"Igorjr19/rinha-de-backend-2026/internal/config"
	"Igorjr19/rinha-de-backend-2026/internal/fraud"
)

func main() {
	cfg := config.Load()
	scorer := fraud.NewScorer()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.NewRouter(scorer),
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    4 << 10,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
