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

# Copy only the built binary
COPY --from=builder /app/bin/laclm ./bin/laclm

# Install bash in case needed
# RUN apt-get update && apt-get install -y bash && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y bash acl && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

# Default command to run your Go app
CMD ["./bin/laclm", "--config", "config.yaml"]
