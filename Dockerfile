FROM alpine:latest

LABEL maintainer="xuanloc0511@gmail.com"

ADD bin/sqlbundle /usr/local/bin

ENTRYPOINT ["/usr/local/bin/sqlbundle"]