# syntax=docker/dockerfile:1

# ------------------------------------------------------------------------------- #
# Stage 1 — Build
# ------------------------------------------------------------------------------- #
FROM golang:1.24 AS build-stage

WORKDIR /app

# Allow automatic toolchain download for newer Go versions.
ENV GOTOOLCHAIN=auto

# Cache dependency downloads in a separate layer.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source tree.
COPY . .

# Compile a fully static Linux binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

# ------------------------------------------------------------------------------- #
# Stage 2 — Test (fails the build if any test fails)
# ------------------------------------------------------------------------------- #
FROM build-stage AS test-stage
RUN go test -v ./...

# ------------------------------------------------------------------------------- #
# Stage 3 — Release (minimal production image)
# ------------------------------------------------------------------------------- #
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

# Copy only the compiled binary from the build stage.
COPY --from=build-stage /server /server

EXPOSE 8080

# Run as non-root for security.
USER nonroot:nonroot

ENTRYPOINT ["/server"]
