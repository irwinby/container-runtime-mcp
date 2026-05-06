# container-runtime-mcp

[![CI](https://github.com/irwinby/container-runtime-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/irwinby/container-runtime-mcp/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/irwinby/container-runtime-mcp/graph/badge.svg)](https://codecov.io/gh/irwinby/container-runtime-mcp)

`container-runtime-mcp` is a Model Context Protocol (MCP) server that exposes container runtime operations as MCP tools.

The current implementation uses Docker-compatible APIs. Docker is supported directly, and compatible runtimes such as Podman may work when exposed through a Docker-compatible API socket.

The server supports two MCP transports:

- **stdio** (default): The MCP client starts the process and communicates over stdin and stdout.
- **HTTP** (streamable): The server listens on a configurable address and path, exposing MCP over HTTP with server-sent events.

## Requirements

- Go 1.26.2 or newer, matching `go.mod`, when building from source.
- A running Docker-compatible container runtime.
- Permission for this process to access the runtime, such as access to the local Docker-compatible socket or an environment configured like the runtime CLI.
- Registry credentials already available to the runtime when using `pull_image` or `push_image` with private registries.
- Docker or Podman, when building or running the container image.

The runtime client is created using the Docker Go SDK, which resolves Docker context and environment variables automatically. Docker Engine API calls use Moby client types. In practice, run this server from an environment where the runtime CLI can list containers for the same user, socket, context, and related environment variables.

## Runtime Support

| Runtime | Support | Notes |
| --- | --- | --- |
| Docker Engine | Supported | Primary supported runtime. Use the same Docker socket, context, and environment where `docker ps` works. |
| Docker Desktop | Supported | Supported through Docker Desktop's Docker Engine API and active Docker context. |
| Colima with Docker runtime | Supported | Supported when Colima is running with Docker runtime and the Docker context/socket points to Colima. |
| Podman | Experimental | May work through Podman's Docker-compatible API socket. Compatibility depends on Podman version and enabled API features. |
| containerd / nerdctl | Not directly supported | Not supported directly because the current provider uses Docker-compatible APIs, not native containerd APIs. |

## Running

Run directly from the repository:

```sh
make run
```

Build a binary:

```sh
make build
```

### Container Image

Build the container image:

```sh
make image-build
```

Run with stdio transport:

```sh
make image-run
```

Run with HTTP transport:

```sh
make image-run-http
```

Run with HTTP transport in read-only mode:

```sh
make image-run-http-read-only
```

Use Podman instead of Docker for image commands:

```sh
make image-build CONTAINER_RUNTIME=podman
```

> **Note:** When running in a container with published ports, set `CONTAINER_RUNTIME_MCP_HTTP_ADDR=0.0.0.0:8080` so the server listens on the container network interface.

> **Security warning:** Mounting `/var/run/docker.sock` gives the container access to the host Docker daemon. Use read-only mode where possible and avoid exposing HTTP transport to untrusted networks.

Use the binary as a stdio MCP server command in your MCP client configuration. For example:

```json
{
  "mcpServers": {
    "container-runtime": {
      "command": "/absolute/path/to/container-runtime-mcp"
    }
  }
}
```

For development, an MCP client can also run the module directly:

```json
{
  "mcpServers": {
    "container-runtime": {
      "command": "go",
      "args": ["run", "/absolute/path/to/container-runtime-mcp"]
    }
  }
}
```

## Transport

### stdio (default)

The MCP client starts the process and communicates with it over stdin and stdout.

### HTTP

Set `CONTAINER_RUNTIME_MCP_TRANSPORT=http` to enable streamable HTTP transport. The server listens on a TCP address and serves MCP over HTTP with SSE.

**Security warning:** HTTP transport exposes destructive runtime operations. The default bind address is `127.0.0.1:8080` to restrict access to localhost.

When the server is in write mode (`CONTAINER_RUNTIME_MCP_READ_ONLY=false`), non-local bind addresses such as `0.0.0.0` require `CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN`. Read-only mode can be used without authentication on any bind address.

Do not expose this server to untrusted networks without authentication and read-only mode.

Local HTTP example:

```sh
CONTAINER_RUNTIME_MCP_TRANSPORT=http CONTAINER_RUNTIME_MCP_HTTP_ADDR=127.0.0.1:3000 make run
```

Or use the provided shortcut:

```sh
make run-http
```

The MCP HTTP endpoint will be available at `http://127.0.0.1:3000/mcp` by default.

## Configuration

Configuration is read from environment variables with the `CONTAINER_RUNTIME_` prefix.

| Environment variable | Default | Description |
| --- | --- | --- |
| `CONTAINER_RUNTIME_MCP_SERVER_NAME` | `Container Runtime` | MCP implementation name reported to clients. |
| `CONTAINER_RUNTIME_MCP_SERVER_TITLE` | empty | MCP implementation title reported to clients. |
| `CONTAINER_RUNTIME_MCP_SERVER_VERSION` | `1.0.0` | MCP implementation version reported to clients. |
| `CONTAINER_RUNTIME_MCP_TRANSPORT` | `stdio` | MCP transport. Supported values: `stdio`, `http`. |
| `CONTAINER_RUNTIME_MCP_HTTP_ADDR` | `127.0.0.1:8080` | HTTP listen address when transport is `http`. |
| `CONTAINER_RUNTIME_MCP_HTTP_PATH` | `/mcp` | HTTP path prefix when transport is `http`. Must start with `/`. |
| `CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT` | `30m` | Idle session timeout for HTTP transport. Parsed as a Go duration. |
| `CONTAINER_RUNTIME_MCP_HTTP_READ_TIMEOUT` | `10s` | HTTP server read timeout. Parsed as a Go duration. |
| `CONTAINER_RUNTIME_MCP_HTTP_IDLE_TIMEOUT` | `120s` | HTTP server idle timeout. Parsed as a Go duration. |
| `CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN` | empty | Bearer token for HTTP transport authentication. Required for non-local addresses when `CONTAINER_RUNTIME_MCP_READ_ONLY=false`. When set, every request must include `Authorization: Bearer <token>`. |
| `CONTAINER_RUNTIME_MCP_READ_ONLY` | `false` | When `true`, only read-only tools are registered and mutating runtime operations are rejected at the service layer. |
| `CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT` | `10m` | Timeout applied to runtime operations. Parsed as a Go duration, such as `30s`, `5m`, or `1h`. |
| `CONTAINER_RUNTIME_LOG_LEVEL` | `info` | Log level. Supported values: `debug`, `info`, `warn`, `error`. |

`CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT` must not be negative. A value of `0` disables the operation timeout and leaves calls bounded only by the MCP request context and runtime client behavior.

`CONTAINER_RUNTIME_MCP_HTTP_SESSION_TIMEOUT` must not be negative.

`CONTAINER_RUNTIME_MCP_HTTP_READ_TIMEOUT` must not be negative.

`CONTAINER_RUNTIME_MCP_HTTP_IDLE_TIMEOUT` must not be negative.

`CONTAINER_RUNTIME_MCP_TRANSPORT` must be `stdio` or `http`.

`CONTAINER_RUNTIME_MCP_HTTP_PATH` must start with `/`.

Example:

```sh
CONTAINER_RUNTIME_MCP_REMOTE_OPERATION_TIMEOUT=2m make run
```

Run with HTTP transport:

```sh
CONTAINER_RUNTIME_MCP_TRANSPORT=http CONTAINER_RUNTIME_MCP_HTTP_ADDR=127.0.0.1:3000 make run
```

Run with HTTP transport and bearer token authentication:

```sh
CONTAINER_RUNTIME_MCP_TRANSPORT=http \
CONTAINER_RUNTIME_MCP_HTTP_ADDR=127.0.0.1:3000 \
CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN=<long-random-token> \
make run
```

Run on a non-local address with authentication (required for write mode):

```sh
CONTAINER_RUNTIME_MCP_TRANSPORT=http \
CONTAINER_RUNTIME_MCP_HTTP_ADDR=0.0.0.0:3000 \
CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN=<long-random-token> \
make run
```

Run in read-only mode over HTTP (authentication optional):

```sh
CONTAINER_RUNTIME_MCP_TRANSPORT=http \
CONTAINER_RUNTIME_MCP_HTTP_ADDR=0.0.0.0:3000 \
CONTAINER_RUNTIME_MCP_READ_ONLY=true \
make run
```

When `CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN` is set, the MCP client must include the token in every request:

```http
Authorization: Bearer <long-random-token>
```

## Runtime Access

This server performs real container runtime operations against the daemon available to the process. Some tools are destructive, including container removal, image removal, stopping containers, restarting containers, executing commands in containers, pushing images, tagging images, creating volumes, and removing volumes.

Set `CONTAINER_RUNTIME_MCP_READ_ONLY=true` to disable all mutating tools. In read-only mode, only `ping`, `info`, `version`, `list_containers`, `inspect_container`, `container_logs`, `list_images`, `inspect_image`, `list_volumes`, and `inspect_volume` are available. This is recommended when exposing the HTTP transport to untrusted clients.

Before configuring an MCP client, verify runtime access from the same environment. For Docker:

```sh
docker ps
```

For Podman:

```sh
podman ps
```

If your runtime setup uses environment variables, contexts, or a custom socket path, configure the MCP client so it starts `container-runtime-mcp` with the same environment.

## Tools

### System

| Tool | Description | Arguments |
| --- | --- | --- |
| `ping` | Ping the runtime daemon to check connectivity. | `{}` |
| `info` | Get runtime system information. | `{}` |
| `version` | Get runtime version information. | `{}` |

### Containers

| Tool | Description | Arguments |
| --- | --- | --- |
| `create_container` | Create a new container. | `name` string, `image` string |
| `list_containers` | List containers. | Optional `all` bool, `limit` number, `size` bool, `latest` bool |
| `inspect_container` | Inspect a container. | `name` string |
| `container_logs` | Get logs from a container. | `name` string, optional `stdout` bool, `stderr` bool, `since` string, `timestamps` bool, `tail` string |
| `start_container` | Start a container. | `name` string |
| `stop_container` | Stop a container. | `name` string, optional `signal` string, optional `timeout_seconds` number |
| `restart_container` | Restart a container. | `name` string, optional `signal` string, optional `timeout_seconds` number |
| `remove_container` | Remove a container. | `name` string, optional `force` bool, `remove_volumes` bool, `remove_links` bool |
| `exec_container` | Execute a command in a running container. | `name` string, `command` array of strings, optional `env` array of strings, `working_dir` string, `user` string, `privileged` bool, `tty` bool, optional `stdin` string |

`name` accepts the container name or ID for tools that operate on an existing container.

For `exec_container`, `command` is an array of strings such as `["sh", "-c", "echo hello"]`. When `tty` is `true`, stdout and stderr are merged into a single stream. Use `stdin` to provide input to the command; the write side is closed after sending so commands expecting EOF (for example `cat` or `read`) will not hang.

For `stop_container` and `restart_container`, `timeout_seconds` follows Docker-compatible API semantics: omit it for the runtime default, use `-1` for indefinite wait, or `0` for immediate termination.

### Images

| Tool | Description | Arguments |
| --- | --- | --- |
| `pull_image` | Pull an image from a registry by reference. | `ref` string, optional `all` bool, optional `platform` object |
| `push_image` | Push an image to a registry by reference. | `ref` string, optional `all` bool, optional `platform` object |
| `list_images` | List images. | Optional `all` bool, `shared_size` bool |
| `inspect_image` | Inspect an image. | `ref` string |
| `remove_image` | Remove an image. | `ref` string, optional `force` bool, `prune_children` bool, `platform` object |
| `tag_image` | Tag an image. | `source` string, `target` string |

`ref` accepts an image reference or image ID depending on the runtime operation. Examples include `nginx:latest`, `alpine:3.20`, `registry.example.com/app:latest`, or an image ID.

### Volumes

| Tool | Description | Arguments |
| --- | --- | --- |
| `list_volumes` | List volumes. | Optional `dangling` bool |
| `inspect_volume` | Inspect a volume. | `name` string |
| `create_volume` | Create a volume. | Optional `name` string, `driver` string, `driver_opts` object, `labels` object |
| `remove_volume` | Remove a volume. | `name` string, optional `force` bool |

## Platform Argument Shape

The `platform` argument is an OCI platform object, not a slash-delimited string.

Example for Linux amd64:

```json
{
  "ref": "nginx:latest",
  "platform": {
    "os": "linux",
    "architecture": "amd64"
  }
}
```

Supported fields are:

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `os` | string | Yes | Operating system, such as `linux` or `windows`. |
| `architecture` | string | Yes | CPU architecture, such as `amd64`, `arm64`, or `arm`. |
| `variant` | string | No | CPU variant, such as `v7` for ARMv7. |
| `os.version` | string | No | OS version, mainly used for Windows images. |
| `os.features` | array of strings | No | Required OS features. |

Example for ARMv7:

```json
{
  "ref": "alpine:latest",
  "platform": {
    "os": "linux",
    "architecture": "arm",
    "variant": "v7"
  }
}
```

## Argument Examples

Create a container:

```json
{
  "name": "web",
  "image": "nginx:latest"
}
```

List all containers, including stopped containers:

```json
{
  "all": true
}
```

Stop a container with a SIGTERM grace period:

```json
{
  "name": "web",
  "signal": "SIGTERM",
  "timeout_seconds": 10
}
```

Remove a container and its anonymous volumes:

```json
{
  "name": "web",
  "force": true,
  "remove_volumes": true
}
```

Pull an image for a specific platform:

```json
{
  "ref": "nginx:latest",
  "platform": {
    "os": "linux",
    "architecture": "amd64"
  }
}
```

Tag an image:

```json
{
  "source": "nginx:latest",
  "target": "registry.example.com/nginx:latest"
}
```

### exec_container

Run a simple command in a running container:

```json
{
  "name": "web",
  "command": ["sh", "-c", "pwd && ls -la"]
}
```

Run a command that reads from stdin:

```json
{
  "name": "web",
  "command": ["sh", "-c", "cat > /tmp/message.txt && wc -c /tmp/message.txt"],
  "stdin": "hello from container-runtime-mcp\n"
}
```

Run a command with environment variables and a working directory:

```json
{
  "name": "web",
  "command": ["sh", "-c", "echo $APP_ENV && pwd"],
  "env": ["APP_ENV=dev"],
  "working_dir": "/app"
}
```

### container_logs

Get the last 100 lines from both stdout and stderr:

```json
{
  "name": "web",
  "tail": "100"
}
```

Get recent logs from the last 10 minutes with timestamps:

```json
{
  "name": "web",
  "since": "10m",
  "timestamps": true
}
```

Get only stdout output, all available lines:

```json
{
  "name": "web",
  "stdout": true,
  "stderr": false,
  "tail": "all"
}
```

### Volume examples

Create a named local volume with labels:

```json
{
  "name": "app-data",
  "driver": "local",
  "labels": {
    "app": "web",
    "environment": "dev"
  }
}
```

List only dangling volumes:

```json
{
  "dangling": true
}
```

Inspect a volume:

```json
{
  "name": "app-data"
}
```

Remove a volume with force:

```json
{
  "name": "app-data",
  "force": true
}
```

## Development

Run tests:

```sh
make test
```

Run tests with the race detector:

```sh
make test-race
```

Format code:

```sh
make fmt
```

Run all linting checks (formatting, imports, vet, and golangci-lint):

```sh
make lint
```

Run only golangci-lint:

```sh
make golangci-lint
```

Tidy module dependencies:

```sh
make tidy
```

Generate mocks:

```sh
make generate
```

Run everything (lint, test, and build):

```sh
make all
```

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE).
