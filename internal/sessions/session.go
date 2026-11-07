package sessions

import (
	"campbell-checkbook/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Session represents a user session, stored as a JWT.
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
}

// GenerateToken creates a signed JWT for the given user ID.
func GenerateToken(cfg *config.Config, userID string) (string, error) {
	// Use HS256 signing method.
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Secret must be provided.
	secret := cfg.JWTSecret
	if secret == "" {
		return "", jwt.ErrTokenMalformed // use generic error
	}
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ValidateToken parses and validates a JWT, returning the user ID if valid.
func ValidateToken(cfg *config.Config, tokenStr string) (string, error) {
	secret := cfg.JWTSecret
	if secret == "" {
		return "", jwt.ErrTokenMalformed
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenInvalidClaims
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
	}
	return "", jwt.ErrTokenInvalidClaims
}
