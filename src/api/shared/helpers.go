package shared

import (
	"net/http"
	"time"

	"holodeck1/server"
)

// Vector3 represents a 3D vector
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// GetClientID extracts client ID from request headers
func GetClientID(r *http.Request) string {
	if clientID := r.Header.Get("X-HD1-ID"); clientID != "" {
		return clientID
	}
	return "api-client-" + time.Now().Format("20060102150405")
}

// GetHubFromContext extracts the hub from request context
func GetHubFromContext(r *http.Request) *server.Hub {
	if hub := r.Context().Value("hub"); hub != nil {
		if h, ok := hub.(*server.Hub); ok {
			return h
		}
	}
	return nil
}