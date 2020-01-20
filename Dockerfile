#FROM golang:latest AS build
FROM golang:alpine as builder
LABEL version="1.0"

RUN mkdir -p /go/scr/app

WORKDIR /go/scr/app
COPY .  /go/scr/app

RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/scr/app/main .
COPY --from=builder /go/scr/app/.env .

EXPOSE 8080

CMD ["./main"]