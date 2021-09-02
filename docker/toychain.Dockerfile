FROM golang:1.17 AS builder
WORKDIR /src
ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

WORKDIR /src/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server

FROM alpine:3.14
RUN apk update && \
  apk upgrade && \
  apk add --no-cache curl tzdata ca-certificates bash nano busybox-extras && \
  rm -rf /var/cache/apk/*

WORKDIR /app
#COPY --from=builder /app/deployments/database /app/deployments/database
COPY --from=builder /src/cmd/server /app/server


RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser && \
  chown -R appuser:appuser /app

USER appuser