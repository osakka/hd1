package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"holodeck1/config"
	"holodeck1/logging"
	"holodeck1/server"
)

func main() {
	// Initialize configuration system - Single Source of Truth
	if err := config.Initialize(); err != nil {
		log.Fatalf("FATAL: Configuration initialization failed: %v", err)
	}

	// Additional flags specific to main only (not part of configuration system)
	var (
		help = flag.Bool("help", false, "Show help message")
	)
	
	// Parse flags after config initialization to allow overrides
	if !flag.Parsed() {
		flag.Parse()
	}

	if *help {
		showHelp()
		return
	}

	// Initialize unified logging system with configuration from config system
	logConfig := &logging.Config{
		Level:        config.Config.Logging.Level,
		TraceModules: config.Config.Logging.TraceModules,
		LogDir:       config.Config.Logging.LogDir,
	}
	if err := logging.ApplyConfig(logConfig); err != nil {
		log.Fatalf("FATAL: Failed to initialize logging: %v", err)
	}

	// Setup legacy logging compatibility if specified
	if config.Config.Logging.LogFile != "" {
		if err := setupLogging(config.Config.Logging.LogFile); err != nil {
			logging.Warn("legacy logging setup failed", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Handle daemon mode
	if config.GetDaemon() {
		if err := becomedaemon(config.GetPIDFile()); err != nil {
			logging.Fatal("failed to daemonize process", map[string]interface{}{
				"pid_file": config.GetPIDFile(),
				"error":    err.Error(),
			})
		}
		defer removePidFile(config.GetPIDFile())
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
	
	// Template-processed JavaScript files
	http.HandleFunc("/static/js/hd1-console.js", server.ServeConsoleJS)
	
	// REVOLUTIONARY: Auto-generated API router from specification
	apiRouter := NewAPIRouter(hub)
	http.Handle("/api/", apiRouter)
	
	// Serve static files with proper cache control headers
	fileServer := http.FileServer(http.Dir(config.GetStaticDir()))
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		"version":     "v3.4.0",
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

func showHelp() {
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

func ensureDirectories() error {
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

func setupLogging(logFile string) error {
	if logFile == "" {
		// Default timestamped log file
		logFile = filepath.Join(config.Config.Paths.LogDir, fmt.Sprintf("hd1_%s.log", 
			server.GetTimestamp()))
	}
	
	return server.SetupFileLogging(logFile)
}

func writePidFile(pidFile string, pid ...int) error {
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

func removePidFile(pidFile string) {
	os.Remove(pidFile)
}

func becomedaemon(pidFile string) error {
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
		if err := writePidFile(pidFile, cmd.Process.Pid); err != nil {
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