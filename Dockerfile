FROM golang:1.11.2-alpine AS builder

ARG DEPLOY_SSH_PRIVATE_KEY
ARG APP_VERSION

ENV EXECUTABLE_NAME=dliver-project-skeleton
ENV ROOT_PACKAGE=gitlab.com/proemergotech/dliver-project-skeleton
ENV DEP_VERSION=0.5.0

RUN apk add --update --no-cache wget openssh-client git
RUN wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/v$DEP_VERSION/dep-linux-amd64 && chmod +x /usr/local/bin/dep

ADD . $GOPATH/src/$ROOT_PACKAGE
WORKDIR $GOPATH/src/$ROOT_PACKAGE

RUN eval $(ssh-agent -s) \
  && echo "${DEPLOY_SSH_PRIVATE_KEY}" | ssh-add - \
  && mkdir ~/.ssh && touch ~/.ssh/known_hosts \
  && ssh-keyscan -t rsa gitlab.com >> ~/.ssh/known_hosts \
  && dep ensure -vendor-only

RUN go build -ldflags "-X gitlab.com/proemergotech/dliver-project-skeleton/app/config.AppVersion=$APP_VERSION" -o "/tmp/$EXECUTABLE_NAME"




FROM alpine:latest

ENV EXECUTABLE_NAME=dliver-project-skeleton

RUN set -eux; \
  apk add --no-cache --virtual ca-certificates

WORKDIR /usr/local/bin/

COPY --from=builder /tmp/$EXECUTABLE_NAME ./$EXECUTABLE_NAME
RUN chmod +x ./$EXECUTABLE_NAME

EXPOSE 80

# can't use variables here, there is no shell to interpret them
CMD ["dliver-project-skeleton"]
