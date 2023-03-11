include .env

VERSION ?= $(shell git rev-parse --short HEAD)
ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

clean:
	rm -rf ./bin/*

tfarm:
	go build \
	-ldflags="-X 'github.com/cbodonnell/tfarm/pkg/version.Version=${VERSION}'" \
	-o ./bin/tfarm ./cmd/tfarm/main.go

tfarm-server-start:
	TFARMD_FRPC_BIN_PATH=${TFARMD_FRPC_BIN_PATH} \
	TFARMD_WORK_DIR=${TFARMD_WORK_DIR} \
	TFARMD_LOG_LEVEL=${TFARMD_LOG_LEVEL} \
	./bin/tfarm server start \
		--frpc-log-level ${TFARMD_LOG_LEVEL}

tfarm-server-start-dev:
	TFARMD_FRPC_BIN_PATH=${TFARMD_FRPC_BIN_PATH} \
	TFARMD_WORK_DIR=${TFARMD_DEV_WORK_DIR} \
	TFARMD_LOG_LEVEL=${TFARMD_LOG_LEVEL} \
	TFARMD_FRPS_TOKEN=${TFARMD_FRPS_TOKEN} \
	./bin/tfarm server start \
		--frps-server-addr=ranch.tunnel.farm \
		--frps-server-port=30070 \
		--frps-token=${TFARMD_FRPS_TOKEN} \
		--frpc-log-level ${TFARMD_LOG_LEVEL}

tfarm-server-configure-dev:
	TFARMD_WORK_DIR=${TFARMD_DEV_WORK_DIR} \
	./bin/tfarm server configure

tfarm-server-certs-regenerate:
	TFARMD_WORK_DIR=${TFARMD_WORK_DIR} \
	./bin/tfarm server certs regenerate

trivy-fs:
	trivy fs --scanners vuln --ignore-unfixed .