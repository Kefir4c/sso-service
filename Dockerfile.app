FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/sso ./cmd/sso/main.go

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/bin/sso /app/sso
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations

ENV CONFIG_PATH=/app/config/remote.yaml
ENV POSTGRES_PASS=postgres
ENV REDIS_PASS=

EXPOSE 44044

CMD ["/app/sso"]