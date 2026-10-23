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
FROM gcr.io/distroless/static:nonroot@sha256:18f4cc26c9ffa7f20a359298be7deadac5c607c8e03a02b8119f1a5fc2f82c99

WORKDIR /
COPY --from=builder /goblog/goblog /goblog
COPY --from=builder ./goblog/sample ./
ENTRYPOINT ["/goblog"]