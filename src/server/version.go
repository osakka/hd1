package server

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	jsVersion     string
	jsVersionTime time.Time
)

func init() {
	updateJSVersion()
}

func updateJSVersion() {
	// API-FIRST DEVELOPMENT: Version driven by API specification (single source of truth)
	apiSpecPath := "api.yaml"
	
	// Include generated files that derive from API spec
	consoleJSPath := "../share/htdocs/static/js/hd1-console.js"
	threeJSPath := "../share/htdocs/static/js/hd1-threejs.js"
	indexPath := "../share/htdocs/index.html"

	// PRIMARY: API specification hash (single source of truth)
	apiSpecHash := getFileHash(apiSpecPath)
	
	// SECONDARY: Generated artifacts that should match API spec
	consoleJSHash := getFileHash(consoleJSPath)
	threeJSHash := getFileHash(threeJSPath)
	indexHash := getFileHash(indexPath)

	// Version format: API-spec-hash + generated-artifacts
	jsVersion = fmt.Sprintf("%s-%s-%s-%s", 
		apiSpecHash[:8],   // API specification drives everything
		consoleJSHash[:8], // Console UI (generated)
		threeJSHash[:8],   // Three.js integration
		indexHash[:8])     // Main UI template
	jsVersionTime = time.Now()
}

func getFileHash(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		return "00000000000000000000000000000000" // 32 chars for safe slicing
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "00000000000000000000000000000000" // 32 chars for safe slicing
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetJSVersion() string {
	// Only update version when explicitly requested or if files actually changed
	// Remove automatic time-based updates to prevent unnecessary refreshes
	return jsVersion
}

func ReplaceVersionPlaceholder(html string) string {
	return strings.ReplaceAll(html, "${JS_VERSION}", GetJSVersion())
}