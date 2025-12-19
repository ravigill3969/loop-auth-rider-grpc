package models

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type CreateCheckoutSessionRequest struct {
	EstimatedPrice       float32     `json:"estimated_price"`
	PickupLocation       string      `json:"pickup_location"`
	DropoffLocation      string      `json:"dropoff_location"`
	EstimatedDistanceKm  float32     `json:"estimated_distance_km"`
	EstimatedDurationMin int64       `json:"estimated_duration_min"`
	PickupCoords         Coordinates `json:"pickup_coords"`
	DropoffCoords        Coordinates `json:"dropoff_coords"`
	RiderName            string      `json:"rider_name"`
	RiderAge             int32       `json:"rider_age"`
	Gender               string      `json:"gender"`
}

type PaymentError struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	StripeCode string `json:"stripe_code,omitempty"`
}

type CreateCheckoutSessionResponse struct {
	Success         bool          `json:"success"`
	CheckoutURL     string        `json:"checkout_url,omitempty"`
	SessionID       string        `json:"session_id,omitempty"`
	PaymentIntentID string        `json:"payment_intent_id,omitempty"`
	Status          string        `json:"status"`
	Error           *PaymentError `json:"error,omitempty"`
}
