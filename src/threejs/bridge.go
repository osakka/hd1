// Package threejs provides Three.js integration bridge
// Pure JSON-based communication with Three.js scene manager
package threejs

import (
	"encoding/json"
	"fmt"
	"sync"

	"holodeck1/logging"
)

// Bridge represents the Three.js integration bridge
// Maintains scene state as JSON and provides serialization
type Bridge struct {
	// Scene state as JSON objects
	entities map[string]interface{}
	avatars  map[string]interface{}
	scene    map[string]interface{}
	
	// Thread safety
	mutex sync.RWMutex
}

// NewBridge creates a new Three.js bridge
func NewBridge() *Bridge {
	return &Bridge{
		entities: make(map[string]interface{}),
		avatars:  make(map[string]interface{}),
		scene: map[string]interface{}{
			"background": "#87CEEB", // Default sky blue
			"fog":        nil,
		},
	}
}

// ApplyOperation applies an operation to the Three.js scene state
func (b *Bridge) ApplyOperation(operation map[string]interface{}) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	opType, ok := operation["type"].(string)
	if !ok {
		return fmt.Errorf("operation missing type")
	}

	data, ok := operation["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("operation missing data")
	}

	switch opType {
	case "avatar_move":
		return b.applyAvatarMove(data)
	case "entity_create":
		return b.applyEntityCreate(data)
	case "entity_update":
		return b.applyEntityUpdate(data)
	case "entity_delete":
		return b.applyEntityDelete(data)
	case "scene_update":
		return b.applySceneUpdate(data)
	default:
		return fmt.Errorf("unknown operation type: %s", opType)
	}
}

// applyAvatarMove handles avatar movement operations
func (b *Bridge) applyAvatarMove(data map[string]interface{}) error {
	sessionID, ok := data["session_id"].(string)
	if !ok {
		return fmt.Errorf("avatar_move missing session_id")
	}

	// Get or create avatar
	avatar, exists := b.avatars[sessionID]
	if !exists {
		avatar = map[string]interface{}{
			"session_id": sessionID,
			"position":   map[string]interface{}{"x": 0, "y": 0, "z": 0},
			"rotation":   map[string]interface{}{"x": 0, "y": 0, "z": 0},
			"animation":  "idle",
		}
		b.avatars[sessionID] = avatar
	}

	avatarData := avatar.(map[string]interface{})

	// Update position
	if position, ok := data["position"].(map[string]interface{}); ok {
		avatarData["position"] = position
	}

	// Update rotation
	if rotation, ok := data["rotation"].(map[string]interface{}); ok {
		avatarData["rotation"] = rotation
	}

	// Update animation
	if animation, ok := data["animation"].(string); ok {
		avatarData["animation"] = animation
	}

	logging.Debug("avatar updated in bridge", map[string]interface{}{
		"session_id": sessionID,
		"position":   avatarData["position"],
	})

	return nil
}

// applyEntityCreate handles entity creation operations
func (b *Bridge) applyEntityCreate(data map[string]interface{}) error {
	entityID, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("entity_create missing id")
	}

	geometry, ok := data["geometry"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("entity_create missing geometry")
	}

	material, ok := data["material"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("entity_create missing material")
	}

	// Create entity
	entity := map[string]interface{}{
		"id":       entityID,
		"geometry": geometry,
		"material": material,
		"position": data["position"], // May be nil
		"rotation": data["rotation"], // May be nil
		"scale":    data["scale"],    // May be nil
		"visible":  true,             // Default visible
	}

	// Set visibility if provided
	if visible, ok := data["visible"].(bool); ok {
		entity["visible"] = visible
	}

	// Set default transforms if not provided
	if entity["position"] == nil {
		entity["position"] = map[string]interface{}{"x": 0, "y": 0, "z": 0}
	}
	if entity["rotation"] == nil {
		entity["rotation"] = map[string]interface{}{"x": 0, "y": 0, "z": 0}
	}
	if entity["scale"] == nil {
		entity["scale"] = map[string]interface{}{"x": 1, "y": 1, "z": 1}
	}

	b.entities[entityID] = entity

	logging.Info("entity created in bridge", map[string]interface{}{
		"entity_id": entityID,
		"geometry":  geometry["type"],
		"material":  material["type"],
	})

	return nil
}

// applyEntityUpdate handles entity update operations
func (b *Bridge) applyEntityUpdate(data map[string]interface{}) error {
	entityID, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("entity_update missing id")
	}

	entity, exists := b.entities[entityID]
	if !exists {
		return fmt.Errorf("entity not found: %s", entityID)
	}

	entityData := entity.(map[string]interface{})

	// Update provided properties
	if position, ok := data["position"]; ok {
		entityData["position"] = position
	}
	if rotation, ok := data["rotation"]; ok {
		entityData["rotation"] = rotation
	}
	if scale, ok := data["scale"]; ok {
		entityData["scale"] = scale
	}
	if visible, ok := data["visible"]; ok {
		entityData["visible"] = visible
	}
	if material, ok := data["material"]; ok {
		entityData["material"] = material
	}

	logging.Debug("entity updated in bridge", map[string]interface{}{
		"entity_id": entityID,
	})

	return nil
}

// applyEntityDelete handles entity deletion operations
func (b *Bridge) applyEntityDelete(data map[string]interface{}) error {
	entityID, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("entity_delete missing id")
	}

	if _, exists := b.entities[entityID]; !exists {
		return fmt.Errorf("entity not found: %s", entityID)
	}

	delete(b.entities, entityID)

	logging.Info("entity deleted from bridge", map[string]interface{}{
		"entity_id": entityID,
	})

	return nil
}

// applySceneUpdate handles scene update operations
func (b *Bridge) applySceneUpdate(data map[string]interface{}) error {
	// Update scene properties
	if background, ok := data["background"].(string); ok {
		b.scene["background"] = background
	}

	if fog, ok := data["fog"].(map[string]interface{}); ok {
		b.scene["fog"] = fog
	}

	logging.Debug("scene updated in bridge", map[string]interface{}{
		"background": b.scene["background"],
	})

	return nil
}

// GetSceneState returns the current Three.js scene state as JSON
func (b *Bridge) GetSceneState() map[string]interface{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Convert entities map to array
	entities := make([]interface{}, 0, len(b.entities))
	for _, entity := range b.entities {
		entities = append(entities, entity)
	}

	// Convert avatars map to array
	avatars := make([]interface{}, 0, len(b.avatars))
	for _, avatar := range b.avatars {
		avatars = append(avatars, avatar)
	}

	// Return complete scene state
	return map[string]interface{}{
		"background": b.scene["background"],
		"fog":        b.scene["fog"],
		"entities":   entities,
		"avatars":    avatars,
	}
}

// GetEntity returns a specific entity by ID
func (b *Bridge) GetEntity(entityID string) (map[string]interface{}, bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	entity, exists := b.entities[entityID]
	if !exists {
		return nil, false
	}

	return entity.(map[string]interface{}), true
}

// GetAvatar returns a specific avatar by session ID
func (b *Bridge) GetAvatar(sessionID string) (map[string]interface{}, bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	avatar, exists := b.avatars[sessionID]
	if !exists {
		return nil, false
	}

	return avatar.(map[string]interface{}), true
}

// ListEntities returns all entities
func (b *Bridge) ListEntities() []map[string]interface{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	entities := make([]map[string]interface{}, 0, len(b.entities))
	for _, entity := range b.entities {
		entities = append(entities, entity.(map[string]interface{}))
	}

	return entities
}

// ListAvatars returns all avatars
func (b *Bridge) ListAvatars() []map[string]interface{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	avatars := make([]map[string]interface{}, 0, len(b.avatars))
	for _, avatar := range b.avatars {
		avatars = append(avatars, avatar.(map[string]interface{}))
	}

	return avatars
}

// GetStats returns bridge statistics
func (b *Bridge) GetStats() map[string]interface{} {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return map[string]interface{}{
		"entity_count": len(b.entities),
		"avatar_count": len(b.avatars),
		"scene_props":  len(b.scene),
	}
}

// SerializeToJSON returns the complete scene state as JSON string
func (b *Bridge) SerializeToJSON() (string, error) {
	sceneState := b.GetSceneState()
	
	jsonData, err := json.Marshal(sceneState)
	if err != nil {
		return "", fmt.Errorf("failed to serialize scene state: %v", err)
	}

	return string(jsonData), nil
}

// Clear resets the bridge to initial state
func (b *Bridge) Clear() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.entities = make(map[string]interface{})
	b.avatars = make(map[string]interface{})
	b.scene = map[string]interface{}{
		"background": "#87CEEB",
		"fog":        nil,
	}

	logging.Info("Three.js bridge cleared", map[string]interface{}{})
}