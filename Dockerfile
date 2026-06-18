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
FROM gcr.io/distroless/static:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240

WORKDIR /
COPY --from=builder /goblog/goblog /goblog
COPY --from=builder ./goblog/sample ./
ENTRYPOINT ["/goblog"]