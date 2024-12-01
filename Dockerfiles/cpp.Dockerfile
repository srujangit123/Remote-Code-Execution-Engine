FROM gcc:4.9

COPY run-code.sh /usr/bin/
RUN chmod +x /usr/bin/run-code.sh