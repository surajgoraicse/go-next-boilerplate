package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	db_sqlc "github.com/surajgoraicse/go-next-boilerplate/internal/db/sqlc"
)

// AuthTokens holds both the access and refresh tokens returned by auth operations.
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

// jwt token payload
type TokenPayload struct {
	UserID          uuid.UUID        `json:"user_id"`
	Email           string           `json:"email"`
	Name            string           `json:"name"`
	Role            db_sqlc.UserRole `json:"role"`
	IsEmailVerified bool             `json:"is_email_verified"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a JWT token with the given payload and expiry time
func GenerateAccessToken(payload TokenPayload, jwtExpiryTime time.Duration, jwtIssuer string) (string, error) {

	payload.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    jwtIssuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpiryTime)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// GenerateSecureToken creates a 32-byte URL-safe random token.
// Use cases : generating refresh token,
func GenerateSecureToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(h[:])
	return raw, hash, nil
}

// HashToken computes the SHA-256 hex digest of a raw token.
func HashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
