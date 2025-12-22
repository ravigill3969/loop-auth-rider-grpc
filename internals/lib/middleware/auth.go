package middleware

import (
	"context"
	"fmt"
	"strings"

	jwtlib "github.com/loop/backend/rider-auth/lib/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	UserIDKey contextKey = "userId"
	EmailKey  contextKey = "email"
)

func AuthInterceptor(secretKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/rider_auth.AuthService/Login" ||
			info.FullMethod == "/rider_auth.AuthService/Register" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		// Extract token from "Bearer <token>" format
		token := strings.TrimSpace(strings.TrimPrefix(authHeader[0], "Bearer"))
		if token == "" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format")
		}

		// Verify token
		claims, err := jwtlib.VerifyToken(token, secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token: %v", err)
		}

		// Add claims to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		fmt.Printf("Authenticated user: %s (ID: %s)\n", claims.Email, claims.UserID)

		return handler(ctx, req)
	}
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func GetEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(EmailKey).(string)
	if !ok || email == "" {
		return "", fmt.Errorf("email not found in context")
	}
	return email, nil
}
