package gui

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/baditaflorin/go_permutation_api/internal/config"
)

// Handler handles GUI HTTP requests
type Handler struct {
	config *config.Config
}

// NewHandler creates a new GUI handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
	}
}

// HandleIndex serves the main GUI page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Permutation API - Configuration</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            overflow: hidden;
        }
        .header {
            background: #4c51bf;
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 { font-size: 2em; margin-bottom: 10px; }
        .header p { opacity: 0.9; }
        .content { padding: 30px; }
        .section {
            background: #f7fafc;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .section h2 {
            color: #2d3748;
            margin-bottom: 15px;
            font-size: 1.3em;
            border-bottom: 2px solid #4c51bf;
            padding-bottom: 10px;
        }
        .form-group {
            margin-bottom: 15px;
        }
        label {
            display: block;
            color: #4a5568;
            font-weight: 600;
            margin-bottom: 5px;
        }
        input, select {
            width: 100%;
            padding: 10px;
            border: 1px solid #cbd5e0;
            border-radius: 5px;
            font-size: 14px;
        }
        input:focus, select:focus {
            outline: none;
            border-color: #4c51bf;
            box-shadow: 0 0 0 3px rgba(76, 81, 191, 0.1);
        }
        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 20px;
        }
        button {
            flex: 1;
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }
        .btn-primary {
            background: #4c51bf;
            color: white;
        }
        .btn-primary:hover {
            background: #434190;
        }
        .btn-secondary {
            background: #48bb78;
            color: white;
        }
        .btn-secondary:hover {
            background: #38a169;
        }
        .message {
            padding: 12px;
            border-radius: 5px;
            margin-bottom: 20px;
            display: none;
        }
        .message.success {
            background: #c6f6d5;
            color: #22543d;
            border: 1px solid #9ae6b4;
        }
        .message.error {
            background: #fed7d7;
            color: #742a2a;
            border: 1px solid #fc8181;
        }
        .info-box {
            background: #bee3f8;
            border-left: 4px solid #3182ce;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 4px;
        }
        .info-box p { color: #2c5282; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Permutation API</h1>
            <p>Configure your permutation generation service</p>
        </div>
        <div class="content">
            <div id="message" class="message"></div>

            <div class="info-box">
                <p><strong>Current API Endpoint:</strong> http://{{.Server.Host}}:{{.Server.Port}}</p>
            </div>

            <div class="section">
                <h2>Server Configuration</h2>
                <div class="form-group">
                    <label>API Port</label>
                    <input type="text" id="server_port" value="{{.Server.Port}}">
                </div>
                <div class="form-group">
                    <label>Host</label>
                    <input type="text" id="server_host" value="{{.Server.Host}}">
                </div>
                <div class="form-group">
                    <label>GUI Port</label>
                    <input type="text" id="gui_port" value="{{.Server.GUIPort}}">
                </div>
            </div>

            <div class="section">
                <h2>Application Settings</h2>
                <div class="form-group">
                    <label>Maximum Elements</label>
                    <input type="number" id="max_elements" value="{{.App.MaxElements}}">
                </div>
                <div class="form-group">
                    <label>Memory Stats Frequency</label>
                    <input type="number" id="memory_stats_freq" value="{{.App.MemoryStatsFreq}}">
                </div>
            </div>

            <div class="section">
                <h2>Database Configuration</h2>
                <div class="form-group">
                    <label>Database Driver</label>
                    <select id="db_driver">
                        <option value="">None</option>
                        <option value="postgres" {{if eq .Database.Driver "postgres"}}selected{{end}}>PostgreSQL</option>
                        <option value="mysql" {{if eq .Database.Driver "mysql"}}selected{{end}}>MySQL</option>
                        <option value="sqlite3" {{if eq .Database.Driver "sqlite3"}}selected{{end}}>SQLite</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Host</label>
                    <input type="text" id="db_host" value="{{.Database.Host}}">
                </div>
                <div class="form-group">
                    <label>Port</label>
                    <input type="text" id="db_port" value="{{.Database.Port}}">
                </div>
                <div class="form-group">
                    <label>Database Name</label>
                    <input type="text" id="db_database" value="{{.Database.Database}}">
                </div>
                <div class="form-group">
                    <label>Table Name</label>
                    <input type="text" id="db_table" value="{{.Database.Table}}">
                </div>
                <div class="form-group">
                    <label>Column Name</label>
                    <input type="text" id="db_column" value="{{.Database.Column}}">
                </div>
            </div>

            <div class="button-group">
                <button class="btn-primary" onclick="saveConfig()">Save Configuration</button>
                <button class="btn-secondary" onclick="downloadConfig()">Download Config File</button>
            </div>
        </div>
    </div>

    <script>
        function showMessage(text, type) {
            const msg = document.getElementById('message');
            msg.textContent = text;
            msg.className = 'message ' + type;
            msg.style.display = 'block';
            setTimeout(() => { msg.style.display = 'none'; }, 5000);
        }

        async function saveConfig() {
            const config = {
                server: {
                    port: document.getElementById('server_port').value,
                    host: document.getElementById('server_host').value,
                    gui_port: document.getElementById('gui_port').value,
                    shutdown_timeout: 5
                },
                database: {
                    driver: document.getElementById('db_driver').value,
                    host: document.getElementById('db_host').value,
                    port: document.getElementById('db_port').value,
                    database: document.getElementById('db_database').value,
                    table: document.getElementById('db_table').value,
                    column: document.getElementById('db_column').value,
                },
                app: {
                    max_elements: parseInt(document.getElementById('max_elements').value),
                    memory_stats_freq: parseInt(document.getElementById('memory_stats_freq').value),
                    enable_cors: true,
                    quiet: false
                }
            };

            try {
                const response = await fetch('/api/config/save', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(config)
                });

                if (response.ok) {
                    showMessage('Configuration saved successfully!', 'success');
                } else {
                    const error = await response.text();
                    showMessage('Failed to save: ' + error, 'error');
                }
            } catch (err) {
                showMessage('Error: ' + err.message, 'error');
            }
        }

        function downloadConfig() {
            const config = {
                server: {
                    port: document.getElementById('server_port').value,
                    host: document.getElementById('server_host').value,
                    gui_port: document.getElementById('gui_port').value,
                    shutdown_timeout: 5
                },
                database: {
                    driver: document.getElementById('db_driver').value,
                    host: document.getElementById('db_host').value,
                    port: document.getElementById('db_port').value,
                    database: document.getElementById('db_database').value,
                    table: document.getElementById('db_table').value,
                    column: document.getElementById('db_column').value,
                },
                app: {
                    max_elements: parseInt(document.getElementById('max_elements').value),
                    memory_stats_freq: parseInt(document.getElementById('memory_stats_freq').value),
                    enable_cors: true,
                    quiet: false
                }
            };

            const blob = new Blob([JSON.stringify(config, null, 2)], { type: 'application/json' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'config.json';
            a.click();
            URL.revokeObjectURL(url);
            showMessage('Configuration file downloaded!', 'success');
        }
    </script>
</body>
</html>`

	t := template.Must(template.New("index").Parse(tmpl))
	t.Execute(w, h.config)
}

// HandleConfig returns the current configuration as JSON
func (h *Handler) HandleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.config)
}

// HandleSaveConfig saves configuration to a file
func (h *Handler) HandleSaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate configuration
	if err := newConfig.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save to file
	configPath := "configs/runtime_config.json"
	if err := newConfig.SaveToFile(configPath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration saved successfully"))
}

// HandleLoadConfig loads configuration from an uploaded file
func (h *Handler) HandleLoadConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("config")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var newConfig config.Config
	if err := json.NewDecoder(file).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid configuration file", http.StatusBadRequest)
		return
	}

	// Validate configuration
	if err := newConfig.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update current configuration
	h.config = &newConfig

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Configuration loaded successfully",
	})
}
