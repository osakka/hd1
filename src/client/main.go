package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	
	"holodeck1/config"
	"holodeck1/logging"
)

func main() {
	// Initialize configuration system
	if err := config.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: Configuration initialization failed: %v\n", err)
		os.Exit(1)
	}
	
	// Initialize logging for client operations
	if err := logging.InitLogger(config.Config.Logging.LogDir, config.Config.Logging.Level, []string{}); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logging: %v\n", err)
		// Continue without logging rather than fail
	}
	
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	
	switch command {
	case "get-version":
		getVersion()
	case "instantiate-prop":
		instantiateProp()
	case "list-scenes":
		listScenes()
	case "load-scene":
		loadScene()
	case "stop-recording":
		stopRecording()
	case "set-camera-position":
		setCameraPosition()
	case "list-sessions":
		listSessions()
	case "create-session":
		createSession()
	case "delete-session":
		deleteSession()
	case "get-session":
		getSession()
	case "list-environments":
		listEnvironments()
	case "apply-environment":
		applyEnvironment()
	case "fork-scene":
		forkScene()
	case "start-recording":
		startRecording()
	case "play-recording":
		playRecording()
	case "get-logging-config":
		getLoggingConfig()
	case "set-logging-config":
		setLoggingConfig()
	case "list-props":
		listProps()
	case "get-recording-status":
		getRecordingStatus()
	case "list-objects":
		listObjects()
	case "create-object":
		createObject()
	case "get-object":
		getObject()
	case "update-object":
		updateObject()
	case "delete-object":
		deleteObject()
	case "set-log-level":
		setLogLevel()
	case "set-trace-modules":
		setTraceModules()
	case "get-logs":
		getLogs()
	case "set-canvas":
		setCanvas()
	case "save-scene-from-session":
		saveSceneFromSession()
	case "force-refresh":
		forceRefresh()
	case "start-camera-orbit":
		startCameraOrbit()
	case "help":
		showHelp()
	default:
		logging.Error("unknown command", map[string]interface{}{
			"command": command,
		})
		fmt.Printf("Unknown command: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("HD1 Client - Auto-generated from API specification")
	fmt.Println("Available commands:")
	fmt.Println("  get-version - GET /version")
	fmt.Println("  instantiate-prop - POST /sessions/{sessionId}/props/{propId}")
	fmt.Println("  list-scenes - GET /scenes")
	fmt.Println("  load-scene - POST /scenes/{sceneId}")
	fmt.Println("  stop-recording - POST /sessions/{sessionId}/recording/stop")
	fmt.Println("  set-camera-position - PUT /sessions/{sessionId}/camera/position")
	fmt.Println("  list-sessions - GET /sessions")
	fmt.Println("  create-session - POST /sessions")
	fmt.Println("  delete-session - DELETE /sessions/{sessionId}")
	fmt.Println("  get-session - GET /sessions/{sessionId}")
	fmt.Println("  list-environments - GET /environments")
	fmt.Println("  apply-environment - POST /environments/{environmentId}")
	fmt.Println("  fork-scene - POST /scenes/{sceneId}/fork")
	fmt.Println("  start-recording - POST /sessions/{sessionId}/recording/start")
	fmt.Println("  play-recording - POST /sessions/{sessionId}/recording/play")
	fmt.Println("  get-logging-config - GET /admin/logging/config")
	fmt.Println("  set-logging-config - POST /admin/logging/config")
	fmt.Println("  list-props - GET /props")
	fmt.Println("  get-recording-status - GET /sessions/{sessionId}/recording/status")
	fmt.Println("  list-objects - GET /sessions/{sessionId}/objects")
	fmt.Println("  create-object - POST /sessions/{sessionId}/objects")
	fmt.Println("  get-object - GET /sessions/{sessionId}/objects/{objectName}")
	fmt.Println("  update-object - PUT /sessions/{sessionId}/objects/{objectName}")
	fmt.Println("  delete-object - DELETE /sessions/{sessionId}/objects/{objectName}")
	fmt.Println("  set-log-level - POST /admin/logging/level")
	fmt.Println("  set-trace-modules - POST /admin/logging/trace")
	fmt.Println("  get-logs - GET /admin/logging/logs")
	fmt.Println("  set-canvas - POST /browser/canvas")
	fmt.Println("  save-scene-from-session - POST /sessions/{sessionId}/scenes/save")
	fmt.Println("  force-refresh - POST /browser/refresh")
	fmt.Println("  start-camera-orbit - POST /sessions/{sessionId}/camera/orbit")

}

func makeRequest(method, path string, body interface{}) {
	url := config.GetAPIBase() + path
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			logging.Error("JSON marshal failed", map[string]interface{}{
				"error": err.Error(),
			})
			fmt.Printf("Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		logging.Error("request creation failed", map[string]interface{}{
			"method": method,
			"url": url,
			"error": err.Error(),
		})
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("HTTP request failed", map[string]interface{}{
			"method": method,
			"url": url,
			"error": err.Error(),
		})
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("response read failed", map[string]interface{}{
			"method": method,
			"url": url,
			"error": err.Error(),
		})
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

// clientError logs error and displays user message
func clientError(command, message string) {
	logging.Error("client command error", map[string]interface{}{
		"command": command,
		"error":   message,
	})
	fmt.Printf("Error: %s\n", message)
	os.Exit(1)
}

func getVersion() {
	makeRequest("GET", "/version", nil)
}

func instantiateProp() {
	if len(os.Args) < 4 {
		clientError("instantiate-prop", "Missing required parameters")
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			clientError("instantiate-prop", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/props/" + os.Args[3] + "", body)
}

func listScenes() {
	makeRequest("GET", "/scenes", nil)
}

func loadScene() {
	if len(os.Args) < 3 {
		clientError("load-scene", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("load-scene", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/scenes/" + os.Args[2] + "", body)
}

func stopRecording() {
	if len(os.Args) < 3 {
		clientError("stop-recording", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("stop-recording", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/recording/stop", body)
}

func setCameraPosition() {
	if len(os.Args) < 3 {
		clientError("set-camera-position", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("set-camera-position", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/camera/position", body)
}

func listSessions() {
	makeRequest("GET", "/sessions", nil)
}

func createSession() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			clientError("create-session", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions", body)
}

func deleteSession() {
	if len(os.Args) < 3 {
		clientError("delete-session", "Missing required parameter")
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "", nil)
}

func getSession() {
	if len(os.Args) < 3 {
		clientError("get-session", "Missing required parameter")
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "", nil)
}

func listEnvironments() {
	makeRequest("GET", "/environments", nil)
}

func applyEnvironment() {
	if len(os.Args) < 3 {
		clientError("apply-environment", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("apply-environment", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/environments/" + os.Args[2] + "", body)
}

func forkScene() {
	if len(os.Args) < 3 {
		clientError("fork-scene", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("fork-scene", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/scenes/" + os.Args[2] + "/fork", body)
}

func startRecording() {
	if len(os.Args) < 3 {
		clientError("start-recording", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("start-recording", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/recording/start", body)
}

func playRecording() {
	if len(os.Args) < 3 {
		clientError("play-recording", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("play-recording", fmt.Sprintf("Error parsing JSON: %v", err))
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
			clientError("set-logging-config", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/admin/logging/config", body)
}

func listProps() {
	makeRequest("GET", "/props", nil)
}

func getRecordingStatus() {
	if len(os.Args) < 3 {
		clientError("get-recording-status", "Missing required parameter")
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/recording/status", nil)
}

func listObjects() {
	if len(os.Args) < 3 {
		clientError("list-objects", "Missing required parameter")
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/objects", nil)
}

func createObject() {
	if len(os.Args) < 3 {
		clientError("create-object", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("create-object", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/objects", body)
}

func getObject() {
	if len(os.Args) < 4 {
		clientError("get-object", "Missing required parameters")
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/objects/" + os.Args[3] + "", nil)
}

func updateObject() {
	if len(os.Args) < 4 {
		clientError("update-object", "Missing required parameters")
	}
	var body interface{}
	if len(os.Args) > 3 {
		if err := json.Unmarshal([]byte(os.Args[4]), &body); err != nil {
			clientError("update-object", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/objects/" + os.Args[3] + "", body)
}

func deleteObject() {
	if len(os.Args) < 4 {
		clientError("delete-object", "Missing required parameters")
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/objects/" + os.Args[3] + "", nil)
}

func setLogLevel() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			clientError("set-log-level", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/admin/logging/level", body)
}

func setTraceModules() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			clientError("set-trace-modules", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/admin/logging/trace", body)
}

func getLogs() {
	makeRequest("GET", "/admin/logging/logs", nil)
}

func setCanvas() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			clientError("set-canvas", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/browser/canvas", body)
}

func saveSceneFromSession() {
	if len(os.Args) < 3 {
		clientError("save-scene-from-session", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("save-scene-from-session", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/scenes/save", body)
}

func forceRefresh() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			clientError("force-refresh", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/browser/refresh", body)
}

func startCameraOrbit() {
	if len(os.Args) < 3 {
		clientError("start-camera-orbit", "Missing required parameter")
	}
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[3]), &body); err != nil {
			clientError("start-camera-orbit", fmt.Sprintf("Error parsing JSON: %v", err))
		}
	}
	makeRequest("POST", "/sessions/" + os.Args[2] + "/camera/orbit", body)
}
