package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
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
		Handler:           api.NewRouter(scorer),
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       2 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    4 << 10,
	}

	ln, err := listen(cfg)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}

func listen(cfg config.Config) (net.Listener, error) {
	if cfg.SocketPath != "" {
		_ = os.Remove(cfg.SocketPath)
		ln, err := net.Listen("unix", cfg.SocketPath)
		if err != nil {
			return nil, err
		}
		if err := os.Chmod(cfg.SocketPath, 0666); err != nil {
			ln.Close()
			return nil, err
		}
		return ln, nil
	}
	return net.Listen("tcp", ":"+cfg.Port)
}
