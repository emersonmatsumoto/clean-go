FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./
COPY api/ ./api/
COPY products/ ./products/
COPY payments/ ./payments/
COPY orders/ ./orders/

RUN go work sync
RUN go build -o /app/server ./api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
