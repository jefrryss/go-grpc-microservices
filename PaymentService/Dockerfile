FROM golang:1.25.1-alpine AS builder

WORKDIR /src

COPY PaymentService/go.mod PaymentService/go.sum ./PaymentService/
COPY shared/go.mod shared/go.sum ./shared/

WORKDIR /src/PaymentService
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

WORKDIR /src
COPY PaymentService ./PaymentService
COPY shared ./shared

WORKDIR /src/PaymentService
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/payment-service ./cmd

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -D -H -G app app

WORKDIR /app

COPY --from=builder --chown=app:app \
    /out/payment-service ./payment-service

USER app

EXPOSE 50053

ENTRYPOINT ["./payment-service"]
