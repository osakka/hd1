package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
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
	// Create logs directory using professional build structure
	logsDir := "../build/logs"
	os.MkdirAll(logsDir, 0755)

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logsDir, fmt.Sprintf("holodeck_%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Failed to create log file: %v", err)
		return &LogManager{}
	}

	log.Printf("Logging to: %s", logPath)
	return &LogManager{logFile: logFile}
}

func (lm *LogManager) Log(level, source, message string, data interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"source":    source,
		"message":   message,
		"data":      data,
	}

	if lm.logFile != nil {
		jsonData, _ := json.Marshal(logEntry)
		lm.logFile.WriteString(string(jsonData) + "\n")
		lm.logFile.Sync()
	}

	// Also log to console
	log.Printf("[%s] %s: %s", level, source, message)
}

func (lm *LogManager) LogClientMessage(msg LogMessage) {
	lm.Log(msg.Level, "CLIENT", msg.Message, map[string]interface{}{
		"url":       msg.URL,
		"userAgent": msg.UserAgent,
		"data":      msg.Data,
		"timestamp": msg.Timestamp,
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
	h.BroadcastMessage(jsonData)
}

// Force reload all clients
func (h *Hub) ForceReload() {
	reloadMsg := map[string]interface{}{
		"type": "reload",
	}

	jsonData, _ := json.Marshal(reloadMsg)
	h.BroadcastMessage(jsonData)
}

// GetTimestamp returns a timestamp string for file naming
func GetTimestamp() string {
	return time.Now().Format("2006-01-02_15-04-05")
}

// SetupFileLogging configures logging to a specific file
func SetupFileLogging(logFile string) error {
	// Ensure directory exists
	dir := filepath.Dir(logFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory %s: %v", dir, err)
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %v", logFile, err)
	}

	// Set log output to file
	log.SetOutput(file)
	log.Printf("THD logging initialized: %s", logFile)
	
	return nil
}