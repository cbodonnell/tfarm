# tfarm

`tfarm` is a tool for creating and managing tunnels. It is a wrapper around [frp](https://github.com/fatedier/frp) that provides a simple interface for managing tunnels. It consists of two components: a server and a client. The server is a wrapper around `frpc` that manages tunnels. The client is a CLI for interfacing with the server to create and manage tunnels.

## Installation

If you are using `go` version 1.17 or later, you can install `frpc` and `tfarm` using the `go install` command.
```bash
go install github.com/fatedier/frp/cmd/frpc@latest
go install github.com/cbodonnell/tfarm/cmd/tfarm@latest
```

*If you are using an older version of `go`, you can install `tfarm` using `go get`.*

Otherwise, you can:
1. Download the latest release of [frp](https://github.com/fatedier/frp/releases) and [tfarm](https://github.com/cbodonnell/tfarm/releases) for your platform.
2. Extract the binaries from the archives.
3. Move the binaries to a directory in your PATH.

## Usage

### Setup the tfarm server

*Note: It is recommended to run the tfarm server using a process manager like `systemd` ([example](https://github.com/cbodonnell/tfarm/blob/main/examples/systemd/tfarm.service)).*

#### Configuration

The following environment variables can **optionally** be used to configure the tfarm server.

Set the `TFARMD_WORK_DIR` environment variable to the path where you want to store the tfarm server state. By default, it will be the current working directory.
```bash
export TFARMD_WORK_DIR=/path/to/work/dir
```

Set the `TFARMD_FRPC_BIN_PATH` environment variable to the path of the `frpc` binary. By default, tfarm will search the current users PATH.
```bash
export TFARMD_FRPC_BIN_PATH=/path/to/frpc
```

#### Start the tfarm server process

Start the tfarmd server.

```bash
tfarm server start
```

In another terminal, create the $HOME/.tfarm directory and copy the server's `client.json` file to it.

```bash
mkdir -p $HOME/.tfarm
cp $TFARMD_WORK_DIR/tls/client.json $HOME/.tfarm
```

Check the status of the tfarm server.

```bash
tfarm status
```

The next step is to configure the tfarm server as a ranch client.

#### Configure the tfarm server as a ranch client

Use the `tfarm ranch` command to interact with the tfarm ranch. The tfarm ranch is the `frps` server that `frpc` connects to. It provides an identity and access layer for `frps`. By default, tfarm will connect to the `tunnel.farm` ranch.

Login to the tfarm ranch.

```bash
tfarm ranch login
```

Create a new ranch client and use it to configure the tfarm server.

```bash
tfarm ranch clients create --credentials | tfarm configure --credentials-stdin
```

### Manage tunnels with the tfarm CLI

Check the status of the tfarm server.

```bash
tfarm status
```

Create a tunnel that forwards traffic to local port 8080.

```bash
tfarm create my-tunnel -p 8080
```

Check the status.

```bash
tfarm status
```

Delete the tunnel.

```bash
tfarm delete my-tunnel
```

## Development

### Dependencies

* Git
* Make
* Go
* Docker
* frpc

### Configuration

Create a `.env` file with the following environment variables.

```bash
TFARMD_FRPC_BIN_PATH=/path/to/frpc
TFARMD_WORK_DIR=/path/to/work/dir
TFARMD_LOG_LEVEL=debug
```

### Build

```bash
make tfarm
```
