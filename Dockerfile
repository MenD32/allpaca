FROM golang:1.24.2 AS builder


WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o allpaca .

FROM ubuntu:latest

WORKDIR /root/
COPY --from=builder /app/allpaca .

ENTRYPOINT ["./allpaca"]