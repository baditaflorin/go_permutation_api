package api

import (
	"encoding/json"
	"net/http"
	"runtime"
)

var (
	// Version information set during build
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
	GoVersion = runtime.Version()
)

// VersionInfo holds version and build information
type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// HandleVersion returns version and build information
func (h *Handler) HandleVersion(w http.ResponseWriter, r *http.Request) {
	info := VersionInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		GoVersion: GoVersion,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
