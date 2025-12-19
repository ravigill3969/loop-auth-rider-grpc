package routes

import (
	"net/http"

	"github.com/loop/backend/rider-auth/rest/internals/handlers"
)

type AuthRoutes struct {
	mux     *http.ServeMux
	handler *handlers.AuthService
}

func NewAuthRoutes(mux *http.ServeMux, handler *handlers.AuthService) *AuthRoutes {
	return &AuthRoutes{
		mux:     mux,
		handler: handler,
	}
}

func (r *AuthRoutes) Register() {

	r.mux.HandleFunc("/api/auth/register", r.handler.RegisterHandler)
	r.mux.HandleFunc("/api/auth/login", r.handler.LoginHandler)
	r.mux.HandleFunc("/api/auth/rider", r.handler.GetRiderDetailsHandler)
}
