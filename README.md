# Fetch MCP

[![GoReleaser](https://github.com/JPratama7/fetch-go/actions/workflows/release.yml/badge.svg)](https://github.com/JPratama7/fetch-go/actions/workflows/release.yml)

`fetch-go` is a simple yet powerful command-line tool written in Go that functions as a Mark3-compatible MCP (Model
Context Protocol) server. It exposes a single tool, `fetch`, which takes a URL, retrieves its content, converts it from
HTML to clean, LLM-friendly Markdown, and returns the result.

This tool is designed to be a building block in larger AI systems, allowing language models to easily access and
understand content from the web.

## Features

- **MCP Server**: Implements the MCP standard for seamless integration with Mark3-powered AI agents.
- **HTML to Markdown**: Converts complex web pages into clean, readable Markdown.
- **Advanced Fetching**: Uses a sophisticated client that can handle gzip compression and impersonates a real browser to
  avoid blocking.
- **Automated Releases**: CI/CD pipeline using GoReleaser and GitHub Actions to automatically build and release binaries
  for multiple platforms.

## Usage

As an MCP server, `fetch-go` is intended to be run as a subprocess by an AI agent or a compatible client. It
communicates over `stdio`.

The server provides the following tool:

**`fetch(url: string, startIndex: number = 0, range: number = 5000)`**

- `url` (required): The URL of the web page to fetch.
- `startIndex` (optional): The starting character index of the content to return. Defaults to `0`.
- `range` (optional): The maximum number of characters to return. Defaults to `5000`.

## Installation

Pre-compiled binaries for Linux, macOS, and Windows are automatically generated for each new release. You can download
the latest version from the [**Releases**](https://github.com/JPratama7/fetch-go/releases) page.

## Building from Source

If you prefer to build the project from source, you'll need to have Go installed (version 1.21 or later).

1. **Clone the repository:**
   ```sh
   git clone https://github.com/JPratama7/fetch-go.git
   cd fetch-go
   ```

2. **Build the binary:**
   ```sh
   go build -o fetch-go main.go
   ```

3. **Run the server:**
   ```sh
   ./fetch-go
   ```

## TODO

- [ ] Add comprehensive unit and integration tests.
- [ ] Implement support for `robots.txt` to ensure the tool is a good web citizen.
