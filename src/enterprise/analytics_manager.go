package enterprise

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"holodeck1/database"
	"holodeck1/logging"
)

type AnalyticsManager struct {
	db *database.DB
}

type AnalyticsEvent struct {
	ID             uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	UserID         *uuid.UUID             `json:"user_id"`
	SessionID      *uuid.UUID             `json:"session_id"`
	EventType      string                 `json:"event_type"`
	EventCategory  string                 `json:"event_category"`
	EventAction    string                 `json:"event_action"`
	EventLabel     string                 `json:"event_label"`
	EventValue     *float64               `json:"event_value"`
	Properties     map[string]interface{} `json:"properties"`
	UserAgent      string                 `json:"user_agent"`
	IPAddress      string                 `json:"ip_address"`
	Referrer       string                 `json:"referrer"`
	Timestamp      time.Time              `json:"timestamp"`
	CreatedAt      time.Time              `json:"created_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type AnalyticsAggregate struct {
	ID             uuid.UUID              `json:"id"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	AggregateType  string                 `json:"aggregate_type"`
	Dimension      string                 `json:"dimension"`
	DimensionValue string                 `json:"dimension_value"`
	Metric         string                 `json:"metric"`
	Value          float64                `json:"value"`
	PeriodType     string                 `json:"period_type"`
	PeriodStart    time.Time              `json:"period_start"`
	PeriodEnd      time.Time              `json:"period_end"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Metadata       map[string]interface{} `json:"metadata"`
}

type EventTrackRequest struct {
	UserID        *uuid.UUID             `json:"user_id"`
	SessionID     *uuid.UUID             `json:"session_id"`
	EventType     string                 `json:"event_type"`
	EventCategory string                 `json:"event_category"`
	EventAction   string                 `json:"event_action"`
	EventLabel    string                 `json:"event_label"`
	EventValue    *float64               `json:"event_value"`
	Properties    map[string]interface{} `json:"properties"`
	UserAgent     string                 `json:"user_agent"`
	IPAddress     string                 `json:"ip_address"`
	Referrer      string                 `json:"referrer"`
}

type AnalyticsQuery struct {
	OrganizationID uuid.UUID  `json:"organization_id"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	EventType      string     `json:"event_type"`
	EventCategory  string     `json:"event_category"`
	UserID         *uuid.UUID `json:"user_id"`
	SessionID      *uuid.UUID `json:"session_id"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
}

type AggregateQuery struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	AggregateType  string    `json:"aggregate_type"`
	Dimension      string    `json:"dimension"`
	Metric         string    `json:"metric"`
	PeriodType     string    `json:"period_type"`
}

type AnalyticsReport struct {
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	StartTime     time.Time                `json:"start_time"`
	EndTime       time.Time                `json:"end_time"`
	Metrics       map[string]float64       `json:"metrics"`
	Dimensions    map[string][]Dimension   `json:"dimensions"`
	TimeSeries    []TimeSeriesPoint        `json:"time_series"`
	TopEvents     []EventSummary           `json:"top_events"`
	UserActivity  []UserActivitySummary    `json:"user_activity"`
}

type Dimension struct {
	Name  string  `json:"name"`
	Value string  `json:"value"`
	Count int     `json:"count"`
}

type TimeSeriesPoint struct {
	Time  time.Time          `json:"time"`
	Value float64            `json:"value"`
	Meta  map[string]float64 `json:"meta"`
}

type EventSummary struct {
	EventType     string  `json:"event_type"`
	EventCategory string  `json:"event_category"`
	EventAction   string  `json:"event_action"`
	Count         int     `json:"count"`
	UniqueUsers   int     `json:"unique_users"`
	AvgValue      float64 `json:"avg_value"`
}

type UserActivitySummary struct {
	UserID       uuid.UUID `json:"user_id"`
	EventCount   int       `json:"event_count"`
	SessionCount int       `json:"session_count"`
	LastActive   time.Time `json:"last_active"`
}

var EventCategories = []string{
	"session",
	"content",
	"collaboration",
	"navigation",
	"interaction",
	"error",
	"performance",
	"security",
}

var AggregateTypes = []string{
	"count",
	"unique_users",
	"unique_sessions",
	"sum",
	"average",
	"min",
	"max",
}

var PeriodTypes = []string{
	"hour",
	"day",
	"week",
	"month",
	"year",
}

func NewAnalyticsManager(db *database.DB) *AnalyticsManager {
	return &AnalyticsManager{
		db: db,
	}
}

func (am *AnalyticsManager) TrackEvent(ctx context.Context, orgID uuid.UUID, req *EventTrackRequest) error {
	// Validate required fields
	if req.EventType == "" {
		return fmt.Errorf("event type is required")
	}
	if req.EventCategory == "" {
		return fmt.Errorf("event category is required")
	}
	if req.EventAction == "" {
		return fmt.Errorf("event action is required")
	}

	// Validate event category
	validCategory := false
	for _, cat := range EventCategories {
		if req.EventCategory == cat {
			validCategory = true
			break
		}
	}
	if !validCategory {
		return fmt.Errorf("invalid event category: %s", req.EventCategory)
	}

	event := &AnalyticsEvent{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         req.UserID,
		SessionID:      req.SessionID,
		EventType:      req.EventType,
		EventCategory:  req.EventCategory,
		EventAction:    req.EventAction,
		EventLabel:     req.EventLabel,
		EventValue:     req.EventValue,
		Properties:     req.Properties,
		UserAgent:      req.UserAgent,
		IPAddress:      req.IPAddress,
		Referrer:       req.Referrer,
		Timestamp:      time.Now(),
		CreatedAt:      time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	if event.Properties == nil {
		event.Properties = make(map[string]interface{})
	}

	query := `
		INSERT INTO analytics_events (id, organization_id, user_id, session_id, event_type, event_category, event_action, event_label, event_value, properties, user_agent, ip_address, referrer, timestamp, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	propertiesJSON, _ := json.Marshal(event.Properties)
	metadataJSON, _ := json.Marshal(event.Metadata)

	_, err := am.db.Conn.ExecContext(ctx, query,
		event.ID, event.OrganizationID, event.UserID, event.SessionID,
		event.EventType, event.EventCategory, event.EventAction, event.EventLabel,
		event.EventValue, propertiesJSON, event.UserAgent, event.IPAddress,
		event.Referrer, event.Timestamp, event.CreatedAt, metadataJSON)
	if err != nil {
		return fmt.Errorf("failed to track event: %w", err)
	}

	// Update aggregates asynchronously
	go am.updateAggregates(context.Background(), event)

	logging.Debug("tracked analytics event", map[string]interface{}{
		"event_id":        event.ID,
		"organization_id": orgID,
		"event_type":      event.EventType,
		"event_category":  event.EventCategory,
		"event_action":    event.EventAction,
		"user_id":         event.UserID,
		"session_id":      event.SessionID,
	})

	return nil
}

func (am *AnalyticsManager) QueryEvents(ctx context.Context, query *AnalyticsQuery) ([]*AnalyticsEvent, error) {
	sqlQuery := `
		SELECT id, organization_id, user_id, session_id, event_type, event_category, event_action, event_label, event_value, properties, user_agent, ip_address, referrer, timestamp, created_at, metadata
		FROM analytics_events
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
	`
	args := []interface{}{query.OrganizationID, query.StartTime, query.EndTime}
	argIndex := 4

	if query.EventType != "" {
		sqlQuery += fmt.Sprintf(" AND event_type = $%d", argIndex)
		args = append(args, query.EventType)
		argIndex++
	}

	if query.EventCategory != "" {
		sqlQuery += fmt.Sprintf(" AND event_category = $%d", argIndex)
		args = append(args, query.EventCategory)
		argIndex++
	}

	if query.UserID != nil {
		sqlQuery += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *query.UserID)
		argIndex++
	}

	if query.SessionID != nil {
		sqlQuery += fmt.Sprintf(" AND session_id = $%d", argIndex)
		args = append(args, *query.SessionID)
		argIndex++
	}

	sqlQuery += " ORDER BY timestamp DESC"

	if query.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, query.Limit)
		argIndex++
	}

	if query.Offset > 0 {
		sqlQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, query.Offset)
	}

	rows, err := am.db.Conn.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []*AnalyticsEvent
	for rows.Next() {
		var event AnalyticsEvent
		var userID, sessionID sql.NullString
		var eventValue sql.NullFloat64
		var propertiesJSON, metadataJSON []byte

		err := rows.Scan(
			&event.ID, &event.OrganizationID, &userID, &sessionID,
			&event.EventType, &event.EventCategory, &event.EventAction, &event.EventLabel,
			&eventValue, &propertiesJSON, &event.UserAgent, &event.IPAddress,
			&event.Referrer, &event.Timestamp, &event.CreatedAt, &metadataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		// Handle nullable fields
		if userID.Valid {
			uid, _ := uuid.Parse(userID.String)
			event.UserID = &uid
		}
		if sessionID.Valid {
			sid, _ := uuid.Parse(sessionID.String)
			event.SessionID = &sid
		}
		if eventValue.Valid {
			event.EventValue = &eventValue.Float64
		}

		// Unmarshal JSON fields
		json.Unmarshal(propertiesJSON, &event.Properties)
		json.Unmarshal(metadataJSON, &event.Metadata)

		events = append(events, &event)
	}

	return events, nil
}

func (am *AnalyticsManager) GenerateReport(ctx context.Context, orgID uuid.UUID, startTime, endTime time.Time) (*AnalyticsReport, error) {
	report := &AnalyticsReport{
		Title:       fmt.Sprintf("Analytics Report - %s to %s", startTime.Format("2006-01-02"), endTime.Format("2006-01-02")),
		Description: "Comprehensive analytics report for the specified period",
		StartTime:   startTime,
		EndTime:     endTime,
		Metrics:     make(map[string]float64),
		Dimensions:  make(map[string][]Dimension),
		TimeSeries:  []TimeSeriesPoint{},
		TopEvents:   []EventSummary{},
		UserActivity: []UserActivitySummary{},
	}

	// Get basic metrics
	metricsQuery := `
		SELECT 
			COUNT(*) as total_events,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT session_id) as unique_sessions,
			COUNT(DISTINCT event_type) as event_types
		FROM analytics_events
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
	`

	var totalEvents, uniqueUsers, uniqueSessions, eventTypes int
	err := am.db.Conn.QueryRowContext(ctx, metricsQuery, orgID, startTime, endTime).Scan(
		&totalEvents, &uniqueUsers, &uniqueSessions, &eventTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	report.Metrics["total_events"] = float64(totalEvents)
	report.Metrics["unique_users"] = float64(uniqueUsers)
	report.Metrics["unique_sessions"] = float64(uniqueSessions)
	report.Metrics["event_types"] = float64(eventTypes)

	// Get top events
	topEventsQuery := `
		SELECT 
			event_type, 
			event_category, 
			event_action,
			COUNT(*) as count,
			COUNT(DISTINCT user_id) as unique_users,
			AVG(event_value) as avg_value
		FROM analytics_events
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
		GROUP BY event_type, event_category, event_action
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err := am.db.Conn.QueryContext(ctx, topEventsQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get top events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var summary EventSummary
		var avgValue sql.NullFloat64

		err := rows.Scan(
			&summary.EventType, &summary.EventCategory, &summary.EventAction,
			&summary.Count, &summary.UniqueUsers, &avgValue)
		if err != nil {
			continue
		}

		if avgValue.Valid {
			summary.AvgValue = avgValue.Float64
		}

		report.TopEvents = append(report.TopEvents, summary)
	}

	// Get dimensions (event categories)
	dimensionsQuery := `
		SELECT event_category, COUNT(*) as count
		FROM analytics_events
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3
		GROUP BY event_category
		ORDER BY count DESC
	`

	rows, err = am.db.Conn.QueryContext(ctx, dimensionsQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get dimensions: %w", err)
	}
	defer rows.Close()

	var categories []Dimension
	for rows.Next() {
		var dim Dimension
		err := rows.Scan(&dim.Value, &dim.Count)
		if err != nil {
			continue
		}
		dim.Name = "Event Category"
		categories = append(categories, dim)
	}
	report.Dimensions["categories"] = categories

	// Get user activity
	userActivityQuery := `
		SELECT 
			user_id,
			COUNT(*) as event_count,
			COUNT(DISTINCT session_id) as session_count,
			MAX(timestamp) as last_active
		FROM analytics_events
		WHERE organization_id = $1 AND timestamp >= $2 AND timestamp <= $3 AND user_id IS NOT NULL
		GROUP BY user_id
		ORDER BY event_count DESC
		LIMIT 20
	`

	rows, err = am.db.Conn.QueryContext(ctx, userActivityQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get user activity: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var summary UserActivitySummary
		var userID string

		err := rows.Scan(&userID, &summary.EventCount, &summary.SessionCount, &summary.LastActive)
		if err != nil {
			continue
		}

		summary.UserID, _ = uuid.Parse(userID)
		report.UserActivity = append(report.UserActivity, summary)
	}

	logging.Info("generated analytics report", map[string]interface{}{
		"organization_id": orgID,
		"start_time":      startTime,
		"end_time":        endTime,
		"total_events":    totalEvents,
		"unique_users":    uniqueUsers,
	})

	return report, nil
}

func (am *AnalyticsManager) updateAggregates(ctx context.Context, event *AnalyticsEvent) {
	// This is a simplified version - in production, you'd want to batch these updates
	// and run them periodically rather than on every event

	// Update hourly aggregates
	hourStart := event.Timestamp.Truncate(time.Hour)
	hourEnd := hourStart.Add(time.Hour)

	// Count aggregate
	am.upsertAggregate(ctx, &AnalyticsAggregate{
		OrganizationID: event.OrganizationID,
		AggregateType:  "count",
		Dimension:      "event_type",
		DimensionValue: event.EventType,
		Metric:         "events",
		Value:          1,
		PeriodType:     "hour",
		PeriodStart:    hourStart,
		PeriodEnd:      hourEnd,
	})

	// Category aggregate
	am.upsertAggregate(ctx, &AnalyticsAggregate{
		OrganizationID: event.OrganizationID,
		AggregateType:  "count",
		Dimension:      "event_category",
		DimensionValue: event.EventCategory,
		Metric:         "events",
		Value:          1,
		PeriodType:     "hour",
		PeriodStart:    hourStart,
		PeriodEnd:      hourEnd,
	})

	// Daily aggregates
	dayStart := event.Timestamp.Truncate(24 * time.Hour)
	dayEnd := dayStart.Add(24 * time.Hour)

	am.upsertAggregate(ctx, &AnalyticsAggregate{
		OrganizationID: event.OrganizationID,
		AggregateType:  "count",
		Dimension:      "event_type",
		DimensionValue: event.EventType,
		Metric:         "events",
		Value:          1,
		PeriodType:     "day",
		PeriodStart:    dayStart,
		PeriodEnd:      dayEnd,
	})
}

func (am *AnalyticsManager) upsertAggregate(ctx context.Context, agg *AnalyticsAggregate) error {
	// Try to update existing aggregate
	updateQuery := `
		UPDATE analytics_aggregates 
		SET value = value + $1, updated_at = NOW()
		WHERE organization_id = $2 AND aggregate_type = $3 AND dimension = $4 
		AND dimension_value = $5 AND metric = $6 AND period_type = $7 AND period_start = $8
	`

	result, err := am.db.Conn.ExecContext(ctx, updateQuery,
		agg.Value, agg.OrganizationID, agg.AggregateType, agg.Dimension,
		agg.DimensionValue, agg.Metric, agg.PeriodType, agg.PeriodStart)
	if err != nil {
		return fmt.Errorf("failed to update aggregate: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// If no rows were updated, insert new aggregate
	if rowsAffected == 0 {
		agg.ID = uuid.New()
		agg.CreatedAt = time.Now()
		agg.UpdatedAt = time.Now()
		agg.Metadata = make(map[string]interface{})

		insertQuery := `
			INSERT INTO analytics_aggregates (id, organization_id, aggregate_type, dimension, dimension_value, metric, value, period_type, period_start, period_end, created_at, updated_at, metadata)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

		metadataJSON, _ := json.Marshal(agg.Metadata)

		_, err = am.db.Conn.ExecContext(ctx, insertQuery,
			agg.ID, agg.OrganizationID, agg.AggregateType, agg.Dimension,
			agg.DimensionValue, agg.Metric, agg.Value, agg.PeriodType,
			agg.PeriodStart, agg.PeriodEnd, agg.CreatedAt, agg.UpdatedAt, metadataJSON)
		if err != nil {
			return fmt.Errorf("failed to insert aggregate: %w", err)
		}
	}

	return nil
}