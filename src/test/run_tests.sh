#!/bin/bash

# HD1 (Holodeck One) Comprehensive Test Suite Runner
# This script runs all test phases for the HD1 universal 3D interface platform

set -e

# Add Go bin directory to PATH if it exists
if [ -d "/usr/local/go/bin" ]; then
    export PATH="/usr/local/go/bin:$PATH"
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${TEST_DIR}/.." && pwd)"
COVERAGE_DIR="${PROJECT_ROOT}/build/coverage"
REPORTS_DIR="${PROJECT_ROOT}/build/reports"
LOG_FILE="${REPORTS_DIR}/test_run.log"

# Test configuration
TIMEOUT=300s
PARALLEL_TESTS=4
COVERAGE_THRESHOLD=80

# Ensure required directories exist
mkdir -p "${COVERAGE_DIR}" "${REPORTS_DIR}"

# Logging function
log() {
    echo -e "${BLUE}[$(date '+%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "${LOG_FILE}"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "${LOG_FILE}"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "${LOG_FILE}"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "${LOG_FILE}"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if Go is installed
    if ! command_exists go; then
        error "Go is not installed. Please install Go 1.19 or later."
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log "Go version: ${GO_VERSION}"
    
    # Check if required tools are installed
    if ! command_exists gofmt; then
        error "gofmt is not available"
        exit 1
    fi
    
    if ! command_exists go; then
        error "go test is not available"
        exit 1
    fi
    
    # Check if gotestsum is installed (optional but recommended)
    if ! command_exists gotestsum; then
        warning "gotestsum not found. Installing..."
        go install gotest.tools/gotestsum@latest
    fi
    
    success "Prerequisites check passed"
}

# Function to setup test environment
setup_test_environment() {
    log "Setting up test environment..."
    
    # Change to project root
    cd "${PROJECT_ROOT}"
    
    # Clean previous test artifacts
    rm -rf "${COVERAGE_DIR}"/*
    rm -rf "${REPORTS_DIR}"/*
    
    # Create test data directory
    mkdir -p "${PROJECT_ROOT}/test_data"
    
    # Generate test fixtures
    log "Generating test fixtures..."
    go run "${TEST_DIR}/fixtures.go" -generate 2>&1 | tee -a "${LOG_FILE}" || true
    
    # Set test environment variables
    export HD1_TEST_MODE=true
    export HD1_TEST_DB_URL="postgres://test:test@localhost:5432/hd1_test?sslmode=disable"
    export HD1_LOG_LEVEL=DEBUG
    export HD1_STATIC_DIR="${PROJECT_ROOT}/share/htdocs/static"
    
    success "Test environment setup complete"
}

# Function to run database migrations for testing
setup_test_database() {
    log "Setting up test database..."
    
    # Check if PostgreSQL is running
    if ! command_exists psql; then
        error "PostgreSQL is not installed or not in PATH"
        exit 1
    fi
    
    # Create test database (ignore errors if it already exists)
    createdb hd1_test 2>/dev/null || true
    
    # Run migrations
    if [ -f "${PROJECT_ROOT}/database/migrations/schema.sql" ]; then
        psql -d hd1_test -f "${PROJECT_ROOT}/database/migrations/schema.sql" 2>&1 | tee -a "${LOG_FILE}" || true
    fi
    
    success "Test database setup complete"
}

# Function to run a specific test phase
run_test_phase() {
    local phase=$1
    local test_file=$2
    local description=$3
    
    log "Running ${phase}: ${description}"
    
    local coverage_file="${COVERAGE_DIR}/${phase}_coverage.out"
    local junit_file="${REPORTS_DIR}/${phase}_junit.xml"
    
    # Run tests with coverage
    if command_exists gotestsum; then
        gotestsum --junitfile="${junit_file}" --format=standard-verbose -- \
            -coverprofile="${coverage_file}" \
            -covermode=atomic \
            -timeout="${TIMEOUT}" \
            -race \
            -v \
            "${PROJECT_ROOT}/test/${test_file}" 2>&1 | tee -a "${LOG_FILE}"
    else
        go test \
            -coverprofile="${coverage_file}" \
            -covermode=atomic \
            -timeout="${TIMEOUT}" \
            -race \
            -v \
            "${PROJECT_ROOT}/test/${test_file}" 2>&1 | tee -a "${LOG_FILE}"
    fi
    
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        success "${phase} tests passed"
        
        # Generate coverage report
        if [ -f "${coverage_file}" ]; then
            local coverage_percent=$(go tool cover -func="${coverage_file}" | grep total | awk '{print $3}' | sed 's/%//')
            log "${phase} coverage: ${coverage_percent}%"
            
            # Generate HTML coverage report
            go tool cover -html="${coverage_file}" -o "${REPORTS_DIR}/${phase}_coverage.html"
        fi
    else
        error "${phase} tests failed with exit code ${exit_code}"
    fi
    
    return $exit_code
}

# Function to run all test phases
run_all_tests() {
    log "Starting comprehensive test suite for HD1 v0.7.0"
    
    local failed_tests=()
    
    # Phase 1: Foundation Tests
    if ! run_test_phase "phase1" "phase1_foundation_test.go" "Foundation - Authentication, Sessions, Services"; then
        failed_tests+=("Phase 1: Foundation")
    fi
    
    # Phase 2: Collaboration Tests
    if ! run_test_phase "phase2" "phase2_collaboration_test.go" "Collaboration - WebRTC, OT, Assets"; then
        failed_tests+=("Phase 2: Collaboration")
    fi
    
    # Phase 3: AI Integration Tests
    if ! run_test_phase "phase3" "phase3_ai_integration_test.go" "AI Integration - LLM Avatars, Content Generation"; then
        failed_tests+=("Phase 3: AI Integration")
    fi
    
    # Phase 4: Universal Platform Tests
    if ! run_test_phase "phase4" "phase4_universal_platform_test.go" "Universal Platform - Clients, Plugins, Enterprise"; then
        failed_tests+=("Phase 4: Universal Platform")
    fi
    
    # Integration Tests
    if ! run_test_phase "integration" "integration_test.go" "Integration - End-to-End Platform Testing"; then
        failed_tests+=("Integration Tests")
    fi
    
    # Report results
    if [ ${#failed_tests[@]} -eq 0 ]; then
        success "All test phases passed!"
        return 0
    else
        error "Failed test phases: ${failed_tests[*]}"
        return 1
    fi
}

# Function to generate combined coverage report
generate_coverage_report() {
    log "Generating combined coverage report..."
    
    local coverage_files=()
    for coverage_file in "${COVERAGE_DIR}"/*_coverage.out; do
        if [ -f "$coverage_file" ]; then
            coverage_files+=("$coverage_file")
        fi
    done
    
    if [ ${#coverage_files[@]} -eq 0 ]; then
        warning "No coverage files found"
        return
    fi
    
    # Combine coverage files
    local combined_coverage="${COVERAGE_DIR}/combined_coverage.out"
    echo "mode: atomic" > "${combined_coverage}"
    
    for coverage_file in "${coverage_files[@]}"; do
        tail -n +2 "$coverage_file" >> "${combined_coverage}"
    done
    
    # Generate combined coverage report
    local total_coverage=$(go tool cover -func="${combined_coverage}" | grep total | awk '{print $3}' | sed 's/%//')
    log "Total coverage: ${total_coverage}%"
    
    # Generate HTML report
    go tool cover -html="${combined_coverage}" -o "${REPORTS_DIR}/combined_coverage.html"
    
    # Check coverage threshold
    if (( $(echo "${total_coverage} < ${COVERAGE_THRESHOLD}" | bc -l) )); then
        warning "Coverage ${total_coverage}% is below threshold ${COVERAGE_THRESHOLD}%"
    else
        success "Coverage ${total_coverage}% meets threshold ${COVERAGE_THRESHOLD}%"
    fi
}

# Function to run benchmark tests
run_benchmarks() {
    log "Running benchmark tests..."
    
    # Find benchmark tests
    local benchmark_files=()
    for test_file in "${TEST_DIR}"/*_test.go; do
        if grep -q "func Benchmark" "$test_file"; then
            benchmark_files+=("$test_file")
        fi
    done
    
    if [ ${#benchmark_files[@]} -eq 0 ]; then
        warning "No benchmark tests found"
        return
    fi
    
    # Run benchmarks
    local benchmark_report="${REPORTS_DIR}/benchmarks.txt"
    go test -bench=. -benchmem -run=^$ "${benchmark_files[@]}" > "${benchmark_report}" 2>&1
    
    success "Benchmark tests completed. Report saved to ${benchmark_report}"
}

# Function to run race condition tests
run_race_tests() {
    log "Running race condition tests..."
    
    # Run all tests with race detection
    go test -race -short "${PROJECT_ROOT}/test/..." 2>&1 | tee "${REPORTS_DIR}/race_test.log"
    
    if [ $? -eq 0 ]; then
        success "No race conditions detected"
    else
        error "Race conditions detected. Check ${REPORTS_DIR}/race_test.log"
    fi
}

# Function to run memory tests
run_memory_tests() {
    log "Running memory tests..."
    
    # Run tests with memory profiling
    go test -memprofile="${REPORTS_DIR}/mem.prof" -short "${PROJECT_ROOT}/test/..." 2>&1 | tee "${REPORTS_DIR}/memory_test.log"
    
    if [ -f "${REPORTS_DIR}/mem.prof" ]; then
        # Generate memory profile report
        go tool pprof -text "${REPORTS_DIR}/mem.prof" > "${REPORTS_DIR}/memory_profile.txt"
        success "Memory profile generated"
    fi
}

# Function to cleanup test environment
cleanup_test_environment() {
    log "Cleaning up test environment..."
    
    # Clean up test database
    dropdb hd1_test 2>/dev/null || true
    
    # Clean up test data
    rm -rf "${PROJECT_ROOT}/test_data"
    
    # Reset environment variables
    unset HD1_TEST_MODE
    unset HD1_TEST_DB_URL
    unset HD1_LOG_LEVEL
    unset HD1_STATIC_DIR
    
    success "Test environment cleanup complete"
}

# Function to generate test report
generate_test_report() {
    log "Generating test report..."
    
    local report_file="${REPORTS_DIR}/test_report.html"
    
    cat > "${report_file}" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>HD1 Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background-color: #f0f0f0; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; }
        .success { color: green; }
        .error { color: red; }
        .warning { color: orange; }
        pre { background-color: #f5f5f5; padding: 10px; border-radius: 3px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="header">
        <h1>HD1 (Holodeck One) Test Report</h1>
        <p>Generated: $(date)</p>
        <p>Version: v0.7.0</p>
    </div>
    
    <div class="section">
        <h2>Test Results</h2>
        <table>
            <tr><th>Phase</th><th>Status</th><th>Coverage</th></tr>
EOF

    # Add test results to report
    for phase in phase1 phase2 phase3 phase4 integration; do
        local coverage_file="${COVERAGE_DIR}/${phase}_coverage.out"
        local coverage="N/A"
        
        if [ -f "$coverage_file" ]; then
            coverage=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
        fi
        
        echo "            <tr><td>$phase</td><td>✓ Passed</td><td>$coverage</td></tr>" >> "$report_file"
    done
    
    cat >> "${report_file}" << EOF
        </table>
    </div>
    
    <div class="section">
        <h2>Coverage Reports</h2>
        <ul>
            <li><a href="combined_coverage.html">Combined Coverage Report</a></li>
            <li><a href="phase1_coverage.html">Phase 1 Coverage</a></li>
            <li><a href="phase2_coverage.html">Phase 2 Coverage</a></li>
            <li><a href="phase3_coverage.html">Phase 3 Coverage</a></li>
            <li><a href="phase4_coverage.html">Phase 4 Coverage</a></li>
        </ul>
    </div>
    
    <div class="section">
        <h2>Additional Reports</h2>
        <ul>
            <li><a href="benchmarks.txt">Benchmark Results</a></li>
            <li><a href="race_test.log">Race Condition Tests</a></li>
            <li><a href="memory_test.log">Memory Tests</a></li>
        </ul>
    </div>
    
    <div class="section">
        <h2>Log Files</h2>
        <pre>$(tail -n 50 "${LOG_FILE}")</pre>
    </div>
</body>
</html>
EOF

    success "Test report generated: ${report_file}"
}

# Main function
main() {
    local start_time=$(date +%s)
    
    log "HD1 (Holodeck One) Comprehensive Test Suite"
    log "=========================================="
    
    # Initialize log file
    echo "HD1 Test Suite Run - $(date)" > "${LOG_FILE}"
    
    # Check command line arguments
    case "${1:-all}" in
        "phase1")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_test_phase "phase1" "phase1_foundation_test.go" "Foundation Tests"
            ;;
        "phase2")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_test_phase "phase2" "phase2_collaboration_test.go" "Collaboration Tests"
            ;;
        "phase3")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_test_phase "phase3" "phase3_ai_integration_test.go" "AI Integration Tests"
            ;;
        "phase4")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_test_phase "phase4" "phase4_universal_platform_test.go" "Universal Platform Tests"
            ;;
        "integration")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_test_phase "integration" "integration_test.go" "Integration Tests"
            ;;
        "benchmarks")
            check_prerequisites
            setup_test_environment
            run_benchmarks
            ;;
        "race")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_race_tests
            ;;
        "memory")
            check_prerequisites
            setup_test_environment
            setup_test_database
            run_memory_tests
            ;;
        "all"|*)
            check_prerequisites
            setup_test_environment
            setup_test_database
            
            if run_all_tests; then
                generate_coverage_report
                run_benchmarks
                run_race_tests
                run_memory_tests
                generate_test_report
                success "All tests completed successfully!"
            else
                error "Some tests failed. Check the reports for details."
                exit 1
            fi
            ;;
    esac
    
    # Cleanup
    cleanup_test_environment
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    log "Test suite completed in ${duration} seconds"
    success "Test artifacts saved to ${REPORTS_DIR}"
}

# Run main function with all arguments
main "$@"