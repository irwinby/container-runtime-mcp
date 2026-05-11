# Container Lifecycle Workflow

This workflow demonstrates a full container lifecycle: pull an image, create a container, start it, run a command, fetch logs, stop it, and remove it.

## Prerequisites

- The MCP server is running with write access (`CONTAINER_RUNTIME_MCP_READ_ONLY=false`).
- A Docker-compatible runtime is available.
- The runtime can pull images from Docker Hub.

## Steps

### 1. Pull an image

**Tool:** `pull_image`
**Arguments:**

```json
{
  "ref": "alpine:3.20"
}
```

**Expected result:** Image is pulled successfully.

### 2. Create a container

**Tool:** `create_container`
**Arguments:**

```json
{
  "name": "mcp-demo-container",
  "image": "alpine:3.20"
}
```

**Expected result:** Container `mcp-demo-container` is created.

### 3. Start the container

**Tool:** `start_container`
**Arguments:**

```json
{
  "name": "mcp-demo-container"
}
```

**Expected result:** Container starts successfully.

### 4. Execute a command

**Tool:** `exec_container`
**Arguments:**

```json
{
  "name": "mcp-demo-container",
  "command": ["sh", "-c", "echo hello from mcp"]
}
```

**Expected result:** Output contains `hello from mcp`.

### 5. Get container logs

**Tool:** `container_logs`
**Arguments:**

```json
{
  "name": "mcp-demo-container",
  "tail": "all"
}
```

**Expected result:** Combined stdout and stderr logs.

### 6. Stop the container

**Tool:** `stop_container`
**Arguments:**

```json
{
  "name": "mcp-demo-container",
  "timeout_seconds": 10
}
```

**Expected result:** Container stops gracefully.

### 7. Remove the container

**Tool:** `remove_container`
**Arguments:**

```json
{
  "name": "mcp-demo-container",
  "force": false,
  "remove_volumes": true
}
```

**Expected result:** Container and any anonymous volumes are removed.

## Cleanup

If any step fails, run the cleanup sequence manually:

1. Stop the container if it is running.
2. Remove the container.
3. Remove the pulled image if you no longer need it:

```json
{
  "ref": "alpine:3.20"
}
```

Use the `remove_image` tool.
