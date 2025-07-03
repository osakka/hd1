package scenegraph

import (
	"encoding/json"
	"fmt"
	"holodeck1/logging"
	"os/exec"
	"strings"
	"time"
)

// PlayCanvasBridge provides integration with PlayCanvas engine via Node.js
type PlayCanvasBridge struct {
	SessionID string
	ScriptPath string
}

// NewPlayCanvasBridge creates a new PlayCanvas integration bridge
func NewPlayCanvasBridge(sessionID string) *PlayCanvasBridge {
	return &PlayCanvasBridge{
		SessionID: sessionID,
		ScriptPath: "/opt/hd1/src/api/scenegraph/playcanvas_runner.js",
	}
}

// ExecutePlayCanvasOperation runs PlayCanvas operations via Node.js
func (pcb *PlayCanvasBridge) ExecutePlayCanvasOperation(operation string, params map[string]interface{}) (map[string]interface{}, error) {
	// Create operation payload
	operationData := map[string]interface{}{
		"operation": operation,
		"sessionId": pcb.SessionID,
		"params":    params,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	// Convert to JSON for Node.js
	payloadBytes, err := json.Marshal(operationData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operation data: %v", err)
	}
	
	logging.Info("executing PlayCanvas operation", map[string]interface{}{
		"operation": operation,
		"sessionId": pcb.SessionID,
		"payload_size": len(payloadBytes),
	})
	
	// Execute via Node.js PlayCanvas runner
	cmd := exec.Command("node", pcb.ScriptPath)
	cmd.Stdin = strings.NewReader(string(payloadBytes))
	
	output, err := cmd.Output()
	if err != nil {
		logging.Error("PlayCanvas operation failed", map[string]interface{}{
			"operation": operation,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("PlayCanvas operation failed: %v", err)
	}
	
	// Parse result
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse PlayCanvas result: %v", err)
	}
	
	logging.Info("PlayCanvas operation completed", map[string]interface{}{
		"operation": operation,
		"success": result["success"],
	})
	
	return result, nil
}

// GetSceneHierarchyFromPlayCanvas retrieves real scene hierarchy
func (pcb *PlayCanvasBridge) GetSceneHierarchyFromPlayCanvas() (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("getSceneHierarchy", map[string]interface{}{})
}

// UpdateSceneHierarchyInPlayCanvas updates scene hierarchy
func (pcb *PlayCanvasBridge) UpdateSceneHierarchyInPlayCanvas(operations []map[string]interface{}) (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("updateSceneHierarchy", map[string]interface{}{
		"operations": operations,
	})
}

// GetSceneStateFromPlayCanvas retrieves real scene state
func (pcb *PlayCanvasBridge) GetSceneStateFromPlayCanvas() (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("getSceneState", map[string]interface{}{})
}

// UpdateSceneStateInPlayCanvas updates scene configuration
func (pcb *PlayCanvasBridge) UpdateSceneStateInPlayCanvas(updates map[string]interface{}) (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("updateSceneState", map[string]interface{}{
		"updates": updates,
	})
}

// CreateEntityInPlayCanvas creates a new PlayCanvas entity
func (pcb *PlayCanvasBridge) CreateEntityInPlayCanvas(entityData map[string]interface{}) (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("createEntity", entityData)
}

// AddComponentToPlayCanvasEntity adds component to PlayCanvas entity
func (pcb *PlayCanvasBridge) AddComponentToPlayCanvasEntity(entityID string, componentType string, componentData map[string]interface{}) (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("addComponent", map[string]interface{}{
		"entityId": entityID,
		"componentType": componentType,
		"componentData": componentData,
	})
}

// ExportSceneFromPlayCanvas exports scene definition
func (pcb *PlayCanvasBridge) ExportSceneFromPlayCanvas(format string, includeAssets bool) (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("exportScene", map[string]interface{}{
		"format": format,
		"includeAssets": includeAssets,
	})
}

// ImportSceneToPlayCanvas imports scene definition
func (pcb *PlayCanvasBridge) ImportSceneToPlayCanvas(sceneData map[string]interface{}, format string, mergeMode string) (map[string]interface{}, error) {
	return pcb.ExecutePlayCanvasOperation("importScene", map[string]interface{}{
		"sceneData": sceneData,
		"format": format,
		"mergeMode": mergeMode,
	})
}

// IsPlayCanvasAvailable checks if PlayCanvas integration is available
func IsPlayCanvasAvailable() bool {
	// Check if Node.js is available
	cmd := exec.Command("node", "--version")
	err := cmd.Run()
	return err == nil
}

// InitializePlayCanvasForSession sets up PlayCanvas application for session
func InitializePlayCanvasForSession(sessionID string) error {
	logging.Info("initializing PlayCanvas for session", map[string]interface{}{
		"sessionId": sessionID,
		"engine": "PlayCanvas v2.8.2",
	})
	
	bridge := NewPlayCanvasBridge(sessionID)
	_, err := bridge.ExecutePlayCanvasOperation("initializeApp", map[string]interface{}{
		"headless": true,
		"width": 1024,
		"height": 768,
	})
	
	return err
}