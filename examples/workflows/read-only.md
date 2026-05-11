# Read-Only Workflow

This workflow uses only safe, non-mutating tools. It is safe to run on any Docker-compatible runtime without creating or removing resources.

## Prerequisites

- The MCP server is running in stdio or HTTP mode.
- The client has access to read-only tools.

## Steps

### 1. Ping the runtime

**Tool:** `ping`
**Arguments:**

```json
{}
```

**Expected result:** The runtime daemon responds with `OK` or similar.

### 2. Get runtime version

**Tool:** `version`
**Arguments:**

```json
{}
```

**Expected result:** Runtime version details (for example, Docker Engine version).

### 3. List containers

**Tool:** `list_containers`
**Arguments:**

```json
{
  "all": true
}
```

**Expected result:** JSON array of containers, including stopped ones.

### 4. Inspect a container

If at least one container exists, inspect it by name or ID.

**Tool:** `inspect_container`
**Arguments:**

```json
{
  "name": "web"
}
```

**Expected result:** Detailed container configuration and state.

### 5. Get container logs

**Tool:** `container_logs`
**Arguments:**

```json
{
  "name": "web",
  "tail": "100"
}
```

**Expected result:** Up to the last 100 lines of logs.

### 6. List images

**Tool:** `list_images`
**Arguments:**

```json
{}
```

**Expected result:** JSON array of images.

### 7. Inspect an image

**Tool:** `inspect_image`
**Arguments:**

```json
{
  "ref": "nginx:latest"
}
```

**Expected result:** Image manifest and configuration details.

### 8. List volumes

**Tool:** `list_volumes`
**Arguments:**

```json
{}
```

**Expected result:** JSON array of volumes.

### 9. Inspect a volume

If a volume exists, inspect it by name.

**Tool:** `inspect_volume`
**Arguments:**

```json
{
  "name": "app-data"
}
```

**Expected result:** Volume driver, mount point, and labels.

## Cleanup

No resources are created, so no cleanup is required.
