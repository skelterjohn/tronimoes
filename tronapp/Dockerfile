# Build stage
FROM cgr.dev/chainguard/node AS builder
WORKDIR /app
USER root 
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Production stage
FROM cgr.dev/chainguard/node AS runner
WORKDIR /app
USER root

ENV NODE_ENV production

# Copy necessary files from builder
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Create cache directory with correct permissions
RUN mkdir -p .next/cache
RUN chown -R node:node .

# Switch back to nonroot user
USER "node"

EXPOSE 3000
ENV PORT 3000
ENV HOSTNAME "0.0.0.0"

# Set both ENTRYPOINT and CMD
ENTRYPOINT ["node"]
CMD ["server.js"]
