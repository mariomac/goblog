FROM golang:1.23 AS builder

WORKDIR /goblog

COPY . .

RUN make compile

FROM ubuntu:latest

WORKDIR /

COPY --from=builder /goblog/bin/goblog /goblog

ENTRYPOINT ["/goblog"]