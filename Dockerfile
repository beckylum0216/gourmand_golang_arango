FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

RUN go build -o server .

FROM alpine:3.18

COPY --from=builder /app/server /app/server

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/server"]