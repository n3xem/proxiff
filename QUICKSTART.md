# Proxiff - Quick Start Guide

This guide will help you get started with Proxiff in 5 minutes.

[日本語版クイックスタートガイド](QUICKSTART-ja.md)

## What is Proxiff?

Proxiff is a Diffy-like HTTP proxy that compares responses from two different servers (newer and current) and logs any differences while returning the current server's response to clients.

## Quick Start

### 1. Run the Demo

The easiest way to see Proxiff in action:

```bash
# Make the demo script executable (if not already)
chmod +x example/run-demo.sh

# Run the demo
./example/run-demo.sh
```

This will start:
- Current server on port 8081
- Newer server on port 8082
- Proxiff proxy on port 8080

### 2. Test It

In another terminal, try these requests:

```bash
# Basic endpoint - will show version differences
curl http://localhost:8080/

# Users endpoint - newer version has email field
curl http://localhost:8080/api/users

# Status endpoint - different status codes (200 vs 201)
curl http://localhost:8080/api/status
```

### 3. Check the Logs

The Proxiff proxy will log all differences it detects:

```
2025/11/23 14:35:35 Difference detected: Body differs: newer="...", current="..."
2025/11/23 14:35:50 Difference detected: Body differs: newer="...", current="..."
2025/11/23 14:35:51 Difference detected: Status code differs: newer=201, current=200
```

## Manual Setup

If you prefer to run things manually:

### 1. Build

```bash
go build -o proxiff ./cmd/proxiff
go build -o sample-server ./example/servers
```

### 2. Start Servers

```bash
# Terminal 1: Current version
./sample-server -port 8081 -version current

# Terminal 2: Newer version
./sample-server -port 8082 -version newer

# Terminal 3: Proxiff
./proxiff -newer http://localhost:8082 -current http://localhost:8081 -port 8080
```

### 3. Send Requests

```bash
curl http://localhost:8080/
```

## Using with Your Own Servers

Replace the sample servers with your actual services:

```bash
./proxiff \
  -newer http://your-new-service.example.com \
  -current http://your-current-service.example.com \
  -port 8080
```

Then point your clients to `http://localhost:8080` instead of your current service.

## Custom Comparison Logic

Want to customize how responses are compared? Create a plugin:

```bash
# Build a plugin (example: status-only comparison)
go build -o plugin-status-only ./example/plugin-status-only

# Use the plugin
./proxiff \
  -newer http://localhost:8082 \
  -current http://localhost:8081 \
  -port 8080 \
  -plugin ./plugin-status-only
```

## Testing

Run all tests:

```bash
go test ./... -v
```

## Next Steps

- Read [README.md](README.md) for detailed documentation
- Check out [example/plugin-status-only/main.go](example/plugin-status-only/main.go) to learn how to write custom comparison plugins
- Implement your own plugin for your specific needs (e.g., ignore timestamps, compare only specific JSON fields, etc.)

## Common Use Cases

1. **Canary Deployments**: Compare production and canary versions
2. **Migration Testing**: Ensure new implementation matches old behavior
3. **A/B Testing**: Compare different algorithm implementations
4. **Regression Testing**: Detect unexpected changes in API responses
