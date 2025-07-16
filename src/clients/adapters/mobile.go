package adapters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"holodeck1/clients"
)

type MobileAdapter struct {
	info *clients.PlatformInfo
}

func NewMobileAdapter() *MobileAdapter {
	return &MobileAdapter{
		info: &clients.PlatformInfo{
			Name:    "mobile",
			Version: "1.0.0",
			Capabilities: []string{
				"3d_rendering",
				"opengles",
				"touch_input",
				"accelerometer",
				"gyroscope",
				"gps",
				"camera",
				"microphone",
				"notifications",
				"vibration",
				"orientation",
				"battery",
				"network_info",
			},
			Requirements: []string{
				"opengles3",
				"vulkan_optional",
				"metal_optional",
			},
			Metadata: map[string]interface{}{
				"description": "Mobile client supporting OpenGL ES 3D rendering",
				"supported_platforms": []string{
					"ios",
					"android",
					"react_native",
					"flutter",
					"unity",
				},
			},
		},
	}
}

func (m *MobileAdapter) GetPlatformInfo() *clients.PlatformInfo {
	return m.info
}

func (m *MobileAdapter) ValidateConfiguration(config map[string]interface{}) error {
	// Validate graphics API
	if graphics, exists := config["graphics_api"]; exists {
		if graphicsStr, ok := graphics.(string); ok {
			validAPIs := []string{"opengles2", "opengles3", "vulkan", "metal"}
			valid := false
			for _, api := range validAPIs {
				if graphicsStr == api {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("unsupported graphics API: %s", graphicsStr)
			}
		}
	}

	// Validate screen configuration
	if screen, exists := config["screen"]; exists {
		if screenMap, ok := screen.(map[string]interface{}); ok {
			if width, exists := screenMap["width"]; exists {
				if widthFloat, ok := width.(float64); ok {
					if widthFloat < 320 || widthFloat > 4096 {
						return fmt.Errorf("invalid screen width: %f", widthFloat)
					}
				}
			}
			if height, exists := screenMap["height"]; exists {
				if heightFloat, ok := height.(float64); ok {
					if heightFloat < 480 || heightFloat > 4096 {
						return fmt.Errorf("invalid screen height: %f", heightFloat)
					}
				}
			}
			if dpi, exists := screenMap["dpi"]; exists {
				if dpiFloat, ok := dpi.(float64); ok {
					if dpiFloat < 72 || dpiFloat > 600 {
						return fmt.Errorf("invalid screen DPI: %f", dpiFloat)
					}
				}
			}
		}
	}

	// Validate performance settings
	if performance, exists := config["performance"]; exists {
		if perfMap, ok := performance.(map[string]interface{}); ok {
			if quality, exists := perfMap["quality"]; exists {
				if qualityStr, ok := quality.(string); ok {
					validQualities := []string{"low", "medium", "high"}
					valid := false
					for _, vq := range validQualities {
						if qualityStr == vq {
							valid = true
							break
						}
					}
					if !valid {
						return fmt.Errorf("invalid quality setting: %s", qualityStr)
					}
				}
			}
			if powerMode, exists := perfMap["power_mode"]; exists {
				if powerStr, ok := powerMode.(string); ok {
					validModes := []string{"battery_saver", "balanced", "performance"}
					valid := false
					for _, vm := range validModes {
						if powerStr == vm {
							valid = true
							break
						}
					}
					if !valid {
						return fmt.Errorf("invalid power mode: %s", powerStr)
					}
				}
			}
		}
	}

	return nil
}

func (m *MobileAdapter) FormatMessage(message *clients.ClientMessage) (interface{}, error) {
	// Format message for mobile client
	mobileMessage := map[string]interface{}{
		"id":        message.ID.String(),
		"type":      message.Type,
		"action":    message.Action,
		"data":      message.Data,
		"timestamp": message.Timestamp.UnixNano() / 1000000, // Convert to milliseconds
	}

	// Add mobile-specific formatting
	switch message.Type {
	case "3d_object":
		return m.format3DObjectMessage(mobileMessage)
	case "ui_update":
		return m.formatUIUpdateMessage(mobileMessage)
	case "input":
		return m.formatInputMessage(mobileMessage)
	case "sensor":
		return m.formatSensorMessage(mobileMessage)
	default:
		return mobileMessage, nil
	}
}

func (m *MobileAdapter) format3DObjectMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Convert coordinate system for mobile (left-handed)
		if position, exists := data["position"]; exists {
			if posMap, ok := position.(map[string]interface{}); ok {
				data["position"] = map[string]interface{}{
					"x": posMap["x"],
					"y": posMap["y"],
					"z": posMap["z"],
				}
			}
		}

		// Ensure OpenGL ES compatible format
		if material, exists := data["material"]; exists {
			if matMap, ok := material.(map[string]interface{}); ok {
				data["material"] = map[string]interface{}{
					"type":     matMap["type"],
					"color":    matMap["color"],
					"texture":  matMap["texture"],
					"opengles": true,
					"quality":  "mobile",
				}
			}
		}

		// Add mobile-specific optimizations
		data["mobile_optimized"] = true
		data["lod_enabled"] = true
	}

	return message, nil
}

func (m *MobileAdapter) formatUIUpdateMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Convert to mobile UI format
		if element, exists := data["element"]; exists {
			if elemMap, ok := element.(map[string]interface{}); ok {
				data["element"] = map[string]interface{}{
					"type":       elemMap["type"],
					"id":         elemMap["id"],
					"style":      elemMap["style"],
					"properties": elemMap["properties"],
					"content":    elemMap["content"],
					"touch_enabled": true,
					"haptic_feedback": true,
				}
			}
		}
	}

	return message, nil
}

func (m *MobileAdapter) formatInputMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Convert to mobile input format
		if input, exists := data["input"]; exists {
			if inputMap, ok := input.(map[string]interface{}); ok {
				data["input"] = map[string]interface{}{
					"type":      inputMap["type"],
					"touches":   inputMap["touches"],
					"gestures":  inputMap["gestures"],
					"pressure":  inputMap["pressure"],
					"timestamp": inputMap["timestamp"],
				}
			}
		}
	}

	return message, nil
}

func (m *MobileAdapter) formatSensorMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Convert to mobile sensor format
		if sensor, exists := data["sensor"]; exists {
			if sensorMap, ok := sensor.(map[string]interface{}); ok {
				data["sensor"] = map[string]interface{}{
					"type":        sensorMap["type"],
					"values":      sensorMap["values"],
					"accuracy":    sensorMap["accuracy"],
					"timestamp":   sensorMap["timestamp"],
					"calibrated":  true,
				}
			}
		}
	}

	return message, nil
}

func (m *MobileAdapter) ParseMessage(data interface{}) (*clients.ClientMessage, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var rawMessage map[string]interface{}
	err = json.Unmarshal(dataBytes, &rawMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	message := &clients.ClientMessage{}

	// Parse ID
	if id, exists := rawMessage["id"]; exists {
		if idStr, ok := id.(string); ok {
			message.ID = uuid.MustParse(idStr)
		}
	}

	// Parse client ID
	if clientID, exists := rawMessage["client_id"]; exists {
		if clientIDStr, ok := clientID.(string); ok {
			message.ClientID = uuid.MustParse(clientIDStr)
		}
	}

	// Parse type
	if msgType, exists := rawMessage["type"]; exists {
		if typeStr, ok := msgType.(string); ok {
			message.Type = typeStr
		}
	}

	// Parse action
	if action, exists := rawMessage["action"]; exists {
		if actionStr, ok := action.(string); ok {
			message.Action = actionStr
		}
	}

	// Parse data
	if msgData, exists := rawMessage["data"]; exists {
		if dataMap, ok := msgData.(map[string]interface{}); ok {
			message.Data = dataMap
		}
	}

	// Parse timestamp
	if timestamp, exists := rawMessage["timestamp"]; exists {
		if timestampFloat, ok := timestamp.(float64); ok {
			// Convert from milliseconds to time.Time
			message.Timestamp = time.Unix(0, int64(timestampFloat)*1000000)
		}
	}

	return message, nil
}

func (m *MobileAdapter) GetRequiredCapabilities() []string {
	return []string{
		"opengles3",
		"touch_input",
		"accelerometer",
		"gyroscope",
	}
}