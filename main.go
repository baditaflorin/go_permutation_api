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
	"sort"
	"strings"
	"syscall"
	"time"
)

// App constants
const (
	defaultPort = "8080"
	maxElements = 12
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		port := getEnv("PORT", defaultPort)
		startServer(port)
	} else {
		runCLI(os.Args[1:])
	}
}

// Helper function to reverse a slice
func reverse(arr []string, start int) {
	end := len(arr) - 1
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}

func generatePermutationsChannel(arr []string) <-chan []string {
	// Use a buffered channel to minimize blocking
	ch := make(chan []string, 10000)
	go func() {
		defer close(ch)
		sort.Strings(arr) // Ensure starting with the smallest lexicographical permutation
		perm := make([]string, len(arr))
		copy(perm, arr)
		ch <- perm

		for {
			if !nextPermutation(arr) {
				break
			}
			// Reuse the same slice to reduce memory allocations
			copy(perm, arr)
			ch <- perm
		}
	}()
	return ch
}

// Efficient in-place generation of the next lexicographical permutation
func nextPermutation(arr []string) bool {
	i := len(arr) - 2
	for i >= 0 && arr[i] >= arr[i+1] {
		i--
	}
	if i < 0 {
		return false
	}

	j := len(arr) - 1
	for arr[j] <= arr[i] {
		j--
	}

	arr[i], arr[j] = arr[j], arr[i]
	reverse(arr, i+1)
	return true
}

// CLI entry point
func runCLI(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: go run main.go [serve [port]] <elements>")
		os.Exit(1)
	}

	elements := args
	if err := validateInput(elements); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	start := time.Now()
	permutationCount := 0

	//for perm := range generatePermutationsChannel(elements) {
	for range generatePermutationsChannel(elements) {
		//fmt.Println(strings.Join(perm, ", "))
		permutationCount++
		//}
	}

	fmt.Printf("Generated %d permutations in %v\n", permutationCount, time.Since(start))
}

// HTTP Server entry point
func startServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/permutations", handlePermutations)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

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
	if len(elements) > maxElements {
		return fmt.Errorf("Too many elements (max %d allowed)", maxElements)
	}
	for _, elem := range elements {
		if strings.TrimSpace(elem) == "" {
			return errors.New("Empty elements are not allowed")
		}
	}
	return nil
}

func writePermutations(w http.ResponseWriter, elements []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	w.Write([]byte("["))
	first := true
	for perm := range generatePermutationsChannel(elements) {
		if !first {
			w.Write([]byte(","))
		}
		if err := enc.Encode(perm); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		first = false
	}
	w.Write([]byte("]"))
}

// Get environment variable or fallback to default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
