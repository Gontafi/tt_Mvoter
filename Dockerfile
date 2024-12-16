FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -ldflags="-w -s" -o server cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /src/.env .
COPY --from=builder /src/server .
COPY --from=builder /src/internal/db/migrations /app/internal/db/migrations

RUN chmod +x /app/server

CMD ["./server"]