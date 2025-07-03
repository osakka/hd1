package scenegraph

import (
	"net/http"
	"strings"
)

// extractSessionId extracts session ID from URL path
func extractSessionId(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "unknown"
}

// extractSceneId extracts scene ID from URL path
func extractSceneId(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")
	for i, part := range parts {
		if part == "scenes" && i+1 < len(parts) {
			sceneIdPart := parts[i+1]
			// Remove any additional path segments after sceneId
			if idx := strings.Index(sceneIdPart, "/"); idx != -1 {
				return sceneIdPart[:idx]
			}
			return sceneIdPart
		}
	}
	return "unknown"
}