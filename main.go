package main

import (
	"compress/gzip"
	"context"
	"fetch-go/robot"
	"fmt"
	"io"
	"net/url"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const DEFAULT_USER_AGENT = "ModelContextProtocol/1.0 (User-Specified; +https://github.com/JPratama7/fetch-go)"

type FetchArgs struct {
	Url        string `json:"url"`
	StartIndex int    `json:"startIndex"`
	Range      int    `json:"range"`
}

func decoder(src io.Reader, dest io.Writer, encoding string) error {
	switch encoding {
	case "gzip":
		gzipReader, err := gzip.NewReader(src)
		if err != nil {
			return err
		}
		_, err = io.Copy(dest, gzipReader)
		return err

	default:
		_, err := io.Copy(dest, src)
		return err
	}
}

func fetchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	payload := new(FetchArgs)
	err := request.BindArguments(payload)
	if err != nil {
		return nil, err
	}

	// Create a Surf client with advanced features
	surfClient := NewClient()

	parsedUrl, err := url.Parse(payload.Url)
	if err != nil {
		return mcp.NewToolResultText("Invalid Url"), nil
	}

	robotRules, err := robot.FromURL(surfClient, parsedUrl)
	if err == nil && !robotRules.IsAllowed(payload.Url) {
		return mcp.NewToolResultText("Url is not allowed by robots.txt"), nil
	}

	resp, err := surfClient.Get(payload.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawContent := new(strings.Builder)

	err = decoder(resp.Body, rawContent, resp.Header.Get("Content-Encoding"))

	if err != nil {
		return nil, err
	}

	res, err := htmltomarkdown.ConvertString(rawContent.String())
	return mcp.NewToolResultText(res), nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Fetch MCP",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool with explicit name
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

	// Add debug logging for tool registration
	fmt.Printf("Registering tool: %s\n", "fetch")

	// Add tool handler
	s.AddTool(tool, fetchHandler)

	// Start the stdio server
	fmt.Println("Starting MCP server...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
