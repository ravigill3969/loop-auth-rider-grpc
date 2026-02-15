FROM golang:1.25.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY ./rider-auth/rest/go.mod ./rider-auth/rest/go.sum ./rider-auth/rest/

COPY libs ./libs
COPY trip ./trip
COPY ./rider-auth/auth-grpc ./rider-auth/auth-grpc

WORKDIR /app/rider-auth/rest

RUN go mod download

COPY ./rider-auth/rest .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rider-rest ./cmd/server

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/rider-auth/rest .

EXPOSE 8081
USER nonroot:nonroot

CMD ["./rider-rest"]
