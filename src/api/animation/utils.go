package animation

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

// extractAnimationId extracts animation ID from URL path
func extractAnimationId(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "animations" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// PlayCanvasBridge type alias for consistency
type PlayCanvasBridge = scenegraph.PlayCanvasBridge

// NewPlayCanvasBridge creates a new PlayCanvas bridge for animation
func NewPlayCanvasBridge(sessionID string) *PlayCanvasBridge {
	return scenegraph.NewPlayCanvasBridge(sessionID)
}

// IsPlayCanvasAvailable checks if PlayCanvas integration is available
func IsPlayCanvasAvailable() bool {
	return scenegraph.IsPlayCanvasAvailable()
}

// Animation helper functions
func GetAnimationsFromPlayCanvas(bridge *PlayCanvasBridge) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("getAnimations", map[string]interface{}{})
}

func CreateAnimationInPlayCanvas(bridge *PlayCanvasBridge, animationData map[string]interface{}) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("createAnimation", animationData)
}

func PlayAnimationInPlayCanvas(bridge *PlayCanvasBridge, animationId string, playParams map[string]interface{}) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"animationId": animationId,
	}
	for k, v := range playParams {
		params[k] = v
	}
	return bridge.ExecutePlayCanvasOperation("playAnimation", params)
}

func StopAnimationInPlayCanvas(bridge *PlayCanvasBridge, animationId string, stopParams map[string]interface{}) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"animationId": animationId,
	}
	for k, v := range stopParams {
		params[k] = v
	}
	return bridge.ExecutePlayCanvasOperation("stopAnimation", params)
}