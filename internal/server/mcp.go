package server

import (
	"context"
	"maps"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chickenzord/goksei"
	"github.com/chickenzord/portosync/internal/version"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCP wraps the MCP SDK server
type MCP struct {
	kseiClients map[string]*goksei.Client
	mcpServer   *mcp.Server
}

// selectKseiClients get clients by multiple names,
// if empty or nil, it will return all clients
func (m *MCP) selectKseiClients(names []string) map[string]*goksei.Client {
	if len(names) == 0 {
		clients := make(map[string]*goksei.Client)
		maps.Copy(clients, m.kseiClients)

		return clients
	}

	// Return clients matching the provided names
	clients := make(map[string]*goksei.Client)

	for _, name := range names {
		if client, ok := m.kseiClients[name]; ok {
			clients[name] = client
		}
	}

	return clients
}

func (m *MCP) getKseiClientNames() []string {
	names := make([]string, 0, len(m.kseiClients))
	for name := range m.kseiClients {
		names = append(names, name)
	}

	return names
}

func (m *MCP) RunHTTP(ctx context.Context, bindAddress string) error {
	httpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return m.mcpServer
	}, nil)

	return http.ListenAndServe(bindAddress, httpHandler)
}

func (m *MCP) RunStdio(ctx context.Context) error {
	return m.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// NewMCP creates a new MCP server using the official MCP Go SDK
func NewMCP(accounts map[string]string, plainPassword bool) *MCP {
	tempDir, err := os.MkdirTemp("", "portosync")
	if err != nil {
		panic(err)
	}

	authStore, err := goksei.NewFileAuthStore(tempDir)
	if err != nil {
		panic(err)
	}

	gokseiClients := make(map[string]*goksei.Client, len(accounts))
	for username, password := range accounts {
		gokseiClients[username] = goksei.NewClient(goksei.ClientOpts{
			Username:      username,
			Password:      password,
			PlainPassword: plainPassword,
			Timeout:       1 * time.Minute,
			AuthStore:     authStore,
		})
	}

	s := &MCP{
		kseiClients: gokseiClients,
	}

	// Create MCP server with implementation info
	versionInfo := version.Get()
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "portosync",
		Version: versionInfo.Version,
		Title:   "PortoSync MCP Server",
	}, nil)

	// Add get_portfolio tool
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_portfolio",
		Description: "Get consolidated portfolio from Portosync",
	}, s.handleGetPortfolio)

	s.mcpServer = mcpServer

	return s
}

// handleGetPortfolio handles the get_portfolio MCP tool
func (m *MCP) handleGetPortfolio(ctx context.Context, req *mcp.CallToolRequest, args GetPortfolioArgs) (*mcp.CallToolResult, GetPortfolioResult, error) {
	result := GetPortfolioResult{}

	clients := m.selectKseiClients(args.AccountNames)
	if len(clients) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Selected accounts not found, available accounts are " + strings.Join(m.getKseiClientNames(), ", "),
				},
			},
			IsError: true,
		}, result, nil
	}

	balances, err := getAllBalances(clients)
	if err != nil {
		return nil, result, err
	}

	if len(balances) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "No portfolio balances found for selected accounts",
				},
			},
			IsError: true,
		}, result, nil
	}

	result.Balances = balances

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: result.Description(),
			},
		},
	}, result, nil
}
