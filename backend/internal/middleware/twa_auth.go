package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type contextKey string

const (
	TWAUserContextKey contextKey = "twa_user"
)

type TWAUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
	IsPremium    bool   `json:"is_premium,omitempty"`
}

type TWAInitData struct {
	QueryID      string  `json:"query_id,omitempty"`
	User         TWAUser `json:"user"`
	AuthDate     int64   `json:"auth_date"`
	Hash         string  `json:"hash"`
	StartParam   string  `json:"start_param,omitempty"`
}

var (
	ErrInvalidInitData = errors.New("invalid init data")
	ErrExpiredInitData = errors.New("expired init data")
	ErrMissingInitData = errors.New("missing init data")
)

// ValidateTWAInitData validates Telegram WebApp initData
func ValidateTWAInitData(initData string, botToken string, maxAge time.Duration) (*TWAInitData, error) {
	if initData == "" {
		return nil, ErrMissingInitData
	}

	// Parse the init data
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrInvalidInitData
	}

	// Extract hash for verification
	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, ErrInvalidInitData
	}

	// Create data check string (sorted key=value pairs, excluding hash)
	var dataCheckParts []string
	for key, vals := range values {
		if key == "hash" {
			continue
		}
		for _, val := range vals {
			dataCheckParts = append(dataCheckParts, key+"="+val)
		}
	}
	sort.Strings(dataCheckParts)
	dataCheckString := strings.Join(dataCheckParts, "\n")

	// Calculate HMAC-SHA256
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))

	dataHash := hmac.New(sha256.New, secretKey.Sum(nil))
	dataHash.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(dataHash.Sum(nil))

	// Verify hash
	if !hmac.Equal([]byte(receivedHash), []byte(calculatedHash)) {
		return nil, ErrInvalidInitData
	}

	// Parse user data
	var result TWAInitData
	userJSON := values.Get("user")
	if userJSON != "" {
		if err := json.Unmarshal([]byte(userJSON), &result.User); err != nil {
			return nil, ErrInvalidInitData
		}
	}

	// Parse auth_date and validate expiry
	authDateStr := values.Get("auth_date")
	if authDateStr != "" {
		var authDate int64
		if _, err := json.Number(authDateStr).Int64(); err == nil {
			result.AuthDate = authDate
		}
	}

	if maxAge > 0 && result.AuthDate > 0 {
		authTime := time.Unix(result.AuthDate, 0)
		if time.Since(authTime) > maxAge {
			return nil, ErrExpiredInitData
		}
	}

	result.Hash = receivedHash
	result.QueryID = values.Get("query_id")
	result.StartParam = values.Get("start_param")

	return &result, nil
}

// TWAAuthMiddleware creates a middleware that validates TWA initData
func TWAAuthMiddleware(botToken string, maxAge time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get init data from Authorization header
			authHeader := r.Header.Get("Authorization")
			initData := strings.TrimPrefix(authHeader, "tma ")

			data, err := ValidateTWAInitData(initData, botToken, maxAge)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), TWAUserContextKey, &data.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetTWAUser extracts TWA user from context
func GetTWAUser(ctx context.Context) *TWAUser {
	if user, ok := ctx.Value(TWAUserContextKey).(*TWAUser); ok {
		return user
	}
	return nil
}
