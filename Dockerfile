ARG EXECUTABLE_NAME=dliver-project-skeleton

FROM golang:1.14-alpine AS builder

ARG APP_VERSION
ARG EXECUTABLE_NAME
ARG GOPROXY
ARG GONOPROXY
ARG GOPRIVATE

ENV ROOT_PACKAGE=gitlab.com/proemergotech/$EXECUTABLE_NAME

ADD . $GOPATH/src/$ROOT_PACKAGE
WORKDIR $GOPATH/src/$ROOT_PACKAGE

RUN GOPROXY=$GOPROXY GONOPROXY=$GONOPROXY GOPRIVATE=$GOPRIVATE go build -mod=readonly -ldflags "-X $ROOT_PACKAGE/app/config.AppVersion=$APP_VERSION" -o "/tmp/$EXECUTABLE_NAME"




FROM alpine:latest

ARG EXECUTABLE_NAME

RUN set -eux; \
  apk add --no-cache ca-certificates

WORKDIR /usr/local/bin/

COPY --from=builder /tmp/$EXECUTABLE_NAME ./$EXECUTABLE_NAME
RUN chmod +x ./$EXECUTABLE_NAME

EXPOSE 80

# can't use variables here, there is no shell to interpret them
CMD ["dliver-project-skeleton"]
