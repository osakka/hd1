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
	case "list-entities":
		listEntities()
	case "create-entity":
		createEntity()
	case "get-scene-state":
		getSceneState()
	case "update-scene-state":
		updateSceneState()
	case "save-scene-state":
		saveSceneState()
	case "play-recording":
		playRecording()
	case "get-logs":
		getLogs()
	case "activate-session-scene":
		activateSessionScene()
	case "activate-entity":
		activateEntity()
	case "apply-environment":
		applyEnvironment()
	case "sync-session-state":
		syncSessionState()
	case "update-entity":
		updateEntity()
	case "delete-entity":
		deleteEntity()
	case "get-entity":
		getEntity()
	case "get-entity-lifecycle-status":
		getEntityLifecycleStatus()
	case "get-session-graph":
		getSessionGraph()
	case "update-session-graph":
		updateSessionGraph()
	case "list-environments":
		listEnvironments()
	case "export-scene-definition":
		exportSceneDefinition()
	case "play-animation":
		playAnimation()
	case "list-audio-sources":
		listAudioSources()
	case "create-audio-source":
		createAudioSource()
	case "stop-recording":
		stopRecording()
	case "set-camera-position":
		setCameraPosition()
	case "start-camera-orbit":
		startCameraOrbit()
	case "create-channel":
		createChannel()
	case "list-channels":
		listChannels()
	case "load-scene-state":
		loadSceneState()
	case "get-recording-status":
		getRecordingStatus()
	case "bulk-component-operation":
		bulkComponentOperation()
	case "bulk-entity-lifecycle-operation":
		bulkEntityLifecycleOperation()
	case "set-log-level":
		setLogLevel()
	case "set-canvas":
		setCanvas()
	case "enable-entity":
		enableEntity()
	case "destroy-entity":
		destroyEntity()
	case "save-scene-from-session":
		saveSceneFromSession()
	case "get-entity-parent":
		getEntityParent()
	case "set-entity-parent":
		setEntityParent()
	case "import-scene-definition":
		importSceneDefinition()
	case "get-session":
		getSession()
	case "delete-session":
		deleteSession()
	case "deactivate-entity":
		deactivateEntity()
	case "list-rigid-bodies":
		listRigidBodies()
	case "join-session-channel":
		joinSessionChannel()
	case "get-entity-transforms":
		getEntityTransforms()
	case "set-entity-transforms":
		setEntityTransforms()
	case "play-audio":
		playAudio()
	case "set-trace-modules":
		setTraceModules()
	case "stop-animation":
		stopAnimation()
	case "get-channel":
		getChannel()
	case "update-channel":
		updateChannel()
	case "delete-channel":
		deleteChannel()
	case "list-sessions":
		listSessions()
	case "create-session":
		createSession()
	case "load-scene":
		loadScene()
	case "force-refresh":
		forceRefresh()
	case "list-scenes":
		listScenes()
	case "get-physics-world":
		getPhysicsWorld()
	case "update-physics-world":
		updatePhysicsWorld()
	case "stop-audio":
		stopAudio()
	case "get-logging-config":
		getLoggingConfig()
	case "set-logging-config":
		setLoggingConfig()
	case "update-component":
		updateComponent()
	case "remove-component":
		removeComponent()
	case "get-component":
		getComponent()
	case "get-entity-children":
		getEntityChildren()
	case "get-scene-hierarchy":
		getSceneHierarchy()
	case "update-scene-hierarchy":
		updateSceneHierarchy()
	case "instantiate-prop":
		instantiateProp()
	case "reset-scene-state":
		resetSceneState()
	case "disable-entity":
		disableEntity()
	case "get-version":
		getVersion()
	case "leave-session-channel":
		leaveSessionChannel()
	case "get-session-channel-status":
		getSessionChannelStatus()
	case "list-session-scenes":
		listSessionScenes()
	case "create-session-scene":
		createSessionScene()
	case "apply-force":
		applyForce()
	case "list-entity-components":
		listEntityComponents()
	case "add-component":
		addComponent()
	case "get-hierarchy-tree":
		getHierarchyTree()
	case "list-animations":
		listAnimations()
	case "create-animation":
		createAnimation()
	case "start-recording":
		startRecording()
	case "list-props":
		listProps()
	case "fork-scene":
		forkScene()
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
	fmt.Println("  list-entities - GET /sessions/{sessionId}/entities")
	fmt.Println("  create-entity - POST /sessions/{sessionId}/entities")
	fmt.Println("  get-scene-state - GET /sessions/{sessionId}/scene/state")
	fmt.Println("  update-scene-state - PUT /sessions/{sessionId}/scene/state")
	fmt.Println("  save-scene-state - POST /sessions/{sessionId}/scene/state/save")
	fmt.Println("  play-recording - POST /sessions/{sessionId}/recording/play")
	fmt.Println("  get-logs - GET /admin/logging/logs")
	fmt.Println("  activate-session-scene - POST /sessions/{sessionId}/scenes/{sceneId}/activate")
	fmt.Println("  activate-entity - POST /sessions/{sessionId}/entities/{entityId}/lifecycle/activate")
	fmt.Println("  apply-environment - POST /environments/{environmentId}")
	fmt.Println("  sync-session-state - POST /sessions/{sessionId}/channel/sync")
	fmt.Println("  update-entity - PUT /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  delete-entity - DELETE /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  get-entity - GET /sessions/{sessionId}/entities/{entityId}")
	fmt.Println("  get-entity-lifecycle-status - GET /sessions/{sessionId}/entities/{entityId}/lifecycle/status")
	fmt.Println("  get-session-graph - GET /sessions/{sessionId}/channel/graph")
	fmt.Println("  update-session-graph - PUT /sessions/{sessionId}/channel/graph")
	fmt.Println("  list-environments - GET /environments")
	fmt.Println("  export-scene-definition - GET /sessions/{sessionId}/scene/export")
	fmt.Println("  play-animation - POST /sessions/{sessionId}/animations/{animationId}/play")
	fmt.Println("  list-audio-sources - GET /sessions/{sessionId}/audio/sources")
	fmt.Println("  create-audio-source - POST /sessions/{sessionId}/audio/sources")
	fmt.Println("  stop-recording - POST /sessions/{sessionId}/recording/stop")
	fmt.Println("  set-camera-position - PUT /sessions/{sessionId}/camera/position")
	fmt.Println("  start-camera-orbit - POST /sessions/{sessionId}/camera/orbit")
	fmt.Println("  create-channel - POST /channels")
	fmt.Println("  list-channels - GET /channels")
	fmt.Println("  load-scene-state - POST /sessions/{sessionId}/scene/state/load")
	fmt.Println("  get-recording-status - GET /sessions/{sessionId}/recording/status")
	fmt.Println("  bulk-component-operation - POST /sessions/{sessionId}/entities/{entityId}/components/bulk")
	fmt.Println("  bulk-entity-lifecycle-operation - POST /sessions/{sessionId}/entities/lifecycle/bulk")
	fmt.Println("  set-log-level - POST /admin/logging/level")
	fmt.Println("  set-canvas - POST /browser/canvas")
	fmt.Println("  enable-entity - PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/enable")
	fmt.Println("  destroy-entity - DELETE /sessions/{sessionId}/entities/{entityId}/lifecycle/destroy")
	fmt.Println("  save-scene-from-session - POST /sessions/{sessionId}/scenes/save")
	fmt.Println("  get-entity-parent - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/parent")
	fmt.Println("  set-entity-parent - PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/parent")
	fmt.Println("  import-scene-definition - POST /sessions/{sessionId}/scene/import")
	fmt.Println("  get-session - GET /sessions/{sessionId}")
	fmt.Println("  delete-session - DELETE /sessions/{sessionId}")
	fmt.Println("  deactivate-entity - POST /sessions/{sessionId}/entities/{entityId}/lifecycle/deactivate")
	fmt.Println("  list-rigid-bodies - GET /sessions/{sessionId}/physics/rigidbodies")
	fmt.Println("  join-session-channel - POST /sessions/{sessionId}/channel/join")
	fmt.Println("  get-entity-transforms - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms")
	fmt.Println("  set-entity-transforms - PUT /sessions/{sessionId}/entities/{entityId}/hierarchy/transforms")
	fmt.Println("  play-audio - POST /sessions/{sessionId}/audio/sources/{audioId}/play")
	fmt.Println("  set-trace-modules - POST /admin/logging/trace")
	fmt.Println("  stop-animation - POST /sessions/{sessionId}/animations/{animationId}/stop")
	fmt.Println("  get-channel - GET /channels/{channelId}")
	fmt.Println("  update-channel - PUT /channels/{channelId}")
	fmt.Println("  delete-channel - DELETE /channels/{channelId}")
	fmt.Println("  list-sessions - GET /sessions")
	fmt.Println("  create-session - POST /sessions")
	fmt.Println("  load-scene - POST /scenes/{sceneId}")
	fmt.Println("  force-refresh - POST /browser/refresh")
	fmt.Println("  list-scenes - GET /scenes")
	fmt.Println("  get-physics-world - GET /sessions/{sessionId}/physics/world")
	fmt.Println("  update-physics-world - PUT /sessions/{sessionId}/physics/world")
	fmt.Println("  stop-audio - POST /sessions/{sessionId}/audio/sources/{audioId}/stop")
	fmt.Println("  get-logging-config - GET /admin/logging/config")
	fmt.Println("  set-logging-config - POST /admin/logging/config")
	fmt.Println("  update-component - PUT /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  remove-component - DELETE /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  get-component - GET /sessions/{sessionId}/entities/{entityId}/components/{componentType}")
	fmt.Println("  get-entity-children - GET /sessions/{sessionId}/entities/{entityId}/hierarchy/children")
	fmt.Println("  get-scene-hierarchy - GET /sessions/{sessionId}/scene/hierarchy")
	fmt.Println("  update-scene-hierarchy - PUT /sessions/{sessionId}/scene/hierarchy")
	fmt.Println("  instantiate-prop - POST /sessions/{sessionId}/props/{propId}")
	fmt.Println("  reset-scene-state - POST /sessions/{sessionId}/scene/state/reset")
	fmt.Println("  disable-entity - PUT /sessions/{sessionId}/entities/{entityId}/lifecycle/disable")
	fmt.Println("  get-version - GET /version")
	fmt.Println("  leave-session-channel - POST /sessions/{sessionId}/channel/leave")
	fmt.Println("  get-session-channel-status - GET /sessions/{sessionId}/channel/status")
	fmt.Println("  list-session-scenes - GET /sessions/{sessionId}/scenes")
	fmt.Println("  create-session-scene - POST /sessions/{sessionId}/scenes")
	fmt.Println("  apply-force - POST /sessions/{sessionId}/physics/rigidbodies/{entityId}/force")
	fmt.Println("  list-entity-components - GET /sessions/{sessionId}/entities/{entityId}/components")
	fmt.Println("  add-component - POST /sessions/{sessionId}/entities/{entityId}/components")
	fmt.Println("  get-hierarchy-tree - GET /sessions/{sessionId}/entities/hierarchy/tree")
	fmt.Println("  list-animations - GET /sessions/{sessionId}/animations")
	fmt.Println("  create-animation - POST /sessions/{sessionId}/animations")
	fmt.Println("  start-recording - POST /sessions/{sessionId}/recording/start")
	fmt.Println("  list-props - GET /props")
	fmt.Println("  fork-scene - POST /scenes/{sceneId}/fork")

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

func getLogs() {
	makeRequest("GET", "/admin/logging/logs", nil)
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

func applyEnvironment() {
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
	makeRequest("POST", "/environments/" + os.Args[2] + "", body)
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
	makeRequest("POST", "/sessions/" + os.Args[2] + "/channel/sync", body)
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

func getEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "", nil)
}

func getEntityLifecycleStatus() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/status", nil)
}

func getSessionGraph() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/channel/graph", nil)
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
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/channel/graph", body)
}

func listEnvironments() {
	makeRequest("GET", "/environments", nil)
}

func exportSceneDefinition() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/scene/export", nil)
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

func createChannel() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/channels", body)
}

func listChannels() {
	makeRequest("GET", "/channels", nil)
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

func getRecordingStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/recording/status", nil)
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

func destroyEntity() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/lifecycle/destroy", nil)
}

func saveSceneFromSession() {
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
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scenes/save", body)
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

func listRigidBodies() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/physics/rigidbodies", nil)
}

func joinSessionChannel() {
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
	makeRequest("POST", "/sessions/" + os.Args[2] + "/channel/join", body)
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

func getChannel() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/channels/" + os.Args[2] + "", nil)
}

func updateChannel() {
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
	makeRequest("PUT", "/channels/" + os.Args[2] + "", body)
}

func deleteChannel() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("DELETE", "/channels/" + os.Args[2] + "", nil)
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

func loadScene() {
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
	makeRequest("POST", "/scenes/" + os.Args[2] + "", body)
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

func listScenes() {
	makeRequest("GET", "/scenes", nil)
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

func getEntityChildren() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/" + os.Args[3] + "/hierarchy/children", nil)
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

func instantiateProp() {
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
	makeRequest("POST", "/sessions/" + os.Args[2] + "/props/" + os.Args[3] + "", body)
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

func getVersion() {
	makeRequest("GET", "/version", nil)
}

func leaveSessionChannel() {
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
	makeRequest("POST", "/sessions/" + os.Args[2] + "/channel/leave", body)
}

func getSessionChannelStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/channel/status", nil)
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

func getHierarchyTree() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/entities/hierarchy/tree", nil)
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

func listProps() {
	makeRequest("GET", "/props", nil)
}

func forkScene() {
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
	makeRequest("POST", "/scenes/" + os.Args[2] + "/fork", body)
}
