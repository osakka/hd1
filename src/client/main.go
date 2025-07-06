package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	
	"holodeck1/config"
)

func main() {
	// Initialize configuration system
	if err := config.Initialize(); err != nil {
		fmt.Printf("FATAL: Configuration initialization failed: %v\n", err)
		os.Exit(1)
	}
	
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	
	switch command {
	case "get-session":
		getSession()
	case "delete-session":
		deleteSession()
	case "disable-entity":
		disableEntity()
	case "get-scene-state":
		getSceneState()
	case "update-scene-state":
		updateSceneState()
	case "save-scene-state":
		saveSceneState()
	case "stop-animation":
		stopAnimation()
	case "set-log-level":
		setLogLevel()
	case "join-session-world":
		joinSessionWorld()
	case "bulk-component-operation":
		bulkComponentOperation()
	case "import-scene-definition":
		importSceneDefinition()
	case "enable-entity":
		enableEntity()
	case "get-world":
		getWorld()
	case "update-world":
		updateWorld()
	case "delete-world":
		deleteWorld()
	case "get-session-world-status":
		getSessionWorldStatus()
	case "get-entity":
		getEntity()
	case "update-entity":
		updateEntity()
	case "delete-entity":
		deleteEntity()
	case "stop-audio":
		stopAudio()
	case "get-recording-status":
		getRecordingStatus()
	case "start-camera-orbit":
		startCameraOrbit()
	case "list-rigid-bodies":
		listRigidBodies()
	case "reset-scene-state":
		resetSceneState()
	case "set-canvas":
		setCanvas()
	case "get-hierarchy-tree":
		getHierarchyTree()
	case "activate-entity":
		activateEntity()
	case "post-client-join-sync":
		postClientJoinSync()
	case "get-avatar-specification":
		getAvatarSpecification()
	case "play-animation":
		playAnimation()
	case "get-logs":
		getLogs()
	case "force-refresh":
		forceRefresh()
	case "play-audio":
		playAudio()
	case "get-world-state-sync":
		getWorldStateSync()
	case "add-component":
		addComponent()
	case "list-entity-components":
		listEntityComponents()
	case "load-scene-state":
		loadSceneState()
	case "update-session-graph":
		updateSessionGraph()
	case "get-session-graph":
		getSessionGraph()
	case "get-physics-world":
		getPhysicsWorld()
	case "update-physics-world":
		updatePhysicsWorld()
	case "get-avatar-asset":
		getAvatarAsset()
	case "list-entities":
		listEntities()
	case "create-entity":
		createEntity()
	case "list-worlds":
		listWorlds()
	case "create-world":
		createWorld()
	case "leave-session-world":
		leaveSessionWorld()
	case "sync-session-state":
		syncSessionState()
	case "deactivate-entity":
		deactivateEntity()
	case "get-entity-parent":
		getEntityParent()
	case "set-entity-parent":
		setEntityParent()
	case "get-entity-children":
		getEntityChildren()
	case "list-sessions":
		listSessions()
	case "create-session":
		createSession()
	case "bulk-entity-lifecycle-operation":
		bulkEntityLifecycleOperation()
	case "stop-recording":
		stopRecording()
	case "set-camera-position":
		setCameraPosition()
	case "get-camera-position":
		getCameraPosition()
	case "update-component":
		updateComponent()
	case "remove-component":
		removeComponent()
	case "get-component":
		getComponent()
	case "export-scene-definition":
		exportSceneDefinition()
	case "start-recording":
		startRecording()
	case "get-entity-transforms":
		getEntityTransforms()
	case "set-entity-transforms":
		setEntityTransforms()
	case "destroy-entity":
		destroyEntity()
	case "update-scene-hierarchy":
		updateSceneHierarchy()
	case "get-scene-hierarchy":
		getSceneHierarchy()
	case "activate-session-scene":
		activateSessionScene()
	case "play-recording":
		playRecording()
	case "get-logging-config":
		getLoggingConfig()
	case "set-logging-config":
		setLoggingConfig()
	case "set-trace-modules":
		setTraceModules()
	case "get-version":
		getVersion()
	case "get-session-avatar":
		getSessionAvatar()
	case "set-session-avatar":
		setSessionAvatar()
	case "change-session-avatar":
		changeSessionAvatar()
	case "get-entity-lifecycle-status":
		getEntityLifecycleStatus()
	case "list-session-scenes":
		listSessionScenes()
	case "create-session-scene":
		createSessionScene()
	case "list-animations":
		listAnimations()
	case "create-animation":
		createAnimation()
	case "apply-force":
		applyForce()
	case "list-avatars":
		listAvatars()
	case "list-audio-sources":
		listAudioSources()
	case "create-audio-source":
		createAudioSource()
	case "help":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("HD1 Client - Auto-generated from API specification")
	fmt.Println("Available commands:")
	fmt.Println("  get-session - GET /sessions/{sessionId}")
	fmt.Println("  delete-session - DELETE /sessions/{sessionId}")
	fmt.Println("  disable-entity - PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable")
	fmt.Println("  get-scene-state - GET /sessions/{sessionId}/scene/state")
	fmt.Println("  update-scene-state - PUT /sessions/{sessionId}/scene/state")
	fmt.Println("  save-scene-state - POST /sessions/{sessionId}/scene/state/save")
	fmt.Println("  stop-animation - POST /sessions/{sessionId}/animations/{animationId}/stop")
	fmt.Println("  set-log-level - POST /admin/logging/level")
	fmt.Println("  join-session-world - POST /sessions/{sessionId}/world/join")
	fmt.Println("  bulk-component-operation - POST /sessions/{sessionId}/entities/{entityId}/components/bulk")
	fmt.Println("  import-scene-definition - POST /sessions/{sessionId}/scene/import")
	fmt.Println("  enable-entity - PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable")
	fmt.Println("  get-world - GET /worlds/{worldId}")
	fmt.Println("  update-world - PUT /worlds/{worldId}")
	fmt.Println("  delete-world - DELETE /worlds/{worldId}")
	fmt.Println("  get-session-world-status - GET /sessions/{sessionId}/world/status")
	fmt.Println("  get-entity - GET /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  update-entity - PUT /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  delete-entity - DELETE /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  stop-audio - POST /sessions/{sessionId}/audio/sources/{audioId}/stop")
	fmt.Println("  get-recording-status - GET /sessions/{sessionId}/recording/status")
	fmt.Println("  start-camera-orbit - POST /sessions/{sessionId}/camera/orbit")
	fmt.Println("  list-rigid-bodies - GET /sessions/{sessionId}/physics/rigidbodies")
	fmt.Println("  reset-scene-state - POST /sessions/{sessionId}/scene/state/reset")
	fmt.Println("  set-canvas - POST /browser/canvas")
	fmt.Println("  get-hierarchy-tree - GET /sessions/{sessionId}/entities/hierarchy/tree")
	fmt.Println("  activate-entity - POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate")
	fmt.Println("  post-client-join-sync - POST /sessions/{sessionId}/sync/client-join")
	fmt.Println("  get-avatar-specification - GET /avatars/{avatarType}")
	fmt.Println("  play-animation - POST /sessions/{sessionId}/animations/{animationId}/play")
	fmt.Println("  get-logs - GET /admin/logging/logs")
	fmt.Println("  force-refresh - POST /browser/refresh")
	fmt.Println("  play-audio - POST /sessions/{sessionId}/audio/sources/{audioId}/play")
	fmt.Println("  get-world-state-sync - GET /sessions/{sessionId}/sync/world-state")
	fmt.Println("  add-component - POST /sessions/{sessionId}/entities/{entityId}/components")
	fmt.Println("  list-entity-components - GET /sessions/{sessionId}/entities/{entityId}/components")
	fmt.Println("  load-scene-state - POST /sessions/{sessionId}/scene/state/load")
	fmt.Println("  update-session-graph - PUT /sessions/{sessionId}/world/graph")
	fmt.Println("  get-session-graph - GET /sessions/{sessionId}/world/graph")
	fmt.Println("  get-physics-world - GET /sessions/{sessionId}/physics/world")
	fmt.Println("  update-physics-world - PUT /sessions/{sessionId}/physics/world")
	fmt.Println("  get-avatar-asset - GET /avatars/{avatarType}/asset")
	fmt.Println("  list-entities - GET /sessions/{sessionId}/entities")
	fmt.Println("  create-entity - POST /sessions/{sessionId}/entities")
	fmt.Println("  list-worlds - GET /worlds")
	fmt.Println("  create-world - POST /worlds")
	fmt.Println("  leave-session-world - POST /sessions/{sessionId}/world/leave")
	fmt.Println("  sync-session-state - POST /sessions/{sessionId}/world/sync")
	fmt.Println("  deactivate-entity - POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate")
	fmt.Println("  get-entity-parent - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent")
	fmt.Println("  set-entity-parent - PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent")
	fmt.Println("  get-entity-children - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children")
	fmt.Println("  list-sessions - GET /sessions")
	fmt.Println("  create-session - POST /sessions")
	fmt.Println("  bulk-entity-lifecycle-operation - POST /sessions/{sessionId}/entities/lifecycle/bulk")
	fmt.Println("  stop-recording - POST /sessions/{sessionId}/recording/stop")
	fmt.Println("  set-camera-position - PUT /sessions/{sessionId}/camera/position")
	fmt.Println("  get-camera-position - GET /sessions/{sessionId}/camera/position")
	fmt.Println("  update-component - PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  remove-component - DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  get-component - GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  export-scene-definition - GET /sessions/{sessionId}/scene/export")
	fmt.Println("  start-recording - POST /sessions/{sessionId}/recording/start")
	fmt.Println("  get-entity-transforms - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms")
	fmt.Println("  set-entity-transforms - PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms")
	fmt.Println("  destroy-entity - DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy")
	fmt.Println("  update-scene-hierarchy - PUT /sessions/{sessionId}/scene/hierarchy")
	fmt.Println("  get-scene-hierarchy - GET /sessions/{sessionId}/scene/hierarchy")
	fmt.Println("  activate-session-scene - POST /sessions/{sessionId}/scenes/{sceneId}/activate")
	fmt.Println("  play-recording - POST /sessions/{sessionId}/recording/play")
	fmt.Println("  get-logging-config - GET /admin/logging/config")
	fmt.Println("  set-logging-config - POST /admin/logging/config")
	fmt.Println("  set-trace-modules - POST /admin/logging/trace")
	fmt.Println("  get-version - GET /version")
	fmt.Println("  get-session-avatar - GET /sessions/{sessionId}/avatar")
	fmt.Println("  set-session-avatar - POST /sessions/{sessionId}/avatar")
	fmt.Println("  change-session-avatar - PUT /sessions/{sessionId}/avatar")
	fmt.Println("  get-entity-lifecycle-status - GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status")
	fmt.Println("  list-session-scenes - GET /sessions/{sessionId}/scenes")
	fmt.Println("  create-session-scene - POST /sessions/{sessionId}/scenes")
	fmt.Println("  list-animations - GET /sessions/{sessionId}/animations")
	fmt.Println("  create-animation - POST /sessions/{sessionId}/animations")
	fmt.Println("  apply-force - POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force")
	fmt.Println("  list-avatars - GET /avatars")
	fmt.Println("  list-audio-sources - GET /sessions/{sessionId}/audio/sources")
	fmt.Println("  create-audio-source - POST /sessions/{sessionId}/audio/sources")

}

func makeRequest(method, path string, body interface{}) {
	url := config.GetAPIBase() + path
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}
	
	// Pretty print JSON
	var jsonData interface{}
	if err := json.Unmarshal(responseBody, &jsonData); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println(string(responseBody))
	}
}


func getSession() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "", nil)
}

func deleteSession() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "", nil)
}

func disableEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/disable", body)
}

func getSceneState() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/scene/state", nil)
}

func updateSceneState() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/scene/state", body)
}

func saveSceneState() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scene/state/save", body)
}

func stopAnimation() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/animations/" + os.Args[3] + "/stop", body)
}

func setLogLevel() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/admin/logging/level", body)
}

func joinSessionWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/world/join", body)
}

func bulkComponentOperation() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components/bulk", body)
}

func importSceneDefinition() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scene/import", body)
}

func enableEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/enable", body)
}

func getWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/worlds/" + os.Args[2] + "", nil)
}

func updateWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/worlds/" + os.Args[2] + "", body)
}

func deleteWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("DELETE", "/worlds/" + os.Args[2] + "", nil)
}

func getSessionWorldStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/world/status", nil)
}

func getEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "", nil)
}

func updateEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "", body)
}

func deleteEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "", nil)
}

func stopAudio() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/audio/sources/" + os.Args[3] + "/stop", body)
}

func getRecordingStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/recording/status", nil)
}

func startCameraOrbit() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/camera/orbit", body)
}

func listRigidBodies() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/physics/rigidbodies", nil)
}

func resetSceneState() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scene/state/reset", body)
}

func setCanvas() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/browser/canvas", body)
}

func getHierarchyTree() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/hierarchy/tree", nil)
}

func activateEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/activate", body)
}

func postClientJoinSync() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/sync/client-join", body)
}

func getAvatarSpecification() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/avatars/" + os.Args[2] + "", nil)
}

func playAnimation() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/animations/" + os.Args[3] + "/play", body)
}

func getLogs() {
	makeRequest("GET", "/admin/logging/logs", nil)
}

func forceRefresh() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/browser/refresh", body)
}

func playAudio() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/audio/sources/" + os.Args[3] + "/play", body)
}

func getWorldStateSync() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/sync/world-state", nil)
}

func addComponent() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components", body)
}

func listEntityComponents() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components", nil)
}

func loadSceneState() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scene/state/load", body)
}

func updateSessionGraph() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/world/graph", body)
}

func getSessionGraph() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/world/graph", nil)
}

func getPhysicsWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/physics/world", nil)
}

func updatePhysicsWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/physics/world", body)
}

func getAvatarAsset() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/avatars/" + os.Args[2] + "/asset", nil)
}

func listEntities() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities", nil)
}

func createEntity() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/entities", body)
}

func listWorlds() {
	makeRequest("GET", "/worlds", nil)
}

func createWorld() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/worlds", body)
}

func leaveSessionWorld() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/world/leave", body)
}

func syncSessionState() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/world/sync", body)
}

func deactivateEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/deactivate", body)
}

func getEntityParent() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/parent", nil)
}

func setEntityParent() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/parent", body)
}

func getEntityChildren() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/children", nil)
}

func listSessions() {
	makeRequest("GET", "/sessions", nil)
}

func createSession() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions", body)
}

func bulkEntityLifecycleOperation() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/entities/lifecycle/bulk", body)
}

func stopRecording() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/recording/stop", body)
}

func setCameraPosition() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/camera/position", body)
}

func getCameraPosition() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/camera/position", nil)
}

func updateComponent() {
	if len(os.Args) < 5 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 4 {
		if err := json.Unmarshal([]byte(os.Args[5]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components/" + os.Args[4] + "", body)
}

func removeComponent() {
	if len(os.Args) < 5 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components/" + os.Args[4] + "", nil)
}

func getComponent() {
	if len(os.Args) < 5 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components/" + os.Args[4] + "", nil)
}

func exportSceneDefinition() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/scene/export", nil)
}

func startRecording() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/recording/start", body)
}

func getEntityTransforms() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/transforms", nil)
}

func setEntityTransforms() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/transforms", body)
}

func destroyEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/destroy", nil)
}

func updateSceneHierarchy() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/scene/hierarchy", body)
}

func getSceneHierarchy() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/scene/hierarchy", nil)
}

func activateSessionScene() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scenes/" + os.Args[3] + "/activate", body)
}

func playRecording() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/recording/play", body)
}

func getLoggingConfig() {
	makeRequest("GET", "/admin/logging/config", nil)
}

func setLoggingConfig() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/admin/logging/config", body)
}

func setTraceModules() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/admin/logging/trace", body)
}

func getVersion() {
	makeRequest("GET", "/version", nil)
}

func getSessionAvatar() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/avatar", nil)
}

func setSessionAvatar() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/avatar", body)
}

func changeSessionAvatar() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/avatar", body)
}

func getEntityLifecycleStatus() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/status", nil)
}

func listSessionScenes() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/scenes", nil)
}

func createSessionScene() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scenes", body)
}

func listAnimations() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/animations", nil)
}

func createAnimation() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/animations", body)
}

func applyForce() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/physics/rigidbodies/" + os.Args[3] + "/force", body)
}

func listAvatars() {
	makeRequest("GET", "/avatars", nil)
}

func listAudioSources() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/audio/sources", nil)
}

func createAudioSource() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/audio/sources", body)
}
