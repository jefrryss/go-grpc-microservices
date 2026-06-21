FROM golang:1.25.1-alpine AS builder

WORKDIR /src

COPY OrderService/go.mod OrderService/go.sum ./OrderService/
COPY shared/go.mod shared/go.sum ./shared/

WORKDIR /src/OrderService

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /src

COPY OrderService ./OrderService
COPY shared ./shared

WORKDIR /src/OrderService

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/order-service ./cmd && \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/migrator ./cmd/migrator

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -D -H -G app app

WORKDIR /app

COPY --from=builder --chown=app:app \
    /out/order-service ./order-service

COPY --from=builder --chown=app:app \
    /out/migrator ./migrator

COPY --chown=app:app OrderService/api ./api
COPY --chown=app:app OrderService/migrations ./migrations

USER app

EXPOSE 50051 8080

ENTRYPOINT ["./order-service"]
