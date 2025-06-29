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

	"holodeck/server"
)

// THD Path Configuration - 100% Absolute Paths
const (
	THD_ROOT_DIR      = "/opt/holo-deck"
	THD_BUILD_DIR     = THD_ROOT_DIR + "/build"
	THD_BIN_DIR       = THD_BUILD_DIR + "/bin"
	THD_LOG_DIR       = THD_BUILD_DIR + "/logs"
	THD_RUNTIME_DIR   = THD_BUILD_DIR + "/runtime"
	THD_SHARE_DIR     = THD_ROOT_DIR + "/share"
	THD_HTDOCS_DIR    = THD_SHARE_DIR + "/htdocs"
	THD_STATIC_DIR    = THD_HTDOCS_DIR + "/static"
	THD_PID_FILE      = THD_RUNTIME_DIR + "/thd.pid"
	THD_DEFAULT_HOST  = "0.0.0.0"
	THD_DEFAULT_PORT  = "8080"
)

func main() {
	// Command line flags - LONG FLAGS ONLY
	var (
		daemonize = flag.Bool("daemon", false, "Run THD as daemon")
		pidFile   = flag.String("pid-file", THD_PID_FILE, "PID file path (absolute)")
		logFile   = flag.String("log-file", "", "Log file path (absolute, defaults to timestamped)")
		host      = flag.String("host", THD_DEFAULT_HOST, "Host to bind to")
		port      = flag.String("port", THD_DEFAULT_PORT, "Port to bind to")
		staticDir = flag.String("static-dir", THD_STATIC_DIR, "Static files directory (absolute)")
		help      = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Ensure all directories exist
	if err := ensureDirectories(); err != nil {
		log.Fatalf("FATAL: Failed to create directories: %v", err)
	}

	// Setup logging with absolute path
	if err := setupLogging(*logFile); err != nil {
		log.Fatalf("FATAL: Failed to setup logging: %v", err)
	}

	// Handle daemon mode
	if *daemonize {
		if err := becomedaemon(*pidFile); err != nil {
			log.Fatalf("FATAL: Failed to daemonize: %v", err)
		}
		defer removePidFile(*pidFile)
	}

	// Validate static directory
	if _, err := os.Stat(*staticDir); os.IsNotExist(err) {
		log.Fatalf("FATAL: Static directory does not exist: %s", *staticDir)
	}

	// Initialize THD
	hub := server.NewHub()
	go hub.Run()

	// WebSocket and static files
	http.HandleFunc("/", server.ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWS(hub, w, r)
	})
	
	// REVOLUTIONARY: Auto-generated API router from specification
	apiRouter := NewAPIRouter(hub)
	http.Handle("/api/", apiRouter)
	
	// Serve static files with proper cache control headers
	fileServer := http.FileServer(http.Dir(*staticDir))
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

	// Startup banner
	log.Println("THD (The Holo-Deck) - Professional Daemon")
	log.Println("SPEC-DRIVEN ARCHITECTURE - Single Source of Truth")
	log.Println("api.yaml drives all routing automatically")
	log.Println("")
	log.Printf("Root Directory: %s", THD_ROOT_DIR)
	log.Printf("Static Directory: %s", *staticDir)
	log.Printf("Log Directory: %s", THD_LOG_DIR)
	log.Printf("Runtime Directory: %s", THD_RUNTIME_DIR)
	if *daemonize {
		log.Printf("PID File: %s", *pidFile)
	}
	log.Println("")
	log.Println("API Endpoints:")
	log.Println("   POST /api/sessions - Create session")
	log.Println("   GET  /api/sessions - List sessions")
	log.Println("   POST /api/sessions/{id}/world - Initialize world")
	log.Println("   POST /api/sessions/{id}/objects - Create objects")
	log.Println("   PUT  /api/sessions/{id}/camera/position - Set camera")
	log.Println("")
	
	bindAddr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("THD Server starting on %s", bindAddr)
	log.Fatal(http.ListenAndServe(bindAddr, nil))
}

func showHelp() {
	fmt.Println("THD (The Holo-Deck) - Professional 3D Visualization Engine")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  thd [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --daemon          Run THD as daemon")
	fmt.Println("  --pid-file PATH   PID file path (absolute)")
	fmt.Println("  --log-file PATH   Log file path (absolute)")
	fmt.Println("  --host HOST       Host to bind to (default: 0.0.0.0)")
	fmt.Println("  --port PORT       Port to bind to (default: 8080)")
	fmt.Println("  --static-dir PATH Static files directory (absolute)")
	fmt.Println("  --help            Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  thd")
	fmt.Println("  thd --daemon --log-file /opt/holo-deck/build/logs/thd.log")
	fmt.Println("  thd --host 127.0.0.1 --port 9090")
	fmt.Println()
	fmt.Printf("DEFAULT PATHS:\n")
	fmt.Printf("  Root: %s\n", THD_ROOT_DIR)
	fmt.Printf("  Static: %s\n", THD_STATIC_DIR)
	fmt.Printf("  Logs: %s\n", THD_LOG_DIR)
	fmt.Printf("  PID: %s\n", THD_PID_FILE)
}

func ensureDirectories() error {
	dirs := []string{THD_BUILD_DIR, THD_BIN_DIR, THD_LOG_DIR, THD_RUNTIME_DIR}
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
		logFile = filepath.Join(THD_LOG_DIR, fmt.Sprintf("thd_%s.log", 
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