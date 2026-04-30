package gui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baditaflorin/go_permutation_api/internal/config"
)

// Server represents the GUI server
type Server struct {
	config     *config.Config
	httpServer *http.Server
	handler    *Handler
}

// NewServer creates a new GUI server
func NewServer(cfg *config.Config) *Server {
	handler := NewHandler(cfg)

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/", handler.HandleIndex)
	mux.HandleFunc("/api/config", handler.HandleConfig)
	mux.HandleFunc("/api/config/save", handler.HandleSaveConfig)
	mux.HandleFunc("/api/config/load", handler.HandleLoadConfig)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.GUIPort),
		Handler: mux,
	}

	return &Server{
		config:     cfg,
		httpServer: httpServer,
		handler:    handler,
	}
}

// Start starts the GUI server with graceful shutdown
func (s *Server) Start() error {
	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting GUI server on %s\n", s.httpServer.Addr)
		log.Printf("Open http://%s in your browser\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("GUI server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down GUI server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(s.config.Server.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("GUI server shutdown failed: %w", err)
	}

	log.Println("GUI server stopped")
	return nil
}

// Stop stops the GUI server
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
