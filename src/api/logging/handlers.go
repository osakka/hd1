package logging

import (
	"encoding/json"
	"net/http"
	"strconv"

	"holodeck/logging"
	"holodeck/server"
)

// GetLoggingConfigHandler returns current logging configuration
func GetLoggingConfigHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	configJSON, err := logging.GetConfigJSON()
	if err != nil {
		http.Error(w, "Failed to get logging configuration", http.StatusInternalServerError)
		logging.Error("failed to get logging configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(configJSON)

	logging.Info("logging configuration retrieved")
}

// SetLoggingConfigHandler updates logging configuration
func SetLoggingConfigHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var configData json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&configData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logging.Error("invalid JSON in logging config request", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := logging.UpdateConfigFromJSON(configData); err != nil {
		http.Error(w, "Failed to update logging configuration", http.StatusBadRequest)
		logging.Error("failed to update logging configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Return updated configuration
	configJSON, err := logging.GetConfigJSON()
	if err != nil {
		http.Error(w, "Failed to get updated configuration", http.StatusInternalServerError)
		logging.Error("failed to get updated logging configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(configJSON)

	logging.Info("logging configuration updated successfully")
}

// SetLogLevelHandler sets the logging level
func SetLogLevelHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Level string `json:"level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logging.Error("invalid JSON in log level request", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := logging.SetLevelFromString(request.Level); err != nil {
		http.Error(w, "Invalid log level", http.StatusBadRequest)
		logging.Error("invalid log level specified", map[string]interface{}{
			"level": request.Level,
			"error": err.Error(),
		})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"level":   request.Level,
		"message": "Log level updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("log level updated", map[string]interface{}{
		"new_level": request.Level,
	})
}

// SetTraceModulesHandler enables/disables tracing for specific modules
func SetTraceModulesHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Enable  []string `json:"enable"`
		Disable []string `json:"disable"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		logging.Error("invalid JSON in trace modules request", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	logger := logging.GetLogger()

	if len(request.Enable) > 0 {
		logger.EnableTrace(request.Enable)
	}

	if len(request.Disable) > 0 {
		logger.DisableTrace(request.Disable)
	}

	response := map[string]interface{}{
		"success": true,
		"enabled": request.Enable,
		"disabled": request.Disable,
		"message": "Trace modules updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("trace modules updated", map[string]interface{}{
		"enabled":  request.Enable,
		"disabled": request.Disable,
	})
}

// GetLogsHandler returns recent log entries
func GetLogsHandler(w http.ResponseWriter, r *http.Request, hub *server.Hub) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get count parameter (default to 100)
	countStr := r.URL.Query().Get("count")
	count := 100 // Default
	if countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 {
			if parsedCount > 1000 {
				count = 1000 // Maximum limit
			} else {
				count = parsedCount
			}
		}
	}

	// Read log entries
	entries, err := logging.ReadLogEntries(count)
	if err != nil {
		http.Error(w, "Failed to read log entries", http.StatusInternalServerError)
		logging.Error("failed to read log entries", map[string]interface{}{
			"error": err.Error(),
			"count": count,
		})
		return
	}

	response := map[string]interface{}{
		"entries": entries,
		"count":   len(entries),
		"requested": count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Debug("log entries retrieved", map[string]interface{}{
		"entries_returned": len(entries),
		"requested_count": count,
	})
}