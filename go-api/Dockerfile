FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go-api-server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go-api-server .

EXPOSE 8080

CMD ["/app/go-api-server"]
