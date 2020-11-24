FROM golang:1.15

#%:{{ `
ARG SERVICE_NAME=project-skeleton
#%: ` | replace "project-skeleton" .ProjectName | trim }}

RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/go-delve/delve/cmd/dlv

RUN mkdir -p /go/src/github.com/proemergotech/$SERVICE_NAME
WORKDIR /go/src/github.com/proemergotech/$SERVICE_NAME

VOLUME /go/src/github.com/proemergotech/$SERVICE_NAME

COPY ./entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["run"]

EXPOSE 80
