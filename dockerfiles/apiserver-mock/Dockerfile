FROM alpine:latest

ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/skpr-apiserver-mock /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/skpr-apiserver-mock"]
