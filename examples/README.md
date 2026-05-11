# Examples

This directory contains practical, copy-paste-ready examples for using `container-runtime-mcp`.

## Prerequisites

- A running Docker-compatible container runtime accessible from your environment.
- Go 1.26.2 or newer if building from source.
- The `container-runtime-mcp` binary built or the module source available.

Before running any example, confirm runtime access:

```sh
docker ps
```

## MCP Client Configuration

The most common way to use this server is through an MCP client that starts the process over **stdio**.

- [mcp-client.json](mcp-client.json) — standard stdio config.
- [mcp-client-read-only.json](mcp-client-read-only.json) — stdio config with read-only mode enabled.

Replace `/absolute/path/to/container-runtime-mcp` with the actual path to the binary or module root.

## Running the Server

### stdio (default)

```sh
make run
```

Or run directly:

```sh
go run .
```

### HTTP transport

Local HTTP without authentication (localhost only):

```sh
source examples/http-local.env
make run-http
```

Read-only HTTP on a non-local address (safe without authentication):

```sh
source examples/http-read-only.env
make run-http
```

Authenticated HTTP with bearer token:

```sh
source examples/http-authenticated.env
make run-http
```

## Workflow Examples

Each workflow is a realistic sequence of MCP tool calls. You can run them via an MCP client or by sending requests to the HTTP endpoint.

- [read-only.md](workflows/read-only.md) — safe, non-mutating operations that work on any runtime.
- [container-lifecycle.md](workflows/container-lifecycle.md) — create, start, exec, stop, and remove a container.
- [images.md](workflows/images.md) — pull, list, inspect, tag, and remove an image.
- [volumes.md](workflows/volumes.md) — create, list, inspect, and remove a volume.

## Verifying Examples

1. Build the project to confirm it compiles:

```sh
make build
```

2. Run tests:

```sh
make test
```

3. Confirm the stdio server starts without error:

```sh
make run
# Stop with Ctrl-C
```

4. Confirm the HTTP server starts without error:

```sh
make run-http
# Stop with Ctrl-C
```

## Notes

- Mutating examples (`container-lifecycle`, `images`, `volumes`) create real resources. Each workflow includes cleanup steps.
- Use `CONTAINER_RUNTIME_MCP_READ_ONLY=true` to prevent accidental mutations.
- When using HTTP transport on a non-local address, always set `CONTAINER_RUNTIME_MCP_HTTP_AUTH_TOKEN` unless read-only mode is enabled.
