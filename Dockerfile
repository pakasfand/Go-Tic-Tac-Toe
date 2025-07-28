# Multi-stage build for Go Tic-Tac-Toe
FROM golang:latest AS builder

# Set working directory
WORKDIR /app

# Copy all source code first (needed for Go workspaces)
COPY . .

# Remove individual go.work files that have incorrect paths
RUN rm -f client/go.work server/go.work client/go.work.sum server/go.work.sum

# Create a root-level go.work file to manage all modules
RUN echo "go 1.23.4" > go.work && \
    echo "" >> go.work && \
    echo "use (" >> go.work && \
    echo "    ./client" >> go.work && \
    echo "    ./server" >> go.work && \
    echo "    ./shared_types" >> go.work && \
    echo ")" >> go.work

# Show the workspace structure for debugging
RUN echo "=== Workspace structure ===" && \
    cat go.work && \
    echo "=== Available directories ===" && \
    ls -la

# Build WebAssembly client
WORKDIR /app/client
ENV GOOS=js
ENV GOARCH=wasm
RUN go build -o main.wasm .

# Build Go server
WORKDIR /app/server
ENV GOOS=linux
ENV GOARCH=amd64
RUN CGO_ENABLED=0 go build -o server .

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built server binary
COPY --from=builder /app/server/server .

# Copy the client files (HTML, WASM, JS, assets)
COPY --from=builder /app/client/ ./client/

# Expose port 8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the server
CMD ["./server"]
