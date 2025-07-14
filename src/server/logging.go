package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"holodeck1/logging"
)

type LogMessage struct {
	Type      string      `json:"type"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
	URL       string      `json:"url,omitempty"`
	UserAgent string      `json:"userAgent,omitempty"`
}

type LogManager struct {
	logFile *os.File
}

func NewLogManager() *LogManager {
	// Create logs directory using standard build structure
	logsDir := "../build/logs"
	os.MkdirAll(logsDir, 0755)

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("hd1_%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logging.Error("failed to create legacy log file", map[string]interface{}{
			"path": logPath,
			"error": err.Error(),
		})
		return &LogManager{}
	}

	logging.Info("legacy log manager initialized", map[string]interface{}{
		"log_path": logPath,
	})
	return &LogManager{logFile: logFile}
}

func (lm *LogManager) Log(level, source, message string, data interface{}) {
	// Route to unified logging system
	dataMap, ok := data.(map[string]interface{})
	if !ok && data != nil {
		dataMap = map[string]interface{}{"legacy_data": data}
	}
	if dataMap == nil {
		dataMap = make(map[string]interface{})
	}
	dataMap["legacy_source"] = source

	// Convert legacy levels to unified logging
	switch level {
	case "info":
		logging.Info(message, dataMap)
	case "warn", "warning":
		logging.Warn(message, dataMap)
	case "error":
		logging.Error(message, dataMap)
	case "debug":
		logging.Debug(message, dataMap)
	default:
		logging.Info(message, dataMap)
	}

	// Legacy file writing for backward compatibility
	if lm.logFile != nil {
		logEntry := map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"level":     level,
			"source":    source,
			"message":   message,
			"data":      data,
		}
		jsonData, _ := json.Marshal(logEntry)
		lm.logFile.WriteString(string(jsonData) + "\n")
		lm.logFile.Sync()
	}
}

func (lm *LogManager) LogClientMessage(msg LogMessage) {
	// Route directly to unified logging
	logging.Debug("client message", map[string]interface{}{
		"level":     msg.Level,
		"message":   msg.Message,
		"url":       msg.URL,
		"userAgent": msg.UserAgent,
		"data":      msg.Data,
		"timestamp": msg.Timestamp,
		"source":    "browser",
	})
}

func (lm *LogManager) Close() {
	if lm.logFile != nil {
		lm.logFile.Close()
	}
}

// Send log message to client
func (h *Hub) SendLog(level, message string, data interface{}) {
	logMsg := map[string]interface{}{
		"type":    "log",
		"level":   level,
		"message": message,
		"data":    data,
	}

	jsonData, _ := json.Marshal(logMsg)
	h.broadcastMessage(jsonData)
}

// Force reload all clients
func (h *Hub) ForceReload() {
	reloadMsg := map[string]interface{}{
		"type": "reload",
	}

	jsonData, _ := json.Marshal(reloadMsg)
	h.broadcastMessage(jsonData)
}

// GetTimestamp returns a timestamp string for file naming
func GetTimestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

// SetupFileLogging configures logging to a specific file (legacy compatibility)
func SetupFileLogging(logFile string) error {
	// Ensure directory exists
	dir := filepath.Dir(logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logging.Error("failed to create legacy log directory", map[string]interface{}{
			"directory": dir,
			"error": err.Error(),
		})
		return fmt.Errorf("failed to create log directory %s: %v", dir, err)
	}

	logging.Info("legacy file logging setup complete", map[string]interface{}{
		"log_file": logFile,
		"note": "unified logging system is primary",
	})
	
	return nil
}