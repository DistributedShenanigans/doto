FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/doto cmd/doto/main.go

FROM alpine:latest AS runtime

WORKDIR /app
COPY --from=builder /app/bin/doto ./
COPY --from=builder /app/config/default-config.yaml ./config/

EXPOSE 8080

CMD ["./doto"]
