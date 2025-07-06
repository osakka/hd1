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
	case "enable-entity":
		enableEntity()
	case "leave-session-world":
		leaveSessionWorld()
	case "list-avatars":
		listAvatars()
	case "get-scene-hierarchy":
		getSceneHierarchy()
	case "update-scene-hierarchy":
		updateSceneHierarchy()
	case "get-recording-status":
		getRecordingStatus()
	case "get-entity-children":
		getEntityChildren()
	case "import-scene-definition":
		importSceneDefinition()
	case "list-sessions":
		listSessions()
	case "create-session":
		createSession()
	case "get-camera-position":
		getCameraPosition()
	case "set-camera-position":
		setCameraPosition()
	case "list-worlds":
		listWorlds()
	case "create-world":
		createWorld()
	case "bulk-component-operation":
		bulkComponentOperation()
	case "list-session-scenes":
		listSessionScenes()
	case "create-session-scene":
		createSessionScene()
	case "play-audio":
		playAudio()
	case "stop-audio":
		stopAudio()
	case "force-refresh":
		forceRefresh()
	case "get-component":
		getComponent()
	case "update-component":
		updateComponent()
	case "remove-component":
		removeComponent()
	case "post-client-join-sync":
		postClientJoinSync()
	case "sync-session-state":
		syncSessionState()
	case "save-scene-state":
		saveSceneState()
	case "list-rigid-bodies":
		listRigidBodies()
	case "play-recording":
		playRecording()
	case "list-entity-components":
		listEntityComponents()
	case "add-component":
		addComponent()
	case "destroy-entity":
		destroyEntity()
	case "list-animations":
		listAnimations()
	case "create-animation":
		createAnimation()
	case "play-animation":
		playAnimation()
	case "get-physics-world":
		getPhysicsWorld()
	case "update-physics-world":
		updatePhysicsWorld()
	case "apply-force":
		applyForce()
	case "get-logging-config":
		getLoggingConfig()
	case "set-logging-config":
		setLoggingConfig()
	case "get-avatar-asset":
		getAvatarAsset()
	case "set-entity-parent":
		setEntityParent()
	case "get-entity-parent":
		getEntityParent()
	case "get-scene-state":
		getSceneState()
	case "update-scene-state":
		updateSceneState()
	case "list-audio-sources":
		listAudioSources()
	case "create-audio-source":
		createAudioSource()
	case "get-hierarchy-tree":
		getHierarchyTree()
	case "deactivate-entity":
		deactivateEntity()
	case "load-scene-state":
		loadSceneState()
	case "reset-scene-state":
		resetSceneState()
	case "set-log-level":
		setLogLevel()
	case "start-camera-orbit":
		startCameraOrbit()
	case "get-session-graph":
		getSessionGraph()
	case "update-session-graph":
		updateSessionGraph()
	case "get-session-avatar":
		getSessionAvatar()
	case "set-session-avatar":
		setSessionAvatar()
	case "change-session-avatar":
		changeSessionAvatar()
	case "get-entity-transforms":
		getEntityTransforms()
	case "set-entity-transforms":
		setEntityTransforms()
	case "disable-entity":
		disableEntity()
	case "bulk-entity-lifecycle-operation":
		bulkEntityLifecycleOperation()
	case "get-logs":
		getLogs()
	case "set-canvas":
		setCanvas()
	case "get-version":
		getVersion()
	case "stop-animation":
		stopAnimation()
	case "get-avatar-specification":
		getAvatarSpecification()
	case "get-session-world-status":
		getSessionWorldStatus()
	case "activate-session-scene":
		activateSessionScene()
	case "stop-recording":
		stopRecording()
	case "get-world":
		getWorld()
	case "update-world":
		updateWorld()
	case "delete-world":
		deleteWorld()
	case "delete-entity":
		deleteEntity()
	case "get-entity":
		getEntity()
	case "update-entity":
		updateEntity()
	case "get-entity-lifecycle-status":
		getEntityLifecycleStatus()
	case "export-scene-definition":
		exportSceneDefinition()
	case "start-recording":
		startRecording()
	case "set-trace-modules":
		setTraceModules()
	case "list-entities":
		listEntities()
	case "create-entity":
		createEntity()
	case "join-session-world":
		joinSessionWorld()
	case "get-world-state-sync":
		getWorldStateSync()
	case "activate-entity":
		activateEntity()
	case "delete-session":
		deleteSession()
	case "get-session":
		getSession()
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
	fmt.Println("  enable-entity - PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable")
	fmt.Println("  leave-session-world - POST /sessions/{sessionId}/world/leave")
	fmt.Println("  list-avatars - GET /avatars")
	fmt.Println("  get-scene-hierarchy - GET /sessions/{sessionId}/scene/hierarchy")
	fmt.Println("  update-scene-hierarchy - PUT /sessions/{sessionId}/scene/hierarchy")
	fmt.Println("  get-recording-status - GET /sessions/{sessionId}/recording/status")
	fmt.Println("  get-entity-children - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children")
	fmt.Println("  import-scene-definition - POST /sessions/{sessionId}/scene/import")
	fmt.Println("  list-sessions - GET /sessions")
	fmt.Println("  create-session - POST /sessions")
	fmt.Println("  get-camera-position - GET /sessions/{sessionId}/camera/position")
	fmt.Println("  set-camera-position - PUT /sessions/{sessionId}/camera/position")
	fmt.Println("  list-worlds - GET /worlds")
	fmt.Println("  create-world - POST /worlds")
	fmt.Println("  bulk-component-operation - POST /sessions/{sessionId}/entities/{entityId}/components/bulk")
	fmt.Println("  list-session-scenes - GET /sessions/{sessionId}/scenes")
	fmt.Println("  create-session-scene - POST /sessions/{sessionId}/scenes")
	fmt.Println("  play-audio - POST /sessions/{sessionId}/audio/sources/{audioId}/play")
	fmt.Println("  stop-audio - POST /sessions/{sessionId}/audio/sources/{audioId}/stop")
	fmt.Println("  force-refresh - POST /browser/refresh")
	fmt.Println("  get-component - GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  update-component - PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  remove-component - DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  post-client-join-sync - POST /sessions/{sessionId}/sync/client-join")
	fmt.Println("  sync-session-state - POST /sessions/{sessionId}/world/sync")
	fmt.Println("  save-scene-state - POST /sessions/{sessionId}/scene/state/save")
	fmt.Println("  list-rigid-bodies - GET /sessions/{sessionId}/physics/rigidbodies")
	fmt.Println("  play-recording - POST /sessions/{sessionId}/recording/play")
	fmt.Println("  list-entity-components - GET /sessions/{sessionId}/entities/{entityId}/components")
	fmt.Println("  add-component - POST /sessions/{sessionId}/entities/{entityId}/components")
	fmt.Println("  destroy-entity - DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy")
	fmt.Println("  list-animations - GET /sessions/{sessionId}/animations")
	fmt.Println("  create-animation - POST /sessions/{sessionId}/animations")
	fmt.Println("  play-animation - POST /sessions/{sessionId}/animations/{animationId}/play")
	fmt.Println("  get-physics-world - GET /sessions/{sessionId}/physics/world")
	fmt.Println("  update-physics-world - PUT /sessions/{sessionId}/physics/world")
	fmt.Println("  apply-force - POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force")
	fmt.Println("  get-logging-config - GET /admin/logging/config")
	fmt.Println("  set-logging-config - POST /admin/logging/config")
	fmt.Println("  get-avatar-asset - GET /avatars/{avatarType}/asset")
	fmt.Println("  set-entity-parent - PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent")
	fmt.Println("  get-entity-parent - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent")
	fmt.Println("  get-scene-state - GET /sessions/{sessionId}/scene/state")
	fmt.Println("  update-scene-state - PUT /sessions/{sessionId}/scene/state")
	fmt.Println("  list-audio-sources - GET /sessions/{sessionId}/audio/sources")
	fmt.Println("  create-audio-source - POST /sessions/{sessionId}/audio/sources")
	fmt.Println("  get-hierarchy-tree - GET /sessions/{sessionId}/entities/hierarchy/tree")
	fmt.Println("  deactivate-entity - POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate")
	fmt.Println("  load-scene-state - POST /sessions/{sessionId}/scene/state/load")
	fmt.Println("  reset-scene-state - POST /sessions/{sessionId}/scene/state/reset")
	fmt.Println("  set-log-level - POST /admin/logging/level")
	fmt.Println("  start-camera-orbit - POST /sessions/{sessionId}/camera/orbit")
	fmt.Println("  get-session-graph - GET /sessions/{sessionId}/world/graph")
	fmt.Println("  update-session-graph - PUT /sessions/{sessionId}/world/graph")
	fmt.Println("  get-session-avatar - GET /sessions/{sessionId}/avatar")
	fmt.Println("  set-session-avatar - POST /sessions/{sessionId}/avatar")
	fmt.Println("  change-session-avatar - PUT /sessions/{sessionId}/avatar")
	fmt.Println("  get-entity-transforms - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms")
	fmt.Println("  set-entity-transforms - PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms")
	fmt.Println("  disable-entity - PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable")
	fmt.Println("  bulk-entity-lifecycle-operation - POST /sessions/{sessionId}/entities/lifecycle/bulk")
	fmt.Println("  get-logs - GET /admin/logging/logs")
	fmt.Println("  set-canvas - POST /browser/canvas")
	fmt.Println("  get-version - GET /version")
	fmt.Println("  stop-animation - POST /sessions/{sessionId}/animations/{animationId}/stop")
	fmt.Println("  get-avatar-specification - GET /avatars/{avatarType}")
	fmt.Println("  get-session-world-status - GET /sessions/{sessionId}/world/status")
	fmt.Println("  activate-session-scene - POST /sessions/{sessionId}/scenes/{sceneId}/activate")
	fmt.Println("  stop-recording - POST /sessions/{sessionId}/recording/stop")
	fmt.Println("  get-world - GET /worlds/{worldId}")
	fmt.Println("  update-world - PUT /worlds/{worldId}")
	fmt.Println("  delete-world - DELETE /worlds/{worldId}")
	fmt.Println("  delete-entity - DELETE /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  get-entity - GET /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  update-entity - PUT /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  get-entity-lifecycle-status - GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status")
	fmt.Println("  export-scene-definition - GET /sessions/{sessionId}/scene/export")
	fmt.Println("  start-recording - POST /sessions/{sessionId}/recording/start")
	fmt.Println("  set-trace-modules - POST /admin/logging/trace")
	fmt.Println("  list-entities - GET /sessions/{sessionId}/entities")
	fmt.Println("  create-entity - POST /sessions/{sessionId}/entities")
	fmt.Println("  join-session-world - POST /sessions/{sessionId}/world/join")
	fmt.Println("  get-world-state-sync - GET /sessions/{sessionId}/sync/world-state")
	fmt.Println("  activate-entity - POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate")
	fmt.Println("  delete-session - DELETE /sessions/{sessionId}")
	fmt.Println("  get-session - GET /sessions/{sessionId}")

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

func listAvatars() {
	makeRequest("GET", "/avatars", nil)
}

func getSceneHierarchy() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/scene/hierarchy", nil)
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

func getRecordingStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/recording/status", nil)
}

func getEntityChildren() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/children", nil)
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

func getCameraPosition() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/camera/position", nil)
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

func getComponent() {
	if len(os.Args) < 5 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components/" + os.Args[4] + "", nil)
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

func listRigidBodies() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/physics/rigidbodies", nil)
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

func listEntityComponents() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/components", nil)
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

func destroyEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/destroy", nil)
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

func getAvatarAsset() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/avatars/" + os.Args[2] + "/asset", nil)
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

func getEntityParent() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/parent", nil)
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

func getHierarchyTree() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/hierarchy/tree", nil)
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

func getSessionGraph() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/world/graph", nil)
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

func getLogs() {
	makeRequest("GET", "/admin/logging/logs", nil)
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

func getVersion() {
	makeRequest("GET", "/version", nil)
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

func getAvatarSpecification() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/avatars/" + os.Args[2] + "", nil)
}

func getSessionWorldStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/world/status", nil)
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

func deleteEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "", nil)
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

func getEntityLifecycleStatus() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/status", nil)
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

func getWorldStateSync() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/sync/world-state", nil)
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

func deleteSession() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "", nil)
}

func getSession() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "", nil)
}
