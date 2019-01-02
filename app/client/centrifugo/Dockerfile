FROM znly/protoc

ENV GOMPLATE_VERSION=2.6.0

#RUN apt-get update && apt-get install -y zip

RUN apk add --no-cache wget ca-certificates bash

RUN wget -O /usr/local/bin/gomplate https://github.com/hairyhenderson/gomplate/releases/download/v$GOMPLATE_VERSION/gomplate_linux-amd64 \
  && chmod +x /usr/local/bin/gomplate

ADD generate.sh .

RUN chmod +x generate.sh

VOLUME /tmp/proto

ENTRYPOINT []

CMD ["./generate.sh"]
