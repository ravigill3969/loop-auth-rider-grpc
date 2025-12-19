package routes

import (
	"net/http"

	"github.com/loop/backend/rider-auth/rest/internals/handlers"
	"github.com/loop/backend/rider-auth/rest/internals/middleware"
)

type PaymentRoutes struct {
	mux       *http.ServeMux
	handler   *handlers.PaymentService
	secretKey string
}

func NewPaymentRoutes(mux *http.ServeMux, handler *handlers.PaymentService, secretKey string) *PaymentRoutes {
	return &PaymentRoutes{
		mux:       mux,
		handler:   handler,
		secretKey: secretKey,
	}
}

func (r *PaymentRoutes) Register() {
	// Wrap handler with JWT middleware
	jwtMiddleware := middleware.JWTVerifyMiddleware(r.secretKey)
	r.mux.Handle("/api/payment/create-checkout-session", jwtMiddleware(http.HandlerFunc(r.handler.CreateCheckoutSessionHandler)))
}
