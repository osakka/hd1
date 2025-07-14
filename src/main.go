// Package main provides the HD1 (Holodeck One) daemon entry point.
// HD1 is an API-first Three.js game engine platform that exposes complete 3D game
// development capabilities through REST endpoints with TCP-simple WebSocket
// synchronization for collaborative 3D environments.
//
// Architecture:
//   - Configuration system: Flags > Environment Variables > .env File > Defaults
//   - Unified logging: Structured JSON logging with module-based tracing
//   - WebSocket hub: TCP-simple sequence-based synchronization
//   - API router: Auto-generated from OpenAPI specification
//   - Three.js integration: Direct WebGL rendering with zero abstraction
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"holodeck1/config"
	"holodeck1/logging"
	"holodeck1/server"
)

// main is the HD1 daemon entry point.
// Initializes configuration, logging, WebSocket hub, and HTTP server
// following the startup sequence: Config → Logging → Hub → Router → Server
func main() {
	// Configuration initialization: Load settings from all sources
	// Priority: Flags > Environment Variables > .env File > Defaults
	if err := config.Initialize(); err != nil {
		// Cannot use structured logging before logging is initialized
		fmt.Fprintf(os.Stderr, "FATAL: Configuration initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Main-specific flags: Local flags not part of global configuration system
	// Used for operational commands like help display
	var (
		help = flag.Bool("help", false, "Show help message")
	)
	
	// Parse flags after config initialization to allow overrides
	if !flag.Parsed() {
		flag.Parse()
	}

	if *help {
		display_help_information()
		return
	}

	// Logging initialization: Setup structured logging with config integration
	// Supports module-based tracing and configurable log levels
	logConfig := &logging.Config{
		Level:        config.Config.Logging.Level,
		TraceModules: config.Config.Logging.TraceModules,
		LogDir:       config.Config.Logging.LogDir,
	}
	if err := logging.ApplyConfig(logConfig); err != nil {
		// Cannot use structured logging before logging is initialized
		fmt.Fprintf(os.Stderr, "FATAL: Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}

	// Setup legacy logging compatibility if specified
	if config.Config.Logging.LogFile != "" {
		if err := configure_file_logging(config.Config.Logging.LogFile); err != nil {
			logging.Warn("legacy logging setup failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Handle daemon mode
	if config.GetDaemon() {
		if err := convert_to_daemon_process(config.GetPIDFile()); err != nil {
			logging.Fatal("failed to daemonize process", map[string]interface{}{
				"pid_file": config.GetPIDFile(),
				"error":    err.Error(),
			})
		}
		defer remove_process_identifier_file(config.GetPIDFile())
	}

	// Validate static directory from configuration
	if _, err := os.Stat(config.GetStaticDir()); os.IsNotExist(err) {
		logging.Fatal("static directory does not exist", map[string]interface{}{
			"static_dir": config.GetStaticDir(),
		})
	}

	// Initialize HD1
	hub := server.NewHub()
	go hub.Run()

	// Initialize template processor with configured static directory
	server.InitializeTemplateProcessor(config.GetStaticDir())
	
	// WebSocket and static files
	http.HandleFunc("/", server.ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	})
	
	// REVOLUTIONARY: Auto-generated API router from specification
	apiRouter := NewAPIRouter(hub)
	http.Handle("/api/", apiRouter)
	
	// Template-processed JavaScript files with API-driven versioning (must be before static handler)
	http.HandleFunc("/static/js/hd1-console.js", server.ServeConsoleJS)
	
	// Serve static files with proper cache control headers
	fileServer := http.FileServer(http.Dir(config.GetStaticDir()))
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip template-processed files (they have their own handlers)
		if r.URL.Path == "/static/js/hd1-console.js" {
			http.NotFound(w, r) // This should never be reached due to HandleFunc precedence
			return
		}
		
		// Set cache control headers for static assets
		if filepath.Ext(r.URL.Path) == ".js" || filepath.Ext(r.URL.Path) == ".css" {
			// For development: no-cache for JS/CSS to avoid cache issues
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		} else {
			// For other static assets (images, etc): cache for 1 hour
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		http.StripPrefix("/static/", fileServer).ServeHTTP(w, r)
	}))

	// Standard startup banner
	logging.Info("HD1 (Holodeck One) daemon starting", map[string]interface{}{
		"version":      config.GetVersion(),
		"architecture": "spec-driven",
	})
	
	logging.Info("directory configuration", map[string]interface{}{
		"root_dir":    config.GetRootDir(),
		"static_dir":  config.GetStaticDir(),
		"log_dir":     config.Config.Paths.LogDir,
		"runtime_dir": config.Config.Paths.RuntimeDir,
	})
	
	if config.GetDaemon() {
		logging.Info("daemon mode enabled", map[string]interface{}{
			"pid_file": config.GetPIDFile(),
		})
	}

	logging.Info("core API endpoints initialized", map[string]interface{}{
		"sessions":    "/api/sessions",
		"objects":     "/api/sessions/{id}/objects", 
		"world":       "/api/sessions/{id}/world",
		"camera":      "/api/sessions/{id}/camera/position",
		"scenes":      "/api/scenes",
		"recording":   "/api/sessions/{id}/recording/*",
		"admin":       "/admin/logging/*",
	})
	
	bindAddr := fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port)
	logging.Info("server binding to address", map[string]interface{}{
		"address": bindAddr,
		"host":    config.Config.Server.Host,
		"port":    config.Config.Server.Port,
	})
	
	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		logging.Fatal("server failed to start", map[string]interface{}{
			"address": bindAddr,
			"error":   err.Error(),
		})
	}
}

func display_help_information() {
	fmt.Println("HD1 (Holodeck One) - Professional 3D Holodeck Platform")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  hd1 [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --daemon          Run HD1 as daemon")
	fmt.Println("  --pid-file PATH   PID file path (absolute)")
	fmt.Println("  --log-file PATH   Log file path (absolute)")
	fmt.Println("  --host HOST       Host to bind to (default: 0.0.0.0)")
	fmt.Println("  --port PORT       Port to bind to (default: 8080)")
	fmt.Println("  --static-dir PATH Static files directory (absolute)")
	fmt.Println("  --help            Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  hd1")
	fmt.Println("  hd1 --daemon --log-file /opt/hd1/build/logs/hd1.log")
	fmt.Println("  hd1 --host 127.0.0.1 --port 9090")
	fmt.Println()
	fmt.Printf("DEFAULT PATHS:\n")
	fmt.Printf("  Root: %s\n", config.GetRootDir())
	fmt.Printf("  Static: %s\n", config.GetStaticDir())
	fmt.Printf("  Logs: %s\n", config.Config.Paths.LogDir)
	fmt.Printf("  PID: %s\n", config.GetPIDFile())
}

func create_required_build_directories() error {
	dirs := []string{
		config.Config.Paths.BuildDir, 
		config.Config.Paths.BinDir, 
		config.Config.Paths.LogDir, 
		config.Config.Paths.RuntimeDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}
	return nil
}

func configure_file_logging(logFile string) error {
	if logFile == "" {
		// Default timestamped log file
		logFile = filepath.Join(config.Config.Paths.LogDir, fmt.Sprintf("hd1_%s.log", 
			server.GetTimestamp()))
	}
	
	return server.SetupFileLogging(logFile)
}

func write_process_identifier_file(pidFile string, pid ...int) error {
	file, err := os.Create(pidFile)
	if err != nil {
		return err
	}
	defer file.Close()
	
	pidToWrite := os.Getpid()
	if len(pid) > 0 {
		pidToWrite = pid[0]
	}
	
	_, err = fmt.Fprintf(file, "%d\n", pidToWrite)
	return err
}

func remove_process_identifier_file(pidFile string) {
	os.Remove(pidFile)
}

func convert_to_daemon_process(pidFile string) error {
	// Fork and exit parent
	if os.Getppid() != 1 {
		// We are in the parent process, need to fork
		executable, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %v", err)
		}
		
		// Get current args excluding --daemon flag for child
		args := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "--daemon" {
				args = append(args, arg)
			}
		}
		
		// Start child process
		cmd := &exec.Cmd{
			Path: executable,
			Args: append([]string{executable}, args...),
			Env:  os.Environ(),
		}
		
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start daemon process: %v", err)
		}
		
		// Write PID file from parent before exiting
		if err := write_process_identifier_file(pidFile, cmd.Process.Pid); err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("failed to write PID file: %v", err)
		}
		
		// Parent exits
		os.Exit(0)
	}
	
	// We are now in the child process
	// Create new session
	if _, err := syscall.Setsid(); err != nil {
		return fmt.Errorf("failed to create new session: %v", err)
	}
	
	// Change working directory to root
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("failed to change working directory: %v", err)
	}
	
	// Close stdin, stdout, stderr
	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)
	
	// Reopen to /dev/null
	devNull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open /dev/null: %v", err)
	}
	syscall.Dup2(int(devNull.Fd()), 0)
	syscall.Dup2(int(devNull.Fd()), 1)
	syscall.Dup2(int(devNull.Fd()), 2)
	devNull.Close()
	
	return nil
}