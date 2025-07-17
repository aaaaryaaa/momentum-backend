# Use the official Golang image
FROM golang:1.24-alpine

# Create app directory
WORKDIR /app

# Install git and certificates for Go modules
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download Go modules
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o post-service .

# Expose the port
EXPOSE 8082

# Run the binary
CMD ["./post-service"]
