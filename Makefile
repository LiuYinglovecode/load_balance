# make file for alb project
# define global variable
SHELL:=/bin/bash
PROJ=alb
ORG_PATH=code.htres.cn/casicloud
REPO_PATH=$(ORG_PATH)/$(PROJ)
export PATH := $(PWD)/bin:$(PATH)
export GOBIN=$(PWD)/bin
# build number
BN=$(shell ./scripts/gen-bn.sh)
VERSION=$(shell cat ALB_VERSION).$(BN)

LD_FLAGS="-w -X $(REPO_PATH)/version.Version=$(VERSION)"

SRCS := $(shell find . -name '*.go'| grep -v vendor)

build: bin/lbagent bin/lbmc bin/adc-cp

bin/lbagent:
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/lbagent

bin/lbmc:
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/lbmc

bin/adc-cp:
	@rm -f vendor\github.com\coreos\etcd\client\keys.generated.go
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/alb-cp

clean:
	@rm -rf bin/

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: proto
proto:
	@protoc -I apis/ apis/agent_apis.proto --go_out=plugins=grpc:apis

test:
	@go test -v ./...

testrace:
	@go test -v --race ./...

.PHONY: lint
lint: 
	@for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done