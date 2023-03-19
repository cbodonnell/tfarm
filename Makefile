include .env

VERSION ?= $(shell git rev-parse --short HEAD)
ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

clean:
	rm -rf ./bin/* ./dist/*

tfarm:
	go build \
	-ldflags="-X 'github.com/cbodonnell/tfarm/pkg/version.Version=${VERSION}'" \
	-o ./bin/tfarm ./cmd/tfarm/main.go

tfarm-server-start: tfarm
	TFARMD_FRPC_BIN_PATH=${TFARMD_FRPC_BIN_PATH} \
	TFARMD_WORK_DIR=${TFARMD_WORK_DIR} \
	TFARMD_LOG_LEVEL=${TFARMD_LOG_LEVEL} \
	./bin/tfarm server start \
		--frps-server-addr=localhost \
		--frps-server-port=7000 \
		--frpc-log-level ${TFARMD_LOG_LEVEL}

tfarm-server-start-dev: tfarm
	TFARMD_FRPC_BIN_PATH=${TFARMD_FRPC_BIN_PATH} \
	TFARMD_WORK_DIR=${TFARMD_DEV_WORK_DIR} \
	TFARMD_LOG_LEVEL=${TFARMD_LOG_LEVEL} \
	./bin/tfarm server start \
		--frpc-log-level ${TFARMD_LOG_LEVEL}

tfarm-server-configure-dev: tfarm
	TFARMD_WORK_DIR=${TFARMD_DEV_WORK_DIR} \
	./bin/tfarm server configure

tfarm-server-certs-regenerate: tfarm
	TFARMD_WORK_DIR=${TFARMD_WORK_DIR} \
	./bin/tfarm server certs regenerate

trivy-fs:
	trivy fs --scanners vuln --ignore-unfixed .

release: clean
	goreleaser release -f ./deploy/.goreleaser.yaml

release-snapshot: clean
	goreleaser release -f ./deploy/.goreleaser.yaml --snapshot