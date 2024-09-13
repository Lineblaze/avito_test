FROM golang:1.22-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /zadanie-6105/backend

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/ ./

WORKDIR /zadanie-6105/backend/cmd/app

RUN go build -o /app/tenders-service

FROM alpine:latest

ENV POSTGRES_CONN=postgres://postgres:963852741@postgres:5432/tenders
ENV SERVER_ADDRESS=0.0.0.0:8080

WORKDIR /app

COPY --from=builder /app/tenders-service .

EXPOSE 8080

CMD ["./tenders-service"]
