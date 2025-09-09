package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chickenzord/portosync/internal/server"
	"github.com/chickenzord/portosync/internal/version"
)

func main() {
	ctx := context.Background()

	bindAddr := os.Getenv("BIND_ADDR")
	kseiAccounts := parseKseiAccountsString(os.Getenv("KSEI_ACCOUNTS"))
	kseiPlainPassword := os.Getenv("KSEI_PLAIN_PASSWORD") != "false" // default to true
	kseiAuthCacheDir := os.Getenv("KSEI_AUTH_CACHE_DIR")

	if kseiAuthCacheDir == "" {
		dir, err := os.MkdirTemp("", "portosync_ksei_auth")
		if err != nil {
			panic(err)
		}

		kseiAuthCacheDir = dir
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands: mcp-stdio, mcp-http, version\n")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "version" {
		versionInfo := version.Get()
		fmt.Println(versionInfo.String())
		os.Exit(0)
	}

	if bindAddr == "" {
		bindAddr = ":8080"
	}

	for name := range kseiAccounts {
		fmt.Printf("Loaded KSEI account: %s\n", name)
	}

	mcpServer := server.NewMCP(kseiAccounts, kseiPlainPassword, kseiAuthCacheDir)

	switch command {
	case "mcp-http":
		fmt.Printf("Starting portosync HTTP server on %s\n", bindAddr)

		if err := mcpServer.RunHTTP(ctx, bindAddr); err != nil {
			fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n", err)
			os.Exit(1)
		}
	case "mcp-stdio":
		if err := mcpServer.RunStdio(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error running MCP server: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command %s\n", command)
		os.Exit(1)
	}
}
