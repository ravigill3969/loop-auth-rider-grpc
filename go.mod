module github.com/loop/backend/rider-auth/rest

go 1.25.1

require ravigill/rider-grpc-server v0.0.0

replace ravigill/rider-grpc-server => ../auth-grpc

require ravigill/loop-auth-utils v0.0.0

replace ravigill/loop-auth-utils => ../../libs/authutils

require ravigill/loop-grpc-trip v0.0.0

replace ravigill/loop-grpc-trip => ../../trip/grpc-trip

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.17.2
	google.golang.org/grpc v1.77.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
