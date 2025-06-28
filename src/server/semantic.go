package server

import (
	"encoding/json"
	"math"
)

// Semantic UI System - AI-friendly interface creation

type SemanticElement struct {
	Type       string                 `json:"type"`       // "title", "button", "card", "dashboard", etc.
	Intent     string                 `json:"intent"`     // "primary", "secondary", "warning", "success"
	Content    string                 `json:"content"`    // Text content
	Size       string                 `json:"size"`       // "small", "medium", "large", "xl"
	Position   string                 `json:"position"`   // "top", "center", "bottom", "left", "right"
	Behavior   string                 `json:"behavior"`   // "static", "floating", "sticky", "animated"
	Data       map[string]interface{} `json:"data"`       // Custom data
	Children   []SemanticElement      `json:"children"`   // Nested elements
	Responsive bool                   `json:"responsive"` // Auto-adapt to screen size
}

type LayoutContext struct {
	ScreenWidth  int  `json:"screen_width"`
	ScreenHeight int  `json:"screen_height"`
	IsMobile     bool `json:"is_mobile"`
	IsTouch      bool `json:"is_touch"`
}

// Convert semantic description to actual 3D+text objects
func (h *Hub) CreateSemanticUI(elements []SemanticElement, context LayoutContext) []byte {
	var objects []map[string]interface{}
	
	// Auto-layout system
	layout := h.calculateLayout(elements, context)
	
	for i, element := range elements {
		pos := layout[i]
		
		switch element.Type {
		case "title":
			objects = append(objects, h.createTitle(element, pos, context))
		case "button":
			objects = append(objects, h.createButton(element, pos, context)...)
		case "card":
			objects = append(objects, h.createCard(element, pos, context))
		case "dashboard":
			objects = append(objects, h.createDashboard(element, pos, context)...)
		case "progress":
			objects = append(objects, h.createProgress(element, pos, context))
		case "notification":
			objects = append(objects, h.createNotification(element, pos, context))
		}
	}
	
	response := map[string]interface{}{
		"type":    "create",
		"objects": objects,
	}
	
	jsonData, _ := json.Marshal(response)
	return jsonData
}

// Intelligent layout calculation
func (h *Hub) calculateLayout(elements []SemanticElement, context LayoutContext) []Position {
	positions := make([]Position, len(elements))
	
	// Responsive grid system
	cols := 1
	if context.ScreenWidth > 768 && !context.IsMobile {
		cols = 3
	} else if context.ScreenWidth > 480 {
		cols = 2
	}
	
	// Auto-spacing based on screen size
	spacing := 2.0
	if context.IsMobile {
		spacing = 1.5
	}
	
	// Calculate positions using golden ratio and design principles
	for i, element := range elements {
		row := i / cols
		col := i % cols
		
		// Base position (flip X to match renderer coordinate system)
		x := -((float64(col) - float64(cols-1)/2) * spacing)
		y := float64(len(elements)/cols-row) * spacing * 0.8
		z := 0.0
		
		// Apply semantic positioning
		switch element.Position {
		case "top":
			y += 3.0
		case "bottom":
			y -= 3.0
		case "left":
			x += 2.0  // Flipped: left means positive X in renderer
		case "right":
			x -= 2.0  // Flipped: right means negative X in renderer
		case "center":
			x, y = 0, 0
		case "floating":
			z = 1.0
		}
		
		// Apply intent-based positioning
		switch element.Intent {
		case "primary":
			// Primary elements get prominence
			y += 0.5
			z += 0.2
		case "warning":
			// Warnings float forward
			z += 0.5
		}
		
		positions[i] = Position{X: x, Y: y, Z: z}
	}
	
	return positions
}

type Position struct {
	X, Y, Z float64
}

// Create semantic components
func (h *Hub) createTitle(element SemanticElement, pos Position, context LayoutContext) map[string]interface{} {
	fontSize := 2.0
	switch element.Size {
	case "small":
		fontSize = 1.2
	case "medium":
		fontSize = 1.8
	case "large":
		fontSize = 2.5
	case "xl":
		fontSize = 3.5
	}
	
	// Responsive sizing
	if context.IsMobile {
		fontSize *= 0.8
	}
	
	color := map[string]float64{"r": 1, "g": 1, "b": 1, "a": 1}
	switch element.Intent {
	case "primary":
		color = map[string]float64{"r": 0, "g": 0.6, "b": 1, "a": 1}
	case "warning":
		color = map[string]float64{"r": 1, "g": 0.8, "b": 0, "a": 1}
	case "success":
		color = map[string]float64{"r": 0, "g": 1, "b": 0.3, "a": 1}
	}
	
	animation := map[string]interface{}{}
	if element.Behavior == "animated" {
		animation = map[string]interface{}{
			"duration": 3,
			"easing":   "ease-in-out",
			"loop":     true,
		}
	}
	
	return map[string]interface{}{
		"id":   "title_" + element.Content,
		"type": "text",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y, "z": pos.Z},
		},
		"color": color,
		"data": map[string]interface{}{
			"text":     element.Content,
			"fontSize": fontSize,
		},
		"animation": animation,
	}
}

func (h *Hub) createButton(element SemanticElement, pos Position, context LayoutContext) []map[string]interface{} {
	width, height := 1.5, 0.5
	if element.Size == "large" {
		width, height = 2.0, 0.7
	}
	
	bgColor := map[string]float64{"r": 0.2, "g": 0.4, "b": 0.8, "a": 0.8}
	switch element.Intent {
	case "primary":
		bgColor = map[string]float64{"r": 0, "g": 0.5, "b": 1, "a": 0.9}
	case "warning":
		bgColor = map[string]float64{"r": 1, "g": 0.6, "b": 0, "a": 0.9}
	case "success":
		bgColor = map[string]float64{"r": 0, "g": 0.8, "b": 0.2, "a": 0.9}
	}
	
	// Create button background
	buttonBg := map[string]interface{}{
		"id":   "button_" + element.Content,
		"type": "cube",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y, "z": pos.Z},
			"scale":    map[string]float64{"x": width, "y": height, "z": 0.1},
		},
		"color": bgColor,
	}
	
	// Create button text
	textColor := map[string]float64{"r": 1, "g": 1, "b": 1, "a": 1}
	if element.Intent == "warning" {
		textColor = map[string]float64{"r": 0, "g": 0, "b": 0, "a": 1}
	}
	
	buttonText := map[string]interface{}{
		"id":   "button_text_" + element.Content,
		"type": "text",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y, "z": pos.Z + 0.05},
		},
		"data": map[string]interface{}{
			"text":     element.Content,
			"fontSize": 0.6,
		},
		"color": textColor,
	}
	
	return []map[string]interface{}{buttonBg, buttonText}
}

func (h *Hub) createCard(element SemanticElement, pos Position, context LayoutContext) map[string]interface{} {
	width, height := 2.0, 1.5
	if context.IsMobile {
		width *= 0.8
	}
	
	return map[string]interface{}{
		"id":   "card_" + element.Content,
		"type": "cube",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y, "z": pos.Z},
			"scale":    map[string]float64{"x": width, "y": height, "z": 0.05},
		},
		"color": map[string]float64{"r": 0.9, "g": 0.9, "b": 0.9, "a": 0.95},
	}
}

func (h *Hub) createDashboard(element SemanticElement, pos Position, context LayoutContext) []map[string]interface{} {
	var objects []map[string]interface{}
	
	// Create dashboard background
	objects = append(objects, map[string]interface{}{
		"id":   "dashboard_bg",
		"type": "plane",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y, "z": pos.Z - 0.1},
			"scale":    map[string]float64{"x": 5, "y": 3, "z": 1},
		},
		"color": map[string]float64{"r": 0.1, "g": 0.1, "b": 0.15, "a": 0.9},
	})
	
	// Create dashboard widgets
	widgets := []string{"CPU", "Memory", "Network", "Storage"}
	for i, widget := range widgets {
		angle := float64(i) * math.Pi / 2
		x := pos.X + 1.5*math.Cos(angle)
		y := pos.Y + 1.5*math.Sin(angle)
		
		objects = append(objects, map[string]interface{}{
			"id":   "widget_" + widget,
			"type": "cube",
			"transform": map[string]interface{}{
				"position": map[string]float64{"x": x, "y": y, "z": pos.Z + 0.1},
				"scale":    map[string]float64{"x": 0.8, "y": 0.8, "z": 0.1},
			},
			"color": map[string]float64{"r": 0, "g": 0.7, "b": 0.9, "a": 0.8},
			"animation": map[string]interface{}{
				"duration": 2 + float64(i)*0.5,
				"easing":   "ease-in-out",
				"loop":     true,
			},
		})
	}
	
	return objects
}

func (h *Hub) createProgress(element SemanticElement, pos Position, context LayoutContext) map[string]interface{} {
	progress := 0.7 // Default 70%
	if val, ok := element.Data["progress"].(float64); ok {
		progress = val
	}
	
	return map[string]interface{}{
		"id":   "progress_" + element.Content,
		"type": "cube",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y, "z": pos.Z},
			"scale":    map[string]float64{"x": 3 * progress, "y": 0.2, "z": 0.1},
		},
		"color": map[string]float64{"r": 0, "g": 1, "b": 0.3, "a": 0.8},
	}
}

func (h *Hub) createNotification(element SemanticElement, pos Position, context LayoutContext) map[string]interface{} {
	color := map[string]float64{"r": 0.2, "g": 0.6, "b": 1, "a": 0.9}
	switch element.Intent {
	case "warning":
		color = map[string]float64{"r": 1, "g": 0.8, "b": 0, "a": 0.9}
	case "error":
		color = map[string]float64{"r": 1, "g": 0.3, "b": 0.3, "a": 0.9}
	case "success":
		color = map[string]float64{"r": 0, "g": 1, "b": 0.4, "a": 0.9}
	}
	
	return map[string]interface{}{
		"id":   "notification_" + element.Content,
		"type": "cube",
		"transform": map[string]interface{}{
			"position": map[string]float64{"x": pos.X, "y": pos.Y + 2.5, "z": pos.Z + 1},
			"scale":    map[string]float64{"x": 2.5, "y": 0.6, "z": 0.1},
		},
		"color": color,
		"animation": map[string]interface{}{
			"duration": 1,
			"easing":   "bounce",
			"loop":     false,
		},
	}
}