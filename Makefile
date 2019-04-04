GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

build: darwin web

all: darwin linux web

linux-full: linux web

darwin-full: darwin web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/submitsrv.darwin backend/submitsrv/*.go

web:
	mkdir -p bin/
	cd frontend/; yarn install; yarn build
	rm -rf bin/public
	mv frontend/dist bin/public

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/submitsrv.linux backend/submitsrv/*.go

clean:
	$(GOCLEAN)
	rm -rf bin
