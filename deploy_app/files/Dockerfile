FROM golang:1.19.8 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./app/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
