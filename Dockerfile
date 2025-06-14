FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY . .

# First build
RUN go build -o honeypot ./cmd/honeypot

# Run analysis and generate binary_info.json
RUN ./honeypot -build

# Rebuild with embedded data
RUN go build -o honeypot ./cmd/honeypot

FROM scratch
COPY --from=builder /build/honeypot /honeypot
