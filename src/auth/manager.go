package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"holodeck1/database"
	"holodeck1/logging"
)

type Manager struct {
	db        *database.DB
	jwtSecret []byte
}

type User struct {
	ID        uuid.UUID              `json:"id"`
	Email     string                 `json:"email"`
	Username  string                 `json:"username"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	LastLogin *time.Time             `json:"last_login"`
	Status    string                 `json:"status"`
	Profile   map[string]interface{} `json:"profile"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type Token struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	TokenType    string     `json:"token_type"`
	ExpiresAt    time.Time  `json:"expires_at"`
	Scope        string     `json:"scope"`
	CreatedAt    time.Time  `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string                 `json:"email"`
	Username string                 `json:"username"`
	Password string                 `json:"password"`
	Profile  map[string]interface{} `json:"profile"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func NewManager(db *database.DB, jwtSecret []byte) *Manager {
	return &Manager{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (m *Manager) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// Check if user already exists
	var exists bool
	err := m.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)", req.Email, req.Username).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "active",
		Profile:   req.Profile,
		Metadata:  make(map[string]interface{}),
	}

	if user.Profile == nil {
		user.Profile = make(map[string]interface{})
	}

	profileJSON, _ := json.Marshal(user.Profile)
	metadataJSON, _ := json.Marshal(user.Metadata)

	query := `
		INSERT INTO users (id, email, username, password_hash, status, profile, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`

	err = m.db.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.Username,
		string(hashedPassword),
		user.Status,
		profileJSON,
		metadataJSON,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		logging.Error("failed to create user", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logging.Info("user registered", map[string]interface{}{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
	})

	return user, nil
}

func (m *Manager) Login(ctx context.Context, req *LoginRequest) (*User, *Token, error) {
	user := &User{}
	var passwordHash string
	var profileJSON, metadataJSON []byte

	query := `
		SELECT id, email, username, password_hash, created_at, updated_at, last_login, status, profile, metadata
		FROM users WHERE email = $1 AND status = 'active'
	`

	err := m.db.QueryRowContext(ctx, query, req.Email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&passwordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.Status,
		&profileJSON,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("invalid credentials")
		}
		logging.Error("failed to get user for login", map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		return nil, nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		logging.Warn("invalid login attempt", map[string]interface{}{
			"email": req.Email,
		})
		return nil, nil, fmt.Errorf("invalid credentials")
	}

	json.Unmarshal(profileJSON, &user.Profile)
	json.Unmarshal(metadataJSON, &user.Metadata)

	// Create tokens
	token, err := m.createTokens(ctx, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	// Update last login
	_, err = m.db.ExecContext(ctx, "UPDATE users SET last_login = NOW() WHERE id = $1", user.ID)
	if err != nil {
		logging.Error("failed to update last login", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		// Don't fail the login for this
	}

	logging.Info("user logged in", map[string]interface{}{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
	})

	return user, token, nil
}

func (m *Manager) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*Token, error) {
	// Get existing token
	var tokenID, userID uuid.UUID
	var expiresAt time.Time

	query := `
		SELECT id, user_id, expires_at 
		FROM oauth_tokens 
		WHERE refresh_token = $1
	`

	err := m.db.QueryRowContext(ctx, query, req.RefreshToken).Scan(&tokenID, &userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Create new tokens
	token, err := m.createTokens(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new tokens: %w", err)
	}

	// Delete old token
	_, err = m.db.ExecContext(ctx, "DELETE FROM oauth_tokens WHERE id = $1", tokenID)
	if err != nil {
		logging.Error("failed to delete old token", map[string]interface{}{
			"token_id": tokenID,
			"error":    err.Error(),
		})
		// Don't fail the refresh for this
	}

	logging.Info("token refreshed", map[string]interface{}{
		"user_id": userID,
	})

	return token, nil
}

func (m *Manager) Logout(ctx context.Context, accessToken string) error {
	query := `DELETE FROM oauth_tokens WHERE access_token = $1`

	result, err := m.db.ExecContext(ctx, query, accessToken)
	if err != nil {
		logging.Error("failed to logout user", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to logout user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("token not found")
	}

	logging.Info("user logged out", nil)
	return nil
}

func (m *Manager) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	var profileJSON, metadataJSON []byte

	query := `
		SELECT id, email, username, created_at, updated_at, last_login, status, profile, metadata
		FROM users WHERE id = $1
	`

	err := m.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
		&user.Status,
		&profileJSON,
		&metadataJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logging.Error("failed to get user by ID", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	json.Unmarshal(profileJSON, &user.Profile)
	json.Unmarshal(metadataJSON, &user.Metadata)

	return user, nil
}

func (m *Manager) ValidateToken(ctx context.Context, tokenString string) (*User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if token exists in database
	var exists bool
	err = m.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM oauth_tokens WHERE access_token = $1)", tokenString).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("token not found")
	}

	// Get user
	user, err := m.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (m *Manager) createTokens(ctx context.Context, userID uuid.UUID) (*Token, error) {
	user, err := m.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate access token
	accessTokenClaims := &JWTClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "hd1-platform",
			Subject:   user.ID.String(),
		},
	}

	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err := accessTokenJWT.SignedString(m.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	token := &Token{
		ID:           uuid.New(),
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // Refresh token expires in 7 days
		Scope:        "read write",
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO oauth_tokens (id, user_id, access_token, refresh_token, token_type, expires_at, scope)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`

	err = m.db.QueryRowContext(ctx, query,
		token.ID,
		token.UserID,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.ExpiresAt,
		token.Scope,
	).Scan(&token.CreatedAt)

	if err != nil {
		logging.Error("failed to create token", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}

func (m *Manager) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM oauth_tokens WHERE expires_at < NOW()`

	result, err := m.db.ExecContext(ctx, query)
	if err != nil {
		logging.Error("failed to cleanup expired tokens", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		logging.Info("cleaned up expired tokens", map[string]interface{}{
			"count": rowsAffected,
		})
	}

	return nil
}

func (m *Manager) StartTokenCleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CleanupExpiredTokens(ctx)
		}
	}
}