package auth

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"pulse-backend/internal/models"
)

var (
	// In production this can be loaded from env JWT_SECRET
	JWTSecret = []byte("pulse-secret-key-32-chars-minimum-token!")

	ErrInvalidToken = errors.New("invalid or expired token")
	ErrRevokedToken = errors.New("token has been revoked/logged out")

	// In-memory sync.Map caching revoked token signatures for zero-latency O(1) checks
	revokedTokensMap = sync.Map{}

	// Single Device / Single Active Session registry: userID -> sessionID
	userActiveSessions = sync.Map{}

	// User Activity registry: userID -> unix timestamp (seconds)
	userLastActive = sync.Map{}
)

type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

// InitRevocationStore preloads unexpired revoked tokens and active sessions into memory
func InitRevocationStore(db *gorm.DB) {
	if db == nil {
		return
	}
	var tokens []models.RevokedToken
	now := time.Now()
	if err := db.Where("expires_at > ?", now).Find(&tokens).Error; err != nil {
		log.Printf("[Auth] Warning: Could not preload revoked tokens: %v", err)
		return
	}

	loaded := 0
	for _, t := range tokens {
		revokedTokensMap.Store(t.Token, t.ExpiresAt.Unix())
		loaded++
	}

	// Also preload active sessions from users
	var users []models.User
	if err := db.Select("id", "current_session_id", "last_active_at").Where("current_session_id != ''").Find(&users).Error; err == nil {
		for _, u := range users {
			userActiveSessions.Store(u.ID, u.CurrentSessionID)
			if !u.LastActiveAt.IsZero() {
				userLastActive.Store(u.ID, u.LastActiveAt.Unix())
			}
		}
	}

	log.Printf("[Auth] Initialized in-memory token revocation cache with %d active revoked sessions", loaded)
}

// Single Active Session helpers
func SetActiveSession(userID string, sessionID string) {
	userActiveSessions.Store(userID, sessionID)
	userLastActive.Store(userID, time.Now().Unix())
}

func GetActiveSession(userID string) (string, bool) {
	if val, ok := userActiveSessions.Load(userID); ok {
		if sID, okStr := val.(string); okStr {
			return sID, true
		}
	}
	return "", false
}

func ClearActiveSession(userID string) {
	userActiveSessions.Delete(userID)
	userLastActive.Delete(userID)
}

func UpdateLastActive(userID string) {
	userLastActive.Store(userID, time.Now().Unix())
}

func GetLastActive(userID string) time.Time {
	if val, ok := userLastActive.Load(userID); ok {
		if u, okInt := val.(int64); okInt {
			return time.Unix(u, 0)
		}
	}
	return time.Time{}
}

// RevokeToken registers a token as revoked in both SQLite and memory cache
func RevokeToken(db *gorm.DB, tokenString string, userID string, expiresAt time.Time) error {
	if tokenString == "" {
		return nil
	}

	// Update in-memory cache immediately
	revokedTokensMap.Store(tokenString, expiresAt.Unix())

	if db != nil {
		rev := models.RevokedToken{
			Token:     tokenString,
			UserID:    userID,
			RevokedAt: time.Now(),
			ExpiresAt: expiresAt,
		}
		// FirstOrCreate or Create to prevent unique index collision
		if err := db.Where(models.RevokedToken{Token: tokenString}).FirstOrCreate(&rev).Error; err != nil {
			log.Printf("[Auth] Error persisting revoked token to DB: %v", err)
			return err
		}
	}

	return nil
}

// IsTokenRevoked verifies whether a token is flagged as revoked
func IsTokenRevoked(tokenString string) bool {
	if val, ok := revokedTokensMap.Load(tokenString); ok {
		expUnix, isInt := val.(int64)
		if isInt && time.Now().Unix() > expUnix {
			revokedTokensMap.Delete(tokenString)
			return false
		}
		return true
	}
	return false
}

// GenerateToken generates a signed JWT token with user_id, role, and unique session_id
func GenerateToken(userID string, email string, role string, sessionID string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pulse-backend",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken parses and validates the token string and checks revocation state
func ValidateToken(tokenString string) (*Claims, error) {
	// 1. Fast revocation check
	if IsTokenRevoked(tokenString) {
		return nil, ErrRevokedToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
