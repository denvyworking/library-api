FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd
RUN CGO_ENABLED=0 go build -o migrate ./cmd/migrate

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

RUN apk --no-cache add ca-certificates bash

EXPOSE 8080

CMD ["./main"]