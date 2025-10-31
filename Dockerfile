# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bookmarks .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bookmarks .

# Copy default bookmarks file (can be overridden with volume mount)
COPY bookmarks.yaml .

# Copy static files
COPY static ./static

# Expose port
EXPOSE 8080

# Run the application
CMD ["./bookmarks"]
