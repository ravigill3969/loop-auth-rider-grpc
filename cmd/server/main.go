package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	pb "ravigill/rider-grpc-server/proto"

	"github.com/loop/backend/rider-auth/rest/internals/configs"
	"github.com/loop/backend/rider-auth/rest/internals/handlers"
	"github.com/loop/backend/rider-auth/rest/internals/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HTTPServer struct {
	mux           *http.ServeMux
	authClient    pb.AuthServiceClient
	paymentClient pb.PaymentServiceClient
}

func NewHTTPServer(authClient pb.AuthServiceClient, paymentClient pb.PaymentServiceClient) *HTTPServer {
	return &HTTPServer{
		mux:           http.NewServeMux(),
		authClient:    authClient,
		paymentClient: paymentClient,
	}
}

func (s *HTTPServer) Start(port string) {

	secretKey := os.Getenv("ACCESS_TOKEN_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("ACCESS_TOKEN_SECRET_KEY environment variable is required")
	}

	authHandler := handlers.NewAuthService(s.authClient)
	authRoutes := routes.NewAuthRoutes(s.mux, authHandler)
	authRoutes.Register()

	paymentHandler := handlers.NewPaymentService(s.paymentClient)
	paymentRoutes := routes.NewPaymentRoutes(s.mux, paymentHandler, secretKey)
	paymentRoutes.Register()

	fmt.Println("Server is running on PORT" + " " + port)

	err := http.ListenAndServe(""+port, corsMiddleware(s.mux))

	fmt.Println(err)
}

func main() {
	err := configs.LoadEnv()

	if err != nil {
		fmt.Println("FNot loading")
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	authConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connect to Auth gRPC:", err)
	}
	authClient := pb.NewAuthServiceClient(authConn)

	paymentConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connect to Payment gRPC:", err)
	}
	paymentClient := pb.NewPaymentServiceClient(paymentConn)

	httpServer := NewHTTPServer(authClient, paymentClient)
	httpServer.Start(port)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
