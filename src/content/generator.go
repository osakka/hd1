package content

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/llm"
	"holodeck1/logging"
)

type Generator struct {
	db         *database.DB
	llmManager *llm.Manager
	templates  map[string]*Template
	jobs       map[uuid.UUID]*GenerationJob
	mu         sync.RWMutex
}

type Template struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Prompt      string                 `json:"prompt"`
	Variables   []TemplateVariable     `json:"variables"`
	Configuration map[string]interface{} `json:"configuration"`
	OutputFormat string                 `json:"output_format"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TemplateVariable struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
	Options     []string    `json:"options,omitempty"`
}

type GenerationJob struct {
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"session_id"`
	UserID       uuid.UUID              `json:"user_id"`
	TemplateID   uuid.UUID              `json:"template_id"`
	Variables    map[string]interface{} `json:"variables"`
	Status       string                 `json:"status"`
	Progress     int                    `json:"progress"`
	Result       *GenerationResult      `json:"result"`
	Error        string                 `json:"error"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	CompletedAt  *time.Time             `json:"completed_at"`
}

type GenerationResult struct {
	Content      string                 `json:"content"`
	Type         string                 `json:"type"`
	Format       string                 `json:"format"`
	Assets       []GeneratedAsset       `json:"assets"`
	Metadata     map[string]interface{} `json:"metadata"`
	Usage        *llm.UsageStats        `json:"usage"`
}

type GeneratedAsset struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	URL      string    `json:"url"`
	MimeType string    `json:"mime_type"`
	Size     int64     `json:"size"`
}

type GenerateRequest struct {
	SessionID  uuid.UUID              `json:"session_id"`
	UserID     uuid.UUID              `json:"user_id"`
	TemplateID uuid.UUID              `json:"template_id"`
	Variables  map[string]interface{} `json:"variables"`
	Options    map[string]interface{} `json:"options"`
}

type CreateTemplateRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Category      string                 `json:"category"`
	Prompt        string                 `json:"prompt"`
	Variables     []TemplateVariable     `json:"variables"`
	Configuration map[string]interface{} `json:"configuration"`
	OutputFormat  string                 `json:"output_format"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type TemplateFilter struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

func NewGenerator(db *database.DB, llmManager *llm.Manager) *Generator {
	return &Generator{
		db:         db,
		llmManager: llmManager,
		templates:  make(map[string]*Template),
		jobs:       make(map[uuid.UUID]*GenerationJob),
	}
}

func (g *Generator) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*Template, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	templateID := uuid.New()
	template := &Template{
		ID:            templateID,
		Name:          req.Name,
		Description:   req.Description,
		Type:          req.Type,
		Category:      req.Category,
		Prompt:        req.Prompt,
		Variables:     req.Variables,
		Configuration: req.Configuration,
		OutputFormat:  req.OutputFormat,
		Metadata:      req.Metadata,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO content_templates (id, name, description, type, category, prompt, variables, configuration, output_format, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	
	variablesJSON, _ := json.Marshal(template.Variables)
	configJSON, _ := json.Marshal(template.Configuration)
	metadataJSON, _ := json.Marshal(template.Metadata)
	
	_, err := g.db.Conn.ExecContext(ctx, query,
		template.ID, template.Name, template.Description, template.Type,
		template.Category, template.Prompt, variablesJSON, configJSON,
		template.OutputFormat, metadataJSON, template.CreatedAt, template.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store template: %w", err)
	}

	g.templates[template.Name] = template

	logging.Info("created content template", map[string]interface{}{
		"template_id": templateID,
		"name":        req.Name,
		"type":        req.Type,
		"category":    req.Category,
	})

	return template, nil
}

func (g *Generator) GetTemplate(ctx context.Context, templateID uuid.UUID) (*Template, error) {
	// Check memory first
	g.mu.RLock()
	for _, template := range g.templates {
		if template.ID == templateID {
			g.mu.RUnlock()
			return template, nil
		}
	}
	g.mu.RUnlock()

	// Load from database
	query := `
		SELECT id, name, description, type, category, prompt, variables, configuration, output_format, metadata, created_at, updated_at
		FROM content_templates
		WHERE id = $1
	`
	
	var template Template
	var variablesJSON, configJSON, metadataJSON []byte
	
	err := g.db.Conn.QueryRowContext(ctx, query, templateID).Scan(
		&template.ID, &template.Name, &template.Description, &template.Type,
		&template.Category, &template.Prompt, &variablesJSON, &configJSON,
		&template.OutputFormat, &metadataJSON, &template.CreatedAt, &template.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Unmarshal JSON fields
	json.Unmarshal(variablesJSON, &template.Variables)
	json.Unmarshal(configJSON, &template.Configuration)
	json.Unmarshal(metadataJSON, &template.Metadata)

	g.mu.Lock()
	g.templates[template.Name] = &template
	g.mu.Unlock()

	return &template, nil
}

func (g *Generator) ListTemplates(ctx context.Context, filter *TemplateFilter) ([]*Template, error) {
	query := `
		SELECT id, name, description, type, category, prompt, variables, configuration, output_format, metadata, created_at, updated_at
		FROM content_templates
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filter.Type)
		argIndex++
	}

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, filter.Category)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := g.db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []*Template
	for rows.Next() {
		var template Template
		var variablesJSON, configJSON, metadataJSON []byte
		
		err := rows.Scan(
			&template.ID, &template.Name, &template.Description, &template.Type,
			&template.Category, &template.Prompt, &variablesJSON, &configJSON,
			&template.OutputFormat, &metadataJSON, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		// Unmarshal JSON fields
		json.Unmarshal(variablesJSON, &template.Variables)
		json.Unmarshal(configJSON, &template.Configuration)
		json.Unmarshal(metadataJSON, &template.Metadata)

		templates = append(templates, &template)
	}

	return templates, nil
}

func (g *Generator) GenerateContent(ctx context.Context, req *GenerateRequest) (*GenerationJob, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Get template
	template, err := g.GetTemplate(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}

	// Validate variables
	err = g.validateVariables(template, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("variable validation failed: %w", err)
	}

	// Create job
	jobID := uuid.New()
	job := &GenerationJob{
		ID:         jobID,
		SessionID:  req.SessionID,
		UserID:     req.UserID,
		TemplateID: req.TemplateID,
		Variables:  req.Variables,
		Status:     "pending",
		Progress:   0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Store in database
	query := `
		INSERT INTO content_generation_jobs (id, session_id, user_id, template_id, variables, status, progress, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	variablesJSON, _ := json.Marshal(job.Variables)
	
	_, err = g.db.Conn.ExecContext(ctx, query,
		job.ID, job.SessionID, job.UserID, job.TemplateID,
		variablesJSON, job.Status, job.Progress, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store job: %w", err)
	}

	g.jobs[jobID] = job

	// Start generation asynchronously
	go g.processGenerationJob(ctx, job, template)

	logging.Info("started content generation job", map[string]interface{}{
		"job_id":      jobID,
		"session_id":  req.SessionID,
		"user_id":     req.UserID,
		"template_id": req.TemplateID,
	})

	return job, nil
}

func (g *Generator) processGenerationJob(ctx context.Context, job *GenerationJob, template *Template) {
	// Update job status
	g.updateJobStatus(ctx, job, "processing", 10)

	// Build prompt from template
	prompt, err := g.buildPrompt(template, job.Variables)
	if err != nil {
		g.updateJobError(ctx, job, fmt.Sprintf("failed to build prompt: %v", err))
		return
	}

	g.updateJobStatus(ctx, job, "processing", 30)

	// Get LLM provider configuration
	config := template.Configuration
	if config == nil {
		config = make(map[string]interface{})
	}

	// Find an available avatar or create a temporary one
	avatarID, err := g.getOrCreateAvatar(ctx, job.SessionID, template)
	if err != nil {
		g.updateJobError(ctx, job, fmt.Sprintf("failed to get avatar: %v", err))
		return
	}

	g.updateJobStatus(ctx, job, "processing", 50)

	// Generate content using LLM
	chatRequest := &llm.ChatRequest{
		AvatarID: avatarID,
		Message:  prompt,
		Context:  "content_generation",
		UserID:   job.UserID,
	}

	chatResponse, err := g.llmManager.Chat(ctx, chatRequest)
	if err != nil {
		g.updateJobError(ctx, job, fmt.Sprintf("failed to generate content: %v", err))
		return
	}

	g.updateJobStatus(ctx, job, "processing", 80)

	// Process the generated content
	result := &GenerationResult{
		Content:  chatResponse.Response,
		Type:     template.Type,
		Format:   template.OutputFormat,
		Assets:   []GeneratedAsset{},
		Metadata: chatResponse.Metadata,
		Usage:    chatResponse.Usage,
	}

	// Post-process content based on template type
	err = g.postProcessContent(ctx, job, template, result)
	if err != nil {
		g.updateJobError(ctx, job, fmt.Sprintf("failed to post-process content: %v", err))
		return
	}

	g.updateJobStatus(ctx, job, "processing", 90)

	// Store result
	job.Result = result
	completedAt := time.Now()
	job.CompletedAt = &completedAt
	job.Status = "completed"
	job.Progress = 100
	job.UpdatedAt = time.Now()

	// Update database
	resultJSON, _ := json.Marshal(result)
	query := `
		UPDATE content_generation_jobs 
		SET result = $1, status = $2, progress = $3, updated_at = $4, completed_at = $5
		WHERE id = $6
	`
	_, err = g.db.Conn.ExecContext(ctx, query, resultJSON, job.Status, job.Progress, job.UpdatedAt, job.CompletedAt, job.ID)
	if err != nil {
		logging.Error("failed to update job result", map[string]interface{}{
			"job_id": job.ID,
			"error":  err.Error(),
		})
	}

	logging.Info("completed content generation job", map[string]interface{}{
		"job_id":        job.ID,
		"template_id":   job.TemplateID,
		"content_length": len(result.Content),
		"processing_time": time.Since(job.CreatedAt).Seconds(),
	})
}

func (g *Generator) buildPrompt(template *Template, variables map[string]interface{}) (string, error) {
	prompt := template.Prompt
	
	// Replace variables in prompt
	for _, variable := range template.Variables {
		placeholder := fmt.Sprintf("{{%s}}", variable.Name)
		value, exists := variables[variable.Name]
		if !exists {
			if variable.Required {
				return "", fmt.Errorf("required variable %s not provided", variable.Name)
			}
			value = variable.Default
		}
		
		valueStr := fmt.Sprintf("%v", value)
		prompt = strings.ReplaceAll(prompt, placeholder, valueStr)
	}
	
	return prompt, nil
}

func (g *Generator) validateVariables(template *Template, variables map[string]interface{}) error {
	for _, variable := range template.Variables {
		value, exists := variables[variable.Name]
		
		if !exists {
			if variable.Required {
				return fmt.Errorf("required variable %s not provided", variable.Name)
			}
			continue
		}
		
		// Type validation
		switch variable.Type {
		case "string":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("variable %s must be a string", variable.Name)
			}
		case "number":
			if _, ok := value.(float64); !ok {
				return fmt.Errorf("variable %s must be a number", variable.Name)
			}
		case "boolean":
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("variable %s must be a boolean", variable.Name)
			}
		case "array":
			if _, ok := value.([]interface{}); !ok {
				return fmt.Errorf("variable %s must be an array", variable.Name)
			}
		}
		
		// Options validation
		if len(variable.Options) > 0 {
			valueStr := fmt.Sprintf("%v", value)
			valid := false
			for _, option := range variable.Options {
				if valueStr == option {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("variable %s must be one of: %v", variable.Name, variable.Options)
			}
		}
	}
	
	return nil
}

func (g *Generator) getOrCreateAvatar(ctx context.Context, sessionID uuid.UUID, template *Template) (uuid.UUID, error) {
	// Look for existing content generation avatar
	filter := &llm.AvatarFilter{
		SessionID: sessionID,
		Type:      "content_generator",
		State:     "active",
		Limit:     1,
	}
	
	avatars, err := g.llmManager.GetAvatarsBySession(ctx, filter)
	if err != nil {
		return uuid.Nil, err
	}
	
	if len(avatars) > 0 {
		return avatars[0].ID, nil
	}
	
	// Create new content generation avatar
	avatarRequest := &llm.CreateAvatarRequest{
		SessionID:     sessionID,
		Name:          "Content Generator",
		Type:          "content_generator",
		Provider:      "openai",
		Model:         "gpt-3.5-turbo",
		Personality:   "I am a helpful AI assistant specialized in generating high-quality content based on templates and user requirements.",
		Knowledge:     "I have extensive knowledge about content creation, writing, and various formats including text, code, and structured data.",
		Capabilities:  []string{"text_generation", "code_generation", "analysis", "creative_writing"},
		Configuration: template.Configuration,
		Position:      &llm.Position3D{X: 0, Y: 0, Z: 0},
		Avatar3D:      &llm.Avatar3DConfig{
			ModelURL: "/api/avatars/content_generator/asset",
			Scale:    1.0,
			Animations: []string{"idle", "thinking", "writing"},
		},
	}
	
	avatar, err := g.llmManager.CreateAvatar(ctx, avatarRequest)
	if err != nil {
		return uuid.Nil, err
	}
	
	return avatar.ID, nil
}

func (g *Generator) postProcessContent(ctx context.Context, job *GenerationJob, template *Template, result *GenerationResult) error {
	// Post-processing based on template type
	switch template.Type {
	case "code":
		// Validate and format code
		return g.postProcessCode(result)
	case "json":
		// Validate and format JSON
		return g.postProcessJSON(result)
	case "markdown":
		// Process markdown
		return g.postProcessMarkdown(result)
	case "html":
		// Process HTML
		return g.postProcessHTML(result)
	default:
		// Default text processing
		return g.postProcessText(result)
	}
}

func (g *Generator) postProcessCode(result *GenerationResult) error {
	// Basic code formatting and validation
	content := strings.TrimSpace(result.Content)
	
	// Remove markdown code blocks if present
	if strings.HasPrefix(content, "```") {
		lines := strings.Split(content, "\n")
		if len(lines) > 2 {
			content = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}
	
	result.Content = content
	return nil
}

func (g *Generator) postProcessJSON(result *GenerationResult) error {
	// Validate and format JSON
	var jsonData interface{}
	err := json.Unmarshal([]byte(result.Content), &jsonData)
	if err != nil {
		return fmt.Errorf("generated content is not valid JSON: %w", err)
	}
	
	// Pretty format JSON
	formatted, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}
	
	result.Content = string(formatted)
	return nil
}

func (g *Generator) postProcessMarkdown(result *GenerationResult) error {
	// Basic markdown processing
	content := strings.TrimSpace(result.Content)
	result.Content = content
	return nil
}

func (g *Generator) postProcessHTML(result *GenerationResult) error {
	// Basic HTML processing
	content := strings.TrimSpace(result.Content)
	result.Content = content
	return nil
}

func (g *Generator) postProcessText(result *GenerationResult) error {
	// Basic text processing
	content := strings.TrimSpace(result.Content)
	result.Content = content
	return nil
}

func (g *Generator) updateJobStatus(ctx context.Context, job *GenerationJob, status string, progress int) {
	job.Status = status
	job.Progress = progress
	job.UpdatedAt = time.Now()
	
	query := `UPDATE content_generation_jobs SET status = $1, progress = $2, updated_at = $3 WHERE id = $4`
	_, err := g.db.Conn.ExecContext(ctx, query, status, progress, job.UpdatedAt, job.ID)
	if err != nil {
		logging.Error("failed to update job status", map[string]interface{}{
			"job_id": job.ID,
			"error":  err.Error(),
		})
	}
}

func (g *Generator) updateJobError(ctx context.Context, job *GenerationJob, errorMsg string) {
	job.Status = "failed"
	job.Error = errorMsg
	job.UpdatedAt = time.Now()
	
	query := `UPDATE content_generation_jobs SET status = $1, error = $2, updated_at = $3 WHERE id = $4`
	_, err := g.db.Conn.ExecContext(ctx, query, job.Status, job.Error, job.UpdatedAt, job.ID)
	if err != nil {
		logging.Error("failed to update job error", map[string]interface{}{
			"job_id": job.ID,
			"error":  err.Error(),
		})
	}
}

func (g *Generator) GetJob(ctx context.Context, jobID uuid.UUID) (*GenerationJob, error) {
	g.mu.RLock()
	job, exists := g.jobs[jobID]
	g.mu.RUnlock()
	
	if exists {
		return job, nil
	}
	
	// Load from database
	query := `
		SELECT id, session_id, user_id, template_id, variables, status, progress, result, error, created_at, updated_at, completed_at
		FROM content_generation_jobs
		WHERE id = $1
	`
	
	var dbJob GenerationJob
	var variablesJSON, resultJSON []byte
	var completedAt sql.NullTime
	
	err := g.db.Conn.QueryRowContext(ctx, query, jobID).Scan(
		&dbJob.ID, &dbJob.SessionID, &dbJob.UserID, &dbJob.TemplateID,
		&variablesJSON, &dbJob.Status, &dbJob.Progress, &resultJSON,
		&dbJob.Error, &dbJob.CreatedAt, &dbJob.UpdatedAt, &completedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}
	
	// Unmarshal JSON fields
	json.Unmarshal(variablesJSON, &dbJob.Variables)
	if len(resultJSON) > 0 {
		json.Unmarshal(resultJSON, &dbJob.Result)
	}
	
	if completedAt.Valid {
		dbJob.CompletedAt = &completedAt.Time
	}
	
	g.mu.Lock()
	g.jobs[jobID] = &dbJob
	g.mu.Unlock()
	
	return &dbJob, nil
}

func (g *Generator) GetJobsBySession(ctx context.Context, sessionID uuid.UUID, limit int) ([]*GenerationJob, error) {
	query := `
		SELECT id, session_id, user_id, template_id, variables, status, progress, result, error, created_at, updated_at, completed_at
		FROM content_generation_jobs
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := g.db.Conn.QueryContext(ctx, query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()
	
	var jobs []*GenerationJob
	for rows.Next() {
		var job GenerationJob
		var variablesJSON, resultJSON []byte
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&job.ID, &job.SessionID, &job.UserID, &job.TemplateID,
			&variablesJSON, &job.Status, &job.Progress, &resultJSON,
			&job.Error, &job.CreatedAt, &job.UpdatedAt, &completedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		
		// Unmarshal JSON fields
		json.Unmarshal(variablesJSON, &job.Variables)
		if len(resultJSON) > 0 {
			json.Unmarshal(resultJSON, &job.Result)
		}
		
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}
		
		jobs = append(jobs, &job)
	}
	
	return jobs, nil
}

func (g *Generator) CleanupCompletedJobs(ctx context.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	// Find jobs completed more than 24 hours ago
	cutoff := time.Now().Add(-24 * time.Hour)
	
	for jobID, job := range g.jobs {
		if job.CompletedAt != nil && job.CompletedAt.Before(cutoff) {
			delete(g.jobs, jobID)
			logging.Info("cleaned up completed job", map[string]interface{}{
				"job_id": jobID,
			})
		}
	}
}

func (g *Generator) StartCleanupWorker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			g.CleanupCompletedJobs(ctx)
		}
	}
}