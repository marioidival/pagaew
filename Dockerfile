FROM golang:latest as builder
ENV CGO_ENABLED 0

COPY . /service

WORKDIR /service
ENV GOBIN /service/bin
RUN go install -ldflags "-X main.build=version" ./cmd/...

FROM alpine:3.14
COPY --from=builder /service/bin /service
COPY --from=builder /service/cmd/dbmigrate/migrations /service/migrations

WORKDIR /service

