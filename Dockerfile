# Stage 1: Build the Go binary
FROM golang:1.22 AS builder

WORKDIR /app

# Copy Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# ✅ Build the binary (Disable CGO and use target platform)
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o cortex-api .

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# ✅ Ensure ca-certificates are installed (fixes TLS issues)
RUN apk --no-cache add ca-certificates

# ✅ Copy the built binary from the builder stage
COPY --from=builder /app/cortex-api /app/cortex-api

# ✅ Ensure the binary has execute permissions
RUN chmod +x /app/cortex-api

# Expose API port
EXPOSE 8080

# ✅ Run the binary
CMD ["/app/cortex-api"]
