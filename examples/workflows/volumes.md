# Volume Workflow

This workflow demonstrates volume operations: create, list, inspect, and remove.

## Prerequisites

- The MCP server is running with write access (`CONTAINER_RUNTIME_MCP_READ_ONLY=false`).
- A Docker-compatible runtime is available.

## Steps

### 1. Create a volume

**Tool:** `create_volume`
**Arguments:**

```json
{
  "name": "mcp-demo-volume",
  "driver": "local",
  "labels": {
    "app": "demo",
    "env": "test"
  }
}
```

**Expected result:** Volume `mcp-demo-volume` is created successfully.

### 2. List volumes

**Tool:** `list_volumes`
**Arguments:**

```json
{}
```

**Expected result:** JSON array of volumes, including `mcp-demo-volume`.

### 3. Inspect the volume

**Tool:** `inspect_volume`
**Arguments:**

```json
{
  "name": "mcp-demo-volume"
}
```

**Expected result:** Volume driver, mount point, and labels.

### 4. Remove the volume

**Tool:** `remove_volume`
**Arguments:**

```json
{
  "name": "mcp-demo-volume",
  "force": false
}
```

**Expected result:** Volume is removed.

## Cleanup

If any step fails and the volume still exists, remove it manually:

```json
{
  "name": "mcp-demo-volume",
  "force": true
}
```
