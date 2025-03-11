# syntax=docker/dockerfile:1.4
# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

ARG TARGETOS
ARG TARGETARCH
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build -ldflags="-w -s" -o /out/app ./cmd/main.go

FROM gcr.io/distroless/static:latest
COPY --from=build /out/app /app
ENTRYPOINT ["/app"]