package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {

	ctx := context.Background()

	mcpClient, err := client.NewStdioMCPClient(
		"go",
		[]string{}, // Empty ENV
		"run",
		"../mcp-stdio-server/main.go",
	)
	if err != nil {
		panic(err)
	}
	defer mcpClient.Close()

	// Create context with timeout
	//ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//defer cancel()

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "budgie-client",
		Version: "1.0.0",
	}

	initResult, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	fmt.Printf(
		"Initialized with server: %s %s\n\n",
		initResult.ServerInfo.Name,
		initResult.ServerInfo.Version,
	)

	// List Tools
	fmt.Println("Listing available tools...")
	toolsRequest := mcp.ListToolsRequest{}
	tools, err := mcpClient.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}
	fmt.Println()

	// NOTE: Call the MCP tool with the arguments
	request := mcp.CallToolRequest{}
	request.Params.Name = "say_hello"
	request.Params.Arguments =  map[string]string{
		"name": "Alice",
	}

	toolResponse, _ := mcpClient.CallTool(ctx, request)
	result := toolResponse.Content[0].(mcp.TextContent).Text

	fmt.Printf("🛠️ Tool Response: %s\n", result)

}

// TODO" WithMCPStdioClient
