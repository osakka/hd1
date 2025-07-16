package adapters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"holodeck1/clients"
)

type WebAdapter struct {
	info *clients.PlatformInfo
}

func NewWebAdapter() *WebAdapter {
	return &WebAdapter{
		info: &clients.PlatformInfo{
			Name:    "web",
			Version: "1.0.0",
			Capabilities: []string{
				"3d_rendering",
				"webgl",
				"webrtc",
				"websocket",
				"file_upload",
				"geolocation",
				"notifications",
				"fullscreen",
				"pointer_lock",
				"gamepad",
				"media_capture",
			},
			Requirements: []string{
				"webgl2",
				"websocket",
				"es2020",
			},
			Metadata: map[string]interface{}{
				"description": "Web browser client supporting WebGL 3D rendering",
				"supported_browsers": []string{
					"chrome",
					"firefox",
					"safari",
					"edge",
				},
			},
		},
	}
}

func (w *WebAdapter) GetPlatformInfo() *clients.PlatformInfo {
	return w.info
}

func (w *WebAdapter) ValidateConfiguration(config map[string]interface{}) error {
	// Validate WebGL support
	if webgl, exists := config["webgl_version"]; exists {
		if webglStr, ok := webgl.(string); ok {
			if webglStr != "1.0" && webglStr != "2.0" {
				return fmt.Errorf("unsupported WebGL version: %s", webglStr)
			}
		}
	}

	// Validate viewport configuration
	if viewport, exists := config["viewport"]; exists {
		if viewportMap, ok := viewport.(map[string]interface{}); ok {
			if width, exists := viewportMap["width"]; exists {
				if widthFloat, ok := width.(float64); ok {
					if widthFloat < 320 || widthFloat > 7680 {
						return fmt.Errorf("invalid viewport width: %f", widthFloat)
					}
				}
			}
			if height, exists := viewportMap["height"]; exists {
				if heightFloat, ok := height.(float64); ok {
					if heightFloat < 240 || heightFloat > 4320 {
						return fmt.Errorf("invalid viewport height: %f", heightFloat)
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
					validQualities := []string{"low", "medium", "high", "ultra"}
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
		}
	}

	return nil
}

func (w *WebAdapter) FormatMessage(message *clients.ClientMessage) (interface{}, error) {
	// Format message for web client
	webMessage := map[string]interface{}{
		"id":        message.ID,
		"type":      message.Type,
		"action":    message.Action,
		"data":      message.Data,
		"timestamp": message.Timestamp.UnixNano() / 1000000, // Convert to milliseconds
	}

	// Add web-specific formatting
	switch message.Type {
	case "3d_object":
		return w.format3DObjectMessage(webMessage)
	case "ui_update":
		return w.formatUIUpdateMessage(webMessage)
	case "media":
		return w.formatMediaMessage(webMessage)
	default:
		return webMessage, nil
	}
}

func (w *WebAdapter) format3DObjectMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Convert coordinate system if needed
		if position, exists := data["position"]; exists {
			if posMap, ok := position.(map[string]interface{}); ok {
				// Web uses right-handed coordinate system
				data["position"] = map[string]interface{}{
					"x": posMap["x"],
					"y": posMap["y"],
					"z": posMap["z"],
				}
			}
		}

		// Ensure WebGL-compatible format
		if material, exists := data["material"]; exists {
			if matMap, ok := material.(map[string]interface{}); ok {
				data["material"] = map[string]interface{}{
					"type":    matMap["type"],
					"color":   matMap["color"],
					"texture": matMap["texture"],
					"webgl":   true,
				}
			}
		}
	}

	return message, nil
}

func (w *WebAdapter) formatUIUpdateMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Convert to DOM-compatible format
		if element, exists := data["element"]; exists {
			if elemMap, ok := element.(map[string]interface{}); ok {
				data["element"] = map[string]interface{}{
					"tag":        elemMap["tag"],
					"id":         elemMap["id"],
					"class":      elemMap["class"],
					"style":      elemMap["style"],
					"attributes": elemMap["attributes"],
					"content":    elemMap["content"],
				}
			}
		}
	}

	return message, nil
}

func (w *WebAdapter) formatMediaMessage(message map[string]interface{}) (interface{}, error) {
	if data, exists := message["data"].(map[string]interface{}); exists {
		// Ensure web-compatible media format
		if media, exists := data["media"]; exists {
			if mediaMap, ok := media.(map[string]interface{}); ok {
				data["media"] = map[string]interface{}{
					"type":     mediaMap["type"],
					"url":      mediaMap["url"],
					"format":   mediaMap["format"],
					"controls": true,
					"autoplay": false,
				}
			}
		}
	}

	return message, nil
}

func (w *WebAdapter) ParseMessage(data interface{}) (*clients.ClientMessage, error) {
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

func (w *WebAdapter) GetRequiredCapabilities() []string {
	return []string{
		"webgl2",
		"websocket",
		"file_api",
		"fullscreen_api",
	}
}