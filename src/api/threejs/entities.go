package threejs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"holodeck1/logging"
	"holodeck1/server"
	"holodeck1/sync"
)

// Vector3 represents a 3D vector
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Geometry represents Three.js geometry
type Geometry struct {
	Type     string  `json:"type"`
	Width    float64 `json:"width,omitempty"`
	Height   float64 `json:"height,omitempty"`
	Depth    float64 `json:"depth,omitempty"`
	Radius   float64 `json:"radius,omitempty"`
	Segments int     `json:"segments,omitempty"`
}

// Material represents Three.js material
type Material struct {
	Type        string  `json:"type"`
	Color       string  `json:"color"`
	Transparent bool    `json:"transparent,omitempty"`
	Opacity     float64 `json:"opacity,omitempty"`
	Metalness   float64 `json:"metalness,omitempty"`
	Roughness   float64 `json:"roughness,omitempty"`
}

// CreateEntityRequest represents the request to create an entity
type CreateEntityRequest struct {
	Geometry Geometry `json:"geometry"`
	Material Material `json:"material"`
	Position *Vector3 `json:"position,omitempty"`
	Rotation *Vector3 `json:"rotation,omitempty"`
	Scale    *Vector3 `json:"scale,omitempty"`
	Visible  *bool    `json:"visible,omitempty"`
}

// CreateEntityResponse represents the response after creating an entity
type CreateEntityResponse struct {
	Success  bool   `json:"success"`
	EntityID string `json:"entity_id"`
	SeqNum   uint64 `json:"seq_num"`
}

// UpdateEntityRequest represents the request to update an entity
type UpdateEntityRequest struct {
	Position *Vector3  `json:"position,omitempty"`
	Rotation *Vector3  `json:"rotation,omitempty"`
	Scale    *Vector3  `json:"scale,omitempty"`
	Visible  *bool     `json:"visible,omitempty"`
	Material *Material `json:"material,omitempty"`
}

// UpdateEntityResponse represents the response after updating an entity
type UpdateEntityResponse struct {
	Success bool   `json:"success"`
	SeqNum  uint64 `json:"seq_num"`
}

// DeleteEntityResponse represents the response after deleting an entity
type DeleteEntityResponse struct {
	Success bool   `json:"success"`
	SeqNum  uint64 `json:"seq_num"`
}

// CreateEntity handles POST /api/threejs/entities
func CreateEntity(w http.ResponseWriter, r *http.Request) {
	var req CreateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate geometry
	if err := validateGeometry(req.Geometry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate material
	if err := validateMaterial(req.Material); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate entity ID
	entityID := generateEntityID()

	// Get client ID
	clientID := getClientID(r)

	// Create operation data
	operationData := map[string]interface{}{
		"id":       entityID,
		"geometry": req.Geometry,
		"material": req.Material,
	}

	// Add optional properties
	if req.Position != nil {
		operationData["position"] = req.Position
	}
	if req.Rotation != nil {
		operationData["rotation"] = req.Rotation
	}
	if req.Scale != nil {
		operationData["scale"] = req.Scale
	}
	if req.Visible != nil {
		operationData["visible"] = *req.Visible
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "entity_create",
		Data:      operationData,
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := CreateEntityResponse{
		Success:  true,
		EntityID: entityID,
		SeqNum:   operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

	logging.Info("entity created via API", map[string]interface{}{
		"entity_id": entityID,
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
	})
}

// UpdateEntity handles PUT /api/threejs/entities/{entityId}
func UpdateEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entityID := vars["entityId"]

	if entityID == "" {
		http.Error(w, "Entity ID required", http.StatusBadRequest)
		return
	}

	var req UpdateEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate material if provided
	if req.Material != nil {
		if err := validateMaterial(*req.Material); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Get client ID
	clientID := getClientID(r)

	// Create operation data
	operationData := map[string]interface{}{
		"id": entityID,
	}

	// Add provided updates
	if req.Position != nil {
		operationData["position"] = req.Position
	}
	if req.Rotation != nil {
		operationData["rotation"] = req.Rotation
	}
	if req.Scale != nil {
		operationData["scale"] = req.Scale
	}
	if req.Visible != nil {
		operationData["visible"] = *req.Visible
	}
	if req.Material != nil {
		operationData["material"] = req.Material
	}

	// Create operation
	operation := &sync.Operation{
		ClientID:  clientID,
		Type:      "entity_update",
		Data:      operationData,
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := UpdateEntityResponse{
		Success: true,
		SeqNum:  operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("entity updated via API", map[string]interface{}{
		"entity_id": entityID,
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
	})
}

// DeleteEntity handles DELETE /api/threejs/entities/{entityId}
func DeleteEntity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	entityID := vars["entityId"]

	if entityID == "" {
		http.Error(w, "Entity ID required", http.StatusBadRequest)
		return
	}

	// Get client ID
	clientID := getClientID(r)

	// Create operation
	operation := &sync.Operation{
		ClientID: clientID,
		Type:     "entity_delete",
		Data: map[string]interface{}{
			"id": entityID,
		},
		Timestamp: time.Now(),
	}

	// Get hub and submit operation
	hub := getHubFromContext(r)
	if hub == nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hub.GetSync().SubmitOperation(operation)

	// Return response
	response := DeleteEntityResponse{
		Success: true,
		SeqNum:  operation.SeqNum,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logging.Info("entity deleted via API", map[string]interface{}{
		"entity_id": entityID,
		"client_id": clientID,
		"seq_num":   operation.SeqNum,
	})
}

// Helper functions
func validateGeometry(geom Geometry) error {
	validTypes := map[string]bool{
		"box":      true,
		"sphere":   true,
		"plane":    true,
		"cylinder": true,
	}

	if !validTypes[geom.Type] {
		return fmt.Errorf("invalid geometry type: %s", geom.Type)
	}

	return nil
}

func validateMaterial(mat Material) error {
	validTypes := map[string]bool{
		"basic":    true,
		"phong":    true,
		"standard": true,
	}

	if !validTypes[mat.Type] {
		return fmt.Errorf("invalid material type: %s", mat.Type)
	}

	if mat.Color == "" {
		return fmt.Errorf("material color is required")
	}

	return nil
}

func generateEntityID() string {
	return "entity-" + time.Now().Format("20060102150405") + "-" + fmt.Sprintf("%d", time.Now().UnixNano()%10000)
}

func getClientID(r *http.Request) string {
	if clientID := r.Header.Get("X-Client-ID"); clientID != "" {
		return clientID
	}
	return "api-client-" + time.Now().Format("20060102150405")
}

func getHubFromContext(r *http.Request) *server.Hub {
	if hub := r.Context().Value("hub"); hub != nil {
		if h, ok := hub.(*server.Hub); ok {
			return h
		}
	}
	return nil
}