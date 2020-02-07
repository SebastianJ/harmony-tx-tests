SHELL := /bin/bash
version := $(shell git rev-list --count HEAD)
commit := $(shell git describe --always --long --dirty)
built_at := $(shell date +%FT%T%z)
built_by := ${USER}

flags := -gcflags="all=-N -l -c 2"
ldflags := -X main.version=v${version} -X main.commit=${commit}
ldflags += -X main.builtAt=${built_at} -X main.builtBy=${built_by}

upload-path-linux := 's3://tools.harmony.one/release/linux-x86_64/harmony-tx-tests'

dist := ./dist/harmony-tx-tests
env := GO111MODULE=on
DIR := ${CURDIR}

all:
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build -o $(dist) -ldflags="$(ldflags)" cmd/main.go

static:
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/mcl
	make -C $(shell go env GOPATH)/src/github.com/harmony-one/bls minimised_static BLS_SWAP_G=1
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build -o $(dist) -ldflags="$(ldflags) -w -extldflags \"-static\"" cmd/main.go

debug:
	source $(shell go env GOPATH)/src/github.com/harmony-one/harmony/scripts/setup_bls_build_flags.sh && $(env) go build $(flags) -o $(dist) -ldflags="$(ldflags)" cmd/main.go

upload-linux:static
	aws s3 cp dist/harmony-tx-sender ${upload-path-linux} --acl public-read

.PHONY:clean upload-linux

clean:
	@rm -f $(dist)
	@rm -rf ./dist
