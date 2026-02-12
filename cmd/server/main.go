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
	"github.com/loop/backend/rider-auth/rest/internals/websocket"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	trippb "ravigill/loop-grpc-trip/proto/trip_proto"
)

type HTTPServer struct {
	mux           *http.ServeMux
	authClient    pb.AuthServiceClient
	paymentClient pb.PaymentServiceClient
	tripClient    trippb.TripServiceClient
}

func NewHTTPServer(authClient pb.AuthServiceClient, paymentClient pb.PaymentServiceClient, tripClient trippb.TripServiceClient) *HTTPServer {
	return &HTTPServer{
		mux:           http.NewServeMux(),
		authClient:    authClient,
		paymentClient: paymentClient,
		tripClient:    tripClient,
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

	tripHandler := handlers.NewTripService(s.tripClient)
	tripRoutes := routes.NewTripRoutes(s.mux, tripHandler, secretKey)
	tripRoutes.Register()

	redis_url := os.Getenv("REDIS_URL")

	if redis_url == "" {
		log.Fatal("REDIS_URL not set")
	}

	redis_client, err := newRedis(redis_url)

	if err != nil {
		log.Fatal("redis conn error", err)
	}

	webSocket := websocket.NewWSService(redis_client)

	s.mux.HandleFunc("/ws", webSocket.WSHandler)

	fmt.Println("Server is running on PORT" + " " + port)

	err = http.ListenAndServe(""+port, corsMiddleware(s.mux))

	fmt.Println(err)
}

func main() {
	err := configs.LoadEnv()

	if err != nil {
		fmt.Println("env file not loading")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	//rider auth
	RIDER_AUTH_GRPC_SERVER_URL := os.Getenv("RIDER_GRPC_AUTH_SERVER")
	fmt.Println(RIDER_AUTH_GRPC_SERVER_URL)

	if RIDER_AUTH_GRPC_SERVER_URL == "" {
		fmt.Println("rider grpc server link is missing")
	}

	authConn, err := grpc.NewClient(RIDER_AUTH_GRPC_SERVER_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connect to Auth gRPC:", err)
	}
	authClient := pb.NewAuthServiceClient(authConn)

	//payment
	PAYMENT_GRPC_SERVER_URL := os.Getenv("PAYMENT_GRPC_SERVER")

	fmt.Println(PAYMENT_GRPC_SERVER_URL)

	if PAYMENT_GRPC_SERVER_URL == "" {
		fmt.Println("payment grpc server link is missing")
	}

	paymentConn, err := grpc.NewClient(PAYMENT_GRPC_SERVER_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connect to Payment gRPC:", err)
	}
	paymentClient := pb.NewPaymentServiceClient(paymentConn)

	//trip
	TRIP_GRPC_SERVER_URL := os.Getenv("TRIP_GRPC_SERVER")
	fmt.Println(TRIP_GRPC_SERVER_URL)

	if TRIP_GRPC_SERVER_URL == "" {
		fmt.Println("payment grpc server link is missing")
	}
	tripConn, err := grpc.NewClient(TRIP_GRPC_SERVER_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connect to Trip gRPC:", err)
	}
	tripClient := trippb.NewTripServiceClient(tripConn)

	httpServer := NewHTTPServer(authClient, paymentClient, tripClient)
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

func newRedis(url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return redis.NewClient(opt), nil
}
