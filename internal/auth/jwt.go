package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// JWTClaims holds the decoded JWT claims
type JWTClaims struct {
	Exp int64 `json:"exp"`
	Iat int64 `json:"iat"`
	Sub string `json:"sub"`
}

// ParseJWTExpiration extracts the expiration time from a JWT token without verification
func ParseJWTExpiration(token string) (time.Time, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, nil
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return time.Time{}, err
	}

	if claims.Exp == 0 {
		return time.Time{}, nil
	}

	return time.Unix(claims.Exp, 0), nil
}

// IsTokenExpired checks if a token is expired or will expire within the given margin
func IsTokenExpired(token string, margin time.Duration) bool {
	exp, err := ParseJWTExpiration(token)
	if err != nil || exp.IsZero() {
		return true
	}

	return time.Now().Add(margin).After(exp)
}
