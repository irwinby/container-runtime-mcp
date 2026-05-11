# Image Workflow

This workflow demonstrates common image operations: pull, list, inspect, tag, and remove.

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

### 2. List images

**Tool:** `list_images`
**Arguments:**

```json
{}
```

**Expected result:** JSON array of images, including `alpine:3.20`.

### 3. Inspect the image

**Tool:** `inspect_image`
**Arguments:**

```json
{
  "ref": "alpine:3.20"
}
```

**Expected result:** Image details such as architecture, size, and layers.

### 4. Tag the image

**Tool:** `tag_image`
**Arguments:**

```json
{
  "source": "alpine:3.20",
  "target": "local/alpine:3.20"
}
```

**Expected result:** New tag `local/alpine:3.20` is created.

### 5. Remove the original image reference

**Tool:** `remove_image`
**Arguments:**

```json
{
  "ref": "alpine:3.20",
  "force": false,
  "prune_children": false
}
```

**Expected result:** The `alpine:3.20` tag is removed. The `local/alpine:3.20` tag still points to the same image ID, so the underlying layers are preserved.

### 6. Remove the tagged image

**Tool:** `remove_image`
**Arguments:**

```json
{
  "ref": "local/alpine:3.20",
  "force": false,
  "prune_children": false
}
```

**Expected result:** The image is fully removed because no remaining tags point to it.

## Cleanup

If any step fails, run `remove_image` for any created references:

```json
{
  "ref": "alpine:3.20"
}
```

```json
{
  "ref": "local/alpine:3.20"
}
```
