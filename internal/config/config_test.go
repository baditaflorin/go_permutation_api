package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Server.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Server.Port)
	}
	if cfg.App.MaxElements != 12 {
		t.Errorf("expected max elements 12, got %d", cfg.App.MaxElements)
	}
	if cfg.Server.ShutdownTimeout != 5 {
		t.Errorf("expected shutdown timeout 5, got %d", cfg.Server.ShutdownTimeout)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("MAX_ELEMENTS", "20")
	os.Setenv("ENABLE_CORS", "false")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("MAX_ELEMENTS")
		os.Unsetenv("ENABLE_CORS")
	}()

	cfg := Default()
	cfg.loadFromEnv()

	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}
	if cfg.Database.Host != "testhost" {
		t.Errorf("expected db host testhost, got %s", cfg.Database.Host)
	}
	if cfg.App.MaxElements != 20 {
		t.Errorf("expected max elements 20, got %d", cfg.App.MaxElements)
	}
	if cfg.App.EnableCORS != false {
		t.Errorf("expected CORS false, got %v", cfg.App.EnableCORS)
	}
}

func TestSaveAndLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.json")

	// Create a config
	cfg := Default()
	cfg.Server.Port = "7070"
	cfg.App.MaxElements = 15

	// Save to file
	err := cfg.SaveToFile(configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Load from file
	loadedCfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loadedCfg.Server.Port != "7070" {
		t.Errorf("expected port 7070, got %s", loadedCfg.Server.Port)
	}
	if loadedCfg.App.MaxElements != 15 {
		t.Errorf("expected max elements 15, got %d", loadedCfg.App.MaxElements)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  Default(),
			wantErr: false,
		},
		{
			name: "empty server port",
			config: &Config{
				Server: ServerConfig{Port: ""},
				App:    AppConfig{MaxElements: 12, MemoryStatsFreq: 1000},
			},
			wantErr: true,
		},
		{
			name: "invalid max elements",
			config: &Config{
				Server: ServerConfig{Port: "8080", GUIPort: "3000"},
				App:    AppConfig{MaxElements: 0, MemoryStatsFreq: 1000},
			},
			wantErr: true,
		},
		{
			name: "database driver without host",
			config: &Config{
				Server:   ServerConfig{Port: "8080", GUIPort: "3000"},
				Database: DatabaseConfig{Driver: "postgres", Host: ""},
				App:      AppConfig{MaxElements: 12, MemoryStatsFreq: 1000},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
