# MCP Gateway (STDIO)

Use the Docker MCP Gateway from go

Repository: https://github.com/docker/mcp-gateway

```bash
# Run the MCP gateway (stdio) - the one installed with docker desktop
docker mcp gateway run
```

Then run:

```bash
go run main.go
```

## Settings

Connection:
```json
{
    "mcpServers": {
        "MCP_DOCKER": {
            "command": "docker",
            "args": ["mcp", "gateway", "run"]
        }
    }
}
```

DuckDuckGo MCP server:
- `fetch_content`: fetch and parse content from a webpage URL.
- `search`: search DuckDuckGo and return formatted results.
