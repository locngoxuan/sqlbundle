.PHONY: sqlbundle

SQLBUNDLE_BUILD=sqlbundle
INSTALL_DIR=/usr/local/bin

#apply for develop
dev:
	env CGO_ENABLED=0 go build -tags='no_oracle' -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd

#apply on release
release:
	env CGO_ENABLED=1 go build -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD}-linux -a ./cmd
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o ./bin/${SQLBUNDLE_BUILD}-wins.exe -a ./cmd
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD}-darwin -a ./cmd

release-pgonly:
	env CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD} -a ./cmd
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD}-linux -a ./cmd
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/${SQLBUNDLE_BUILD}-wins.exe -a ./cmd
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/${SQLBUNDLE_BUILD}-darwin -a ./cmd

#apply on release
compress:
	upx --brute ./bin/${SQLBUNDLE_BUILD}
	upx --brute ./bin/${SQLBUNDLE_BUILD}-linux
	upx --brute ./bin/${SQLBUNDLE_BUILD}-darwin

install:
	chmod 755 ./bin/${SQLBUNDLE_BUILD}
	cp -r ./bin/${SQLBUNDLE_BUILD} ${INSTALL_DIR}/${SQLBUNDLE_BUILD}

clean:
	rm -rf ./bin