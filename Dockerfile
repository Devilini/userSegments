FROM golang:1.21 AS builder

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
