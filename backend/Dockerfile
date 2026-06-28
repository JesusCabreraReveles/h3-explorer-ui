# syntax=docker/dockerfile:1

# h3-go links the H3 C library via cgo, so the build stage needs a C toolchain
# (CGO_ENABLED=1) and the runtime image needs a glibc base.

# --- Build stage ----------------------------------------------------------
FROM golang:1.26-bookworm AS build

WORKDIR /src

# Cache dependency downloads across builds.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/server ./cmd/server

# --- Runtime stage --------------------------------------------------------
# distroless/base ships glibc + CA certs and runs as a non-root user.
FROM gcr.io/distroless/base-debian12 AS runtime

COPY --from=build /out/server /server

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/server"]
