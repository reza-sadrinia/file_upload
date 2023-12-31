# Stage 1: Build the Go binary
FROM golang:1.21.5 AS builder

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/src/app/main .
COPY --from=builder /go/src/app/templates /app/templates

EXPOSE 8080

CMD ["./main"]
