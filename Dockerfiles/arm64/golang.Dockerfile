FROM arm64v8/alpine:latest

RUN apk update && apk add --no-cache \
    go \
    libc-dev \
    bash

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY run-code.sh /usr/bin/
RUN chmod +x /usr/bin/run-code.sh