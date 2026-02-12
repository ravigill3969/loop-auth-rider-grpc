package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	trippb "ravigill/loop-grpc-trip/proto/trip_proto"

	"github.com/loop/backend/rider-auth/rest/internals/middleware"
	"github.com/loop/backend/rider-auth/rest/internals/models"
	"google.golang.org/grpc/metadata"
)

type TripClient struct {
	trip_client trippb.TripServiceClient
}

func NewTripService(tripClient trippb.TripServiceClient) *TripClient {
	return &TripClient{
		trip_client: tripClient,
	}
}

func (t *TripClient) CancelRide(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", "Only POST method is accepted")
		return
	}

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

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)

	rider_id, ok := r.Context().Value(middleware.RiderIDKey).(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", "Please login to perform this action.")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body", err.Error())
		return
	}
	defer r.Body.Close()

	var req models.CancelRideS
	if err := json.Unmarshal(body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	fmt.Println("driver_id", req.DriverID)

	resp, err := t.trip_client.RiderCancelTripHandler(ctx, &trippb.RiderCancelTripHandlerRequest{
		TripId:   req.TripID,
		DriverId: req.DriverID,
		Reason:   req.Reason,
		RiderId:  rider_id,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to cancel ride", err.Error())
	}

	if !resp.Status {
		respondWithError(w, http.StatusBadRequest, "Failed to cancel ride", err.Error())
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (t *TripClient) RetrieveActiveRide(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {

		cookie, err := r.Cookie("access_token")
		if err == nil {
			authHeader = cookie.Value
		}
	}

	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token", "Authorization header or cookie is required")
		return
	}

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)

	riderID, ok := middleware.GetRiderIDFromContext(r.Context())

	if ok != nil {
		respondWithError(w, http.StatusUnauthorized, "Please login or register", "Unauthorized")
		return
	}

	resp, err := t.trip_client.RetrieveTripBasedOnStatusForRider(ctx, &trippb.RetrieveTripBasedOnStatusForRiderRequest{
		RiderId: riderID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), "Get active ride error")
		return
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (t *TripClient) RetrieveActiveRideWithID(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {

		cookie, err := r.Cookie("access_token")
		if err == nil {
			authHeader = cookie.Value
		}
	}

	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token", "Authorization header or cookie is required")
		return
	}

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)

	_, ok := middleware.GetRiderIDFromContext(r.Context())

	if ok != nil {
		respondWithError(w, http.StatusUnauthorized, "Please login or register", "Unauthorized")
		return
	}

	tripID := r.URL.Query().Get("tid")

	if tripID == "" {
		respondWithError(w, http.StatusBadRequest, "What are you looking for dawg?", "tid is necessary!")
		return
	}

	resp, err := t.trip_client.GetTripWithTripID(ctx, &trippb.GetTripWithTripIDRequest{
		TripId: tripID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), "Get active rider error")
		return
	}

	respondWithJSON(w, http.StatusOK, resp)

}
