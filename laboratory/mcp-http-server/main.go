package main

import (
	"context"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-http-tests",
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

	// Start the HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	).Start(":" + httpPort)
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
