FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o smart-cashier-backend ./cmd/main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/smart-cashier-backend .
RUN mkdir -p uploads data

ENV GIN_MODE=release
ENV PORT=8080
ENV DB_PATH=/app/data/smart_cashier.db

EXPOSE 8080

CMD ["./smart-cashier-backend"]
