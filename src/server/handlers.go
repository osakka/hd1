package server

import (
	"net/http"
	"holodeck1/logging"
)

var templateProcessor *TemplateProcessor

// InitializeTemplateProcessor sets up the template processor
func InitializeTemplateProcessor(staticDir string) {
	templateProcessor = NewTemplateProcessor(staticDir)
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check if template processor is initialized
	if templateProcessor == nil {
		logging.Error("template processor not initialized", nil)
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}
	
	// Serve template-processed index.html
	if err := templateProcessor.ServeIndex(w, r); err != nil {
		logging.Error("failed to serve index template", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Template processing failed", http.StatusInternalServerError)
		return
	}
}

// ServeConsoleJS serves the hd1-console.js with template processing
func ServeConsoleJS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check if template processor is initialized
	if templateProcessor == nil {
		logging.Error("template processor not initialized", nil)
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}
	
	// Serve template-processed console JavaScript
	if err := templateProcessor.ServeTemplate(w, r, "static/js/hd1-console.js", "application/javascript"); err != nil {
		logging.Error("failed to serve console template", map[string]interface{}{
			"error": err.Error(),
		})
		http.Error(w, "Template processing failed", http.StatusInternalServerError)
		return
	}
}
