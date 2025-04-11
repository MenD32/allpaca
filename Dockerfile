FROM golang:1.24.2 AS builder

WORKDIR /app
RUN mkdir dist
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM ubuntu:latest

WORKDIR /root/
COPY --from=builder /app/dist/allpaca .

ENTRYPOINT ["./allpaca"]