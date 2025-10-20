# Portosync MCP Server

![MCP Server](https://badge.mcpx.dev?type=server)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/chickenzord/portosync)
[![Go Report Card](https://goreportcard.com/badge/github.com/chickenzord/portosync)](https://goreportcard.com/report/github.com/chickenzord/portosync)
[![codecov](https://codecov.io/github/chickenzord/portosync/graph/badge.svg?token=wmr7FbJadF)](https://codecov.io/github/chickenzord/portosync)
![Go Build](https://github.com/chickenzord/portosync/actions/workflows/go.yml/badge.svg?branch=main)
![Docker Build](https://github.com/chickenzord/portosync/actions/workflows/docker.yml/badge.svg?branch=main)
![Code License](https://img.shields.io/github/license/chickenzord/portosync)

A Model Context Protocol (MCP) server that provides seamless integration with financial portfolio data from multiple sources. This server enables AI assistants like Claude to interact with your financial portfolio balances through standardized MCP tools, with initial support for KSEI (Indonesian central custodian holding data across all securities).

## Features

- 🏦 **KSEI Integration** - Fetch portfolio data from KSEI AKSES using goksei library
- 🚀 **Dual Mode Support** - Run via stdio (MCP standard) or HTTP server
- 🐳 **Docker Ready** - Multi-stage Alpine-based container
- 🔒 **Secure** - Runs as non-root user with minimal dependencies
- 📝 **Self-Describing** - Rich tool descriptions with clear intents, parameter schemas, and behavior annotations

Upcoming features

- 📊 **Portfolio Tracking** - Store and query financial balances as time series data
- 🔄 **Background Jobs** - Periodic data fetching with jitter for reliability
- 🗄️ **SQLite Database** - Lightweight, self-contained database storage

## Tools Available

The server provides self-describing tools with comprehensive metadata including clear descriptions, parameter schemas, and behavior annotations (read-only, idempotent, open-world hints). This enables AI assistants to understand tool capabilities and make informed decisions about tool usage.

### `get_portfolio`
**Title:** Get Portfolio Balances

Retrieves current investment portfolio balances from KSEI (Indonesian Central Securities Depository) accounts. Returns detailed information about holdings including asset symbols, names, quantities, values, and currencies.

**Parameters:**
- `account_names` (array of strings, optional): List of specific account names to retrieve portfolio data from. Each name must match a configured account. If empty or omitted, returns portfolio data from all configured accounts. Use the `list_account_names` tool to discover available account names.

**Returns:** Array of balance objects with fields:
- `source_type`, `source_account`: Source information
- `asset_symbol`, `asset_name`, `asset_type`, `asset_sub_type`: Asset identification
- `units_amount`, `units_value`, `units_currency`: Quantity and value data

**Behavior Annotations:**
- ✓ Read-only (does not modify data)
- ✗ Non-idempotent (data changes daily during settlement hours)
- ✗ Closed-world (accesses only your private configured accounts)

### `list_account_names`
**Title:** List Available Account Names

Lists all account names that are currently configured in the server. Each account name represents a separate KSEI AKSES account connection.

**Parameters:**
- None (returns all configured accounts)

**Returns:**
- `account_names` (array of strings): List of configured account names that can be used with `get_portfolio`

**Behavior Annotations:**
- ✓ Read-only (does not modify data)
- ✓ Idempotent (always returns same list)
- ✗ Closed-world (queries internal server configuration)

## Installation

### Prerequisites

- KSEI AKSES account credentials
- Go 1.25+ (if building from source)

### Option 1: Docker (Recommended)

```bash
# Pull the image
docker pull ghcr.io/chickenzord/portosync:latest

# Run in stdio mode (for MCP clients)
docker run --rm -i \
  -e KSEI_ACCOUNTS="personal:email1@example.com:pass1,business:email2@example.com:pass2" \
  -e KSEI_AUTH_CACHE_DIR="/tmp/ksei_cache" \
  ghcr.io/chickenzord/portosync:latest mcp-stdio

# Run in HTTP mode
docker run --rm -p 8080:8080 \
  -e KSEI_ACCOUNTS="personal:user1@example.com:pass1,business:user2@example.com:pass2" \
  -e KSEI_AUTH_CACHE_DIR="/tmp/ksei_cache" \
  -e BIND_ADDR=":8080" \
  ghcr.io/chickenzord/portosync:latest mcp-http
```

### Option 2: Go Install

```bash
go install github.com/chickenzord/portosync/cmd/portosync@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/chickenzord/portosync.git
cd portosync
go build -o portosync ./cmd/portosync
```

## Configuration

### Environment Variables

- `KSEI_ACCOUNTS` (required): KSEI account configurations in format "name:username:password,name2:username2:password2"
- `KSEI_AUTH_CACHE_DIR` (optional): Directory to cache KSEI authentication tokens (default: temp directory)
- `KSEI_PLAIN_PASSWORD` (optional): Set to "false" to use encrypted passwords (default: true)
- `BIND_ADDR` (optional): HTTP server bind address (default: ":8080")

### KSEI Account Configuration

The `KSEI_ACCOUNTS` environment variable should contain comma-separated account configurations:

```
KSEI_ACCOUNTS="personal:your.email@example.com:yourpassword,business:business.email@example.com:businesspassword"
```

Each account configuration follows the format: `name:username:password`
- `name`: A friendly name for the account (e.g., "personal", "business", "family")
- `username`: Your KSEI AKSES email/username
- `password`: Your KSEI AKSES password

## MCP Client Configuration

### Claude Desktop

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "portosync": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-e", "KSEI_ACCOUNTS=personal:your.email@example.com:yourpassword",
        "-e", "KSEI_AUTH_CACHE_DIR=/tmp/ksei_cache",
        "ghcr.io/chickenzord/portosync:latest",
        "mcp-stdio"
      ]
    }
  }
}
```

Or if using a local binary:

```json
{
  "mcpServers": {
    "portosync": {
      "command": "portosync",
      "args": ["mcp-stdio"],
      "env": {
        "KSEI_ACCOUNTS": "personal:your.email@example.com:yourpassword",
        "KSEI_AUTH_CACHE_DIR": "/tmp/ksei_cache"
      }
    }
  }
}
```

### VS Code with MCP Extension

```json
{
  "mcp.servers": [
    {
      "name": "portosync",
      "command": "portosync",
      "args": ["mcp-stdio"],
      "env": {
        "KSEI_ACCOUNTS": "personal:your.email@example.com:yourpassword",
        "KSEI_AUTH_CACHE_DIR": "/tmp/ksei_cache"
      }
    }
  ]
}
```

## Usage Examples

Once configured with an MCP client, you can use natural language commands:

- *"Show me my current portfolio balances"*
- *"What's the history of my BBCA holdings in KSEI?"*
- *"Show balances for my personal account only"*
- *"What stocks do I own and their current values?"*
- *"What accounts are available?"*
- *"List all my configured accounts"*

## Development

### Running Locally

```bash
# Set environment variables
export KSEI_ACCOUNTS="personal:your.email@example.com:yourpassword"
export KSEI_AUTH_CACHE_DIR="/tmp/ksei_cache"

# Run in stdio mode
go run ./cmd/portosync mcp-stdio

# Run in HTTP mode
go run ./cmd/portosync mcp-http
```

### Testing

```bash
# Run tests
go test ./...

# Test with a real KSEI account
go run ./cmd/portosync mcp-stdio < test-requests.json
```

### Building

```bash
# Build binary
go build -o portosync ./cmd/portosync

# Build Docker image
docker build -t portosync .
```

## Docker Compose Example

```yaml
version: '3.8'
services:
  portosync:
    image: ghcr.io/chickenzord/portosync:latest
    environment:
      - KSEI_ACCOUNTS=personal:your.email@example.com:yourpassword
      - KSEI_AUTH_CACHE_DIR=/app/cache
      - BIND_ADDR=:8080
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./cache:/app/cache
    command: ["mcp-http"]
```

## API Reference

The HTTP mode exposes MCP endpoints at:
- `POST /mcp/v1/initialize` - Initialize MCP session
- `POST /mcp/v1/tools/list` - List available tools  
- `POST /mcp/v1/tools/call` - Call a tool

See the [MCP specification](https://spec.modelcontextprotocol.io/) for detailed API documentation.

## Troubleshooting

### Common Issues

**Authentication Error**: Verify your KSEI credentials in `KSEI_ACCOUNTS` are correct and the account is active.

**Connection Error**: Ensure KSEI AKSES is accessible and your network allows connections to their servers.

**Cache Issues**: Clear the `KSEI_AUTH_CACHE_DIR` directory if experiencing persistent authentication problems.

### Logs

Container logs:
```bash
docker logs <container-id>
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)  
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [goksei](https://github.com/chickenzord/goksei) - Go library for KSEI AKSES integration
- [MCP Specification](https://spec.modelcontextprotocol.io/) - Model Context Protocol specification
- [Claude Desktop](https://claude.ai/desktop) - AI assistant with MCP support

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/chickenzord/portosync/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/chickenzord/portosync/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/chickenzord/portosync/wiki)

---

Built with ❤️ in Indonesia for better financial literacy and MCP ecosystem
