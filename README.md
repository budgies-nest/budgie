# budgie

> Only a simple Agent pattern with the [OpenAI Golang SDK](https://github.com/openai/openai-go) and [mcp-go](https://github.com/mark3labs/mcp-go).

## 📚 Documentation

### Core Features
- **[Completion Handlers](./docs/COMPLETION_HANDLERS.md)** - Hook into the completion lifecycle to modify inputs and outputs
- **[Agent Creation](./docs/01-create-an-agent.md)** - Getting started with creating agents
- **[Chat Completion](./docs/02-chat-completion.md)** - Basic text completions
- **[Stream Completion](./docs/03-chat-stream-completion.md)** - Streaming responses
- **[Tool Calls](./docs/04-simple-tool-call.md)** - Function calling capabilities
- **[Multiple Tool Calls](./docs/05-multiple-tool-calls.md)** - Handling multiple tools
- **[Alternative Tool Detection](./docs/06-alternative-tool-call-detection.md)** - Alternative tool calling method
- **[MCP Integration](./docs/07-mcp-stdio-tool-call.md)** - Model Control Protocol support
- **[RAG Memory](./docs/09-rag-in-memory.md)** - Retrieval-Augmented Generation
- **[Logging](./docs/14-logging.md)** - Comprehensive logging system

### Advanced Features
- **[Agent as MCP Server](./docs/12-agent-as-a-mcp-server.md)** - Expose agents as MCP servers
- **[Agent as REST API](./docs/13-agent-as-a-rest-api-server.md)** - HTTP API endpoints
- **[Vector Stores](./docs/10-in-memory-vectore-store.md)** - In-memory vector storage
- **[Text Chunking](./docs/11-chunking.md)** - Document processing utilities

## 🍳 Cookbook Examples

See the [cookbook](./cookbook/) directory for practical examples:
- **[Response Modification](./cookbook/14-response-modification/)** - Complete guide to using completion handlers
- **[Basic Chat](./cookbook/01-chat-completion/)** - Simple chat completion
- **[Tool Usage](./cookbook/03-one-tool-call/)** - Function calling examples
- **[MCP Integration](./cookbook/07-mcp-stdio/)** - MCP client examples
- **[RAG Implementation](./cookbook/09-rag-memory/)** - RAG memory examples