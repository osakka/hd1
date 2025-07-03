package physics

import (
	"net/http"
	"strings"
	"holodeck1/api/scenegraph"
)

// extractSessionId extracts session ID from URL path
func extractSessionId(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "sessions" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// extractEntityId extracts entity ID from URL path
func extractEntityId(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "rigidbodies" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// PlayCanvasBridge type alias for consistency
type PlayCanvasBridge = scenegraph.PlayCanvasBridge

// NewPlayCanvasBridge creates a new PlayCanvas bridge for physics
func NewPlayCanvasBridge(sessionID string) *PlayCanvasBridge {
	return scenegraph.NewPlayCanvasBridge(sessionID)
}

// IsPlayCanvasAvailable checks if PlayCanvas integration is available
func IsPlayCanvasAvailable() bool {
	return scenegraph.IsPlayCanvasAvailable()
}

// Physics helper functions
func GetPhysicsWorldFromPlayCanvas(bridge *PlayCanvasBridge) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("getPhysicsWorld", map[string]interface{}{})
}

func UpdatePhysicsWorldInPlayCanvas(bridge *PlayCanvasBridge, worldData map[string]interface{}) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("updatePhysicsWorld", worldData)
}

func GetRigidBodiesFromPlayCanvas(bridge *PlayCanvasBridge) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("getRigidBodies", map[string]interface{}{})
}

func ApplyForceInPlayCanvas(bridge *PlayCanvasBridge, entityId string, forceData map[string]interface{}) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"entityId": entityId,
	}
	for k, v := range forceData {
		params[k] = v
	}
	return bridge.ExecutePlayCanvasOperation("applyForce", params)
}