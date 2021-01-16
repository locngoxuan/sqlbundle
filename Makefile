.PHONY: sqlbundle

SQLBUNDLE_BUILD=sqlbundle
INSTALL_DIR=/usr/local/bin
VERSION=1.3.0

prepare:
	@export GOPROXY=direct
	@export GOSUMDB=off
	go get -v .

#apply for develop
dev:
	env CGO_ENABLED=0 go build -tags='no_oracle' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

dev-nodb:
	env CGO_ENABLED=0 go build -tags='no_oracle no_postgres' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

#apply on release
release:
	env CGO_ENABLED=1 go build -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

release-pgonly:
	env CGO_ENABLED=0 go build -tags='no_oracle' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags='no_oracle' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD}-linux -a ./cmd
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags='no_oracle' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD}-darwin -a ./cmd
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags='no_oracle' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD}-wins.exe -a ./cmd

release-nodb:
	env CGO_ENABLED=0 go build -tags='no_oracle no_postgres' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags='no_oracle no_postgres' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD}-linux -a ./cmd
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -tags='no_oracle no_postgres' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD}-darwin -a ./cmd
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags='no_oracle no_postgres' -ldflags="-s -w -X main.version=${VERSION}" -o ./bin/${SQLBUNDLE_BUILD}-wins.exe -a ./cmd

#apply on release
compress:
	upx --brute ./bin/${SQLBUNDLE_BUILD}
	upx --brute ./bin/${SQLBUNDLE_BUILD}-linux
	upx --brute ./bin/${SQLBUNDLE_BUILD}-darwin
	upx --brute ./bin/${SQLBUNDLE_BUILD}-wins.exe

install:
	chmod 755 ./bin/${SQLBUNDLE_BUILD}
	cp -r ./bin/${SQLBUNDLE_BUILD} ${INSTALL_DIR}/${SQLBUNDLE_BUILD}

clean:
	rm -rf ./bin