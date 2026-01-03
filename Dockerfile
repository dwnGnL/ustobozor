FROM golang:1.25.1 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/api ./cmd/api

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app \
    && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY migrations /app/migrations

RUN mkdir -p /app/uploads

ENV APP_PORT=8080
EXPOSE 8080

USER app
ENTRYPOINT ["/app/api"]
