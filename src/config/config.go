package config

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// HD1Config represents the complete HD1 configuration system
// Priority: Flags > Environment Variables > Config File > Defaults
type HD1Config struct {
	Server    ServerConfig    `json:"server"`
	Paths     PathsConfig     `json:"paths"`
	Logging   LoggingConfig   `json:"logging"`
	Client    ClientConfig    `json:"client"`
	WebSocket WebSocketConfig `json:"websocket"`
	Session   SessionConfig   `json:"session"`
	Worlds    WorldsConfig    `json:"worlds"`
	Avatars   AvatarsConfig   `json:"avatars"`
	Sync      SyncConfig      `json:"sync"`
}

type ServerConfig struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	APIBase         string `json:"api_base"`
	InternalAPIBase string `json:"internal_api_base"`
	StaticDir       string `json:"static_dir"`
	Daemon          bool   `json:"daemon"`
	Version         string `json:"version"`
}

type PathsConfig struct {
	RootDir      string `json:"root_dir"`
	BuildDir     string `json:"build_dir"`
	BinDir       string `json:"bin_dir"`
	LogDir       string `json:"log_dir"`
	RuntimeDir   string `json:"runtime_dir"`
	ShareDir     string `json:"share_dir"`
	HtDocsDir    string `json:"htdocs_dir"`
	PIDFile      string `json:"pid_file"`
	WorldsDir    string `json:"worlds_dir"`
	AvatarsDir   string `json:"avatars_dir"`
	RecordingsDir string `json:"recordings_dir"`
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

// WebSocketConfig contains WebSocket-specific configuration
type WebSocketConfig struct {
	WriteTimeout        time.Duration `json:"write_timeout"`
	PongTimeout         time.Duration `json:"pong_timeout"`
	PingPeriod          time.Duration `json:"ping_period"`
	MaxMessageSize      int64         `json:"max_message_size"`
	ReadBufferSize      int           `json:"read_buffer_size"`
	WriteBufferSize     int           `json:"write_buffer_size"`
	ClientWorldBuffer int           `json:"client_world_buffer"`
}

// SessionConfig contains session management configuration
type SessionConfig struct {
	CleanupInterval     time.Duration `json:"cleanup_interval"`
	InactivityTimeout   time.Duration `json:"inactivity_timeout"`
	HTTPClientTimeout   time.Duration `json:"http_client_timeout"`
	DefaultSessionID    string        `json:"default_session_id"`
}

// WorldsConfig contains world system configuration
type WorldsConfig struct {
	ConfigFile       string   `json:"config_file"`
	DefaultWorld     string   `json:"default_world"`
	ProtectedList    []string `json:"protected_list"`
	AutoJoinOnCreate bool     `json:"auto_join_on_create"`
	SyncOnJoin       bool     `json:"sync_on_join"`
}

// AvatarsConfig contains avatar system configuration
type AvatarsConfig struct {
	ConfigFile              string        `json:"config_file"`
	MaxConcurrentCreations  int           `json:"max_concurrent_creations"`
	HealthCheckInterval     time.Duration `json:"health_check_interval"`
	PositionUpdateThrottle  time.Duration `json:"position_update_throttle"`
	MaxReconnectAttempts    int           `json:"max_reconnect_attempts"`
	ReconnectDelay          time.Duration `json:"reconnect_delay"`
	MaxReconnectDelay       time.Duration `json:"max_reconnect_delay"`
	HeartbeatFrequency      time.Duration `json:"heartbeat_frequency"`
}

// SyncConfig contains HD1-VSC synchronization protocol configuration
type SyncConfig struct {
	Protocol                string        `json:"protocol"`                 // HD1-VSC protocol version
	SyncInterval            time.Duration `json:"sync_interval"`            // Sync broadcast interval
	MaxDeltaLog            int           `json:"max_delta_log"`            // Maximum delta operations to keep
	ChecksumAlgorithm      string        `json:"checksum_algorithm"`       // Checksum algorithm (sha256, md5)
	CausalityTimeout       time.Duration `json:"causality_timeout"`        // Timeout for out-of-order operations
	DeltaQueueSize         int           `json:"delta_queue_size"`         // Size of delta operation queue
	AvatarRegistrySize     int           `json:"avatar_registry_size"`     // Initial avatar registry capacity
	BroadcastWorldBuffer int           `json:"broadcast_world_buffer"` // Broadcast world buffer size
	WorldStateCompressionEnabled bool    `json:"world_state_compression_enabled"` // Enable world state compression
	PerformanceMetricsEnabled bool      `json:"performance_metrics_enabled"`     // Enable sync performance metrics
	VectorClockPrecision   int           `json:"vector_clock_precision"`   // Vector clock precision bits
}

// Global configuration instance - Single Source of Truth
var Config *HD1Config

// create_unique_session_identifier creates a unique session identifier
func create_unique_session_identifier() string {
	// Generate 8 random bytes for session ID
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return fmt.Sprintf("session-%d", time.Now().Unix())
	}
	
	// Convert to lowercase hex string
	id := fmt.Sprintf("session-%x", bytes)
	return id
}

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
	c.Server.APIBase = "http://0.0.0.0:8080/api"
	c.Server.InternalAPIBase = "http://localhost:8080/api"
	c.Server.Version = "v0.7.0"
	
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
	c.Paths.WorldsDir = filepath.Join(rootDir, "share", "worlds")
	c.Paths.AvatarsDir = filepath.Join(rootDir, "share", "avatars")
	c.Paths.RecordingsDir = filepath.Join(rootDir, "recordings")
	c.Server.StaticDir = filepath.Join(rootDir, "share", "htdocs", "static")
	
	// Logging defaults
	c.Logging.Level = "INFO"
	c.Logging.TraceModules = []string{}
	c.Logging.LogDir = c.Paths.LogDir
	
	// WebSocket defaults (based on current hardcoded values)
	c.WebSocket.WriteTimeout = 10 * time.Second
	c.WebSocket.PongTimeout = 60 * time.Second
	c.WebSocket.PingPeriod = 54 * time.Second // (60 * 9) / 10
	c.WebSocket.MaxMessageSize = 1048576  // 1MB for GLB assets
	c.WebSocket.ReadBufferSize = 1048576  // 1MB read buffer
	c.WebSocket.WriteBufferSize = 1048576 // 1MB write buffer
	c.WebSocket.ClientWorldBuffer = 256
	
	// Session defaults (based on current hardcoded values)
	c.Session.CleanupInterval = 2 * time.Minute
	c.Session.InactivityTimeout = 10 * time.Minute
	c.Session.HTTPClientTimeout = 5 * time.Second
	c.Session.DefaultSessionID = create_unique_session_identifier()
	
	// Worlds defaults
	c.Worlds.ConfigFile = "config.yaml"
	c.Worlds.DefaultWorld = "world_one"
	c.Worlds.ProtectedList = []string{"world_one", "world_two"}
	c.Worlds.AutoJoinOnCreate = true
	c.Worlds.SyncOnJoin = true
	
	// Avatars defaults (based on current hardcoded values)
	c.Avatars.ConfigFile = "config.yaml"
	c.Avatars.MaxConcurrentCreations = 2
	c.Avatars.HealthCheckInterval = 5 * time.Second
	c.Avatars.PositionUpdateThrottle = 16 * time.Millisecond // ~60fps
	c.Avatars.MaxReconnectAttempts = 99
	c.Avatars.ReconnectDelay = 1 * time.Second
	c.Avatars.MaxReconnectDelay = 30 * time.Second
	c.Avatars.HeartbeatFrequency = 5 * time.Second
	
	// Sync protocol defaults (eliminating hardcoded values)
	c.Sync.Protocol = "HD1-VSC-v1.0"
	c.Sync.SyncInterval = 16 * time.Millisecond  // ~60fps sync rate
	c.Sync.MaxDeltaLog = 10000                   // Store 10k delta operations
	c.Sync.ChecksumAlgorithm = "sha256"          // Cryptographic integrity
	c.Sync.CausalityTimeout = 5 * time.Second    // Timeout for out-of-order ops
	c.Sync.DeltaQueueSize = 1000                 // Queue size for causality resolution
	c.Sync.AvatarRegistrySize = 100              // Initial avatar registry capacity
	c.Sync.BroadcastWorldBuffer = 1024         // Configurable broadcast buffer
	c.Sync.WorldStateCompressionEnabled = true   // Enable compression for performance
	c.Sync.PerformanceMetricsEnabled = false     // Disable metrics by default
	c.Sync.VectorClockPrecision = 64             // 64-bit vector clock precision
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
	if internalAPIBase := os.Getenv("HD1_INTERNAL_API_BASE"); internalAPIBase != "" {
		c.Server.InternalAPIBase = internalAPIBase
	}
	if version := os.Getenv("HD1_VERSION"); version != "" {
		c.Server.Version = version
	}
	if daemon := os.Getenv("HD1_DAEMON"); daemon == "true" || daemon == "1" {
		c.Server.Daemon = true
	}
	
	// Path configuration
	if rootDir := os.Getenv("HD1_ROOT_DIR"); rootDir != "" {
		c.Paths.RootDir = rootDir
		// Recompute derived paths
		c.calculate_dependent_directory_paths()
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
	
	if worldsDir := os.Getenv("HD1_WORLDS_DIR"); worldsDir != "" {
		c.Paths.WorldsDir = worldsDir
	}
	if avatarsDir := os.Getenv("HD1_AVATARS_DIR"); avatarsDir != "" {
		c.Paths.AvatarsDir = avatarsDir
	}
	if recordingsDir := os.Getenv("HD1_RECORDINGS_DIR"); recordingsDir != "" {
		c.Paths.RecordingsDir = recordingsDir
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
	
	// WebSocket configuration
	if writeTimeout := os.Getenv("HD1_WEBSOCKET_WRITE_TIMEOUT"); writeTimeout != "" {
		if timeout, err := time.ParseDuration(writeTimeout); err == nil {
			c.WebSocket.WriteTimeout = timeout
		}
	}
	if pongTimeout := os.Getenv("HD1_WEBSOCKET_PONG_TIMEOUT"); pongTimeout != "" {
		if timeout, err := time.ParseDuration(pongTimeout); err == nil {
			c.WebSocket.PongTimeout = timeout
		}
	}
	if pingPeriod := os.Getenv("HD1_WEBSOCKET_PING_PERIOD"); pingPeriod != "" {
		if period, err := time.ParseDuration(pingPeriod); err == nil {
			c.WebSocket.PingPeriod = period
		}
	}
	if maxMessageSize := os.Getenv("HD1_WEBSOCKET_MAX_MESSAGE_SIZE"); maxMessageSize != "" {
		if size, err := strconv.ParseInt(maxMessageSize, 10, 64); err == nil {
			c.WebSocket.MaxMessageSize = size
		}
	}
	if readBufferSize := os.Getenv("HD1_WEBSOCKET_READ_BUFFER_SIZE"); readBufferSize != "" {
		if size, err := strconv.Atoi(readBufferSize); err == nil {
			c.WebSocket.ReadBufferSize = size
		}
	}
	if writeBufferSize := os.Getenv("HD1_WEBSOCKET_WRITE_BUFFER_SIZE"); writeBufferSize != "" {
		if size, err := strconv.Atoi(writeBufferSize); err == nil {
			c.WebSocket.WriteBufferSize = size
		}
	}
	if worldBuffer := os.Getenv("HD1_WEBSOCKET_CLIENT_WORLD_BUFFER"); worldBuffer != "" {
		if size, err := strconv.Atoi(worldBuffer); err == nil {
			c.WebSocket.ClientWorldBuffer = size
		}
	}
	
	// Session configuration
	if cleanupInterval := os.Getenv("HD1_SESSION_CLEANUP_INTERVAL"); cleanupInterval != "" {
		if interval, err := time.ParseDuration(cleanupInterval); err == nil {
			c.Session.CleanupInterval = interval
		}
	}
	if inactivityTimeout := os.Getenv("HD1_SESSION_INACTIVITY_TIMEOUT"); inactivityTimeout != "" {
		if timeout, err := time.ParseDuration(inactivityTimeout); err == nil {
			c.Session.InactivityTimeout = timeout
		}
	}
	if httpClientTimeout := os.Getenv("HD1_SESSION_HTTP_CLIENT_TIMEOUT"); httpClientTimeout != "" {
		if timeout, err := time.ParseDuration(httpClientTimeout); err == nil {
			c.Session.HTTPClientTimeout = timeout
		}
	}
	if defaultSessionID := os.Getenv("HD1_SESSION_DEFAULT_ID"); defaultSessionID != "" {
		c.Session.DefaultSessionID = defaultSessionID
	}
	
	// Worlds configuration
	if configFile := os.Getenv("HD1_WORLDS_CONFIG_FILE"); configFile != "" {
		c.Worlds.ConfigFile = configFile
	}
	if defaultWorld := os.Getenv("HD1_WORLDS_DEFAULT_WORLD"); defaultWorld != "" {
		c.Worlds.DefaultWorld = defaultWorld
	}
	if autoJoin := os.Getenv("HD1_WORLDS_AUTO_JOIN_ON_CREATE"); autoJoin == "true" || autoJoin == "1" {
		c.Worlds.AutoJoinOnCreate = true
	} else if autoJoin == "false" || autoJoin == "0" {
		c.Worlds.AutoJoinOnCreate = false
	}
	if syncOnJoin := os.Getenv("HD1_WORLDS_SYNC_ON_JOIN"); syncOnJoin == "true" || syncOnJoin == "1" {
		c.Worlds.SyncOnJoin = true
	} else if syncOnJoin == "false" || syncOnJoin == "0" {
		c.Worlds.SyncOnJoin = false
	}
	if protectedList := os.Getenv("HD1_WORLDS_PROTECTED_LIST"); protectedList != "" {
		c.Worlds.ProtectedList = strings.Split(protectedList, ",")
	}
	
	// Avatars configuration
	if configFile := os.Getenv("HD1_AVATARS_CONFIG_FILE"); configFile != "" {
		c.Avatars.ConfigFile = configFile
	}
	if maxConcurrent := os.Getenv("HD1_AVATARS_MAX_CONCURRENT_CREATIONS"); maxConcurrent != "" {
		if max, err := strconv.Atoi(maxConcurrent); err == nil {
			c.Avatars.MaxConcurrentCreations = max
		}
	}
	if healthCheck := os.Getenv("HD1_AVATARS_HEALTH_CHECK_INTERVAL"); healthCheck != "" {
		if interval, err := time.ParseDuration(healthCheck); err == nil {
			c.Avatars.HealthCheckInterval = interval
		}
	}
	if positionThrottle := os.Getenv("HD1_AVATARS_POSITION_UPDATE_THROTTLE"); positionThrottle != "" {
		if throttle, err := time.ParseDuration(positionThrottle); err == nil {
			c.Avatars.PositionUpdateThrottle = throttle
		}
	}
	if maxReconnect := os.Getenv("HD1_AVATARS_MAX_RECONNECT_ATTEMPTS"); maxReconnect != "" {
		if max, err := strconv.Atoi(maxReconnect); err == nil {
			c.Avatars.MaxReconnectAttempts = max
		}
	}
	if reconnectDelay := os.Getenv("HD1_AVATARS_RECONNECT_DELAY"); reconnectDelay != "" {
		if delay, err := time.ParseDuration(reconnectDelay); err == nil {
			c.Avatars.ReconnectDelay = delay
		}
	}
	if maxReconnectDelay := os.Getenv("HD1_AVATARS_MAX_RECONNECT_DELAY"); maxReconnectDelay != "" {
		if delay, err := time.ParseDuration(maxReconnectDelay); err == nil {
			c.Avatars.MaxReconnectDelay = delay
		}
	}
	if heartbeat := os.Getenv("HD1_AVATARS_HEARTBEAT_FREQUENCY"); heartbeat != "" {
		if frequency, err := time.ParseDuration(heartbeat); err == nil {
			c.Avatars.HeartbeatFrequency = frequency
		}
	}
	
	// Sync protocol configuration
	if protocol := os.Getenv("HD1_SYNC_PROTOCOL"); protocol != "" {
		c.Sync.Protocol = protocol
	}
	if syncInterval := os.Getenv("HD1_SYNC_INTERVAL"); syncInterval != "" {
		if interval, err := time.ParseDuration(syncInterval); err == nil {
			c.Sync.SyncInterval = interval
		}
	}
	if maxDeltaLog := os.Getenv("HD1_SYNC_MAX_DELTA_LOG"); maxDeltaLog != "" {
		if max, err := strconv.Atoi(maxDeltaLog); err == nil {
			c.Sync.MaxDeltaLog = max
		}
	}
	if checksumAlgorithm := os.Getenv("HD1_SYNC_CHECKSUM_ALGORITHM"); checksumAlgorithm != "" {
		c.Sync.ChecksumAlgorithm = checksumAlgorithm
	}
	if causalityTimeout := os.Getenv("HD1_SYNC_CAUSALITY_TIMEOUT"); causalityTimeout != "" {
		if timeout, err := time.ParseDuration(causalityTimeout); err == nil {
			c.Sync.CausalityTimeout = timeout
		}
	}
	if deltaQueueSize := os.Getenv("HD1_SYNC_DELTA_QUEUE_SIZE"); deltaQueueSize != "" {
		if size, err := strconv.Atoi(deltaQueueSize); err == nil {
			c.Sync.DeltaQueueSize = size
		}
	}
	if avatarRegistrySize := os.Getenv("HD1_SYNC_AVATAR_REGISTRY_SIZE"); avatarRegistrySize != "" {
		if size, err := strconv.Atoi(avatarRegistrySize); err == nil {
			c.Sync.AvatarRegistrySize = size
		}
	}
	if broadcastBuffer := os.Getenv("HD1_SYNC_BROADCAST_WORLD_BUFFER"); broadcastBuffer != "" {
		if buffer, err := strconv.Atoi(broadcastBuffer); err == nil {
			c.Sync.BroadcastWorldBuffer = buffer
		}
	}
	if compression := os.Getenv("HD1_SYNC_WORLD_STATE_COMPRESSION_ENABLED"); compression == "true" || compression == "1" {
		c.Sync.WorldStateCompressionEnabled = true
	} else if compression == "false" || compression == "0" {
		c.Sync.WorldStateCompressionEnabled = false
	}
	if metrics := os.Getenv("HD1_SYNC_PERFORMANCE_METRICS_ENABLED"); metrics == "true" || metrics == "1" {
		c.Sync.PerformanceMetricsEnabled = true
	} else if metrics == "false" || metrics == "0" {
		c.Sync.PerformanceMetricsEnabled = false
	}
	if precision := os.Getenv("HD1_SYNC_VECTOR_CLOCK_PRECISION"); precision != "" {
		if prec, err := strconv.Atoi(precision); err == nil {
			c.Sync.VectorClockPrecision = prec
		}
	}
}

// loadFlags reads configuration from command line flags
func (c *HD1Config) loadFlags() {
	// Only parse flags if not already parsed
	if !flag.Parsed() {
		// Define flags with current config values as defaults
		// Long and short flag combinations for essential operations
		host := flag.String("host", c.Server.Host, "Host to bind to")
		hostShort := flag.String("h", c.Server.Host, "Host to bind to (short)")
		port := flag.String("port", c.Server.Port, "Port to bind to") 
		portShort := flag.String("p", c.Server.Port, "Port to bind to (short)")
		apiBase := flag.String("api-base", c.Server.APIBase, "API base URL")
		internalAPIBase := flag.String("internal-api-base", c.Server.InternalAPIBase, "Internal API base URL for server communications")
		version := flag.String("version", c.Server.Version, "HD1 version identifier")
		versionShort := flag.String("v", c.Server.Version, "HD1 version identifier (short)")
		daemon := flag.Bool("daemon", c.Server.Daemon, "Run in daemon mode")
		daemonShort := flag.Bool("d", c.Server.Daemon, "Run in daemon mode (short)")
		rootDir := flag.String("root-dir", c.Paths.RootDir, "HD1 root directory (absolute path)")
		buildDir := flag.String("build-dir", c.Paths.BuildDir, "Build directory (absolute path)")
		logDir := flag.String("log-dir", c.Paths.LogDir, "Log directory (absolute path)")
		staticDir := flag.String("static-dir", c.Server.StaticDir, "Static files directory (absolute path)")
		pidFile := flag.String("pid-file", c.Paths.PIDFile, "PID file path (absolute)")
		logFile := flag.String("log-file", c.Logging.LogFile, "Log file path (absolute)")
		logLevel := flag.String("log-level", c.Logging.Level, "Logging level (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
		traceModules := flag.String("trace-modules", strings.Join(c.Logging.TraceModules, ","), "Comma-separated trace modules")
		protectedWorlds := flag.String("protected-worlds", strings.Join(c.Worlds.ProtectedList, ","), "Comma-separated list of protected worlds")
		
		// Extended flags for complete configuration coverage
		worldsDir := flag.String("worlds-dir", c.Paths.WorldsDir, "Worlds configuration directory")
		avatarsDir := flag.String("avatars-dir", c.Paths.AvatarsDir, "Avatars configuration directory")
		recordingsDir := flag.String("recordings-dir", c.Paths.RecordingsDir, "Recordings directory")
		defaultWorld := flag.String("default-world", c.Worlds.DefaultWorld, "Default world identifier")
		autoJoinOnCreate := flag.Bool("auto-join-on-create", c.Worlds.AutoJoinOnCreate, "Auto-join world on session create")
		syncOnJoin := flag.Bool("sync-on-join", c.Worlds.SyncOnJoin, "Sync world state on join")
		
		// WebSocket configuration flags
		writeTimeout := flag.Duration("websocket-write-timeout", c.WebSocket.WriteTimeout, "WebSocket write timeout")
		pongTimeout := flag.Duration("websocket-pong-timeout", c.WebSocket.PongTimeout, "WebSocket pong timeout") 
		pingPeriod := flag.Duration("websocket-ping-period", c.WebSocket.PingPeriod, "WebSocket ping period")
		maxMessageSize := flag.Int64("websocket-max-message-size", c.WebSocket.MaxMessageSize, "WebSocket max message size")
		readBufferSize := flag.Int("websocket-read-buffer-size", c.WebSocket.ReadBufferSize, "WebSocket read buffer size")
		writeBufferSize := flag.Int("websocket-write-buffer-size", c.WebSocket.WriteBufferSize, "WebSocket write buffer size")
		
		// Session configuration flags
		cleanupInterval := flag.Duration("session-cleanup-interval", c.Session.CleanupInterval, "Session cleanup interval")
		inactivityTimeout := flag.Duration("session-inactivity-timeout", c.Session.InactivityTimeout, "Session inactivity timeout")
		httpClientTimeout := flag.Duration("session-http-client-timeout", c.Session.HTTPClientTimeout, "HTTP client timeout")
		
		// Avatar configuration flags
		maxConcurrentCreations := flag.Int("avatars-max-concurrent-creations", c.Avatars.MaxConcurrentCreations, "Max concurrent avatar creations")
		healthCheckInterval := flag.Duration("avatars-health-check-interval", c.Avatars.HealthCheckInterval, "Avatar health check interval")
		positionUpdateThrottle := flag.Duration("avatars-position-update-throttle", c.Avatars.PositionUpdateThrottle, "Avatar position update throttle")
		maxReconnectAttempts := flag.Int("avatars-max-reconnect-attempts", c.Avatars.MaxReconnectAttempts, "Max avatar reconnect attempts")
		reconnectDelay := flag.Duration("avatars-reconnect-delay", c.Avatars.ReconnectDelay, "Avatar reconnect delay")
		maxReconnectDelay := flag.Duration("avatars-max-reconnect-delay", c.Avatars.MaxReconnectDelay, "Max avatar reconnect delay")
		heartbeatFrequency := flag.Duration("avatars-heartbeat-frequency", c.Avatars.HeartbeatFrequency, "Avatar heartbeat frequency")
		
		// Sync protocol configuration flags
		syncProtocol := flag.String("sync-protocol", c.Sync.Protocol, "HD1-VSC sync protocol version")
		syncInterval := flag.Duration("sync-interval", c.Sync.SyncInterval, "Sync broadcast interval")
		maxDeltaLog := flag.Int("sync-max-delta-log", c.Sync.MaxDeltaLog, "Max delta operations to keep")
		checksumAlgorithm := flag.String("sync-checksum-algorithm", c.Sync.ChecksumAlgorithm, "Checksum algorithm")
		causalityTimeout := flag.Duration("sync-causality-timeout", c.Sync.CausalityTimeout, "Causality timeout")
		deltaQueueSize := flag.Int("sync-delta-queue-size", c.Sync.DeltaQueueSize, "Delta operation queue size")
		avatarRegistrySize := flag.Int("sync-avatar-registry-size", c.Sync.AvatarRegistrySize, "Avatar registry size")
		broadcastWorldBuffer := flag.Int("sync-broadcast-world-buffer", c.Sync.BroadcastWorldBuffer, "Broadcast world buffer size")
		worldStateCompression := flag.Bool("sync-world-state-compression", c.Sync.WorldStateCompressionEnabled, "Enable world state compression")
		performanceMetrics := flag.Bool("sync-performance-metrics", c.Sync.PerformanceMetricsEnabled, "Enable sync performance metrics")
		vectorClockPrecision := flag.Int("sync-vector-clock-precision", c.Sync.VectorClockPrecision, "Vector clock precision bits")
		
		flag.Parse()
		
		// Apply flag values (short flags take precedence over long flags)
		if *hostShort != c.Server.Host {
			c.Server.Host = *hostShort
		} else {
			c.Server.Host = *host
		}
		if *portShort != c.Server.Port {
			c.Server.Port = *portShort
		} else {
			c.Server.Port = *port
		}
		if *daemonShort != c.Server.Daemon {
			c.Server.Daemon = *daemonShort
		} else {
			c.Server.Daemon = *daemon
		}
		if *apiBase != "" {
			c.Server.APIBase = *apiBase
			c.Client.APIBase = *apiBase
		}
		if *internalAPIBase != "" {
			c.Server.InternalAPIBase = *internalAPIBase
		}
		if *versionShort != c.Server.Version {
			c.Server.Version = *versionShort
		} else if *version != "" {
			c.Server.Version = *version
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
		if *protectedWorlds != "" {
			c.Worlds.ProtectedList = strings.Split(*protectedWorlds, ",")
		}
		
		// Apply extended configuration flags
		c.Paths.WorldsDir = *worldsDir
		c.Paths.AvatarsDir = *avatarsDir
		c.Paths.RecordingsDir = *recordingsDir
		c.Worlds.DefaultWorld = *defaultWorld
		c.Worlds.AutoJoinOnCreate = *autoJoinOnCreate
		c.Worlds.SyncOnJoin = *syncOnJoin
		
		// Apply WebSocket configuration
		c.WebSocket.WriteTimeout = *writeTimeout
		c.WebSocket.PongTimeout = *pongTimeout
		c.WebSocket.PingPeriod = *pingPeriod
		c.WebSocket.MaxMessageSize = *maxMessageSize
		c.WebSocket.ReadBufferSize = *readBufferSize
		c.WebSocket.WriteBufferSize = *writeBufferSize
		
		// Apply Session configuration
		c.Session.CleanupInterval = *cleanupInterval
		c.Session.InactivityTimeout = *inactivityTimeout
		c.Session.HTTPClientTimeout = *httpClientTimeout
		
		// Apply Avatar configuration
		c.Avatars.MaxConcurrentCreations = *maxConcurrentCreations
		c.Avatars.HealthCheckInterval = *healthCheckInterval
		c.Avatars.PositionUpdateThrottle = *positionUpdateThrottle
		c.Avatars.MaxReconnectAttempts = *maxReconnectAttempts
		c.Avatars.ReconnectDelay = *reconnectDelay
		c.Avatars.MaxReconnectDelay = *maxReconnectDelay
		c.Avatars.HeartbeatFrequency = *heartbeatFrequency
		
		// Apply Sync protocol configuration
		c.Sync.Protocol = *syncProtocol
		c.Sync.SyncInterval = *syncInterval
		c.Sync.MaxDeltaLog = *maxDeltaLog
		c.Sync.ChecksumAlgorithm = *checksumAlgorithm
		c.Sync.CausalityTimeout = *causalityTimeout
		c.Sync.DeltaQueueSize = *deltaQueueSize
		c.Sync.AvatarRegistrySize = *avatarRegistrySize
		c.Sync.BroadcastWorldBuffer = *broadcastWorldBuffer
		c.Sync.WorldStateCompressionEnabled = *worldStateCompression
		c.Sync.PerformanceMetricsEnabled = *performanceMetrics
		c.Sync.VectorClockPrecision = *vectorClockPrecision
		
		// Recompute derived paths if root changed
		c.calculate_dependent_directory_paths()
	}
}

// calculate_dependent_directory_paths calculates dependent paths from root directory
func (c *HD1Config) calculate_dependent_directory_paths() {
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
	if c.Paths.WorldsDir == "" || strings.HasPrefix(c.Paths.WorldsDir, "/opt/hd1") {
		c.Paths.WorldsDir = filepath.Join(c.Paths.ShareDir, "worlds")
	}
	if c.Paths.AvatarsDir == "" || strings.HasPrefix(c.Paths.AvatarsDir, "/opt/hd1") {
		c.Paths.AvatarsDir = filepath.Join(c.Paths.ShareDir, "avatars")
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

// GetWorldsDir returns the configured worlds directory
func GetWorldsDir() string {
	if Config != nil {
		return Config.Paths.WorldsDir
	}
	return "/opt/hd1/share/worlds" // fallback
}

// GetAvatarsDir returns the configured avatars directory
func GetAvatarsDir() string {
	if Config != nil {
		return Config.Paths.AvatarsDir
	}
	return "/opt/hd1/share/avatars" // fallback
}

// GetRecordingsDir returns the configured recordings directory
func GetRecordingsDir() string {
	if Config != nil {
		return Config.Paths.RecordingsDir
	}
	return "/opt/hd1/recordings" // fallback
}

// GetWorldsConfigFile returns the configured worlds config file path
func GetWorldsConfigFile() string {
	if Config != nil {
		return filepath.Join(Config.Paths.WorldsDir, Config.Worlds.ConfigFile)
	}
	return "/opt/hd1/share/worlds/config.yaml" // fallback
}

// GetAvatarsConfigFile returns the configured avatars config file path
func GetAvatarsConfigFile() string {
	if Config != nil {
		return filepath.Join(Config.Paths.AvatarsDir, Config.Avatars.ConfigFile)
	}
	return "/opt/hd1/share/avatars/config.yaml" // fallback
}

// WebSocket configuration getters
func GetWebSocketWriteTimeout() time.Duration {
	if Config != nil {
		return Config.WebSocket.WriteTimeout
	}
	return 10 * time.Second // fallback
}

func GetWebSocketPongTimeout() time.Duration {
	if Config != nil {
		return Config.WebSocket.PongTimeout
	}
	return 60 * time.Second // fallback
}

func GetWebSocketPingPeriod() time.Duration {
	if Config != nil {
		return Config.WebSocket.PingPeriod
	}
	return 54 * time.Second // fallback
}

func GetWebSocketMaxMessageSize() int64 {
	if Config != nil {
		return Config.WebSocket.MaxMessageSize
	}
	return 512 // fallback
}

func GetWebSocketReadBufferSize() int {
	if Config != nil {
		return Config.WebSocket.ReadBufferSize
	}
	return 1024 // fallback
}

func GetWebSocketWriteBufferSize() int {
	if Config != nil {
		return Config.WebSocket.WriteBufferSize
	}
	return 1024 // fallback
}

func GetWebSocketClientWorldBuffer() int {
	if Config != nil {
		return Config.WebSocket.ClientWorldBuffer
	}
	return 256 // fallback
}

// Session configuration getters
func GetSessionCleanupInterval() time.Duration {
	if Config != nil {
		return Config.Session.CleanupInterval
	}
	return 2 * time.Minute // fallback
}

func GetSessionInactivityTimeout() time.Duration {
	if Config != nil {
		return Config.Session.InactivityTimeout
	}
	return 10 * time.Minute // fallback
}

func GetSessionHTTPClientTimeout() time.Duration {
	if Config != nil {
		return Config.Session.HTTPClientTimeout
	}
	return 5 * time.Second // fallback
}

func GetSessionDefaultID() string {
	if Config != nil {
		return Config.Session.DefaultSessionID
	}
	return create_unique_session_identifier() // fallback - generate unique ID
}

// Worlds configuration getters
func GetWorldsDefaultWorld() string {
	if Config != nil {
		return Config.Worlds.DefaultWorld
	}
	return "world_one" // fallback
}

func GetWorldsAutoJoinOnCreate() bool {
	if Config != nil {
		return Config.Worlds.AutoJoinOnCreate
	}
	return true // fallback
}

func GetWorldsSyncOnJoin() bool {
	if Config != nil {
		return Config.Worlds.SyncOnJoin
	}
	return true // fallback
}

// GetWorldsProtectedList returns the list of protected worlds
func GetWorldsProtectedList() []string {
	if Config != nil {
		return Config.Worlds.ProtectedList
	}
	return []string{"world_one", "world_two"} // fallback
}

// GetInternalAPIBase returns the configured internal API base URL
func GetInternalAPIBase() string {
	if Config != nil {
		return Config.Server.InternalAPIBase
	}
	return "http://localhost:8080/api" // fallback
}

// GetVersion returns the configured HD1 version
func GetVersion() string {
	if Config != nil {
		return Config.Server.Version
	}
	return "v0.7.0" // fallback
}

// Avatars configuration getters
func GetAvatarsMaxConcurrentCreations() int {
	if Config != nil {
		return Config.Avatars.MaxConcurrentCreations
	}
	return 2 // fallback
}

func GetAvatarsHealthCheckInterval() time.Duration {
	if Config != nil {
		return Config.Avatars.HealthCheckInterval
	}
	return 5 * time.Second // fallback
}

func GetAvatarsPositionUpdateThrottle() time.Duration {
	if Config != nil {
		return Config.Avatars.PositionUpdateThrottle
	}
	return 16 * time.Millisecond // fallback
}

func GetAvatarsMaxReconnectAttempts() int {
	if Config != nil {
		return Config.Avatars.MaxReconnectAttempts
	}
	return 99 // fallback
}

func GetAvatarsReconnectDelay() time.Duration {
	if Config != nil {
		return Config.Avatars.ReconnectDelay
	}
	return 1 * time.Second // fallback
}

func GetAvatarsMaxReconnectDelay() time.Duration {
	if Config != nil {
		return Config.Avatars.MaxReconnectDelay
	}
	return 30 * time.Second // fallback
}

func GetAvatarsHeartbeatFrequency() time.Duration {
	if Config != nil {
		return Config.Avatars.HeartbeatFrequency
	}
	return 5 * time.Second // fallback
}

// Sync protocol configuration getters
func GetSyncProtocol() string {
	if Config != nil {
		return Config.Sync.Protocol
	}
	return "HD1-VSC-v1.0" // fallback
}

func GetSyncInterval() time.Duration {
	if Config != nil {
		return Config.Sync.SyncInterval
	}
	return 16 * time.Millisecond // fallback
}

func GetSyncMaxDeltaLog() int {
	if Config != nil {
		return Config.Sync.MaxDeltaLog
	}
	return 10000 // fallback
}

func GetSyncChecksumAlgorithm() string {
	if Config != nil {
		return Config.Sync.ChecksumAlgorithm
	}
	return "sha256" // fallback
}

func GetSyncCausalityTimeout() time.Duration {
	if Config != nil {
		return Config.Sync.CausalityTimeout
	}
	return 5 * time.Second // fallback
}

func GetSyncDeltaQueueSize() int {
	if Config != nil {
		return Config.Sync.DeltaQueueSize
	}
	return 1000 // fallback
}

func GetSyncAvatarRegistrySize() int {
	if Config != nil {
		return Config.Sync.AvatarRegistrySize
	}
	return 100 // fallback
}

func GetSyncBroadcastWorldBuffer() int {
	if Config != nil {
		return Config.Sync.BroadcastWorldBuffer
	}
	return 1024 // fallback
}

func GetSyncWorldStateCompressionEnabled() bool {
	if Config != nil {
		return Config.Sync.WorldStateCompressionEnabled
	}
	return true // fallback
}

func GetSyncPerformanceMetricsEnabled() bool {
	if Config != nil {
		return Config.Sync.PerformanceMetricsEnabled
	}
	return false // fallback
}

func GetSyncVectorClockPrecision() int {
	if Config != nil {
		return Config.Sync.VectorClockPrecision
	}
	return 64 // fallback
}

// GetString returns a configuration value as string (used by database package)
func GetString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}