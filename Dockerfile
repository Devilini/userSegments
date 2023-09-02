FROM golang:alpine as builder

WORKDIR /usr/local/src

COPY go.mod go.sum /
RUN go mod download

COPY . .
COPY .env /
RUN go build -o /bin/app cmd/main.go

FROM alpine
COPY --from=builder bin/app /
COPY --from=builder .env /
CMD ["/app"]
