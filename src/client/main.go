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
	case "update-entity":
		updateEntity()
	case "delete-entity":
		deleteEntity()
	case "move-avatar":
		moveAvatar()
	case "update-scene":
		updateScene()
	case "get-scene":
		getScene()
	case "get-version":
		getVersion()
	case "submit-operation":
		submitOperation()
	case "get-missing-operations":
		getMissingOperations()
	case "get-full-sync":
		getFullSync()
	case "get-sync-stats":
		getSyncStats()
	case "create-entity":
		createEntity()
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
	fmt.Println("  update-entity - PUT /threejs/entities/{entityId}")
	fmt.Println("  delete-entity - DELETE /threejs/entities/{entityId}")
	fmt.Println("  move-avatar - POST /threejs/avatars/{sessionId}/move")
	fmt.Println("  update-scene - PUT /threejs/scene")
	fmt.Println("  get-scene - GET /threejs/scene")
	fmt.Println("  get-version - GET /system/version")
	fmt.Println("  submit-operation - POST /sync/operations")
	fmt.Println("  get-missing-operations - GET /sync/missing/{from}/{to}")
	fmt.Println("  get-full-sync - GET /sync/full")
	fmt.Println("  get-sync-stats - GET /sync/stats")
	fmt.Println("  create-entity - POST /threejs/entities")

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


func updateEntity() {
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
	makeRequest("PUT", "/threejs/entities/" + os.Args[2] + "", body)
}

func deleteEntity() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing required parameter")
		os.Exit(1)
	}
	makeRequest("DELETE", "/threejs/entities/" + os.Args[2] + "", nil)
}

func moveAvatar() {
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
	makeRequest("POST", "/threejs/avatars/" + os.Args[2] + "/move", body)
}

func updateScene() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("PUT", "/threejs/scene", body)
}

func getScene() {
	makeRequest("GET", "/threejs/scene", nil)
}

func getVersion() {
	makeRequest("GET", "/system/version", nil)
}

func submitOperation() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/sync/operations", body)
}

func getMissingOperations() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required parameters")
		os.Exit(1)
	}
	makeRequest("GET", "/sync/missing/" + os.Args[2] + "/" + os.Args[3] + "", nil)
}

func getFullSync() {
	makeRequest("GET", "/sync/full", nil)
}

func getSyncStats() {
	makeRequest("GET", "/sync/stats", nil)
}

func createEntity() {
	var body interface{}
	if len(os.Args) > 2 {
		if err := json.Unmarshal([]byte(os.Args[2]), &body); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	makeRequest("POST", "/threejs/entities", body)
}
