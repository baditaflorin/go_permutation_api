package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"
)

// App constants
const (
	defaultPort  = "8080"
	maxElements  = 12
	defaultQuiet = false
)

// Config holds the application configuration
type Config struct {
	Port  string
	Quiet bool
	Args  []string
}

func main() {
	config := parseFlags(os.Args[1:])

	if config.Port != "" {
		startServer(config.Port)
	} else {
		runCLI(config.Args, config.Quiet)
	}
}

// parseFlags parses command-line flags and arguments
func parseFlags(args []string) Config {
	var cfg Config
	flagSet := flag.NewFlagSet("permutation-generator", flag.ExitOnError)

	quiet := flagSet.Bool("quiet", defaultQuiet, "Suppress permutation output")
	serve := flagSet.Bool("serve", false, "Start HTTP server")
	port := flagSet.String("port", defaultPort, "Port for the HTTP server")

	flagSet.Parse(args)

	if *serve {
		cfg.Port = *port
	} else {
		cfg.Args = flagSet.Args()
		cfg.Quiet = *quiet
	}

	return cfg
}

// runCLI handles the CLI functionality
func runCLI(args []string, quiet bool) {
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	if err := validateInput(args); err != nil {
		logErrorAndExit(err)
	}

	startTime := time.Now()
	permutationCount := 0

	minMem, maxMem := getInitialMemory()

	// Prepare and sort the initial permutation
	perm := make([]string, len(args))
	copy(perm, args)
	sort.Strings(perm)

	// Output the first permutation
	if !quiet {
		fmt.Println(strings.Join(perm, ", "))
	}
	permutationCount++

	// Generate subsequent permutations
	for next := nextPermutation(perm); next; next = nextPermutation(perm) {
		if !quiet {
			fmt.Println(strings.Join(perm, ", "))
		}
		permutationCount++

		// Update memory stats every 1000 permutations
		if permutationCount%1000 == 0 {
			updateMemoryStats(&minMem, &maxMem)
		}
	}

	duration := time.Since(startTime)
	minMemory, maxMemory := getFinalMemoryStats(minMem, maxMem)

	printCLIResults(permutationCount, duration, minMemory, maxMemory)
}

// printUsage displays the usage information
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  serve [--port PORT]          Start HTTP server")
	fmt.Println("  [--quiet] <elements>         Generate permutations")
	flag.PrintDefaults()
}

// logErrorAndExit logs the error message and exits the program
func logErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// getInitialMemory retrieves the initial memory usage
func getInitialMemory() (minMem, maxMem uint64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	minMem = memStats.Alloc
	maxMem = memStats.Alloc
	return
}

// updateMemoryStats updates the minimum and maximum memory usage
func updateMemoryStats(minMem, maxMem *uint64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	if memStats.Alloc < *minMem {
		*minMem = memStats.Alloc
	}
	if memStats.Alloc > *maxMem {
		*maxMem = memStats.Alloc
	}
}

// getFinalMemoryStats formats the memory statistics
func getFinalMemoryStats(minMem, maxMem uint64) (uint64, uint64) {
	return minMem / 1024, maxMem / 1024
}

// printCLIResults outputs the permutation generation results
func printCLIResults(count int, duration time.Duration, minMem, maxMem uint64) {
	fmt.Printf("Generated %d permutations in %v\n", count, duration)
	fmt.Printf("Memory usage: Min = %d KB, Max = %d KB\n", minMem, maxMem)
}

// startServer initializes and starts the HTTP server
func startServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlePermutations)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go handleShutdown(server, idleConnsClosed)

	log.Printf("Starting server on port %s...", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

// handleShutdown gracefully shuts down the server on interrupt
func handleShutdown(server *http.Server, done chan struct{}) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	log.Println("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	close(done)
}

// handlePermutations routes the request based on HTTP method
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

// handleGet processes GET requests for permutations
func handleGet(w http.ResponseWriter, r *http.Request) {
	elements, err := parseElementsFromQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writePermutations(w, elements)
}

// handlePost processes POST requests for permutations
func handlePost(w http.ResponseWriter, r *http.Request) {
	elements, err := parseElementsFromBody(r)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateInput(elements); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writePermutations(w, elements)
}

// parseElementsFromQuery extracts elements from URL query parameters
func parseElementsFromQuery(r *http.Request) ([]string, error) {
	query := r.URL.Query().Get("elements")
	if query == "" {
		return nil, errors.New("Missing 'elements' parameter")
	}
	elements := strings.Split(query, ",")
	if err := validateInput(elements); err != nil {
		return nil, err
	}
	return elements, nil
}

// parseElementsFromBody extracts elements from JSON request body
func parseElementsFromBody(r *http.Request) ([]string, error) {
	var elements []string
	if err := json.NewDecoder(r.Body).Decode(&elements); err != nil {
		return nil, err
	}
	if err := validateInput(elements); err != nil {
		return nil, err
	}
	return elements, nil
}

// validateInput checks if the input elements meet the constraints
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

// writePermutations writes the permutations to the HTTP response in streaming fashion
func writePermutations(w http.ResponseWriter, elements []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	w.Write([]byte("["))

	// Sort first
	sort.Strings(elements)
	// Write the first permutation
	if err := enc.Encode(elements); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	for next := nextPermutation(elements); next; next = nextPermutation(elements) {
		w.Write([]byte(","))
		if err := enc.Encode(elements); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("]"))
}

// nextPermutation generates the next lexicographical permutation in-place
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

// reverse reverses a slice of strings in-place starting from the given index
func reverse(arr []string, start int) {
	end := len(arr) - 1
	for start < end {
		arr[start], arr[end] = arr[end], arr[start]
		start++
		end--
	}
}
