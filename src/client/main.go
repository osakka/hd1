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
	case "get-version":
		getVersion()
	case "create-session":
		createSession()
	case "list-sessions":
		listSessions()
	case "save-scene-from-session":
		saveSceneFromSession()
	case "play-recording":
		playRecording()
	case "list-objects":
		listObjects()
	case "create-object":
		createObject()
	case "update-object":
		updateObject()
	case "delete-object":
		deleteObject()
	case "get-object":
		getObject()
	case "set-trace-modules":
		setTraceModules()
	case "list-environments":
		listEnvironments()
	case "stop-recording":
		stopRecording()
	case "set-log-level":
		setLogLevel()
	case "set-canvas":
		setCanvas()
	case "start-camera-orbit":
		startCameraOrbit()
	case "instantiate-prop":
		instantiateProp()
	case "load-scene":
		loadScene()
	case "fork-scene":
		forkScene()
	case "start-recording":
		startRecording()
	case "get-logging-config":
		getLoggingConfig()
	case "set-logging-config":
		setLoggingConfig()
	case "get-logs":
		getLogs()
	case "get-session":
		getSession()
	case "delete-session":
		deleteSession()
	case "apply-environment":
		applyEnvironment()
	case "list-props":
		listProps()
	case "list-scenes":
		listScenes()
	case "get-recording-status":
		getRecordingStatus()
	case "force-refresh":
		forceRefresh()
	case "set-camera-position":
		setCameraPosition()
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
	fmt.Println("  get-version - GET /version")
	fmt.Println("  create-session - POST /sessions")
	fmt.Println("  list-sessions - GET /sessions")
	fmt.Println("  save-scene-from-session - POST /sessions/{sessionId}/scenes/save")
	fmt.Println("  play-recording - POST /sessions/{sessionId}/recording/play")
	fmt.Println("  list-objects - GET /sessions/{sessionId}/objects")
	fmt.Println("  create-object - POST /sessions/{sessionId}/objects")
	fmt.Println("  update-object - PUT /sessions/{sessionId}/objects/{objectName}")
	fmt.Println("  delete-object - DELETE /sessions/{sessionId}/objects/{objectName}")
	fmt.Println("  get-object - GET /sessions/{sessionId}/objects/{objectName}")
	fmt.Println("  set-trace-modules - POST /admin/logging/trace")
	fmt.Println("  list-environments - GET /environments")
	fmt.Println("  stop-recording - POST /sessions/{sessionId}/recording/stop")
	fmt.Println("  set-log-level - POST /admin/logging/level")
	fmt.Println("  set-canvas - POST /browser/canvas")
	fmt.Println("  start-camera-orbit - POST /sessions/{sessionId}/camera/orbit")
	fmt.Println("  instantiate-prop - POST /sessions/{sessionId}/props/{propId}")
	fmt.Println("  load-scene - POST /scenes/{sceneId}")
	fmt.Println("  fork-scene - POST /scenes/{sceneId}/fork")
	fmt.Println("  start-recording - POST /sessions/{sessionId}/recording/start")
	fmt.Println("  get-logging-config - GET /admin/logging/config")
	fmt.Println("  set-logging-config - POST /admin/logging/config")
	fmt.Println("  get-logs - GET /admin/logging/logs")
	fmt.Println("  get-session - GET /sessions/{sessionId}")
	fmt.Println("  delete-session - DELETE /sessions/{sessionId}")
	fmt.Println("  apply-environment - POST /environments/{environmentId}")
	fmt.Println("  list-props - GET /props")
	fmt.Println("  list-scenes - GET /scenes")
	fmt.Println("  get-recording-status - GET /sessions/{sessionId}/recording/status")
	fmt.Println("  force-refresh - POST /browser/refresh")
	fmt.Println("  set-camera-position - PUT /sessions/{sessionId}/camera/position")

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


func getVersion() {
	makeRequest("GET", "/version", nil)
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

func listSessions() {
	makeRequest("GET", "/sessions", nil)
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

func listObjects() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/objects", nil)
}

func createObject() {
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
	makeRequest("POST", "/sessions/" + os.Args[2] + "/objects", body)
}

func updateObject() {
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
	makeRequest("PUT", "/sessions/" + os.Args[2] + "/objects/" + os.Args[3] + "", body)
}

func deleteObject() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("DELETE", "/sessions/" + os.Args[2] + "/objects/" + os.Args[3] + "", nil)
}

func getObject() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/objects/" + os.Args[3] + "", nil)
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

func listEnvironments() {
	makeRequest("GET", "/environments", nil)
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

func getLogs() {
	makeRequest("GET", "/admin/logging/logs", nil)
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

func listProps() {
	makeRequest("GET", "/props", nil)
}

func listScenes() {
	makeRequest("GET", "/scenes", nil)
}

func getRecordingStatus() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("GET", "/sessions/" + os.Args[2] + "/recording/status", nil)
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
