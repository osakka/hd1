package system

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/yaml.v3"
	"holodeck1/logging"
	"holodeck1/server"
)

type VersionResponse struct {
	APIVersion     string    `json:"api_version"`
	JSVersion      string    `json:"js_version"`
	BuildTimestamp time.Time `json:"build_timestamp"`
	Title          string    `json:"title"`
}

type APISpec struct {
	Info struct {
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
	} `yaml:"info"`
}

// GetVersionHandler - GET /version
func GetVersionHandler(w http.ResponseWriter, r *http.Request, hub interface{}) {
	// Cast hub to proper type (not needed for version endpoint but follows pattern)
	_, ok := hub.(*server.Hub)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	logging.Info("version endpoint called", map[string]interface{}{
		"method":   r.Method,
		"endpoint": "/version",
	})

	// Read API specification to get version
	apiVersion := "1.0.0" // fallback
	apiTitle := "HD1 (Holodeck One) API" // fallback
	
	if specData, err := ioutil.ReadFile("api.yaml"); err == nil {
		var spec APISpec
		if err := yaml.Unmarshal(specData, &spec); err == nil {
			apiVersion = spec.Info.Version
			apiTitle = spec.Info.Title
		}
	}

	// Get JS version hash
	jsVersion := server.GetJSVersion()

	response := VersionResponse{
		APIVersion:     apiVersion,
		JSVersion:      jsVersion,
		BuildTimestamp: time.Now(),
		Title:          apiTitle,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logging.Error("failed to encode version response", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logging.Info("version information served", map[string]interface{}{
		"api_version": apiVersion,
		"js_version":  jsVersion[:8], // log first 8 chars
		"title":       apiTitle,
	})
}