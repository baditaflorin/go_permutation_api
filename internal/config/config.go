package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	App      AppConfig      `json:"app"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            string `json:"port"`
	Host            string `json:"host"`
	GUIPort         string `json:"gui_port"`
	ShutdownTimeout int    `json:"shutdown_timeout"` // seconds
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Driver   string `json:"driver"`   // postgres, mysql, sqlite
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Table    string `json:"table"`
	Column   string `json:"column"`
	SSLMode  string `json:"ssl_mode"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	MaxElements        int  `json:"max_elements"`
	Quiet              bool `json:"quiet"`
	MemoryStatsFreq    int  `json:"memory_stats_freq"`    // Update memory stats every N permutations
	EnableCORS         bool `json:"enable_cors"`
	EnableMetrics      bool `json:"enable_metrics"`
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            "8080",
			Host:            "localhost",
			GUIPort:         "3000",
			ShutdownTimeout: 5,
		},
		Database: DatabaseConfig{
			Driver:  "", // No database by default
			Host:    "localhost",
			Port:    "5432",
			SSLMode: "disable",
		},
		App: AppConfig{
			MaxElements:     12,
			Quiet:           false,
			MemoryStatsFreq: 1000,
			EnableCORS:      true,
			EnableMetrics:   false,
		},
	}
}

// LoadFromFile loads configuration from a JSON file
func LoadFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := Default()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}

// SaveToFile saves configuration to a JSON file
func (c *Config) SaveToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// Load loads configuration from environment variables and file
// Priority: Environment Variables > Config File > Defaults
func Load() (*Config, error) {
	config := Default()

	// Try to load from file if CONFIG_FILE env var is set
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		fileConfig, err := LoadFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
		config = fileConfig
	}

	// Override with environment variables
	config.loadFromEnv()

	return config, nil
}

// loadFromEnv loads configuration from environment variables
func (c *Config) loadFromEnv() {
	// Server configuration
	if port := os.Getenv("SERVER_PORT"); port != "" {
		c.Server.Port = port
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		c.Server.Host = host
	}
	if guiPort := os.Getenv("GUI_PORT"); guiPort != "" {
		c.Server.GUIPort = guiPort
	}
	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			c.Server.ShutdownTimeout = val
		}
	}

	// Database configuration
	if driver := os.Getenv("DB_DRIVER"); driver != "" {
		c.Database.Driver = driver
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		c.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		c.Database.Port = port
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		c.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		c.Database.Password = password
	}
	if database := os.Getenv("DB_DATABASE"); database != "" {
		c.Database.Database = database
	}
	if table := os.Getenv("DB_TABLE"); table != "" {
		c.Database.Table = table
	}
	if column := os.Getenv("DB_COLUMN"); column != "" {
		c.Database.Column = column
	}
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		c.Database.SSLMode = sslMode
	}

	// App configuration
	if maxElements := os.Getenv("MAX_ELEMENTS"); maxElements != "" {
		if val, err := strconv.Atoi(maxElements); err == nil {
			c.App.MaxElements = val
		}
	}
	if quiet := os.Getenv("QUIET"); quiet != "" {
		c.App.Quiet = quiet == "true" || quiet == "1"
	}
	if memStatsFreq := os.Getenv("MEMORY_STATS_FREQ"); memStatsFreq != "" {
		if val, err := strconv.Atoi(memStatsFreq); err == nil {
			c.App.MemoryStatsFreq = val
		}
	}
	if enableCORS := os.Getenv("ENABLE_CORS"); enableCORS != "" {
		c.App.EnableCORS = enableCORS == "true" || enableCORS == "1"
	}
	if enableMetrics := os.Getenv("ENABLE_METRICS"); enableMetrics != "" {
		c.App.EnableMetrics = enableMetrics == "true" || enableMetrics == "1"
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}
	if c.Server.GUIPort == "" {
		return fmt.Errorf("GUI port cannot be empty")
	}
	if c.App.MaxElements < 1 {
		return fmt.Errorf("max elements must be at least 1")
	}
	if c.App.MemoryStatsFreq < 1 {
		return fmt.Errorf("memory stats frequency must be at least 1")
	}

	// Validate database config if driver is specified
	if c.Database.Driver != "" {
		if c.Database.Host == "" {
			return fmt.Errorf("database host cannot be empty when driver is specified")
		}
		if c.Database.Database == "" {
			return fmt.Errorf("database name cannot be empty when driver is specified")
		}
	}

	return nil
}
