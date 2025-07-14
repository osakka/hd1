# HD1 (Holodeck One) - Development Build System
# Single source of truth: api.yaml drives Three.js transformation

.PHONY: all clean generate build test run validate client logs status start stop restart daemon-start daemon-stop daemon-status

# Build directory structure - Configuration-driven
BUILD_DIR = $(shell test -n "$$HD1_BUILD_DIR" && echo "$$HD1_BUILD_DIR" || echo "../build")
BIN_DIR = $(BUILD_DIR)/bin
LOG_DIR = $(shell test -n "$$HD1_LOG_DIR" && echo "$$HD1_LOG_DIR" || echo "$(BUILD_DIR)/logs")
RUNTIME_DIR = $(shell test -n "$$HD1_RUNTIME_DIR" && echo "$$HD1_RUNTIME_DIR" || echo "$(BUILD_DIR)/runtime")
SHARE_DIR = $(shell test -n "$$HD1_STATIC_DIR" && dirname "$$HD1_STATIC_DIR" || echo "../share")
CHANNELS_DIR = $(shell test -n "$$HD1_CHANNELS_DIR" && echo "$$HD1_CHANNELS_DIR" || echo "../share/channels")
AVATARS_DIR = $(shell test -n "$$HD1_AVATARS_DIR" && echo "$$HD1_AVATARS_DIR" || echo "../share/avatars")

# Default target
all: validate generate build

# Generate router from Three.js specification
generate:
	@echo "GENERATING THREE.JS CODE FROM SPECIFICATION..."
	go run codegen/generator.go
	@echo "Three.js auto-router generated from api.yaml"
	@echo "DOWNLOADING THREE.JS LIBRARY..."
	@mkdir -p $(SHARE_DIR)/htdocs/static/vendor/threejs
	@if [ ! -f $(SHARE_DIR)/htdocs/static/vendor/threejs/three.min.js ]; then \
		curl -L -o $(SHARE_DIR)/htdocs/static/vendor/threejs/three.min.js https://cdn.jsdelivr.net/npm/three@0.170.0/build/three.min.js; \
		echo "Three.js r170 downloaded to static directory"; \
	else \
		echo "Three.js already exists in static directory"; \
	fi

# Validate Three.js specification
validate:
	@echo "VALIDATING THREE.JS API SPECIFICATION..."
	@if [ ! -f api.yaml ]; then echo "FATAL: api.yaml missing!"; exit 1; fi
	@echo "Three.js specification found"

# Build server with artifact management
build: generate
	@echo "BUILDING HD1 SERVER..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/hd1 .
	@if [ -f $(BIN_DIR)/hd1 ]; then echo "HD1 server built -> $(BIN_DIR)/hd1"; else echo "Build failed"; exit 1; fi

# Test API endpoints
test: build
	@echo "TESTING API ENDPOINTS..."
	@mkdir -p $(LOG_DIR)
	@$(BIN_DIR)/hd1 > $(LOG_DIR)/test.log 2>&1 & \
	SERVER_PID=$$!; \
	sleep 3; \
	API_BASE=$${HD1_API_BASE:-http://$${HD1_HOST:-0.0.0.0}:$${HD1_PORT:-8080}/api}; \
	curl -s "$$API_BASE/sessions" | jq '.' > /dev/null && echo "GET /api/sessions - OK" || echo "GET /api/sessions - FAILED"; \
	curl -s -X POST "$$API_BASE/sessions" | jq '.' > /dev/null && echo "POST /api/sessions - OK" || echo "POST /api/sessions - FAILED"; \
	kill $$SERVER_PID 2>/dev/null || true; \
	echo "API tests complete"

# Run server with logging
run: build
	@echo "STARTING HD1 (Holodeck One)..."
	@mkdir -p $(LOG_DIR) $(RUNTIME_DIR)
	$(BIN_DIR)/hd1

# Create client wrapper with paths
client:
	@echo "CREATING HD1 API CLIENT..."
	@mkdir -p $(BIN_DIR)
	@echo '#!/bin/bash' > $(BIN_DIR)/hd1-client
	@echo '# HD1 (Holodeck One) API Client' >> $(BIN_DIR)/hd1-client
	@echo 'API_BASE="$${HD1_API_BASE:-http://$${HD1_HOST:-0.0.0.0}:$${HD1_PORT:-8080}/api}"' >> $(BIN_DIR)/hd1-client
	@echo 'case "$${1:-help}" in' >> $(BIN_DIR)/hd1-client
	@echo '  "sessions") curl -s "$$API_BASE/sessions" | jq "." ;;' >> $(BIN_DIR)/hd1-client
	@echo '  "create-session") curl -s -X POST "$$API_BASE/sessions" | jq ".session_id" ;;' >> $(BIN_DIR)/hd1-client
	@echo '  "list") curl -s "$$API_BASE/sessions" | jq ".sessions[].id" ;;' >> $(BIN_DIR)/hd1-client
	@echo '  *) echo "Usage: hd1-client sessions|create-session|list" ;;' >> $(BIN_DIR)/hd1-client
	@echo 'esac' >> $(BIN_DIR)/hd1-client
	@chmod +x $(BIN_DIR)/hd1-client
	@echo "HD1 client created -> $(BIN_DIR)/hd1-client"

# Setup Three.js web resources
web:
	@echo "SETTING UP THREE.JS WEB RESOURCES..."
	@mkdir -p $(SHARE_DIR)/htdocs/static
	@$(MAKE) generate
	@echo "Three.js web resources organized in $(SHARE_DIR)/htdocs/"
	@echo "Note: Three.js files are served directly from $(SHARE_DIR)/htdocs/static/"

# Development status
build-status:
	@echo "HD1 BUILD STATUS"
	@echo "================"
	@if [ -f $(BIN_DIR)/hd1 ]; then echo "HD1 Binary: EXISTS ($(BIN_DIR)/hd1)"; else echo "HD1 Binary: MISSING"; fi
	@if [ -f $(BIN_DIR)/hd1-client ]; then echo "HD1 Client: EXISTS ($(BIN_DIR)/hd1-client)"; else echo "HD1 Client: MISSING"; fi
	@if [ -f api.yaml ]; then echo "Three.js API Spec: EXISTS"; else echo "Three.js API Spec: MISSING"; fi
	@if [ -f auto_router.go ]; then echo "Three.js Auto Router: GENERATED"; else echo "Three.js Auto Router: NOT GENERATED"; fi
	@echo "Build artifacts in: $(BUILD_DIR)/"

# Show recent logs
logs:
	@echo "HD1 RECENT LOGS"
	@echo "==============="
	@if [ -f $(LOG_DIR)/test.log ]; then echo "Test logs:"; tail -10 $(LOG_DIR)/test.log; fi
	@if [ -d $(LOG_DIR) ]; then echo "Available logs:"; ls -la $(LOG_DIR)/; fi

# Clean all build artifacts
clean:
	@echo "CLEANING HD1 THREE.JS BUILD ARTIFACTS..."
	@rm -rf $(BUILD_DIR)/bin/hd1 $(BUILD_DIR)/bin/hd1-client
	@rm -f auto_router.go
	@echo "Clean complete"

# Deep clean - remove all build directories
deep-clean:
	@echo "DEEP CLEANING HD1 THREE.JS WORKSPACE..."
	@rm -rf $(BUILD_DIR)
	@rm -f auto_router.go
	@echo "Deep clean complete"

# Daemon control
start: build
	@echo "Starting HD1 daemon..."
	@mkdir -p $(RUNTIME_DIR)
	@if [ -f $(RUNTIME_DIR)/hd1.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/hd1.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "ERROR: HD1 daemon already running (PID: $$PID)"; \
			exit 1; \
		else \
			echo "Cleaning stale PID file"; \
			rm -f $(RUNTIME_DIR)/hd1.pid; \
		fi \
	fi
	@HD1_HOST=$${HD1_HOST:-0.0.0.0}; HD1_PORT=$${HD1_PORT:-8080}; \
	echo "Checking port $$HD1_PORT..."; \
	if netstat -tuln 2>/dev/null | grep -q ":$$HD1_PORT "; then \
		echo "ERROR: Port $$HD1_PORT already in use"; \
		echo "Process using port $$HD1_PORT:"; \
		lsof -ti:$$HD1_PORT | head -5; \
		echo "Run 'make force-stop' to kill all processes on port $$HD1_PORT"; \
		exit 1; \
	fi; \
	echo "Port $$HD1_PORT available"
	@HD1_HOST=$${HD1_HOST:-0.0.0.0} HD1_PORT=$${HD1_PORT:-8080} $(BIN_DIR)/hd1 --daemon
	@sleep 2
	@if [ -f $(RUNTIME_DIR)/hd1.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/hd1.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "HD1 daemon started successfully (PID: $$PID)"; \
		else \
			echo "ERROR: HD1 daemon failed to start"; \
			exit 1; \
		fi \
	else \
		echo "ERROR: HD1 daemon failed to create PID file"; \
		exit 1; \
	fi

stop:
	@echo "Stopping HD1 daemon..."
	@if [ -f $(RUNTIME_DIR)/hd1.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/hd1.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			kill $$PID; \
			sleep 2; \
			if kill -0 $$PID 2>/dev/null; then \
				echo "Daemon still running, forcing termination..."; \
				kill -9 $$PID; \
				sleep 1; \
			fi; \
			echo "HD1 daemon stopped (PID: $$PID)"; \
		else \
			echo "WARNING: Process not running (stale PID file)"; \
		fi; \
		rm -f $(RUNTIME_DIR)/hd1.pid; \
	else \
		echo "WARNING: No PID file found - daemon not running"; \
	fi

force-stop:
	@echo "Force stopping all HD1 processes..."
	@pkill -f hd1 || echo "No HD1 processes found"
	@HD1_HOST=$${HD1_HOST:-0.0.0.0}; HD1_PORT=$${HD1_PORT:-8080}; \
	lsof -ti:$$HD1_PORT | xargs -r kill -9 || echo "No processes on port $$HD1_PORT"
	@rm -f $(RUNTIME_DIR)/hd1.pid
	@echo "Force stop complete"

restart: stop
	@sleep 1
	@$(MAKE) start
	@echo "HD1 daemon restarted"

status:
	@echo "HD1 DAEMON STATUS"
	@echo "================="
	@if [ -f $(RUNTIME_DIR)/hd1.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/hd1.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "STATUS: HD1 daemon is running (PID: $$PID)"; \
			HD1_HOST=$${HD1_HOST:-0.0.0.0}; HD1_PORT=$${HD1_PORT:-8080}; \
			echo "PORT: Checking port $$HD1_PORT..."; \
			if netstat -tuln 2>/dev/null | grep -q ":$$HD1_PORT "; then \
				echo "PORT: Listening on port $$HD1_PORT"; \
			else \
				echo "WARNING: Not listening on port $$HD1_PORT"; \
			fi \
		else \
			echo "ERROR: HD1 daemon is not running (stale PID file)"; \
			rm -f $(RUNTIME_DIR)/hd1.pid; \
		fi \
	else \
		echo "STATUS: HD1 daemon is not running"; \
	fi
	@echo "PROCESSES:"
	@ps aux | grep -v grep | grep hd1 || echo "No HD1 processes found"

# Show all targets
help:
	@echo "HD1 (Holodeck One) Build System"
	@echo "==========================================="
	@echo "Core targets:"
	@echo "  make all       - Complete build pipeline"
	@echo "  make build     - Build HD1 server"
	@echo "  make run       - Start HD1 server (foreground)"
	@echo "  make test      - Test API endpoints"
	@echo ""
	@echo "Daemon control:"
	@echo "  make start     - Start HD1 daemon"
	@echo "  make stop      - Stop HD1 daemon"
	@echo "  make restart   - Restart HD1 daemon"
	@echo "  make status    - Show daemon status"
	@echo ""
	@echo "Development targets:"
	@echo "  make generate  - Generate Three.js router from api.yaml"
	@echo "  make client    - Create HD1 API client"
	@echo "  make web       - Setup web resources"
	@echo "  make build-status - Show build status"
	@echo "  make logs      - Show recent logs"
	@echo ""
	@echo "Maintenance targets:"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deep-clean- Remove all build directories"
	@echo "  make help      - Show this help"