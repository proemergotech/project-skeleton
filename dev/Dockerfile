FROM golang:1.13.3

ARG SERVICE_NAME=dliver-project-skeleton

RUN go get github.com/githubnemo/CompileDaemon

RUN mkdir -p /go/src/gitlab.com/proemergotech/$SERVICE_NAME
WORKDIR /go/src/gitlab.com/proemergotech/$SERVICE_NAME

VOLUME /go/src/gitlab.com/proemergotech/$SERVICE_NAME

COPY ./entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run"]

EXPOSE 80
