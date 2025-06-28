# THD (The Holo-Deck) - Professional Build System
# Single source of truth: api.yaml drives everything

.PHONY: all clean generate build test run validate client logs status start stop restart daemon-start daemon-stop daemon-status

# Professional build directory structure
BUILD_DIR = ../build
BIN_DIR = $(BUILD_DIR)/bin
LOG_DIR = $(BUILD_DIR)/logs
RUNTIME_DIR = $(BUILD_DIR)/runtime
SHARE_DIR = ../share

# Default target
all: validate generate build

# Generate router from specification
generate:
	@echo "GENERATING CODE FROM SPECIFICATION..."
	go run codegen/generator.go
	@echo "Auto-router generated from api.yaml"

# Validate specification
validate:
	@echo "VALIDATING API SPECIFICATION..."
	@if [ ! -f api.yaml ]; then echo "FATAL: api.yaml missing!"; exit 1; fi
	@echo "Specification found"

# Build server with professional artifact management
build: generate
	@echo "BUILDING THD SERVER..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/thd .
	@if [ -f $(BIN_DIR)/thd ]; then echo "THD server built -> $(BIN_DIR)/thd"; else echo "Build failed"; exit 1; fi

# Test API endpoints
test: build
	@echo "TESTING API ENDPOINTS..."
	@mkdir -p $(LOG_DIR)
	@$(BIN_DIR)/thd > $(LOG_DIR)/test.log 2>&1 & \
	SERVER_PID=$$!; \
	sleep 3; \
	curl -s http://localhost:8080/api/sessions | jq '.' > /dev/null && echo "GET /api/sessions - OK" || echo "GET /api/sessions - FAILED"; \
	curl -s -X POST http://localhost:8080/api/sessions | jq '.' > /dev/null && echo "POST /api/sessions - OK" || echo "POST /api/sessions - FAILED"; \
	kill $$SERVER_PID 2>/dev/null || true; \
	echo "API tests complete"

# Run server with professional logging
run: build
	@echo "STARTING THD (The Holo-Deck)..."
	@mkdir -p $(LOG_DIR) $(RUNTIME_DIR)
	$(BIN_DIR)/thd

# Create client wrapper with professional paths
client:
	@echo "CREATING THD API CLIENT..."
	@mkdir -p $(BIN_DIR)
	@echo '#!/bin/bash' > $(BIN_DIR)/thd-client
	@echo '# THD (The Holo-Deck) API Client' >> $(BIN_DIR)/thd-client
	@echo 'API_BASE="http://localhost:8080/api"' >> $(BIN_DIR)/thd-client
	@echo 'case "$${1:-help}" in' >> $(BIN_DIR)/thd-client
	@echo '  "sessions") curl -s "$$API_BASE/sessions" | jq "." ;;' >> $(BIN_DIR)/thd-client
	@echo '  "create-session") curl -s -X POST "$$API_BASE/sessions" | jq ".session_id" ;;' >> $(BIN_DIR)/thd-client
	@echo '  "list") curl -s "$$API_BASE/sessions" | jq ".sessions[].id" ;;' >> $(BIN_DIR)/thd-client
	@echo '  *) echo "Usage: thd-client sessions|create-session|list" ;;' >> $(BIN_DIR)/thd-client
	@echo 'esac' >> $(BIN_DIR)/thd-client
	@chmod +x $(BIN_DIR)/thd-client
	@echo "THD client created -> $(BIN_DIR)/thd-client"

# Setup web resources
web:
	@echo "SETTING UP WEB RESOURCES..."
	@mkdir -p $(SHARE_DIR)/htdocs/static
	@if [ -d renderer/static ]; then cp -r renderer/static/* $(SHARE_DIR)/htdocs/static/; fi
	@echo "Web resources organized in $(SHARE_DIR)/htdocs/"

# Professional development status
build-status:
	@echo "THD BUILD STATUS"
	@echo "================"
	@if [ -f $(BIN_DIR)/thd ]; then echo "THD Binary: EXISTS ($(BIN_DIR)/thd)"; else echo "THD Binary: MISSING"; fi
	@if [ -f $(BIN_DIR)/thd-client ]; then echo "THD Client: EXISTS ($(BIN_DIR)/thd-client)"; else echo "THD Client: MISSING"; fi
	@if [ -f api.yaml ]; then echo "API Spec: EXISTS"; else echo "API Spec: MISSING"; fi
	@if [ -f auto_router.go ]; then echo "Auto Router: GENERATED"; else echo "Auto Router: NOT GENERATED"; fi
	@echo "Build artifacts in: $(BUILD_DIR)/"

# Show recent logs
logs:
	@echo "THD RECENT LOGS"
	@echo "==============="
	@if [ -f $(LOG_DIR)/test.log ]; then echo "Test logs:"; tail -10 $(LOG_DIR)/test.log; fi
	@if [ -d $(LOG_DIR) ]; then echo "Available logs:"; ls -la $(LOG_DIR)/; fi

# Clean all build artifacts
clean:
	@echo "CLEANING THD BUILD ARTIFACTS..."
	@rm -rf $(BUILD_DIR)/bin/thd $(BUILD_DIR)/bin/thd-client
	@rm -f auto_router.go
	@echo "Clean complete"

# Deep clean - remove all build directories
deep-clean:
	@echo "DEEP CLEANING THD WORKSPACE..."
	@rm -rf $(BUILD_DIR)
	@rm -f auto_router.go
	@echo "Deep clean complete"

# Professional daemon control
start: build
	@echo "Starting THD daemon..."
	@mkdir -p $(RUNTIME_DIR)
	@if [ -f $(RUNTIME_DIR)/thd.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/thd.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "ERROR: THD daemon already running (PID: $$PID)"; \
			exit 1; \
		else \
			echo "Cleaning stale PID file"; \
			rm -f $(RUNTIME_DIR)/thd.pid; \
		fi \
	fi
	@echo "Checking port 8080..."
	@if netstat -tuln 2>/dev/null | grep -q ":8080 "; then \
		echo "ERROR: Port 8080 already in use"; \
		echo "Process using port 8080:"; \
		lsof -ti:8080 | head -5; \
		echo "Run 'make force-stop' to kill all processes on port 8080"; \
		exit 1; \
	fi
	@echo "Port 8080 available"
	@$(BIN_DIR)/thd --daemon
	@sleep 2
	@if [ -f $(RUNTIME_DIR)/thd.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/thd.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "THD daemon started successfully (PID: $$PID)"; \
		else \
			echo "ERROR: THD daemon failed to start"; \
			exit 1; \
		fi \
	else \
		echo "ERROR: THD daemon failed to create PID file"; \
		exit 1; \
	fi

stop:
	@echo "Stopping THD daemon..."
	@if [ -f $(RUNTIME_DIR)/thd.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/thd.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			kill $$PID; \
			sleep 2; \
			if kill -0 $$PID 2>/dev/null; then \
				echo "Daemon still running, forcing termination..."; \
				kill -9 $$PID; \
				sleep 1; \
			fi; \
			echo "THD daemon stopped (PID: $$PID)"; \
		else \
			echo "WARNING: Process not running (stale PID file)"; \
		fi; \
		rm -f $(RUNTIME_DIR)/thd.pid; \
	else \
		echo "WARNING: No PID file found - daemon not running"; \
	fi

force-stop:
	@echo "Force stopping all THD processes..."
	@pkill -f thd || echo "No THD processes found"
	@lsof -ti:8080 | xargs -r kill -9 || echo "No processes on port 8080"
	@rm -f $(RUNTIME_DIR)/thd.pid
	@echo "Force stop complete"

restart: stop
	@sleep 1
	@$(MAKE) start
	@echo "THD daemon restarted"

status:
	@echo "THD DAEMON STATUS"
	@echo "================="
	@if [ -f $(RUNTIME_DIR)/thd.pid ]; then \
		PID=$$(cat $(RUNTIME_DIR)/thd.pid); \
		if kill -0 $$PID 2>/dev/null; then \
			echo "STATUS: THD daemon is running (PID: $$PID)"; \
			echo "PORT: Checking port 8080..."; \
			if netstat -tuln 2>/dev/null | grep -q ":8080 "; then \
				echo "PORT: Listening on port 8080"; \
			else \
				echo "WARNING: Not listening on port 8080"; \
			fi \
		else \
			echo "ERROR: THD daemon is not running (stale PID file)"; \
			rm -f $(RUNTIME_DIR)/thd.pid; \
		fi \
	else \
		echo "STATUS: THD daemon is not running"; \
	fi
	@echo "PROCESSES:"
	@ps aux | grep -v grep | grep thd || echo "No THD processes found"

# Show all targets
help:
	@echo "THD (The Holo-Deck) Build System"
	@echo "==========================================="
	@echo "Core targets:"
	@echo "  make all       - Complete build pipeline"
	@echo "  make build     - Build THD server"
	@echo "  make run       - Start THD server (foreground)"
	@echo "  make test      - Test API endpoints"
	@echo ""
	@echo "Daemon control:"
	@echo "  make start     - Start THD daemon"
	@echo "  make stop      - Stop THD daemon"
	@echo "  make restart   - Restart THD daemon"
	@echo "  make daemon-status - Show daemon status"
	@echo ""
	@echo "Development targets:"
	@echo "  make generate  - Generate router from api.yaml"
	@echo "  make client    - Create THD API client"
	@echo "  make web       - Setup web resources"
	@echo "  make status    - Show build status"
	@echo "  make logs      - Show recent logs"
	@echo ""
	@echo "Maintenance targets:"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make deep-clean- Remove all build directories"
	@echo "  make help      - Show this help"