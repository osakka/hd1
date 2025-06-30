package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// TemplateProcessor handles server-side template variable replacement
type TemplateProcessor struct {
	staticDir string
	htdocsDir string
}

// NewTemplateProcessor creates a new template processor
func NewTemplateProcessor(staticDir string) *TemplateProcessor {
	htdocsDir := filepath.Join(staticDir, "..")
	return &TemplateProcessor{
		staticDir: staticDir,
		htdocsDir: htdocsDir,
	}
}

// ProcessTemplate reads a file and processes template variables
func (tp *TemplateProcessor) ProcessTemplate(filePath string) (string, error) {
	// Read the template file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %v", filePath, err)
	}

	// Process template variables
	processed := string(content)
	processed = strings.ReplaceAll(processed, "${JS_VERSION}", GetJSVersion())
	
	return processed, nil
}

// ServeTemplate serves a processed template file with proper headers
func (tp *TemplateProcessor) ServeTemplate(w http.ResponseWriter, r *http.Request, templatePath string, contentType string) error {
	// Get the full path to the template file
	fullPath := filepath.Join(tp.htdocsDir, templatePath)
	
	// Process the template
	content, err := tp.ProcessTemplate(fullPath)
	if err != nil {
		return err
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", contentType)
	
	// For HTML templates, always set no-cache for development
	if contentType == "text/html" {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}

	// Write the processed content
	w.Write([]byte(content))
	return nil
}

// ServeIndex serves the main index.html template with template processing
func (tp *TemplateProcessor) ServeIndex(w http.ResponseWriter, r *http.Request) error {
	return tp.ServeTemplate(w, r, "index.html", "text/html")
}