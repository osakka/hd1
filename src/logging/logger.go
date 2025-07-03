package logging

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogLevel represents logging levels
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// Log rotation constants
const (
	DefaultMaxLogSize = 10 * 1024 * 1024 // 10MB
	DefaultMaxRotations = 3              // Keep 3 rotated logs
)

var levelNames = map[LogLevel]string{
	TRACE: "TRACE",
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

var levelFromString = map[string]LogLevel{
	"TRACE": TRACE,
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"FATAL": FATAL,
}

// Logger provides unified logging for HD1 system
type Logger struct {
	level       LogLevel
	traceModules map[string]bool
	file        *os.File
	mu          sync.RWMutex
	processID   int
	logPath     string
	maxSize     int64 // Maximum log file size in bytes
	maxRotations int   // Maximum number of rotated log files
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	ProcessID int                    `json:"process_id"`
	ThreadID  string                 `json:"thread_id"`
	Level     string                 `json:"level"`
	Function  string                 `json:"function"`
	File      string                 `json:"file"`
	Line      int                    `json:"line"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// InitLogger initializes the global logger
func InitLogger(logDir string, level LogLevel, traceModules []string) error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewLogger(logDir, level, traceModules)
	})
	return err
}

// NewLogger creates a new logger instance
func NewLogger(logDir string, level LogLevel, traceModules []string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := filepath.Join(logDir, "hd1.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	traceMap := make(map[string]bool)
	for _, module := range traceModules {
		traceMap[strings.ToLower(module)] = true
	}

	return &Logger{
		level:        level,
		traceModules: traceMap,
		file:         file,
		processID:    os.Getpid(),
		logPath:      logFile,
		maxSize:      DefaultMaxLogSize,
		maxRotations: DefaultMaxRotations,
	}, nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if defaultLogger == nil {
		// Fallback to stderr if not initialized
		logger, _ := NewLogger("", INFO, nil)
		return logger
	}
	return defaultLogger
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetLevelFromString sets the logging level from string
func (l *Logger) SetLevelFromString(levelStr string) error {
	level, exists := levelFromString[strings.ToUpper(levelStr)]
	if !exists {
		return fmt.Errorf("invalid log level: %s", levelStr)
	}
	l.SetLevel(level)
	return nil
}

// EnableTrace enables tracing for specific modules
func (l *Logger) EnableTrace(modules []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, module := range modules {
		l.traceModules[strings.ToLower(module)] = true
	}
}

// DisableTrace disables tracing for specific modules
func (l *Logger) DisableTrace(modules []string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, module := range modules {
		delete(l.traceModules, strings.ToLower(module))
	}
}

// log is the core logging function
func (l *Logger) log(level LogLevel, message string, data map[string]interface{}) {
	// Thread-safe level check with single lock acquisition
	l.mu.RLock()
	currentLevel := l.level
	enabled := level >= currentLevel
	l.mu.RUnlock()

	// Zero-overhead quick return if level not enabled
	if !enabled {
		return
	}

	// Get caller information
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	funcName := "unknown"
	if fn := runtime.FuncForPC(pc); fn != nil {
		funcName = filepath.Base(fn.Name())
	}

	fileName := filepath.Base(file)
	fileNameNoExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		ProcessID: l.processID,
		ThreadID:  getThreadID(),
		Level:     levelNames[level],
		Function:  funcName,
		File:      fileNameNoExt,
		Line:      line,
		Message:   message,
		Data:      data,
	}

	l.writeEntry(entry, level)
}

// Trace logs trace level messages for specific modules
func (l *Logger) Trace(module, message string, data ...map[string]interface{}) {
	l.mu.RLock()
	enabled := l.traceModules[strings.ToLower(module)]
	l.mu.RUnlock()

	if !enabled {
		return
	}

	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	if dataMap == nil {
		dataMap = make(map[string]interface{})
	}
	dataMap["trace_module"] = module

	l.log(TRACE, message, dataMap)
}

// Debug logs debug level messages
func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.log(DEBUG, message, dataMap)
}

// Info logs info level messages
func (l *Logger) Info(message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.log(INFO, message, dataMap)
}

// Warn logs warning level messages
func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.log(WARN, message, dataMap)
}

// Error logs error level messages
func (l *Logger) Error(message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.log(ERROR, message, dataMap)
}

// Fatal logs fatal level messages and exits
func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	var dataMap map[string]interface{}
	if len(data) > 0 {
		dataMap = data[0]
	}
	l.log(FATAL, message, dataMap)
	os.Exit(1)
}

// writeEntry writes the log entry to file and console
func (l *Logger) writeEntry(entry LogEntry, level LogLevel) {
	// Format for console (human readable)
	consoleMsg := fmt.Sprintf("%s [%d:%s] [%s] %s.%s:%d %s",
		entry.Timestamp[:19], // Truncate nanoseconds for console
		entry.ProcessID,
		entry.ThreadID,
		entry.Level,
		entry.Function,
		entry.File,
		entry.Line,
		entry.Message,
	)

	// Add data if present
	if len(entry.Data) > 0 {
		dataStr, _ := json.Marshal(entry.Data)
		consoleMsg += " " + string(dataStr)
	}

	// Write to console (stderr for errors, stdout for others)
	if level >= ERROR {
		fmt.Fprintln(os.Stderr, consoleMsg)
	} else {
		fmt.Fprintln(os.Stdout, consoleMsg)
	}

	// Write JSON to file if available
	if l.file != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		if jsonData, err := json.Marshal(entry); err == nil {
			l.file.Write(jsonData)
			l.file.Write([]byte("\n"))
			
			// Check if log rotation is needed
			l.checkRotation()
		}
	}
}

// getThreadID returns the current goroutine ID
func getThreadID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	
	// Parse goroutine ID from stack trace: "goroutine 1 [running]:"
	stack := string(buf[:n])
	if idx := strings.Index(stack, " "); idx > 0 && idx >= 10 {
		if gid := stack[10:idx]; gid != "" {
			return gid
		}
	}
	
	// Fallback to "main" if parsing fails
	return "main"
}

// Close closes the logger
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Zero-overhead convenience functions using default logger
func Trace(module, message string, data ...map[string]interface{}) {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.traceModules[strings.ToLower(module)]
	logger.mu.RUnlock()
	
	if enabled {
		logger.Trace(module, message, data...)
	}
}

func Debug(message string, data ...map[string]interface{}) {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.level <= DEBUG
	logger.mu.RUnlock()
	
	if enabled {
		logger.Debug(message, data...)
	}
}

func Info(message string, data ...map[string]interface{}) {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.level <= INFO
	logger.mu.RUnlock()
	
	if enabled {
		logger.Info(message, data...)
	}
}

func Warn(message string, data ...map[string]interface{}) {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.level <= WARN
	logger.mu.RUnlock()
	
	if enabled {
		logger.Warn(message, data...)
	}
}

func Error(message string, data ...map[string]interface{}) {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.level <= ERROR
	logger.mu.RUnlock()
	
	if enabled {
		logger.Error(message, data...)
	}
}

func Fatal(message string, data ...map[string]interface{}) {
	// FATAL always logs regardless of level
	GetLogger().Fatal(message, data...)
}

func SetLevel(level LogLevel) {
	GetLogger().SetLevel(level)
}

func SetLevelFromString(levelStr string) error {
	return GetLogger().SetLevelFromString(levelStr)
}

// Zero-overhead level checking functions for conditional logging
func IsTraceEnabled(module string) bool {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.traceModules[strings.ToLower(module)]
	logger.mu.RUnlock()
	return enabled
}

func IsDebugEnabled() bool {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.level <= DEBUG
	logger.mu.RUnlock()
	return enabled
}

func IsInfoEnabled() bool {
	logger := GetLogger()
	logger.mu.RLock()
	enabled := logger.level <= INFO
	logger.mu.RUnlock()
	return enabled
}

func EnableTrace(modules []string) {
	GetLogger().EnableTrace(modules)
}

func DisableTrace(modules []string) {
	GetLogger().DisableTrace(modules)
}

// checkRotation checks if log rotation is needed and performs it
func (l *Logger) checkRotation() {
	if l.file == nil || l.logPath == "" {
		return
	}

	// Get current file size
	stat, err := l.file.Stat()
	if err != nil {
		return
	}

	// Check if rotation is needed
	if stat.Size() >= l.maxSize {
		l.rotateLog()
	}
}

// rotateLog performs log rotation
func (l *Logger) rotateLog() {
	// Close current file
	l.file.Close()

	// Rotate existing log files
	for i := l.maxRotations; i > 1; i-- {
		oldPath := fmt.Sprintf("%s.%d", l.logPath, i-1)
		newPath := fmt.Sprintf("%s.%d", l.logPath, i)
		
		// Remove the oldest log if it exists
		if i == l.maxRotations {
			os.Remove(newPath)
		}
		
		// Move log files
		os.Rename(oldPath, newPath)
	}

	// Move current log to .1
	os.Rename(l.logPath, l.logPath+".1")

	// Create new log file
	file, err := os.OpenFile(l.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// Fallback to stderr if we can't create new log
		l.file = nil
		return
	}

	l.file = file

	// Log rotation event
	l.Info("log rotation completed", map[string]interface{}{
		"max_size_mb": l.maxSize / (1024 * 1024),
		"max_rotations": l.maxRotations,
	})
}

// ReadLogEntries reads the last N log entries from the log file
func ReadLogEntries(count int) ([]LogEntry, error) {
	logger := GetLogger()
	if logger.logPath == "" {
		return nil, fmt.Errorf("log file path not configured")
	}

	// Read from current log file
	entries := make([]LogEntry, 0, count)
	entriesRead := 0

	// Try to read from current log file first
	if fileEntries, err := readEntriesFromFile(logger.logPath, count); err == nil {
		entries = append(entries, fileEntries...)
		entriesRead = len(fileEntries)
	}

	// If we need more entries, read from rotated logs
	for i := 1; i <= logger.maxRotations && entriesRead < count; i++ {
		rotatedPath := fmt.Sprintf("%s.%d", logger.logPath, i)
		remaining := count - entriesRead
		
		if fileEntries, err := readEntriesFromFile(rotatedPath, remaining); err == nil {
			// Prepend older entries
			entries = append(fileEntries, entries...)
			entriesRead += len(fileEntries)
		}
	}

	// Return the last N entries
	if len(entries) > count {
		entries = entries[len(entries)-count:]
	}

	return entries, nil
}

// readEntriesFromFile reads log entries from a specific file
func readEntriesFromFile(filePath string, maxCount int) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	
	// Read all lines and keep only the last maxCount
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Take the last maxCount lines
	startIndex := 0
	if len(lines) > maxCount {
		startIndex = len(lines) - maxCount
	}

	// Parse JSON entries
	for i := startIndex; i < len(lines); i++ {
		var entry LogEntry
		if err := json.Unmarshal([]byte(lines[i]), &entry); err == nil {
			entries = append(entries, entry)
		}
	}

	return entries, scanner.Err()
}