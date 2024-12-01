FROM amd64/golang:1.23.3-alpine3.20

COPY run-code.sh /usr/bin/
RUN chmod +x /usr/bin/run-code.sh