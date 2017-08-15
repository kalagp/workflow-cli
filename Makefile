#
# Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
# Dell EMC Confidential/Proprietary Information
#
#

ORGANIZATION = cpsd
PROJECT = workflow-cli
BINARYNAME = workflow-cli
GOOUT = ./bin
SHELL = /bin/bash

# variable definitions
PKGPATH = github.com/dellemc-symphony/workflow-cli
COMMITHASH = $(shell git describe --tags --always --dirty)
BUILDDATE = $(shell date -u)
GOVERSION = $(shell go version)
ifndef BUILD_ID
	RELEASEVERSION = v0.0.1-dev
else
	RELEASEVERSION = v0.0.1-${BUILD_ID}
endif

#Flags to pass to main.go
LDFLAGS = -ldflags "-X '${PKGPATH}/cmd.binaryName=${BINARYNAME}' \
	  -X '${PKGPATH}/cmd.buildDate=${BUILDDATE}' \
	  -X '${PKGPATH}/cmd.commitHash=${COMMITHASH}' \
	  -X '${PKGPATH}/cmd.goVersion=${GOVERSION}' \
	  -X '${PKGPATH}/cmd.releaseVersion=${RELEASEVERSION}' "


default: deps build test

creds:
	@$(eval CREDS = $(subst :, ,$(GIT_CREDS)))
	@$(eval GIT_USER = $(word 1, $(CREDS)))  # this variable is set by Jenkinsfile
	@$(eval GIT_PASS = $(word 2, ,$(CREDS)))
	@git config --global user.name $(GIT_USER)
	@echo "machine github.com\n login $(GIT_USER)\n password $(GIT_PASS)" > ~/.netrc
	@chmod 600 ~/.netrc

deps:
	go get github.com/Masterminds/glide
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	glide install
	env GOOS=windows go get -d ./...

mock: build-linux
	go run mock/main/main.go

integration-test: build-linux
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=65 --focus="\bINTEGRATION\b" -- --https=false
	mv cmd/cmd.coverprofile INTEGRATION_http.coverprofile

	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=65 --focus="\bINTEGRATION\b" -- --https=true
	mv cmd/cmd.coverprofile INTEGRATION_https.coverprofile

	cp resources/endpoint_template.yaml $(GOOUT)/linux/endpoint.yaml

unit-test: build-linux
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=65 --focus="\bUNIT\b" -- --https=false
	mv cmd/cmd.coverprofile UNIT_http.coverprofile

	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=65 --focus="\bUNIT\b" -- --https=true
	mv cmd/cmd.coverprofile UNIT_https.coverprofile

	cp resources/endpoint_template.yaml $(GOOUT)/linux/endpoint.yaml

test: build-linux
	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=65 -- --https=false
	mv cmd/cmd.coverprofile http.coverprofile

	ginkgo -r -race -trace -cover -randomizeAllSpecs --slowSpecThreshold=65 -- --https=true
	mv cmd/cmd.coverprofile https.coverprofile

	cp resources/endpoint_template.yaml $(GOOUT)/linux/endpoint.yaml

cover-cmd: test
	go tool cover -html=cmd/cmd.coverprofile

coverage:
	./coverage.sh

build: build-linux build-mac build-windows

build-linux:
	env GOOS=linux go build -o $(GOOUT)/linux/$(BINARYNAME) $(LDFLAGS)
	cp resources/endpoint_template.yaml $(GOOUT)/linux/endpoint.yaml

build-mac:
	env GOOS=darwin go build -o $(GOOUT)/darwin/$(BINARYNAME) $(LDFLAGS)
	cp resources/endpoint_template.yaml $(GOOUT)/darwin/endpoint.yaml

build-windows:
	env GOOS=windows go build -o $(GOOUT)/windows/$(BINARYNAME).exe $(LDFLAGS)
	cp resources/endpoint_template.yaml $(GOOUT)/windows/endpoint.yaml
