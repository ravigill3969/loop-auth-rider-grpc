package models

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	BirthMonth  string `json:"birth_month"`
	BirthYear   int64  `json:"birth_year"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents the user information in responses
type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	BirthMonth  string `json:"birth_month"`
	BirthYear   int64  `json:"birth_year"`
	UpdatedAt   int64  `json:"updated_at"`
	CreatedAt   int64  `json:"created_at"`
}

// Tokens represents authentication tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
}

// AuthResponse represents the response for registration
type AuthResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Status  int64   `json:"status"`
	User    *User   `json:"user,omitempty"`
	Token   *Tokens `json:"token,omitempty"`
}

// LoginResponse represents the response for login
type LoginResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Status  int64   `json:"status"`
	User    *User   `json:"user,omitempty"`
	Token   *Tokens `json:"token,omitempty"`
}

// GetRiderDetailsResponse represents the response for getting rider details
type GetRiderDetailsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int64  `json:"status"`
	User    *User  `json:"user,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  int64  `json:"status"`
	Error   string `json:"error,omitempty"`
}
