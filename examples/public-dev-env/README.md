# Public Dev Environment

This is an example using tunnel.farm to create a public development environment.

The development environment is a simple web server that proxies requests to local services.

## Requirements

1. docker

## Usage

1. Start an nginx container that proxies requests to local services:

```bash
make dev
```

In this example, an nginx server is listening on port 8080. Requests to the /api/ path are proxied to a backend server listening on port 5555 and all other requests are proxied to a frontend server listening on port 3000.

2. Use `tfarm` to expose the nginx server to the internet:

```bash
tfarm create dev -p 8080
```

3. Visit the URL provided by `tfarm status`.

```bash
tfarm status

# Name  Type  Status   LocalAddr       Plugin  RemoteAddr                  Error  
# dev   HTTP  running  127.0.0.1:8080          http://subdomain.tunnel.farm         
```
