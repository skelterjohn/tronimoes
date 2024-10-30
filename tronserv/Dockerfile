# Build stage
FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o tronserv

# Production stage
FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY --from=builder /app/tronserv .

EXPOSE 8080
USER nonroot

ENTRYPOINT ["/app/tronserv"]
