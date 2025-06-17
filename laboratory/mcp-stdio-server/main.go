package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-stdio-tests",
		"0.0.0",
	)
	// Add a tool
	sayHello := mcp.NewTool("say_hello",
		mcp.WithDescription(`Say hello to the given person name`),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the person to greet."),
		),
	)
	s.AddTool(sayHello, sayHelloHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		log.Fatalln("Failed to start server:", err)
		return
	}

}

func sayHelloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	// Check if the argument is provided
	if len(args) == 0 {
		return mcp.NewToolResultText("Hello John Doe"), nil
	}
	var content = "Hello John Doe"
	if name, ok := args["name"]; ok {
		content = "Hello " + name.(string) // use the provided name
	}
	return mcp.NewToolResultText(content), nil

}
