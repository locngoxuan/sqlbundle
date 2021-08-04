.PHONY: sqlbundle

SQLBUNDLE_BUILD=sqlbundle
INSTALL_DIR=/usr/local/bin
VERSION=1.6.0

prepare:
	@export GOPROXY=direct
	@export GOSUMDB=off
	go get -v .

#apply for develop
dev:
	env CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

#apply on release
release:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

#apply on release
compress:
	upx --brute ./bin/${SQLBUNDLE_BUILD}

install:
	chmod 755 ./bin/${SQLBUNDLE_BUILD}
	cp -r ./bin/${SQLBUNDLE_BUILD} ${INSTALL_DIR}/${SQLBUNDLE_BUILD}

docker: release compress
	docker build --force-rm -t xuanloc0511/sqlbundle:${VERSION} -f Dockerfile .
	docker tag xuanloc0511/sqlbundle:${VERSION}  xuanloc0511/sqlbundle:latest 

clean:
	rm -rf ./bin