package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	pb "ravigill/rider-grpc-server/proto"

	"github.com/loop/backend/rider-auth/rest/internals/middleware"
	"github.com/loop/backend/rider-auth/rest/internals/models"
	"google.golang.org/grpc/metadata"
)

type PaymentService struct {
	paymentClient pb.PaymentServiceClient
}

func NewPaymentService(paymentClient pb.PaymentServiceClient) *PaymentService {
	return &PaymentService{
		paymentClient: paymentClient,
	}
}

func (p *PaymentService) CreateCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is accepted")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		cookie, err := r.Cookie("access_token")
		if err == nil {
			authHeader = cookie.Value
		}
	}

	rider_id, ok := r.Context().Value(middleware.RiderIDKey).(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Please login to perform this action.")
	}

	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token", "Authorization header or cookie is required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}
	defer r.Body.Close()

	var req models.CreateCheckoutSessionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	if req.EstimatedPrice <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid estimated_price", "Must be greater than 0")
		return
	}

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)

	grpcReq := &pb.CreateCheckOutSessionRequest{
		RiderId:              rider_id,
		RiderName:            req.RiderName,
		RiderAge:             req.RiderAge,
		Gender:               req.Gender,
		EstimatedPrice:       req.EstimatedPrice,
		PickupLocation:       req.PickupLocation,
		DropoffLocation:      req.DropoffLocation,
		EstimatedDistanceKm:  req.EstimatedDistanceKm,
		EstimatedDurationMin: req.EstimatedDurationMin,
		PickupCoordsLatLng: &pb.Coordinates{
			Lat: req.PickupCoords.Lat,
			Lng: req.PickupCoords.Lng,
		},
		DropoffCoordsLatLng: &pb.Coordinates{
			Lat: req.DropoffCoords.Lat,
			Lng: req.DropoffCoords.Lng,
		},
	}

	grpcResp, err := p.paymentClient.CreateCheckOutSession(ctx, grpcReq)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create checkout session", err.Error())
		return
	}

	resp := models.CreateCheckoutSessionResponse{
		Success:         grpcResp.Success,
		CheckoutURL:     grpcResp.CheckoutUrl,
		SessionID:       grpcResp.SessionId,
		PaymentIntentID: grpcResp.PaymentIntentId,
		Status:          grpcResp.Status,
	}

	if grpcResp.Error != nil {
		resp.Error = &models.PaymentError{
			Code:       grpcResp.Error.Code,
			Message:    grpcResp.Error.Message,
			StripeCode: grpcResp.Error.StripeCode,
		}
	}

	statusCode := http.StatusOK
	if !grpcResp.Success {
		statusCode = http.StatusBadRequest
	}

	respondWithJSON(w, statusCode, resp)
}
