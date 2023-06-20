FROM golang:1.20 as builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /usr/src/app

COPY . .
RUN go mod download

WORKDIR /usr/src/app/cmd/app
RUN go build -o /usr/local/bin/auth -buildvcs=false

FROM alpine:latest

COPY --from=builder /usr/local/bin/auth ./
CMD ["./auth"]
EXPOSE 9090
EXPOSE 50051