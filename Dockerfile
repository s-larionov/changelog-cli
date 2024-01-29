FROM golang:1.21.6-alpine3.19 AS builder

LABEL maintainer="Sergey Larionov <sergey@larionov.it>"

RUN apk update && \
    apk add bash ca-certificates git gcc g++ libc-dev binutils file

WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /opt/changelog-cli .

FROM alpine:3.19 as production

ENV COMMAND "diff"
ENV FILE "CHANGELOG.md"
ENV FROM "latest"
ENV TO "Unreleased"
ENV VERSION ""
ENV BUMP "auto"

RUN apk update && \
    apk add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*
RUN echo "Europe/Moscow" >  /etc/timezone && cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime

WORKDIR /opt
COPY . /opt

CMD ["sh", "-c", "./changelog-cli -command=${COMMAND} -file=${FILE} -from=${FROM} -to=${TO} -bump=${BUMP} -version=${VERSION}"]
