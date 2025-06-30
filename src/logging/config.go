package logging

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
)

// Config holds logging configuration
type Config struct {
	Level        string   `json:"level"`
	TraceModules []string `json:"trace_modules"`
	LogDir       string   `json:"log_dir"`
}

// LoadConfig loads logging configuration from environment, flags, and defaults
func LoadConfig() *Config {
	config := &Config{
		Level:        "INFO",
		TraceModules: []string{},
		LogDir:       "/opt/hd1/build/logs",
	}

	// Load from environment variables
	if level := os.Getenv("HD1_LOG_LEVEL"); level != "" {
		config.Level = strings.ToUpper(level)
	}

	if modules := os.Getenv("HD1_TRACE_MODULES"); modules != "" {
		config.TraceModules = strings.Split(modules, ",")
		for i := range config.TraceModules {
			config.TraceModules[i] = strings.TrimSpace(config.TraceModules[i])
		}
	}

	if logDir := os.Getenv("HD1_LOG_DIR"); logDir != "" {
		config.LogDir = logDir
	}

	// Override with command line flags if provided
	loadFromFlags(config)

	return config
}

// loadFromFlags loads configuration from command line flags
func loadFromFlags(config *Config) {
	logLevel := flag.String("log-level", config.Level, "Logging level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	traceModules := flag.String("trace-modules", strings.Join(config.TraceModules, ","), "Comma-separated list of modules to trace")
	logDir := flag.String("log-dir", config.LogDir, "Directory for log files")

	// Only parse if not already parsed
	if !flag.Parsed() {
		flag.Parse()
	}

	config.Level = strings.ToUpper(*logLevel)
	if *traceModules != "" {
		config.TraceModules = strings.Split(*traceModules, ",")
		for i := range config.TraceModules {
			config.TraceModules[i] = strings.TrimSpace(config.TraceModules[i])
		}
	}
	config.LogDir = *logDir
}

// ApplyConfig applies the configuration to the logger
func ApplyConfig(config *Config) error {
	// Convert level string to LogLevel
	level, exists := levelFromString[config.Level]
	if !exists {
		level = INFO // Default fallback
	}

	// Initialize logger
	if err := InitLogger(config.LogDir, level, config.TraceModules); err != nil {
		return err
	}

	return nil
}

// GetConfigJSON returns current configuration as JSON
func GetConfigJSON() ([]byte, error) {
	logger := GetLogger()
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	levelName := "UNKNOWN"
	if name, exists := levelNames[logger.level]; exists {
		levelName = name
	}

	traceModules := make([]string, 0, len(logger.traceModules))
	for module := range logger.traceModules {
		traceModules = append(traceModules, module)
	}

	config := Config{
		Level:        levelName,
		TraceModules: traceModules,
	}

	return json.Marshal(config)
}

// UpdateConfigFromJSON updates configuration from JSON
func UpdateConfigFromJSON(jsonData []byte) error {
	var config Config
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return err
	}

	logger := GetLogger()

	// Update level
	if config.Level != "" {
		if err := logger.SetLevelFromString(config.Level); err != nil {
			return err
		}
	}

	// Update trace modules
	if len(config.TraceModules) > 0 {
		// Clear existing and set new
		logger.mu.Lock()
		logger.traceModules = make(map[string]bool)
		for _, module := range config.TraceModules {
			logger.traceModules[strings.ToLower(module)] = true
		}
		logger.mu.Unlock()
	}

	return nil
}