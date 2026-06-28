# syntax=docker/dockerfile:1

# --- Build stage ----------------------------------------------------------
FROM node:22-bookworm-slim AS build

WORKDIR /app

# Install dependencies against the lockfile for reproducible builds.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build && npm prune --omit=dev

# --- Runtime stage --------------------------------------------------------
FROM node:22-bookworm-slim AS runtime

WORKDIR /app
ENV NODE_ENV=production
ENV PORT=3000
# Same-origin API proxy target inside the compose network.
ENV API_PROXY_TARGET=http://backend:8080

COPY --from=build /app/build ./build
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./package.json

EXPOSE 3000
USER node

CMD ["node", "build"]
