include .env

VERSION ?= $(shell git rev-parse --short HEAD)
ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

clean:
	rm -rf ./bin/*

tfarmd:
	go build \
	-ldflags="-X 'github.com/cbodonnell/tfarm/cmd/tfarmd/commands.version=${VERSION}'" \
	-o ./bin/tfarmd ./cmd/tfarmd/main.go

tfarm:
	go build \
	-ldflags="-X 'github.com/cbodonnell/tfarm/cmd/tfarm/commands.version=${VERSION}'" \
	-o ./bin/tfarm ./cmd/tfarm/main.go

start-tfarmd:
	TFARMD_FRPC_BIN_PATH=${TFARMD_FRPC_BIN_PATH} \
	TFARMD_WORK_DIR=${TFARMD_WORK_DIR} \
	TFARMD_LOG_LEVEL=${TFARMD_LOG_LEVEL} \
	./bin/tfarmd start \
		--frpc-log-level ${TFARMD_LOG_LEVEL}

start-tfarmd-dev:
	TFARMD_FRPC_BIN_PATH=${TFARMD_FRPC_BIN_PATH} \
	TFARMD_WORK_DIR=${TFARMD_DEV_WORK_DIR} \
	TFARMD_LOG_LEVEL=${TFARMD_LOG_LEVEL} \
	TFARMD_FRPS_TOKEN=${TFARMD_FRPS_TOKEN} \
	./bin/tfarmd start \
		--frps-server-addr=ranch.tunnel.farm \
		--frps-server-port=30070 \
		--frps-token=${TFARMD_FRPS_TOKEN} \
		--frpc-log-level ${TFARMD_LOG_LEVEL}

tfarmd-certs:
	TFARMD_WORK_DIR=${TFARMD_WORK_DIR} \
	./bin/tfarmd certs
