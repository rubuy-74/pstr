FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/pstr ./cmd/pstr/

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /bin/pstr .

EXPOSE 3000
ENTRYPOINT ["./server"]
