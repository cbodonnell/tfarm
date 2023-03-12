# tfarm

tfarm is a tool for creating and managing tunnels. It is a wrapper around [frp](https://github.com/fatedier/frp) that provides a simple interface for managing tunnels. It consists of two components: a server and a client. The server is a wrapper around `frpc` that manages tunnels. The client is a CLI for interfacing with the server to create and manage tunnels.

## Installation

Install `frpc`.
```bash
go install github.com/fatedier/frp/cmd/frpc@latest
```

Install `tfarm`.
```bash
go install github.com/cbodonnell/tfarm/cmd/tfarm@latest
```

## Usage

### Setup the tfarm server

*Note: It is recommended to run the tfarm server using a process manager like `systemd`.*

#### Start the tfarm server process

Set the `TFARMD_WORK_DIR` environment variable to the path where you want to store the tfarm server state.
```bash
export TFARMD_WORK_DIR=/path/to/work/dir
export TFARMD_FRPC_BIN_PATH=$(which frpc)
```

Start the tfarmd server.

```bash
tfarm server start
```

#### Configure the tfarm server as a ranch client

Use the `tfarm ranch` command to interact with the tfarm ranch. The tfarm ranch is a the `frps` server that the tfarm server `frpc` connects to. It acts an identity and access layer for `frps`.

Login to the tfarm ranch.

```bash
tfarm ranch login
```

Create a new ranch client.

```bash
tfarm ranch clients create
```

Configure the tfarm server with the client credentials, where `$RANCH_CLIENT` is the id of the client you just created.

*Note: Use `tfarm ranch clients list` to list ranch clients.*

```bash
tfarm ranch clients get $RANCH_CLIENT --credentials | tfarm configure --credentials-stdin
```

### Manage tunnels with the tfarm client

Once the tfarm server is running, you can use the tfarm client to manage tunnels.

#### Validate that the tfarm server is running

Check the status of the tfarm server.

```bash
tfarm status
```

#### Create a tunnel

Create a tunnel that forwards traffic from the tunnel to local port 8080.

```bash
tfarm create my-tunnel -p 8080
```

Check the status.

```bash
tfarm status
```

#### Delete a tunnel

Delete the tunnel.

```bash
tfarm delete my-tunnel
```
