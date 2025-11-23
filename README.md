# Proxiff

A Diffy-like HTTP proxy tool written in Go that compares responses from two different servers (newer and current).

[日本語README](README-ja.md) | [日本語クイックスタート](QUICKSTART-ja.md)

## Features

- Forward HTTP requests to two different servers (newer and current)
- Compare responses and log differences
- Pluggable comparison logic through gRPC plugin system
- Return the current server's response to clients
- Built with TDD (Test-Driven Development)

## Installation

```bash
go build -o proxiff ./cmd/proxiff
```

## Usage

Basic usage:

```bash
./proxiff start --newer http://localhost:8082 --current http://localhost:8081
```

### Available Commands

- `proxiff start`: Start the proxy server
- `proxiff version`: Print version information
- `proxiff help`: Show help message

### Flags (for start command)

- `--newer`: URL of the newer server (required)
- `--current`: URL of the current server (required)
- `--port`: Port to listen on (default: 8080)
- `--plugin`: Path to comparator plugin binary (optional)

## Example

1. Start two sample servers:

```bash
# Terminal 1: Current server
go run ./example/servers -port 8081 -version current

# Terminal 2: Newer server
go run ./example/servers -port 8082 -version newer
```

2. Start proxiff:

```bash
# Terminal 3: Proxiff proxy
go run ./cmd/proxiff start --newer http://localhost:8082 --current http://localhost:8081
```

3. Send requests to the proxy:

```bash
# This will return the current server's response
# but log any differences with the newer server
curl http://localhost:8080/
curl http://localhost:8080/api/users
curl http://localhost:8080/api/status
```

## Difference Log Examples

Proxiff uses the [google/go-cmp](https://github.com/google/go-cmp) library to compare responses from two servers. When differences are detected, they are logged in an easy-to-read format.

### Field Addition Difference (/api/users)

When the newer version adds an `email` field:

```
2025/11/23 15:07:51 Difference detected:   &comparator.Response{
  	StatusCode: 200,
  	Headers: http.Header{
- 		"Content-Length": {"78"},
+ 		"Content-Length": {"130"},
  		"Content-Type":   {"application/json"},
  		"Date":           {"Sun, 23 Nov 2025 15:07:51 GMT"},
  	},
  	Body: bytes.Join({
  		`{"users":[{"`,
- 		`id":1,"name":"Alice"},{`,
+ 		`email":"alice@example.com","id":1,"name":"Alice"},{"email":"bob@`,
+ 		`example.com",`,
  		`"id":2,"name":"Bob"}],"version":"`,
- 		"current",
+ 		"newer",
  		"\"}\n",
  	}, ""),
  }
```

### Status Code Difference (/api/status)

When the current server returns HTTP 200 and the newer server returns HTTP 201:

```
2025/11/23 15:07:53 Difference detected:   &comparator.Response{
- 	StatusCode: 200,
+ 	StatusCode: 201,
  	Headers: http.Header{
- 		"Content-Length": {"36"},
+ 		"Content-Length": {"34"},
  		"Content-Type":   {"application/json"},
  		"Date":           {"Sun, 23 Nov 2025 15:07:53 GMT"},
  	},
  	Body: bytes.Join({
  		`{"status":"ok","version":"`,
- 		"current",
+ 		"newer",
  		"\"}\n",
  	}, ""),
  }
```

### No Differences

When responses match completely:

```
2025/11/23 15:04:00 Responses match
```

### Symbol Meanings

- `-` symbol: Content from the current server (removed parts)
- `+` symbol: Content from the newer server (added parts)

## Architecture

Proxiff uses a gRPC-based plugin system for comparison logic. All comparators are implemented as plugins. When no plugin is specified, a builtin SimpleComparator runs in the same process for better performance.

See the [Plugin System](#plugin-system-grpc) section below for details.

## Testing

Run all tests:

```bash
go test ./... -v
```

Run tests for a specific package:

```bash
go test ./plugin/builtin/... -v
go test ./proxy/... -v
```

## Project Structure

```
proxiff/
├── cmd/
│   └── proxiff/        # Main CLI application
│       └── main.go
├── comparator/         # Comparison logic types
│   └── comparator.go   # Interface and type definitions
├── proxy/              # Proxy core functionality
│   ├── proxy.go
│   └── proxy_test.go
├── plugin/             # Plugin system
│   ├── builtin/        # Builtin plugins
│   │   ├── simple.go   # Default SimpleComparator implementation
│   │   └── simple_test.go
│   ├── proto/          # gRPC protobuf definitions
│   ├── interface.go    # Plugin interface
│   ├── grpc_client.go  # gRPC client implementation
│   ├── grpc_server.go  # gRPC server implementation
│   ├── client.go       # Plugin loader
│   └── builtin.go      # Builtin plugin loader
└── example/
    ├── deployment/     # Sample deployment configurations
    │   ├── docker/     # Docker Compose example with Nginx
    │   └── nginx/      # Nginx configuration sample
    ├── servers/        # Sample servers for testing
    │   └── main.go
    └── plugin-status-only/  # Status-only comparison plugin example
        └── main.go
```

## Plugin System (gRPC)

Proxiff uses [HashiCorp go-plugin](https://github.com/hashicorp/go-plugin) for a gRPC-based plugin system, providing:

- **Language Agnostic**: Implement plugins in any language via gRPC
- **Process Isolation**: Plugin crashes don't affect the main process
- **Simple Implementation**: Just implement the `Comparator` interface
- **Battle-Tested**: Used in Terraform, Vault, and other HashiCorp products

### Builtin Plugin (Default)

When no `--plugin` flag is specified, Proxiff automatically uses a builtin SimpleComparator that compares:
- HTTP status codes
- Response headers
- Response bodies

The builtin plugin runs in the same process for better performance while maintaining the same plugin interface.

### Creating a Custom Plugin

Implement the `Comparator` interface:

```go
package main

import (
    "github.com/hashicorp/go-hclog"
    "github.com/hashicorp/go-plugin"
    "github.com/n3xem/proxiff/comparator"
    pluginpkg "github.com/n3xem/proxiff/plugin"
)

type MyComparator struct{}

func (m *MyComparator) Compare(newer, current *comparator.Response) *comparator.Result {
    // Custom comparison logic
    return &comparator.Result{
        Match:      newer.StatusCode == current.StatusCode,
        Newer:      newer,
        Current:    current,
        Difference: "custom logic",
    }
}

func main() {
    logger := hclog.New(&hclog.LoggerOptions{
        Level:  hclog.Error,
        Output: nil,
    })

    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: pluginpkg.Handshake,
        Plugins: map[string]plugin.Plugin{
            "comparator": &pluginpkg.ComparatorPlugin{Impl: &MyComparator{}},
        },
        GRPCServer: plugin.DefaultGRPCServer,
        Logger:     logger,
    })
}
```

Build and use:

```bash
# Build the plugin
go build -o my-plugin ./my-plugin

# Use the plugin
./proxiff start --newer http://localhost:8082 --current http://localhost:8081 --plugin ./my-plugin
```

See `example/plugin-status-only/` for a working example that only compares HTTP status codes.

## Deployment Example with Nginx Mirror Module

The `example/deployment/` directory contains sample configurations showing how Proxiff can be combined with Nginx's [mirror module](https://nginx.org/en/docs/http/ngx_http_mirror_module.html) to compare newer and current versions without affecting production traffic.

### Architecture

```
Client
  ↓
Nginx
  ├─> Production Server (returns response to client)
  └─> Proxiff (mirror, response ignored)
        ├─> Newer Server
        └─> Current Server
              ↓
        Difference detection and logging
```

### Docker Compose Sample Setup

The `example/deployment/docker/` directory provides a ready-to-use example environment:

```bash
cd example/deployment/docker

# Start the sample environment
docker compose up -d

# View logs
docker compose logs -f proxiff

# Send test requests
curl http://localhost:8000/api/users

# Cleanup
docker compose down -v
```

### Run Integration Tests

```bash
cd example/deployment/docker
./test-integration.sh
```

See [example/deployment/docker/README.md](example/deployment/docker/README.md) for details on customizing the setup for your needs.

### Benefits

1. **No Production Impact**: Mirrored traffic responses are ignored, so production is unaffected
2. **Real Traffic**: Compare versions using actual production traffic
3. **Timeout Isolation**: Proxiff timeouts don't affect the production service
4. **Gradual Verification**: Validate new versions before production deployment

## License

MIT
