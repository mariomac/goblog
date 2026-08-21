# First stage: Compile Go appllication
FROM golang:1.26 AS builder

ARG TARGETARCH

ENV CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH
WORKDIR /goblog
COPY go.mod go.sum ./
COPY ./vendor ./vendor
COPY ./src ./src
COPY ./sample ./sample
RUN go build -o /goblog/goblog ./src

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot@sha256:1c2c046bc09ed40fad370b599a0b1ae7987f55b01e247cf27a7c27cd97e5bbc7

WORKDIR /
COPY --from=builder /goblog/goblog /goblog
COPY --from=builder ./goblog/sample ./
ENTRYPOINT ["/goblog"]