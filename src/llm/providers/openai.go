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

type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewOpenAIProvider(apiKey string, baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *OpenAIProvider) GenerateResponse(ctx context.Context, prompt string, config map[string]interface{}) (*llm.LLMResponse, error) {
	startTime := time.Now()

	// Parse configuration
	model := p.getConfigString(config, "model", "gpt-3.5-turbo")
	maxTokens := p.getConfigInt(config, "max_tokens", 2048)
	temperature := p.getConfigFloat(config, "temperature", 0.7)
	topP := p.getConfigFloat(config, "top_p", 1.0)

	// Build request
	request := OpenAIRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
		Stream:      false,
	}

	// Serialize request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	// Make request
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var openAIResponse OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(openAIResponse.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	processingTime := time.Since(startTime).Seconds()

	// Build LLM response
	llmResponse := &llm.LLMResponse{
		Content:   openAIResponse.Choices[0].Message.Content,
		Type:      "text",
		Metadata: map[string]interface{}{
			"model":         openAIResponse.Model,
			"finish_reason": openAIResponse.Choices[0].FinishReason,
			"id":            openAIResponse.ID,
		},
		Usage: &llm.UsageStats{
			InputTokens:     openAIResponse.Usage.PromptTokens,
			OutputTokens:    openAIResponse.Usage.CompletionTokens,
			TotalTokens:     openAIResponse.Usage.TotalTokens,
			ProcessingTime:  processingTime,
			Cost:            p.calculateCost(model, openAIResponse.Usage),
		},
		Timestamp: time.Now(),
	}

	logging.Debug("generated OpenAI response", map[string]interface{}{
		"model":           model,
		"input_tokens":    openAIResponse.Usage.PromptTokens,
		"output_tokens":   openAIResponse.Usage.CompletionTokens,
		"processing_time": processingTime,
	})

	return llmResponse, nil
}

func (p *OpenAIProvider) GetCapabilities() []string {
	return []string{
		"text_generation",
		"conversation",
		"code_generation",
		"analysis",
		"translation",
		"summarization",
	}
}

func (p *OpenAIProvider) ValidateConfig(config map[string]interface{}) error {
	// Validate model
	if model, exists := config["model"]; exists {
		if modelStr, ok := model.(string); ok {
			validModels := []string{
				"gpt-3.5-turbo",
				"gpt-3.5-turbo-16k",
				"gpt-4",
				"gpt-4-32k",
				"gpt-4-1106-preview",
				"gpt-4-vision-preview",
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
			if tempFloat < 0 || tempFloat > 2 {
				return fmt.Errorf("temperature must be between 0 and 2")
			}
		} else {
			return fmt.Errorf("temperature must be a number")
		}
	}

	// Validate max_tokens
	if maxTokens, exists := config["max_tokens"]; exists {
		if tokensInt, ok := maxTokens.(float64); ok {
			if tokensInt < 1 || tokensInt > 4096 {
				return fmt.Errorf("max_tokens must be between 1 and 4096")
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

	return nil
}

func (p *OpenAIProvider) getConfigString(config map[string]interface{}, key string, defaultValue string) string {
	if value, exists := config[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

func (p *OpenAIProvider) getConfigInt(config map[string]interface{}, key string, defaultValue int) int {
	if value, exists := config[key]; exists {
		if num, ok := value.(float64); ok {
			return int(num)
		}
	}
	return defaultValue
}

func (p *OpenAIProvider) getConfigFloat(config map[string]interface{}, key string, defaultValue float64) float64 {
	if value, exists := config[key]; exists {
		if num, ok := value.(float64); ok {
			return num
		}
	}
	return defaultValue
}

func (p *OpenAIProvider) calculateCost(model string, usage Usage) float64 {
	// Pricing per 1K tokens (as of 2024)
	var inputCost, outputCost float64
	
	switch model {
	case "gpt-3.5-turbo":
		inputCost = 0.001
		outputCost = 0.002
	case "gpt-3.5-turbo-16k":
		inputCost = 0.003
		outputCost = 0.004
	case "gpt-4":
		inputCost = 0.03
		outputCost = 0.06
	case "gpt-4-32k":
		inputCost = 0.06
		outputCost = 0.12
	case "gpt-4-1106-preview":
		inputCost = 0.01
		outputCost = 0.03
	case "gpt-4-vision-preview":
		inputCost = 0.01
		outputCost = 0.03
	default:
		inputCost = 0.001
		outputCost = 0.002
	}

	return (float64(usage.PromptTokens)/1000)*inputCost + (float64(usage.CompletionTokens)/1000)*outputCost
}