GOARCH=amd64

.PHONY: build
build: linux windows darwin

.PHONY: linux
linux: linux-client linux-server

.PHONY: windows
windows: windows-client windows-server

.PHONY: darwin
darwin: darwin-client darwin-server

.PHONY: client
client: linux-client windows-client

.PHONY: linux-client
linux-client:
	cd client/go && \
	GOOS=linux go build -ldflags="-s -w" -o ../../bin/client-linux-${GOARCH}

.PHONY: windows-client
windows-client:
	cd client/go && \
	GOOS=windows go build -ldflags="-s -w" -o ../../bin/client-windows-${GOARCH}.exe

.PHONY: darwin-client
darwin-client:
	cd client/go && \
	GOOS=darwin go build -ldflags="-s -w" -o ../../bin/client-darwin-${GOARCH}

.PHONY: server
server: linux-server windows-server darwin-server

.PHONY: linux-server
linux-server:
	cd server && \
	GOOS=linux go build -ldflags="-s -w" -o ../bin/server-linux-${GOARCH}

.PHONY: windows-server
windows-server:
	cd server && \
	GOOS=windows go build -ldflags="-s -w" -o ../bin/server-windows-${GOARCH}.exe

.PHONY: darwin-server
darwin-server:
	cd server && \
	GOOS=darwin go build -ldflags="-s -w" -o ../bin/server-darwin-${GOARCH}

.PHONY: lint
lint:
	golangci-lint run ./client/... ./server/...

.PHONY: test
test:
	go test -v -race -count=1 ./...

.PHONY: deps
deps: deps-client deps-server

.PHONY: deps-client
deps-client:
	cd client/go && \
	go mod verify && \
	go mod tidy

.PHONY: deps-server
deps-server:
	cd server && \
	go mod verify && \
	go mod tidy