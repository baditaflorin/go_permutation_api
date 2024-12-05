package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// App constants
const (
	defaultPort = "8080"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		port := getEnv("PORT", defaultPort)
		startServer(port)
	} else {
		runCLI(os.Args[1:])
	}
}

// Functional implementation of permutations in Go
func permutations(arr []string) [][]string {
	if len(arr) == 0 {
		return [][]string{}
	}
	if len(arr) == 1 {
		return [][]string{arr}
	}

	result := [][]string{}
	for i, val := range arr {
		rest := append(append([]string{}, arr[:i]...), arr[i+1:]...)
		for _, perm := range permutations(rest) {
			result = append(result, append([]string{val}, perm...))
		}
	}
	return result
}

// CLI entry point
func runCLI(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: go run main.go [serve [port]] <elements>")
		os.Exit(1)
	}

	elements := args
	perms := permutations(elements)
	printPermutations(perms)
}

// HTTP Server Entry Point
func startServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/permutations", handlePermutations)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Starting server on port %s...", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

// HTTP handler for permutations
func handlePermutations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("elements")
	if query == "" {
		http.Error(w, "Missing 'elements' parameter", http.StatusBadRequest)
		return
	}

	elements := strings.Split(query, ",")
	if err := validateInput(elements); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writePermutations(w, elements)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var elements []string
	if err := json.NewDecoder(r.Body).Decode(&elements); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateInput(elements); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writePermutations(w, elements)
}

func validateInput(elements []string) error {
	if len(elements) > 10 {
		return errors.New("Too many elements (max 10 allowed)")
	}
	for _, elem := range elements {
		if strings.TrimSpace(elem) == "" {
			return errors.New("Empty elements are not allowed")
		}
	}
	return nil
}

func writePermutations(w http.ResponseWriter, elements []string) {
	perms := permutations(elements)
	response, err := json.Marshal(perms)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Utility function for printing permutations to the console
func printPermutations(perms [][]string) {
	for i, perm := range perms {
		fmt.Printf("%d: %s\n", i+1, strings.Join(perm, ", "))
	}
}

// Get environment variable or fallback to default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
