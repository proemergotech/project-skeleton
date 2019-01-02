FROM golang:1.11.2

RUN go get github.com/githubnemo/CompileDaemon

RUN mkdir -p /go/src/gitlab.com/proemergotech/dliver-project-skeleton
WORKDIR /go/src/gitlab.com/proemergotech/dliver-project-skeleton

VOLUME /go/src/gitlab.com/proemergotech/dliver-project-skeleton

COPY ./entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run"]

EXPOSE 80
