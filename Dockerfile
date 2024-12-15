FROM golang:1.23-alpine3.19 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -ldflags="-w -s" -o app cmd/main.go

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /src/app .
COPY --from=builder /src/.env .

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser:appgroup

EXPOSE 8080

ENTRYPOINT ["./app"]