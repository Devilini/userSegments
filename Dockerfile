FROM golang:1.18 AS builder

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy
