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
	// Create version hash based on all frontend files
	rendererPath := "../share/htdocs/static/js/renderer.js"
	matrixPath := "../share/htdocs/static/js/gl-matrix.js"
	aframePath := "../share/htdocs/static/js/hd1-aframe.js"
	consoleJSPath := "../share/htdocs/static/js/hd1-console.js"
	consoleCSSPath := "../share/htdocs/static/css/hd1-console.css"
	indexPath := "../share/htdocs/index.html"
	handlersPath := "server/handlers.go"

	rendererHash := getFileHash(rendererPath)
	matrixHash := getFileHash(matrixPath)
	aframeHash := getFileHash(aframePath)
	consoleJSHash := getFileHash(consoleJSPath)
	consoleCSSHash := getFileHash(consoleCSSPath)
	indexHash := getFileHash(indexPath)
	handlersHash := getFileHash(handlersPath)

	jsVersion = fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s", 
		rendererHash[:8], 
		matrixHash[:8],
		aframeHash[:8],
		consoleJSHash[:8],
		consoleCSSHash[:8],
		indexHash[:8],
		handlersHash[:8])
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