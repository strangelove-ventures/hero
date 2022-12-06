FROM --platform=$BUILDPLATFORM golang:1.18-alpine as BUILD

WORKDIR /hero

RUN apk update && \
  apk --no-cache add make git

# Copy go.mod and go.sum first and download for caching go modules
COPY go.mod go.sum ./

RUN go mod download

# Copy the files from host
COPY . .

ARG TARGETARCH TARGETOS
ENV GOOS=${TARGETOS} GOARCH=${TARGETARCH} LEDGER_ENABLED=false BUILD_TAGS=muslc
RUN make build

FROM alpine:latest

ENV USER_HOME /hero

RUN addgroup --gid 1025 herouser && adduser --uid 1025 -S -G herouser herouser -h "$USER_HOME"

USER herouser

WORKDIR $USER_HOME

COPY --from=BUILD /hero/bin/herod /usr/bin/herod