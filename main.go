package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func fetchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fmt.Printf("URL is %s\n", request.GetString("url", ""))
	fmt.Printf("Start index is %d\n", request.GetInt("startIndex", 0))
	fmt.Printf("Range is %d\n", request.GetInt("range", 0))

	return mcp.NewToolResultText("Fetching..."), nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool
	tool := mcp.NewTool("fetch",
		mcp.WithDescription("Fetch and parse into LLM Friendly Markdown"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to fetch"),
		),
		mcp.WithNumber("startIndex",
			mcp.Description("Start index of the content to fetch"),
			mcp.DefaultNumber(0),
		),
		mcp.WithNumber("range",
			mcp.Description("Range of the content to fetch in characters"),
			mcp.DefaultNumber(5000),
		),
	)

	// Add tool handler
	s.AddTool(tool, fetchHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
