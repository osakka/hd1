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
	// Create version hash based on JS file contents and timestamp
	rendererPath := "../share/htdocs/static/js/renderer.js"
	matrixPath := "../share/htdocs/static/js/gl-matrix.js"

	rendererHash := getFileHash(rendererPath)
	matrixHash := getFileHash(matrixPath)

	jsVersion = fmt.Sprintf("%s-%s-%d", 
		rendererHash[:8], 
		matrixHash[:8], 
		time.Now().Unix())
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
	// Update version if files changed recently
	if time.Since(jsVersionTime) > time.Minute {
		updateJSVersion()
	}
	return jsVersion
}

func ReplaceVersionPlaceholder(html string) string {
	return strings.ReplaceAll(html, "${JS_VERSION}", GetJSVersion())
}