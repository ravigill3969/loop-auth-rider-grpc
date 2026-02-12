package routes

import (
	"net/http"

	"github.com/loop/backend/rider-auth/rest/internals/handlers"
	"github.com/loop/backend/rider-auth/rest/internals/middleware"
)

type TripRoutes struct {
	mux       *http.ServeMux
	handler   *handlers.TripClient
	secretkey string
}

func NewTripRoutes(mux *http.ServeMux, handler *handlers.TripClient, secretkey string) *TripRoutes {
	return &TripRoutes{
		mux:       mux,
		handler:   handler,
		secretkey: secretkey,
	}
}

func (r *TripRoutes) Register() {

	jwtMiddleware := middleware.JWTVerifyMiddleware(r.secretkey)
	r.mux.Handle("POST /api/trip/cancel-ride", jwtMiddleware(http.HandlerFunc(r.handler.CancelRide)))
	r.mux.Handle("GET /api/trip/get-active-ride-id", jwtMiddleware(http.HandlerFunc(r.handler.RetrieveActiveRide)))
	r.mux.Handle("GET /api/trip/get-active-ride/", jwtMiddleware(http.HandlerFunc(r.handler.RetrieveActiveRideWithID)))
}
