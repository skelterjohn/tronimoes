# Build stage
FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go install .
RUN CGO_ENABLED=0 go test ./...

# Production stage
FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY --from=builder /go/bin/tronserv .

EXPOSE 8080
ENV PORT 8080
USER nonroot

ENTRYPOINT ["/app/tronserv"]
CMD ["--addr", "0.0.0.0", "--env", "prod"]