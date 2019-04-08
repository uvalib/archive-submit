GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

build: darwin deploy-templates web

all: darwin linux deploy-templates web

linux-full: linux deploy-templates web

darwin-full: darwin deploy-templates web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/submitsrv.darwin backend/submitsrv/*.go

deploy-templates:
	mkdir -p bin/
	rm -rf bin/templates
	mkdir -p bin/templates
	cp ./templates/* bin/templates

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
