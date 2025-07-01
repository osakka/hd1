package config

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HD1Config represents the complete HD1 configuration system
// Priority: Flags > Environment Variables > Config File > Defaults
type HD1Config struct {
	Server  ServerConfig  `json:"server"`
	Paths   PathsConfig   `json:"paths"`
	Logging LoggingConfig `json:"logging"`
	Client  ClientConfig  `json:"client"`
}

type ServerConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	APIBase   string `json:"api_base"`
	StaticDir string `json:"static_dir"`
	Daemon    bool   `json:"daemon"`
}

type PathsConfig struct {
	RootDir    string `json:"root_dir"`
	BuildDir   string `json:"build_dir"`
	BinDir     string `json:"bin_dir"`
	LogDir     string `json:"log_dir"`
	RuntimeDir string `json:"runtime_dir"`
	ShareDir   string `json:"share_dir"`
	HtDocsDir  string `json:"htdocs_dir"`
	PIDFile    string `json:"pid_file"`
}

type LoggingConfig struct {
	Level       string   `json:"level"`
	TraceModules []string `json:"trace_modules"`
	LogFile     string   `json:"log_file"`
	LogDir      string   `json:"log_dir"`
}

type ClientConfig struct {
	APIBase string `json:"api_base"`
}

// Global configuration instance - Single Source of Truth
var Config *HD1Config

// Initialize loads configuration from all sources with proper priority
func Initialize() error {
	config := &HD1Config{}
	
	// Load defaults first
	config.loadDefaults()
	
	// Load .env file if it exists
	config.loadEnvFile()
	
	// Override with environment variables
	config.loadEnvironmentVariables()
	
	// Override with command line flags (highest priority)
	config.loadFlags()
	
	// Validate and compute derived paths
	if err := config.validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %v", err)
	}
	
	// Set global configuration
	Config = config
	return nil
}

// loadDefaults sets sensible default values
func (c *HD1Config) loadDefaults() {
	// Server defaults
	c.Server.Host = "0.0.0.0"
	c.Server.Port = "8080"
	
	// Path defaults - configurable root directory
	rootDir := "/opt/hd1"
	c.Paths.RootDir = rootDir
	c.Paths.BuildDir = filepath.Join(rootDir, "build")
	c.Paths.BinDir = filepath.Join(rootDir, "build", "bin")
	c.Paths.LogDir = filepath.Join(rootDir, "build", "logs")
	c.Paths.RuntimeDir = filepath.Join(rootDir, "build", "runtime")
	c.Paths.ShareDir = filepath.Join(rootDir, "share")
	c.Paths.HtDocsDir = filepath.Join(rootDir, "share", "htdocs")
	c.Paths.PIDFile = filepath.Join(rootDir, "build", "runtime", "hd1.pid")
	c.Server.StaticDir = filepath.Join(rootDir, "share", "htdocs", "static")
	
	// Logging defaults
	c.Logging.Level = "INFO"
	c.Logging.TraceModules = []string{}
	c.Logging.LogDir = c.Paths.LogDir
}

// loadEnvFile reads configuration from .env file if it exists
func (c *HD1Config) loadEnvFile() {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return // .env file doesn't exist, skip
	}
	
	file, err := os.Open(envFile)
	if err != nil {
		return // Can't open .env file, skip
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse KEY=VALUE format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Remove quotes if present
		value = strings.Trim(value, "\"'")
		
		// Set environment variable (only if not already set)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// loadEnvironmentVariables reads configuration from environment
func (c *HD1Config) loadEnvironmentVariables() {
	// Server configuration
	if host := os.Getenv("HD1_HOST"); host != "" {
		c.Server.Host = host
	}
	if port := os.Getenv("HD1_PORT"); port != "" {
		c.Server.Port = port
	}
	if apiBase := os.Getenv("HD1_API_BASE"); apiBase != "" {
		c.Server.APIBase = apiBase
		c.Client.APIBase = apiBase
	}
	if daemon := os.Getenv("HD1_DAEMON"); daemon == "true" || daemon == "1" {
		c.Server.Daemon = true
	}
	
	// Path configuration
	if rootDir := os.Getenv("HD1_ROOT_DIR"); rootDir != "" {
		c.Paths.RootDir = rootDir
		// Recompute derived paths
		c.computeDerivedPaths()
	}
	if buildDir := os.Getenv("HD1_BUILD_DIR"); buildDir != "" {
		c.Paths.BuildDir = buildDir
	}
	if logDir := os.Getenv("HD1_LOG_DIR"); logDir != "" {
		c.Paths.LogDir = logDir
		c.Logging.LogDir = logDir
	}
	if staticDir := os.Getenv("HD1_STATIC_DIR"); staticDir != "" {
		c.Server.StaticDir = staticDir
	}
	
	// Logging configuration
	if level := os.Getenv("HD1_LOG_LEVEL"); level != "" {
		c.Logging.Level = level
	}
	if modules := os.Getenv("HD1_TRACE_MODULES"); modules != "" {
		c.Logging.TraceModules = strings.Split(modules, ",")
	}
	if logFile := os.Getenv("HD1_LOG_FILE"); logFile != "" {
		c.Logging.LogFile = logFile
	}
}

// loadFlags reads configuration from command line flags
func (c *HD1Config) loadFlags() {
	// Only parse flags if not already parsed
	if !flag.Parsed() {
		// Define flags with current config values as defaults
		host := flag.String("host", c.Server.Host, "Host to bind to")
		port := flag.String("port", c.Server.Port, "Port to bind to") 
		apiBase := flag.String("api-base", c.Server.APIBase, "API base URL")
		daemon := flag.Bool("daemon", c.Server.Daemon, "Run in daemon mode")
		rootDir := flag.String("root-dir", c.Paths.RootDir, "HD1 root directory (absolute path)")
		buildDir := flag.String("build-dir", c.Paths.BuildDir, "Build directory (absolute path)")
		logDir := flag.String("log-dir", c.Paths.LogDir, "Log directory (absolute path)")
		staticDir := flag.String("static-dir", c.Server.StaticDir, "Static files directory (absolute path)")
		pidFile := flag.String("pid-file", c.Paths.PIDFile, "PID file path (absolute)")
		logFile := flag.String("log-file", c.Logging.LogFile, "Log file path (absolute)")
		logLevel := flag.String("log-level", c.Logging.Level, "Logging level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
		traceModules := flag.String("trace-modules", strings.Join(c.Logging.TraceModules, ","), "Comma-separated trace modules")
		
		flag.Parse()
		
		// Apply flag values
		c.Server.Host = *host
		c.Server.Port = *port
		c.Server.Daemon = *daemon
		if *apiBase != "" {
			c.Server.APIBase = *apiBase
			c.Client.APIBase = *apiBase
		}
		c.Paths.RootDir = *rootDir
		c.Paths.BuildDir = *buildDir
		c.Paths.LogDir = *logDir
		c.Logging.LogDir = *logDir
		c.Server.StaticDir = *staticDir
		c.Paths.PIDFile = *pidFile
		c.Logging.LogFile = *logFile
		c.Logging.Level = *logLevel
		if *traceModules != "" {
			c.Logging.TraceModules = strings.Split(*traceModules, ",")
		}
		
		// Recompute derived paths if root changed
		c.computeDerivedPaths()
	}
}

// computeDerivedPaths calculates dependent paths from root directory
func (c *HD1Config) computeDerivedPaths() {
	if c.Paths.BuildDir == "" || strings.HasPrefix(c.Paths.BuildDir, "/opt/hd1") {
		c.Paths.BuildDir = filepath.Join(c.Paths.RootDir, "build")
	}
	if c.Paths.BinDir == "" || strings.HasPrefix(c.Paths.BinDir, "/opt/hd1") {
		c.Paths.BinDir = filepath.Join(c.Paths.BuildDir, "bin")
	}
	if c.Paths.LogDir == "" || strings.HasPrefix(c.Paths.LogDir, "/opt/hd1") {
		c.Paths.LogDir = filepath.Join(c.Paths.BuildDir, "logs")
		c.Logging.LogDir = c.Paths.LogDir
	}
	if c.Paths.RuntimeDir == "" || strings.HasPrefix(c.Paths.RuntimeDir, "/opt/hd1") {
		c.Paths.RuntimeDir = filepath.Join(c.Paths.BuildDir, "runtime")
	}
	if c.Paths.ShareDir == "" || strings.HasPrefix(c.Paths.ShareDir, "/opt/hd1") {
		c.Paths.ShareDir = filepath.Join(c.Paths.RootDir, "share")
	}
	if c.Paths.HtDocsDir == "" || strings.HasPrefix(c.Paths.HtDocsDir, "/opt/hd1") {
		c.Paths.HtDocsDir = filepath.Join(c.Paths.ShareDir, "htdocs")
	}
	if c.Paths.PIDFile == "" || strings.HasPrefix(c.Paths.PIDFile, "/opt/hd1") {
		c.Paths.PIDFile = filepath.Join(c.Paths.RuntimeDir, "hd1.pid")
	}
	if c.Server.StaticDir == "" || strings.HasPrefix(c.Server.StaticDir, "/opt/hd1") {
		c.Server.StaticDir = filepath.Join(c.Paths.HtDocsDir, "static")
	}
}

// validate ensures configuration is valid and complete
func (c *HD1Config) validate() error {
	// Validate required paths are absolute
	if !filepath.IsAbs(c.Paths.RootDir) {
		return fmt.Errorf("root directory must be absolute path: %s", c.Paths.RootDir)
	}
	
	// Compute API base if not set
	if c.Server.APIBase == "" {
		c.Server.APIBase = fmt.Sprintf("http://%s:%s/api", c.Server.Host, c.Server.Port)
	}
	if c.Client.APIBase == "" {
		c.Client.APIBase = c.Server.APIBase
	}
	
	// Ensure all directories exist (create if needed)
	dirs := []string{
		c.Paths.BuildDir,
		c.Paths.BinDir,
		c.Paths.LogDir,
		c.Paths.RuntimeDir,
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	
	return nil
}

// GetAPIBase returns the configured API base URL
func GetAPIBase() string {
	if Config != nil {
		return Config.Client.APIBase
	}
	return "http://localhost:8080/api" // fallback
}

// GetRootDir returns the configured root directory
func GetRootDir() string {
	if Config != nil {
		return Config.Paths.RootDir
	}
	return "/opt/hd1" // fallback
}

// GetStaticDir returns the configured static files directory
func GetStaticDir() string {
	if Config != nil {
		return Config.Server.StaticDir
	}
	return "/opt/hd1/share/htdocs/static" // fallback
}

// GetPIDFile returns the configured PID file path
func GetPIDFile() string {
	if Config != nil {
		return Config.Paths.PIDFile
	}
	return "/opt/hd1/build/runtime/hd1.pid" // fallback
}

// GetHost returns the configured server host
func GetHost() string {
	if Config != nil {
		return Config.Server.Host
	}
	return "0.0.0.0" // fallback
}

// GetPort returns the configured server port
func GetPort() string {
	if Config != nil {
		return Config.Server.Port
	}
	return "8080" // fallback
}

// GetLogDir returns the configured log directory
func GetLogDir() string {
	if Config != nil {
		return Config.Paths.LogDir
	}
	return "/opt/hd1/build/logs" // fallback
}

// GetDaemon returns the daemon mode setting
func GetDaemon() bool {
	if Config != nil {
		return Config.Server.Daemon
	}
	return false // fallback
}