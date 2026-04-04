package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/viral32111/glow-cad-exporter/internal/config"
	"github.com/viral32111/glow-cad-exporter/internal/redis"
)

type Server struct {
	HTTPMultiplexer *http.ServeMux
}

func New(config config.Config, redisClient *redis.Client) (*Server, error) {
	server := &Server{
		HTTPMultiplexer: http.NewServeMux(),
	}

	server.RegisterRoutes()
	slog.Debug("Registered HTTP routes")

	return server, nil
}

func (server *Server) Start(ctx context.Context, address string) error {
	slog.Debug("Starting HTTP server", "address", address)

	httpServer := &http.Server{
		Addr:    address,
		Handler: server.HTTPMultiplexer,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Error stopping HTTP server", "error", err)
		}
	}()

	slog.Info("HTTP server listening", "address", address)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("Unable to start HTTP server on '%s': %w", address, err)
	}

	slog.Info("HTTP server stopped")

	return nil
}

func (server *Server) RegisterRoutes() {
	// server.HTTPMultiplexer.Handle("GET /example", routes.ExampleRouteHandler(server.RouteContext))

	server.HTTPMultiplexer.Handle("/", http.HandlerFunc(server.NotFoundHandler))
}

func (server *Server) NotFoundHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(404)
	json.NewEncoder(response).Encode(map[string]any{
		"error": "Unknown API route",
	})
}
