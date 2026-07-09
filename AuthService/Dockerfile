FROM golang:1.25.1-alpine AS builder

WORKDIR /src

COPY AuthService/go.mod AuthService/go.sum ./AuthService/
COPY shared/go.mod shared/go.sum ./shared/
COPY platform/go.mod platform/go.sum ./platform/

WORKDIR /src/AuthService

RUN --mount=type=cache,target=/go/pkg/mod go mod download

WORKDIR /src

COPY AuthService ./AuthService
COPY shared ./shared
COPY platform ./platform

WORKDIR /src/AuthService

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/auth-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -D -H -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /out/auth-service ./auth-service
COPY --chown=app:app AuthService/migrations ./migrations

USER app

EXPOSE 50054 8084

ENTRYPOINT ["./auth-service"]
