# Build stage
FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go install . ./agent
RUN CGO_ENABLED=0 go test ./game
RUN mkdir /app/bin
RUN cp $HOME/go/bin/* /app/bin

# Production stage
FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY --from=builder /app/bin/tronserv .
COPY --from=builder /app/bin/agent .

EXPOSE 8080
ENV PORT 8080
USER nonroot

ENTRYPOINT ["/app/tronserv"]
CMD ["--addr", "0.0.0.0", "--env", "prod", "--agent-spawner", "local"]