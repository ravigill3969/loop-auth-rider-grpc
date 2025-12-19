package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	pb "ravigill/rider-grpc-server/proto"

	"github.com/loop/backend/rider-auth/rest/internals/models"
	"google.golang.org/grpc/metadata"
)

type AuthService struct {
	authClient pb.AuthServiceClient
}

func NewAuthService(authClient pb.AuthServiceClient) *AuthService {

	return &AuthService{
		authClient: authClient,
	}
}

// RegisterHandler handles user registration
func (a *AuthService) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is accepted")
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}
	defer r.Body.Close()

	var req models.RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields", "email, password, and full_name are required")
		return
	}

	// Call gRPC service
	grpcReq := &pb.RegisterRequest{
		User: &pb.User{
			Email:       req.Email,
			FullName:    req.FullName,
			Password:    req.Password,
			PhoneNumber: req.PhoneNumber,
			BirthMonth:  req.BirthMonth,
			BirthYear:   req.BirthYear,
		},
	}

	grpcResp, err := a.authClient.Register(context.Background(), grpcReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to register user", err.Error())
		return
	}

	if !grpcResp.Success {
		respondWithError(w, http.StatusCreated, grpcResp.Message, grpcResp.Message)
		return
	}

	// Convert gRPC response to REST response
	resp := models.AuthResponse{
		Success: grpcResp.Success,
		Message: grpcResp.Message,
		Status:  grpcResp.Status,
	}

	if grpcResp.User != nil {
		resp.User = &models.User{
			ID:          grpcResp.User.Id,
			Email:       grpcResp.User.Email,
			FullName:    grpcResp.User.FullName,
			PhoneNumber: grpcResp.User.PhoneNumber,
			BirthMonth:  grpcResp.User.BirthMonth,
			BirthYear:   grpcResp.User.BirthYear,
			UpdatedAt:   grpcResp.User.UpdatedAt,
			CreatedAt:   grpcResp.User.CreatedAt,
		}
	}

	access_cookie := http.Cookie{
		Name:     "access_token",
		Value:    grpcResp.Token.TokenType + " " + grpcResp.Token.AccessToken,
		Expires:  time.Now().Add(24 * time.Hour * 3),
		MaxAge:   60 * 60 * 24 * 3,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	refresh_cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    grpcResp.Token.TokenType + " " + grpcResp.Token.RefreshToken,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		MaxAge:   60 * 60 * 24 * 7,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)

	respondWithJSON(w, int(grpcResp.Status), resp)
}

func (a *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is accepted")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}
	defer r.Body.Close()

	var req models.LoginRequest
	if err := json.Unmarshal(body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields", "email and password are required")
		return
	}

	grpcReq := &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	grpcResp, err := a.authClient.Login(context.Background(), grpcReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to login", err.Error())
		return
	}

	if !grpcResp.Success {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour * 24 * 3),
			MaxAge:   -60 * 60 * 24 * 3,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour * 24 * 7),
			MaxAge:   -60 * 60 * 24 * 7,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		respondWithError(w, int(grpcResp.Status), grpcResp.Message, grpcResp.Message)
		return
	}

	resp := models.LoginResponse{
		Success: grpcResp.Success,
		Message: grpcResp.Message,
		Status:  grpcResp.Status,
	}

	if grpcResp.User != nil {
		resp.User = &models.User{
			ID:          grpcResp.User.Id,
			Email:       grpcResp.User.Email,
			FullName:    grpcResp.User.FullName,
			PhoneNumber: grpcResp.User.PhoneNumber,
			BirthMonth:  grpcResp.User.BirthMonth,
			BirthYear:   grpcResp.User.BirthYear,
			UpdatedAt:   grpcResp.User.UpdatedAt,
			CreatedAt:   grpcResp.User.CreatedAt,
		}
	}

	access_cookie := http.Cookie{
		Name:     "access_token",
		Value:    grpcResp.Token.TokenType + " " + grpcResp.Token.AccessToken,
		Expires:  time.Now().Add(24 * time.Hour * 3),
		MaxAge:   60 * 60 * 24 * 3, // 3 days in seconds
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	refresh_cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    grpcResp.Token.TokenType + " " + grpcResp.Token.RefreshToken,
		Expires:  time.Now().Add(24 * time.Hour * 7),
		MaxAge:   60 * 60 * 24 * 7, // 7 days in seconds
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, &access_cookie)
	http.SetCookie(w, &refresh_cookie)

	respondWithJSON(w, int(grpcResp.Status), resp)
}

func (a *AuthService) GetRiderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only GET method is accepted")
		return
	}

	// Extract token from Authorization header or Cookie
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Try to get from cookie
		cookie, err := r.Cookie("access_token")
		if err == nil {
			authHeader = cookie.Value
		}
	}

	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token", "Authorization header or cookie is required")
		return
	}

	// Create context with authorization metadata for gRPC
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)

	// Call gRPC service (ID will come from the token context)
	grpcReq := &pb.GetRiderDetailsRequest{
		Id: "", // Not used anymore, comes from auth context
	}

	grpcResp, err := a.authClient.GetRiderDetails(ctx, grpcReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get rider details", err.Error())
		return
	}

	// Convert gRPC response to REST response
	resp := models.GetRiderDetailsResponse{
		Success: grpcResp.Success,
		Message: grpcResp.Message,
		Status:  grpcResp.Status,
	}

	if grpcResp.User != nil {
		resp.User = &models.User{
			ID:          grpcResp.User.Id,
			Email:       grpcResp.User.Email,
			FullName:    grpcResp.User.FullName,
			PhoneNumber: grpcResp.User.PhoneNumber,
			BirthMonth:  grpcResp.User.BirthMonth,
			BirthYear:   grpcResp.User.BirthYear,
			UpdatedAt:   grpcResp.User.UpdatedAt,
			CreatedAt:   grpcResp.User.CreatedAt,
		}
	}

	respondWithJSON(w, int(grpcResp.Status), resp)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string, errorDetail string) {
	errResp := models.ErrorResponse{
		Success: false,
		Message: message,
		Status:  int64(statusCode),
		Error:   errorDetail,
	}
	respondWithJSON(w, statusCode, errResp)
}
