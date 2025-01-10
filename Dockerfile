# First stage: Compile Go appllication
FROM golang:1.23 AS builder

ARG TARGETARCH

ENV CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH
WORKDIR /goblog
COPY go.mod go.sum ./
RUN go mod download
COPY ./src ./src
COPY ./sample ./sample
RUN go build -o /goblog/goblog ./src

# Second stage: Create minimal image with compilled binary
FROM alpine:latest

WORKDIR /
COPY --from=builder /goblog/goblog /goblog
COPY --from=builder ./goblog/sample ./
ENTRYPOINT ["/goblog"]