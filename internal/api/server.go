package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gowebsocket "golang.org/x/net/websocket"

	"github.com/baditaflorin/go_permutation_api/internal/config"
	"github.com/baditaflorin/go_permutation_api/internal/response/reqid"
	"github.com/baditaflorin/go_permutation_api/internal/security"
	"github.com/baditaflorin/go_permutation_api/internal/websocket"
)

// Server represents the HTTP API server
type Server struct {
	config     *config.Config
	httpServer *http.Server
	handler    *Handler
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) *Server {
	handler := NewHandler(cfg)

	wsHandler := websocket.New(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HandlePermutations)
	mux.HandleFunc("/health", handler.HandleHealth)
	mux.HandleFunc("/version", handler.HandleVersion)
	mux.Handle("/ws", gowebsocket.Handler(wsHandler.ServeWS))
	if cfg.App.EnableMetrics {
		mux.HandleFunc("/metrics", handler.HandleMetrics)
	}
	if cfg.App.EnablePprof {
		RegisterPprofHandlers(mux)
	}

	// Apply middleware (innermost → outermost, so outermost wraps first)
	var finalHandler http.Handler = mux
	finalHandler = recoveryMiddleware(finalHandler)
	finalHandler = security.BodyLimit(finalHandler)
	finalHandler = security.Headers(finalHandler)
	finalHandler = security.TrustedProxies(finalHandler)
	if cfg.App.EnableCORS {
		finalHandler = corsMiddleware(finalHandler)
	}
	finalHandler = compressionMiddleware(finalHandler)
	if cfg.App.EnableMetrics {
		finalHandler = metricsMiddleware(finalHandler)
	}
	finalHandler = reqid.Middleware(finalHandler)
	finalHandler = loggingMiddleware(finalHandler)

	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:        finalHandler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   60 * time.Second, // permutations can stream for longer
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		config:     cfg,
		httpServer: httpServer,
		handler:    handler,
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting API server on %s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down API server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(s.config.Server.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("API server stopped")
	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Handler returns the underlying HTTP handler for testing.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}
