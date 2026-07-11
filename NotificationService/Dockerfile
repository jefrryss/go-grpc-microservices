FROM golang:1.25.1-alpine AS builder

WORKDIR /src

COPY NotificationService/go.mod NotificationService/go.sum ./NotificationService/
COPY shared/go.mod shared/go.sum ./shared/
COPY platform/go.mod platform/go.sum ./platform/

WORKDIR /src/NotificationService

RUN --mount=type=cache,target=/go/pkg/mod go mod download

WORKDIR /src

COPY NotificationService ./NotificationService
COPY shared ./shared
COPY platform ./platform

WORKDIR /src/NotificationService

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/notification-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -D -H -G app app

WORKDIR /app

COPY --from=builder --chown=app:app /out/notification-service ./notification-service

USER app

ENTRYPOINT ["./notification-service"]
