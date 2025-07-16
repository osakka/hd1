package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"holodeck1/llm"
	"holodeck1/logging"
)

type ClaudeProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Messages    []ClaudeMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
	TopK        int             `json:"top_k,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Role         string        `json:"role"`
	Content      []ClaudeContent `json:"content"`
	Model        string        `json:"model"`
	StopReason   string        `json:"stop_reason"`
	StopSequence string        `json:"stop_sequence"`
	Usage        ClaudeUsage   `json:"usage"`
}

type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func NewClaudeProvider(apiKey string, baseURL string) *ClaudeProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &ClaudeProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *ClaudeProvider) GenerateResponse(ctx context.Context, prompt string, config map[string]interface{}) (*llm.LLMResponse, error) {
	startTime := time.Now()

	// Parse configuration
	model := p.getConfigString(config, "model", "claude-3-haiku-20240307")
	maxTokens := p.getConfigInt(config, "max_tokens", 2048)
	temperature := p.getConfigFloat(config, "temperature", 0.7)
	topP := p.getConfigFloat(config, "top_p", 1.0)
	topK := p.getConfigInt(config, "top_k", 0)

	// Build request
	request := ClaudeRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: temperature,
		TopP:        topP,
		TopK:        topK,
		Stream:      false,
	}

	// Serialize request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/messages", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Make request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Claude API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var claudeResponse ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(claudeResponse.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	processingTime := time.Since(startTime).Seconds()

	// Build LLM response
	llmResponse := &llm.LLMResponse{
		Content:   claudeResponse.Content[0].Text,
		Type:      "text",
		Metadata: map[string]interface{}{
			"model":         claudeResponse.Model,
			"stop_reason":   claudeResponse.StopReason,
			"stop_sequence": claudeResponse.StopSequence,
			"id":            claudeResponse.ID,
		},
		Usage: &llm.UsageStats{
			InputTokens:     claudeResponse.Usage.InputTokens,
			OutputTokens:    claudeResponse.Usage.OutputTokens,
			TotalTokens:     claudeResponse.Usage.InputTokens + claudeResponse.Usage.OutputTokens,
			ProcessingTime:  processingTime,
			Cost:            p.calculateCost(model, claudeResponse.Usage),
		},
		Timestamp: time.Now(),
	}

	logging.Debug("generated Claude response", map[string]interface{}{
		"model":           model,
		"input_tokens":    claudeResponse.Usage.InputTokens,
		"output_tokens":   claudeResponse.Usage.OutputTokens,
		"processing_time": processingTime,
	})

	return llmResponse, nil
}

func (p *ClaudeProvider) GetCapabilities() []string {
	return []string{
		"text_generation",
		"conversation",
		"code_generation",
		"analysis",
		"translation",
		"summarization",
		"reasoning",
		"math",
	}
}

func (p *ClaudeProvider) ValidateConfig(config map[string]interface{}) error {
	// Validate model
	if model, exists := config["model"]; exists {
		if modelStr, ok := model.(string); ok {
			validModels := []string{
				"claude-3-haiku-20240307",
				"claude-3-sonnet-20240229",
				"claude-3-opus-20240229",
				"claude-2.1",
				"claude-2.0",
				"claude-instant-1.2",
			}
			valid := false
			for _, validModel := range validModels {
				if modelStr == validModel {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid model: %s", modelStr)
			}
		} else {
			return fmt.Errorf("model must be a string")
		}
	}

	// Validate temperature
	if temp, exists := config["temperature"]; exists {
		if tempFloat, ok := temp.(float64); ok {
			if tempFloat < 0 || tempFloat > 1 {
				return fmt.Errorf("temperature must be between 0 and 1")
			}
		} else {
			return fmt.Errorf("temperature must be a number")
		}
	}

	// Validate max_tokens
	if maxTokens, exists := config["max_tokens"]; exists {
		if tokensInt, ok := maxTokens.(float64); ok {
			if tokensInt < 1 || tokensInt > 8192 {
				return fmt.Errorf("max_tokens must be between 1 and 8192")
			}
		} else {
			return fmt.Errorf("max_tokens must be a number")
		}
	}

	// Validate top_p
	if topP, exists := config["top_p"]; exists {
		if topPFloat, ok := topP.(float64); ok {
			if topPFloat < 0 || topPFloat > 1 {
				return fmt.Errorf("top_p must be between 0 and 1")
			}
		} else {
			return fmt.Errorf("top_p must be a number")
		}
	}

	// Validate top_k
	if topK, exists := config["top_k"]; exists {
		if topKInt, ok := topK.(float64); ok {
			if topKInt < 0 || topKInt > 500 {
				return fmt.Errorf("top_k must be between 0 and 500")
			}
		} else {
			return fmt.Errorf("top_k must be a number")
		}
	}

	return nil
}

func (p *ClaudeProvider) getConfigString(config map[string]interface{}, key string, defaultValue string) string {
	if value, exists := config[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (p *ClaudeProvider) getConfigInt(config map[string]interface{}, key string, defaultValue int) int {
	if value, exists := config[key]; exists {
		if num, ok := value.(float64); ok {
			return int(num)
		}
	}
	return defaultValue
}

func (p *ClaudeProvider) getConfigFloat(config map[string]interface{}, key string, defaultValue float64) float64 {
	if value, exists := config[key]; exists {
		if num, ok := value.(float64); ok {
			return num
		}
	}
	return defaultValue
}

func (p *ClaudeProvider) calculateCost(model string, usage ClaudeUsage) float64 {
	// Pricing per 1K tokens (as of 2024)
	var inputCost, outputCost float64
	
	switch model {
	case "claude-3-haiku-20240307":
		inputCost = 0.00025
		outputCost = 0.00125
	case "claude-3-sonnet-20240229":
		inputCost = 0.003
		outputCost = 0.015
	case "claude-3-opus-20240229":
		inputCost = 0.015
		outputCost = 0.075
	case "claude-2.1":
		inputCost = 0.008
		outputCost = 0.024
	case "claude-2.0":
		inputCost = 0.008
		outputCost = 0.024
	case "claude-instant-1.2":
		inputCost = 0.0008
		outputCost = 0.0024
	default:
		inputCost = 0.00025
		outputCost = 0.00125
	}

	return (float64(usage.InputTokens)/1000)*inputCost + (float64(usage.OutputTokens)/1000)*outputCost
}