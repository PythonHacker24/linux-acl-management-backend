# --- Build stage ---
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make

# --- Final image ---
FROM debian:latest

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/bin/laclm ./bin/laclm
COPY config.yaml .
COPY .env .

# Install `bash` to source the env file
RUN apt-get update && apt-get install -y bash && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

# Use a wrapper script to load env vars and run the binary
COPY docker-entrypoint.sh .

RUN chmod +x docker-entrypoint.sh

ENTRYPOINT ["./docker-entrypoint.sh"]
