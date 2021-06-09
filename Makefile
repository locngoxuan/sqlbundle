.PHONY: sqlbundle

SQLBUNDLE_BUILD=sqlbundle
INSTALL_DIR=/usr/local/bin
VERSION=1.4.0

prepare:
	@export GOPROXY=direct
	@export GOSUMDB=off
	go get -v .

#apply for develop
dev:
	go build -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

dev-nodb:
	env CGO_ENABLED=0 go build -tags='no_oracle no_postgres no_sqlite' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

#apply on release
release:
	go build -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

#apply on release
compress:
	upx --brute ./bin/${SQLBUNDLE_BUILD}

install:
	chmod 755 ./bin/${SQLBUNDLE_BUILD}
	cp -r ./bin/${SQLBUNDLE_BUILD} ${INSTALL_DIR}/${SQLBUNDLE_BUILD}

clean:
	rm -rf ./bin