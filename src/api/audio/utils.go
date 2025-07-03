package audio

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

// extractAudioId extracts audio ID from URL path
func extractAudioId(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "sources" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// PlayCanvasBridge type alias for consistency
type PlayCanvasBridge = scenegraph.PlayCanvasBridge

// NewPlayCanvasBridge creates a new PlayCanvas bridge for audio
func NewPlayCanvasBridge(sessionID string) *PlayCanvasBridge {
	return scenegraph.NewPlayCanvasBridge(sessionID)
}

// IsPlayCanvasAvailable checks if PlayCanvas integration is available
func IsPlayCanvasAvailable() bool {
	return scenegraph.IsPlayCanvasAvailable()
}

// Audio helper functions
func GetAudioSourcesFromPlayCanvas(bridge *PlayCanvasBridge) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("getAudioSources", map[string]interface{}{})
}

func CreateAudioSourceInPlayCanvas(bridge *PlayCanvasBridge, audioData map[string]interface{}) (map[string]interface{}, error) {
	return bridge.ExecutePlayCanvasOperation("createAudioSource", audioData)
}

func PlayAudioInPlayCanvas(bridge *PlayCanvasBridge, audioId string, playParams map[string]interface{}) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"audioId": audioId,
	}
	for k, v := range playParams {
		params[k] = v
	}
	return bridge.ExecutePlayCanvasOperation("playAudio", params)
}

func StopAudioInPlayCanvas(bridge *PlayCanvasBridge, audioId string, stopParams map[string]interface{}) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"audioId": audioId,
	}
	for k, v := range stopParams {
		params[k] = v
	}
	return bridge.ExecutePlayCanvasOperation("stopAudio", params)
}