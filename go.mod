module github.com/loop/backend/rider-auth/rest

go 1.25.1

require ravigill/rider-grpc-server v0.0.0

replace ravigill/rider-grpc-server => ../auth-grpc
replace github.com/loop/backend/rider-auth/lib => ../lib

require github.com/loop/backend/rider-auth/lib v0.0.1


require (
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.77.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251124214823-79d6a2a48846 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
