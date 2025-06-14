FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .

# First build
RUN go build -o honeypot cmd/honeypot/main.go

# Run analysis and generate binary_info.json
RUN ./honeypot -build

# Rebuild with embedded data
RUN go build -o honeypot cmd/honeypot/main.go

FROM scratch
COPY --from=builder /build/honeypot /honeypot