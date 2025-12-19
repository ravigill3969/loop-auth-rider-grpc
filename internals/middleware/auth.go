package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	jwtlib "github.com/loop/backend/rider-auth/lib/jwt"
)

type contextKey string

const (
	RiderIDKey contextKey = "riderId"
	EmailKey   contextKey = "email"
)

func JWTVerifyMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header or cookie
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				cookie, err := r.Cookie("access_token")
				if err == nil {
					authHeader = cookie.Value
				}
			}

			if authHeader == "" {
				http.Error(w, "Missing authorization token", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
			if token == "" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			claims, err := jwtlib.VerifyToken(token, secretKey)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid or expired token: %v", err), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), RiderIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRiderIDFromContext extracts rider ID from request context
func GetRiderIDFromContext(ctx context.Context) (string, error) {
	riderID, ok := ctx.Value(RiderIDKey).(string)
	if !ok || riderID == "" {
		return "", fmt.Errorf("rider ID not found in context")
	}
	return riderID, nil
}

// GetEmailFromContext extracts email from request context
func GetEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(EmailKey).(string)
	if !ok || email == "" {
		return "", fmt.Errorf("email not found in context")
	}
	return email, nil
}
