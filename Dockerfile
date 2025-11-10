# TURN Server Dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./
COPY cmd/ ./cmd/

# Build
RUN go build -o turn-server ./cmd/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/turn-server .

# TURN server uses UDP port 3478 and relay ports 60000-65535
EXPOSE 3478/udp
EXPOSE 60000-65535/udp

# Environment variables (set these in Railway)
# TURN_USERNAME - TURN server username
# TURN_PASSWORD - TURN server password
# PUBLIC_IP - Your server's public IP
# TURN_PORT - TURN port (default: 3478)
# WORKER_ID - Optional Supabase worker ID
# TM_PROJECT - Optional Supabase project URL
# TM_ANONKEY - Optional Supabase anon key

CMD ["./turn-server"]
