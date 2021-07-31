FROM alpine:3.13.5

LABEL maintainer="xuanloc0511@gmail.com"

ADD bin/sqlbundle /usr/local/bin

ENTRYPOINT ["/usr/local/bin/sqlbundle"]