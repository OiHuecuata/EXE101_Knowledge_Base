FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY src/backend/go.mod src/backend/go.sum ./

RUN go mod download

COPY src/backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]