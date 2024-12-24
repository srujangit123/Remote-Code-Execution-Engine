FROM arm64v8/alpine:latest

RUN apk update && apk add --no-cache \
    gcc \
    g++ \
    libc-dev \
    bash

COPY run-code.sh /usr/bin/
RUN chmod +x /usr/bin/run-code.sh