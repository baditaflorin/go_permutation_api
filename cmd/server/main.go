package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/baditaflorin/go_permutation_api/internal/api"
	"github.com/baditaflorin/go_permutation_api/internal/config"
	"github.com/baditaflorin/go_permutation_api/internal/datasource"
	"github.com/baditaflorin/go_permutation_api/internal/gui"
	"github.com/baditaflorin/go_permutation_api/internal/permutation"
	"github.com/baditaflorin/go_permutation_api/pkg/validator"
)

func main() {
	// Parse command-line flags
	serveAPI := flag.Bool("serve", false, "Start the API server")
	serveGUI := flag.Bool("gui", false, "Start the GUI configuration server")
	quiet := flag.Bool("quiet", false, "Suppress permutation output (CLI mode only)")
	csvFile := flag.String("csv", "", "Path to CSV file for input data")
	tsvFile := flag.String("tsv", "", "Path to TSV file for input data")
	column := flag.Int("column", 0, "Column index for CSV/TSV files (0-based)")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v\n", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v\n", err)
	}

	// Override quiet setting from command line
	if *quiet {
		cfg.App.Quiet = true
	}

	// Determine mode of operation
	switch {
	case *serveAPI:
		runAPIServer(cfg)
	case *serveGUI:
		runGUIServer(cfg)
	case *csvFile != "":
		runFromCSV(cfg, *csvFile, *column)
	case *tsvFile != "":
		runFromTSV(cfg, *tsvFile, *column)
	default:
		runCLI(cfg, flag.Args())
	}
}

// runAPIServer starts the API server
func runAPIServer(cfg *config.Config) {
	server := api.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("API server error: %v\n", err)
	}
}

// runGUIServer starts the GUI configuration server
func runGUIServer(cfg *config.Config) {
	server := gui.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("GUI server error: %v\n", err)
	}
}

// runFromCSV loads elements from a CSV file and generates permutations
func runFromCSV(cfg *config.Config, filePath string, column int) {
	source, err := datasource.NewCSVSource(filePath, column)
	if err != nil {
		log.Fatalf("Failed to create CSV source: %v\n", err)
	}
	defer source.Close()

	elements, err := source.Load()
	if err != nil {
		log.Fatalf("Failed to load CSV data: %v\n", err)
	}

	runPermutations(cfg, elements)
}

// runFromTSV loads elements from a TSV file and generates permutations
func runFromTSV(cfg *config.Config, filePath string, column int) {
	source, err := datasource.NewTSVSource(filePath, column)
	if err != nil {
		log.Fatalf("Failed to create TSV source: %v\n", err)
	}
	defer source.Close()

	elements, err := source.Load()
	if err != nil {
		log.Fatalf("Failed to load TSV data: %v\n", err)
	}

	runPermutations(cfg, elements)
}

// runCLI runs the CLI mode with elements from command-line arguments
func runCLI(cfg *config.Config, args []string) {
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	runPermutations(cfg, args)
}

// runPermutations generates and prints permutations
func runPermutations(cfg *config.Config, elements []string) {
	// Sanitize and validate input
	elements = validator.SanitizeElements(elements)
	if err := validator.ValidateElements(elements, cfg.App.MaxElements); err != nil {
		log.Fatalf("Invalid input: %v\n", err)
	}

	// Track memory usage
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	initialMemory := memStats.Alloc

	minMemory := initialMemory
	maxMemory := initialMemory

	// Sort elements for lexicographical order
	sort.Strings(elements)

	// Generate permutations
	start := time.Now()
	gen := permutation.New(elements)
	count := 0

	for {
		perm, ok := gen.Next()
		if !ok {
			break
		}

		count++

		// Print permutation if not in quiet mode
		if !cfg.App.Quiet {
			fmt.Println(formatPermutation(perm))
		}

		// Update memory stats periodically
		if count%cfg.App.MemoryStatsFreq == 0 {
			runtime.ReadMemStats(&memStats)
			currentMemory := memStats.Alloc
			if currentMemory < minMemory {
				minMemory = currentMemory
			}
			if currentMemory > maxMemory {
				maxMemory = currentMemory
			}
		}
	}

	elapsed := time.Since(start)

	// Print statistics
	printStats(count, elapsed, initialMemory, minMemory, maxMemory)
}

// formatPermutation formats a permutation for output
func formatPermutation(perm []string) string {
	result := ""
	for i, elem := range perm {
		if i > 0 {
			result += ","
		}
		result += elem
	}
	return result
}

// printStats prints execution statistics
func printStats(count int, duration time.Duration, initial, min, max uint64) {
	fmt.Printf("\nGenerated %d permutations in %v\n", count, duration)
	fmt.Printf("Memory: Initial=%dKB Min=%dKB Max=%dKB\n",
		initial/1024, min/1024, max/1024)
}

// printUsage prints usage information
func printUsage() {
	fmt.Println("Permutation API - Generate all lexicographical permutations")
	fmt.Println("\nUsage:")
	fmt.Println("  CLI mode:     permutation-api [options] element1 element2 element3 ...")
	fmt.Println("  CSV mode:     permutation-api --csv=file.csv [--column=0]")
	fmt.Println("  TSV mode:     permutation-api --tsv=file.tsv [--column=0]")
	fmt.Println("  API server:   permutation-api --serve")
	fmt.Println("  GUI server:   permutation-api --gui")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("  CONFIG_FILE          Path to configuration file")
	fmt.Println("  SERVER_PORT          API server port (default: 8080)")
	fmt.Println("  GUI_PORT             GUI server port (default: 3000)")
	fmt.Println("  MAX_ELEMENTS         Maximum number of elements (default: 12)")
	fmt.Println("  DB_DRIVER            Database driver (postgres, mysql, sqlite3)")
	fmt.Println("  DB_HOST              Database host")
	fmt.Println("  DB_PORT              Database port")
	fmt.Println("  DB_DATABASE          Database name")
	fmt.Println("  DB_TABLE             Database table name")
	fmt.Println("  DB_COLUMN            Database column name")
	fmt.Println("\nExamples:")
	fmt.Println("  permutation-api a b c")
	fmt.Println("  permutation-api --quiet 1 2 3 4")
	fmt.Println("  permutation-api --csv=data.csv --column=1")
	fmt.Println("  permutation-api --serve")
	fmt.Println("  permutation-api --gui")
}
